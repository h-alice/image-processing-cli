package main

import (
	"bytes"
	"flag"
	"fmt"
	"imagetools/config"
	"io"
	"log"
	"os"
	"path/filepath"

	op "imagecore/operation"
)

func ProcessFile(profile config.ProcessProfileConfig, out io.Writer, in io.Reader) error {

	// Procedure: Decode -> image ops -> encode -> segment ops -> write out

	working_image, err := op.CreateImageFromReader(in)
	if err != nil {
		log.Printf("[x] Error while creating image: %v", err)
		return err
	}

	// Do crop
	//img, err = profile.DoCrop(img)
	//if err != nil {
	//	log.Printf("[x] Error while cropping image: %v", err)
	//	return err
	//}

	output_image := working_image.
		Then(op.Decode()).
		ThenIf(profile.Resize.Factor != 0.0, op.ResizeImageByFactor(profile.Resize.Algorithm, profile.Resize.Factor)).
		ThenIf((profile.Resize.Factor == 0.0) && (profile.Resize.Width != 0), op.ResizeImageByWidth(profile.Resize.Algorithm, profile.Resize.Width)).
		ThenIf((profile.Resize.Factor == 0.0) && (profile.Resize.Width == 0) && (profile.Resize.Height != 0), op.ResizeImageByHeight(profile.Resize.Algorithm, profile.Resize.Height)).
		Then(op.Encode(profile.Output.Format, (*op.EncoderOption)(profile.Output.Options))).
		ThenIf((profile.Output.Format == "jpeg" || profile.Output.Format == "jpg") && profile.ICC != "", op.EmbedProfile(profile.ICC)) // Supports jpeg only for now, will be extended to other formats.

	if output_image.LastError() != nil {
		log.Printf("[x] Error while processing image: %v", output_image.LastError())
		return output_image.LastError()
	}

	// Write output to writer.
	output_image.Then(op.WriteImageToWriter(out))
	if output_image.LastError() != nil {
		log.Printf("[x] Error while writing image: %v", output_image.LastError())
		return output_image.LastError()
	}

	return nil
}

func main() {

	// Parse command line.
	ptr_config_path := flag.String("config", "", "config file path")
	flag.Parse()

	// If path not specified, load defaultprofile from home directory.
	config_path := *ptr_config_path
	if *ptr_config_path == "" {

		log.Printf("[!] Using default config file.")

		var err error
		config_path, err = getProfileFromHomeDir("default", true)
		if err != nil {
			log.Fatalf("[x] Error while getting default config: %s\n", err)
		}
	}

	conf, err := config.LoadConfigFromFile(config_path)
	if err != nil {
		log.Fatalf("[x] Cannot load config file: %s\n", err)
	}

	for _, f := range flag.Args() { // Iterate through input images.

		raw_bytes, err := os.ReadFile(f)
		if err != nil {
			log.Printf("[x] Error while reading file: %s\n", err)
			continue
		}

		tasks := make(chan struct{}, len(conf.Profiles))
		for _, pf := range conf.Profiles { // Apply all profile to input image.

			go func(pf config.ProcessProfileConfig) {

				// TODO: Output dir.
				output_dir := filepath.Dir(f)
				if pf.Output == nil {
					log.Println("[x] No output section.")
					return
				}

				outfile_name := pf.Output.GenerateFileName(f)

				outputbuf := bytes.NewBuffer([]byte{})
				output_full_path := filepath.Join(output_dir, outfile_name)

				err = ProcessFile(pf, outputbuf, bytes.NewBuffer(raw_bytes))
				if err != nil {
					log.Printf("[x] An error occurred while processing image: %s\n", err)
					return
				}

				ofp, err := os.OpenFile(output_full_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
				defer func() {
					ofp.Close()
				}()

				if err != nil {
					log.Fatalln(err)
				}
				fmt.Printf("Writing output %s -> %s [%s]\n", f, output_full_path, pf.ProfileName)

				output_length := outputbuf.Len()

				written, err := io.Copy(ofp, outputbuf)

				if err != nil {
					log.Fatalln(err)
				} else if written != int64(output_length) {
					err = fmt.Errorf("written length mismatch")
					log.Fatal(err)
				}

				tasks <- struct{}{}
			}(pf)

		}

		for i := 0; i < len(conf.Profiles); i++ {
			log.Println(i)
			<-tasks
		}
	}
}
