// Description: Main entry point for the CLI application.
package main

import (
	"context"
	"fmt"
	"imagetools/config"
	"log"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

// Get profile from home directory.
// This function will trying to get profile from `.imgtools` directory in home directory.
//
// profile_name: Profile name.
// create_dir: Create directory if not exists.
func getProfileFromHomeDir(profile_name string, create_dir bool) (path string, err error) {

	// If profile name is empty, use `default``.
	if profile_name == "" {
		profile_name = "default"
	}

	profile_name = fmt.Sprintf("%s.yaml", profile_name) // Append `.yaml` extension.

	home, err := os.UserHomeDir() // Get home directory.
	if err != nil {
		return "", err
	}

	profile_dir := filepath.Join(home, ".imgtools")                // Profile directory.
	profile_file := filepath.Join(home, ".imgtools", profile_name) // Profile file.

	_, err = os.Stat(profile_dir) // Check profile directory.
	if err != nil {

		// If directory not exists and create_dir is true, create directory.
		if os.IsNotExist(err) {

			if create_dir { // Making profile directory if not exists.
				err = os.Mkdir(profile_dir, 0777)
				if err != nil {
					return "", err
				}

				// Writing default profile.
				err = os.WriteFile(profile_file, []byte(config.GenerateDefaultConfig().ToYaml()), 0777)
				if err != nil {
					return "", err
				}
			} else { // If create_dir is false, return error.
				return "", err
			}
		} else {
			return "", err
		}
	}

	return profile_file, err
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

				log.Printf("[!] No profile specified. Trying to load default profile from home directory.\n")

				config_path, err := getProfileFromHomeDir("default", true)
				if err != nil {
					log.Fatalf("[x] Cannot load default config file: %s\n", err)
				}
				config_root, err = config.LoadConfigFromFile(config_path) // Load default config file.
				if err != nil {
					log.Fatalf("[x] Cannot load default config file: %s\n", err)
				}
			}

			// Iterate through input images.
			for _, f := range c.Args().Slice() {
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

				// Dispatch goroutine for each profile.
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

	log.Printf("[+] All images processed.\n")
}
