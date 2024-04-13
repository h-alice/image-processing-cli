package main

import (
	"bytes"
	"flag"
	"imagetools/config"
	"io"
	"log"
	"os"

	op "imagecore/operation"
)

func ProcessFile(profile config.ProcessProfileConfig, in io.Reader) error {

	// Procedure: Decode -> image ops -> encode -> segment ops -> write out

	working_image, err := op.CreateImageFromReader(in)
	if err != nil {
		log.Printf("[x] Error while creating image: %v", err)
		return err
	}

	// Create image processing pipeline.
	for index, pb := range profile.PipelineBlocks {
		log.Printf("Processing Operation #%d: %s", index, pb.Operation)
		working_image = working_image.Then(config.PipelineBlockToOperation(pb))
		if working_image.LastError() != nil {
			log.Printf("[x] Error while processing image: %v", working_image.LastError())
			return working_image.LastError()
		}
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

		// Currently, this function will only affect the `fileName` field in `write` block.
		// This is a temporary solution to the issue which "write" block cannot get original input file name.
		conf.AssignInputFile(f)

		raw_bytes, err := os.ReadFile(f)
		if err != nil {
			log.Printf("[x] Error while reading file: %s\n", err)
			continue
		}

		tasks := make(chan struct{}, len(conf.Profiles))
		for _, pf := range conf.Profiles { // Apply all profile to input image.

			go func(pf config.ProcessProfileConfig) {

				err = ProcessFile(pf, bytes.NewBuffer(raw_bytes))
				if err != nil {
					log.Printf("[x] An error occurred while processing image: %s\n", err)
					return
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
