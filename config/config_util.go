package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

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

// Generate output file name.
func (ocf OutputConfig) GenerateFileName() string {
	original_dir := filepath.Dir(ocf.origFileName)               // Get original file directory.
	original_ext := filepath.Ext(ocf.origFileName)               // Get file extension.
	original_name := filepath.Base(ocf.origFileName)             // Get file name.
	stem := original_name[:len(original_name)-len(original_ext)] // Get file name w/o extension.

	fileSuffix := ""

	switch strings.ToLower(ocf.Format) {
	case "jpeg":
		fileSuffix = ".jpg" // Use JPG instead of JPEG.
	case "":
		// Output format not specified: keep original extension.
		// NOTE: It's not guanteed that the original extension matches the output format.
		fileSuffix = original_ext
	default:
		fileSuffix = "." + ocf.Format // Use specified output format.

	}
	full_file := ocf.NamePrefix + stem + ocf.NameSuffix + fileSuffix

	return filepath.Join(original_dir, full_file)
}

// Check the integrity of pipeline block.
//
// pb: Pipeline block to check.
// Each block must associate with a valid operation.
// If the operation is `resize`, the block must have a valid `ResizeConfig`.
//
// Returns error if pipeline block is invalid.
func checkPipelineBlock(pb PipelineBlock) error {

	switch pb.Operation {
	case OperationDecode: // Decode block.
		// Decode operation does not require additional configuration.
		break
	case OperationResize: // Resize block.
		if pb.Resize == nil {
			return ErrInvalidResizeBlock
		}
	case OperationEncode:
		if pb.Encode == nil { // Encode block.
			return ErrInvalidEncodeBlock
		}
	case OperationCrop: // Crop block.
		if pb.Crop == nil {
			return ErrInvalidCropBlock
		}
	case OperationWrite: // File output block.
		if pb.Write == nil {
			return ErrInvalidWriteBlock
		}
	case OperationIccEmbed: // ICC embedding block.
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

// Assign input file to current config.
//
// This is a temporary solution to the issue which "write" block cannot get original input file name.
// With this helper function, all pipeline block can have similar signature.
// TODO: Find another solution ;)
func (profile_root *ConfigFileRoot) AssignInputFile(input_file string) {

	for _, profile := range profile_root.Profiles {
		for _, pb := range profile.PipelineBlocks {
			if pb.Operation == OperationWrite {
				pb.Write.origFileName = input_file
			}
		}
	}
}

// Generate a config that does nothing to input image.
func GenerateDefaultConfig() ConfigFileRoot {

	// Default config.
	default_config := ConfigFileRoot{
		Profiles: []ProcessProfileConfig{
			{
				ProfileName: "SampleProfile",
				PipelineBlocks: []PipelineBlock{
					{
						Operation: OperationDecode,
					},
					{
						Operation: OperationEncode,
						Encode: &EncodeConfig{
							Format: "jpeg",
						},
					},
					{
						Operation: OperationWrite,
						Write: &OutputConfig{
							NameSuffix: "_output",
						},
					},
				},
			},
		},
	}

	return default_config
}
