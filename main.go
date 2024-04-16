package main

import (
	"context"
	"flag"
	"imagetools/config"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
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
func ProcessFile(profile config.ImageProcessingProfile) error {

	working_image, err := profile.CreateImageFile()
	if err != nil {
		log.Printf("[x] Error while creating image: %v", err)
		return err
	}

	// Create image processing pipeline.
	for _, pb := range profile.PipelineBlocks {
		// log.Printf("Processing Operation #%d: %s", index, pb.Operation)
		working_image = working_image.Then(config.PipelineBlockToOperation(pb))
		if working_image.LastError() != nil {
			log.Printf("[x] Error while processing image: %v", working_image.LastError())
			return working_image.LastError()
		}
	}

	return nil
}

// Main worker.
func mainWorker(ctx context.Context, profile config.ImageProcessingProfile, result_chan chan<- error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get only file name from the path.
	input_image_name := profile.GetAssignedFilePath()
	input_image_name = filepath.Base(input_image_name)

	log.Printf("[.] Processing image [%s] with profile [%s]\n", input_image_name, profile.ProfileName)

	// Process image.
	for {
		select {
		case <-ctx.Done(): // Check if context is cancelled.
			log.Printf("[!] Context cancelled. Exiting...\n")
			result_chan <- nil
			return // Terminate goroutine.
		default:
			err := ProcessFile(profile) // Process image.
			if err != nil {
				log.Printf("[x] Error while processing image: %v\n", err)
				result_chan <- err
				return // Terminate goroutine.
			}
			result_chan <- nil
			return // Goroutine finished.
		}
	}
}

// Main function, defines arguments and flags.
func main() {

	// Context for main worker.
	ctx := context.Background()

	app := &cli.App{
		Name:  "Image Processing CLI",
		Usage: "Batch process images",
		Authors: []*cli.Author{
			{
				Name:  "h-alice",
				Email: "admin@halice.art",
			},
			{
				Name: "natlee",
			},
		},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:  "f",
				Usage: "Config file path (can be specified multiple times)",
			},
		},
		Action: func(c *cli.Context) error {

			// Placeholder for input config file paths.
			loaded_configs := make([]config.ProfileRoot, 0) // Placeholder for loaded configs.

			// Placeholder for config root.
			config_root := config.ProfileRoot{}

			for _, path := range c.StringSlice("f") { // Iterate through input config file paths.
				conf, err := config.LoadConfigFromFile(path) // Load config file.
				if err != nil {
					log.Printf("[!] Error (%s) while loading config file: %s The config file will be ignored.\n", err, path)
				}
				loaded_configs = append(loaded_configs, conf) // Append to loaded configs.
			}

			config_root = config.MergeConfigFiles(loaded_configs...) // Merge all loaded configs.

			if len(config_root.Profiles) == 0 {
				log.Fatalf("[x] No profile found in config file.\n")
			}

			// Iterate through input images.
			for _, f := range flag.Args() { // Iterate through input images.
				// Create result channel to capture return from goroutine.
				result_chan := make(chan error) // Result channel.

				// Currently, this function will only affect the `fileName` field in `write` block.
				// This is a temporary solution to the issue which "write" block cannot get original input file name.
				config_root.AssignInputFile(f)

				// Check if the file exists.
				_, err := os.Stat(f) // Check if file exists.
				if err != nil {
					log.Printf("[!] Input file not found: %s\n", f)
					continue // Skip to next file.
				}

				for _, pf := range config_root.Profiles { // Apply all profile to input image.
					// Process image in goroutine.
					go mainWorker(ctx, pf, result_chan)
				}

				// Wait for all goroutines to finish.
				for range config_root.Profiles {
					<-result_chan
				}
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("[x] Error: %v\n", err)
	}

	// Merge all config files, if any specified.
	//if len(*config_paths) != 0 {

	// If path not specified, load defaultprofile from home directory.

	//	log.Printf("[!] Using default config file.\n")
	//
	//	config_path, err := getProfileFromHomeDir("default", true)
	//
	//	if err != nil {
	//		log.Fatalf("[x] Cannot load default config file: %s\n", err)
	//	}
	//
	//	config_root, err = config.LoadConfigFromFile(config_path) // Load default config file.
	//	if err != nil {
	//		log.Fatalf("[x] Cannot load default config file: %s\n", err)
	//	}
	//

	log.Printf("[+] All images processed.\n")
}
