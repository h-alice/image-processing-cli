package config

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Generate output file name.
func (ocf OutputConfig) GenerateFileName(input_name string) string {

	original_ext := filepath.Ext(input_name)                     // Get file extension.
	original_name := filepath.Base(input_name)                   // Get file name.
	stem := original_name[:len(original_name)-len(original_ext)] // Get file name w/o extension.

	fileSuffix := ""

	switch strings.ToLower(ocf.Format) {
	case "jpeg":
		fileSuffix = ".jpg" // Use JPG instead of JPEG.
	case "":
		fileSuffix = original_ext // Output format not specified: keep original extension.
	default:
		fileSuffix = "." + ocf.Format // Use specified output format.

	}
	full_file := ocf.NamePrefix + stem + ocf.NameSuffix + fileSuffix

	return full_file
}

// Load config file from path.
//
// config_path: Path to config file.
func LoadConfigFromFile(config_path string) (*ConfigFileRoot, error) {

	raw_config, err := os.ReadFile(config_path) // Read raw config file.
	if err != nil {
		return nil, err
	}

	// Converting JSON to config structure.
	var conf ConfigFileRoot                 // Parsed config placeholder.
	err = yaml.Unmarshal(raw_config, &conf) // Convert JSON to structure.
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

/*
// Pretty print config file.
func (pf ConfigFileRoot) PrettyPrint() string {

	var output string = "" // Placeholder for output.

	ident := "  " // Indentation.

	for _, pf := range pf.Profiles { // Iterate through profiles.

		output += "Profile Name: " + pf.ProfileName + "\n"
		output += ident + "ICC Profile: " + pf.ICC + "\n"

		if pf.Resize != nil {
			output += ident + "Resizing Configuration:\n"
			output += ident + ident + "Resize Width: " + fmt.Sprintf("%d", pf.Resize.Width) + "\n"
			output += ident + ident + "Resize Height: " + fmt.Sprintf("%d", pf.Resize.Height) + "\n"
			output += ident + ident + "Resize Factor: " + fmt.Sprintf("%f.2", pf.Resize.Factor) + "\n"
			output += ident + ident + "Resize Algorithm: " + pf.Resize.Algorithm + "\n"
		}

		if pf.Output != nil {
			output += ident + "Output Configuration:\n"
			output += ident + ident + "Output Format: " + pf.Output.Format + "\n"
			output += ident + ident + "Output Name Prefix: " + pf.Output.NamePrefix + "\n"
			output += ident + ident + "Output Name Suffix: " + pf.Output.NameSuffix + "\n"
			if pf.Output.Options != nil {
				output += ident + ident + "Encoder Quality: " + fmt.Sprintf("%d", pf.Output.Options.Quality) + "\n"
			}
		}
		output += "\n"
	}
	return output
}

// Profile instance to yaml string.
func (profile_root ConfigFileRoot) ToYaml() string {

	// Convert to yaml.
	yaml_str, err := yaml.Marshal(profile_root)
	if err != nil {
		return ""
	}

	return string(yaml_str)
}

// Generate a config that does nothing to input image.
func GenerateDefaultConfig() string {

	// Default config.
	default_config := ConfigFileRoot{
		Profiles: []ProcessProfileConfig{
			{
				ProfileName: "SampleProfile",
				ICC:         "",
				Resize: &ResizeConfig{
					Factor:    1.0,
					Algorithm: "nearestneighbor",
				},
				Output: &OutputConfig{
					Format:     "jpg",
					NameSuffix: "",
					NamePrefix: "",
					Options: &OutputOptionConfig{
						Quality: 100,
					},
				},
			},
		},
	}

	return default_config.ToYaml()
}
*/
