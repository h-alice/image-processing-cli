package main

import (
	"bytes"
	"flag"
	op "imagecore/operation"
	"imagetools/config"
	"io"
	"log"
	"os"
	"strings"
)

// This is a list which holds the input config file paths.
type inputConfigFilePaths []string

// Implementing the flag.Value interface.
func (m *inputConfigFilePaths) String() string {
	return strings.Join(*m, ", ")
}

// Implementing the flag.Value interface.
func (m *inputConfigFilePaths) Set(value string) error {
	*m = append(*m, value)
	return nil
}

// Process image file with profile.
func ProcessFile(profile config.ProcessProfileConfig, in io.Reader) error {

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

	config_paths := new(inputConfigFilePaths)

	// Parse command line.
	flag.Var(config_paths, "f", "Config file path (can be specified multiple times)")
	flag.Parse()

	config_root := config.ConfigFileRoot{}

	// Merge all config files, if any specified.
	if len(*config_paths) != 0 {
		loaded_configs := make([]config.ConfigFileRoot, 0) // Placeholder for loaded configs.
		for _, path := range *config_paths {
			conf, err := config.LoadConfigFromFile(path) // Load config file.
			if err != nil {
				log.Printf("[!] Error (%s) while loading config file: %s The config file will be ignored.\n", err, path)
			}
			loaded_configs = append(loaded_configs, conf) // Append to loaded configs.
		}
		config_root = config.MergeConfigFiles(loaded_configs...) // Merge all loaded configs.
	} else { // If path not specified, load defaultprofile from home directory.

		log.Printf("[!] Using default config file.")

		config_path, err := getProfileFromHomeDir("default", true)

		if err != nil {
			log.Fatalf("[x] Cannot load default config file: %s\n", err)
		}

		config_root, err = config.LoadConfigFromFile(config_path) // Load default config file.
		if err != nil {
			log.Fatalf("[x] Cannot load default config file: %s\n", err)
		}
	}

	for _, f := range flag.Args() { // Iterate through input images.

		// Currently, this function will only affect the `fileName` field in `write` block.
		// This is a temporary solution to the issue which "write" block cannot get original input file name.
		config_root.AssignInputFile(f)

		raw_bytes, err := os.ReadFile(f)
		if err != nil {
			log.Printf("[x] Error while reading file: %s\n", err)
			continue
		}

		tasks := make(chan struct{}, len(config_root.Profiles))
		for _, pf := range config_root.Profiles { // Apply all profile to input image.

			go func(pf config.ProcessProfileConfig) {

				err = ProcessFile(pf, bytes.NewBuffer(raw_bytes))
				if err != nil {
					log.Printf("[x] An error occurred while processing image: %s\n", err)
					return
				}

				tasks <- struct{}{}
			}(pf)
		}

		for i := 0; i < len(config_root.Profiles); i++ {
			log.Println(i)
			<-tasks
		}
	}
}
