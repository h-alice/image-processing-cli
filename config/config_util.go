package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	ErrInvalidPipelineBlock     = errors.New("malfomed pipeline block")
	ErrInvalidPipelineBlockType = errors.New("unsupported pipeline block type in config")
	ErrInvalidResizeBlock       = errors.New("resize block provided but no additional configuration")
	ErrInvalidCropBlock         = errors.New("crop block provided but no additional configuration")
	ErrInvalidEncodeBlock       = errors.New("encode block provided but no additional configuration")
	ErrInvalidWriteBlock        = errors.New("file output block provided but no additional configuration")
	ErrInvalidIccBlock          = errors.New("icc embedding block provided but no additional configuration")
)

/*
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
*/

// Check the integrity of pipeline block.
//
// pb: Pipeline block to check.
// Each block must associate with a valid operation.
// If the operation is `resize`, the block must have a valid `ResizeConfig`.
//
// Returns error if pipeline block is invalid.
func checkPipelineBlock(pb PipelineBlock) error {

	switch pb.Operation {
	case "resize":
		if pb.Resize == nil {
			return ErrInvalidResizeBlock
		}
	case "encode":
		if pb.Encode == nil {
			return ErrInvalidEncodeBlock
		}
	case "crop":
		if pb.Crop == nil {
			return ErrInvalidCropBlock
		}
	case "write":
		if pb.Write == nil {
			return ErrInvalidWriteBlock
		}
	case "icc_embed":
		if pb.ICCEmbedProfile == nil {
			return ErrInvalidIccBlock
		}
	default:
		return ErrInvalidPipelineBlockType
	}

	return nil
}

// Check the integrity of pipeline block list.
//
// pbs: List of pipeline blocks to check.
func checkPipelineBlockList(pbs []PipelineBlock) error {

	for _, pb := range pbs {
		err := checkPipelineBlock(pb)
		if err != nil {
			return err
		}
	}

	return nil
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

	// Iterate through profiles.
	for _, profile := range conf.Profiles {

		// Check pipeline blocks.
		err = checkPipelineBlockList(profile.PipelineBlocks)
		if err != nil {
			return nil, err
		}
	}

	return &conf, nil
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

/*
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
