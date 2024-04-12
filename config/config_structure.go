package config

import (
	"errors"

	op "imagecore/operation" // Grab `EncoderOption` from operation package.
)

var ErrNotImplemented = errors.New("operation not implemented")

// OutputOptionConfig structure for output file.
type OutputOptionConfig op.EncoderOption

// Config structure for output file.
//
// Format: Output file format. Either `jpeg` or `png`.
//
// NameSuffix: Suffix of the output file name.
//
// NamePrefix: Prefix of the output file name.
//
// Options: Encoder option. For jpeg use and supports only `Quality` option.
type OutputConfig struct {
	Format     string `yaml:"format"` // Output file format
	NameSuffix string `yaml:"suffix"` // Output file name suffix
	NamePrefix string `yaml:"prefix"` // Output file name prefix
	fileName   string // This is used to store the file name of input image, hence no need to serialize this field.
}

// Config structure for resizing image.
//
// Width: Output image width.
//
// Height: Output image height.
//
// Factor: Resize factor.
//
// Algorithm: Resize algorithm. Either `nearestneighbor`, `catmullrom`, or `approxbilinear`.
//
// NOTE: The `Factor` is prioritized over `Width` and `Height`.
type ResizeConfig struct {
	Width     int     `yaml:"width"`     // Output image width
	Height    int     `yaml:"height"`    // Output image height
	Factor    float32 `yaml:"factor"`    // Resize factor
	Algorithm string  `yaml:"algorithm"` // Resize algorithm
}

type EncodeConfig struct {
	Format  string              `yaml:"format"`  // Output file format
	Options *OutputOptionConfig `yaml:"options"` // Encoder option
}

type IccEmbedConfig struct {
	ProfileName string `yaml:"icc_name"` // Profile name

}

// Config structure for cropping image.
//
// Alignment: Crop alignment. One of `center` or `topleft`.
//
// Width: Crop width.
//
// Height: Crop height.
type CropConfig struct {
	Alignment string `yaml:"alignment"` // Crop alignment
	Width     int    `yaml:"width"`     // Crop width
	Height    int    `yaml:"height"`    // Crop height
}

// Config structure for processing profile.
//
// ProfileName: Profile identifier.
//
// ICC: ICC profile to embed.
//
// Resize: Resize option.
//
// Output: Output file configuration.
type ProcessProfileConfig struct {
	ProfileName    string          `yaml:"profile_name"` // Profile identifier
	PipelineBlocks []PipelineBlock `yaml:"pipeline"`     // Pipeline blocks
}

// Currently not used.
type OutputDirConfig struct {
	DirName string `yaml:"dirName"` // Output directory name
}

// Operation block structure.
//
// Operation: Operation name.
// NOTE: The operation type is one of the following:
// - `decode`
// - `crop`
// - `resize`
// - `icc_embed`
// - `encode`
// - `write`
type PipelineBlock struct {
	Operation       string          `yaml:"operation"`               // Operation name.
	Crop            *CropConfig     `yaml:"crop_config,omitempty"`   // Crop configuration.
	Resize          *ResizeConfig   `yaml:"resize_config,omitempty"` // Resize configuration.
	ICCEmbedProfile *IccEmbedConfig `yaml:"icc_config,omitempty"`    // Embed profile configuration.
	Encode          *EncodeConfig   `yaml:"encode_config,omitempty"` // Encode configuration.
	Write           *OutputConfig   `yaml:"write_config,omitempty"`  // Write configuration.
}

// Config structure for config file.
//
// Profiles: List of profile configurations.
type ConfigFileRoot struct {
	Profiles []ProcessProfileConfig `yaml:"profiles"` // List of profile configurations
}
