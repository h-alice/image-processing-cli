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
	Format     string              `yaml:"format"`     // Output file format
	NameSuffix string              `yaml:"nameSuffix"` // Output file name suffix
	NamePrefix string              `yaml:"namePrefix"` // Output file name prefix
	Options    *OutputOptionConfig `yaml:"options"`    // Encoder option
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
	ProfileName string        `yaml:"profileName"` // Profile identifier
	ICC         string        `yaml:"icc"`         // ICC profile to embed
	Resize      *ResizeConfig `yaml:"resize"`      // Resize option
	Output      *OutputConfig `yaml:"output"`      // Output file configuration
}

// Currently not used.
type OutputDirConfig struct {
	DirName string `yaml:"dirName"` // Output directory name
}

// Config structure for config file.
//
// Profiles: List of profile configurations.
type ConfigFileRoot struct {
	Profiles []ProcessProfileConfig `yaml:"profiles"` // List of profile configurations
}
