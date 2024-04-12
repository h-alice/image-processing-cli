package main

import (
	"fmt"
	"imagetools/config"
	"os"
	"path/filepath"
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
