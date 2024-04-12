package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	config, err := LoadConfigFromFile("test_resources/test_full_conf.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	for _, pf := range config.Profiles {
		t.Logf("Profile: %s", pf.ProfileName)
		for _, pb := range pf.PipelineBlocks {
			t.Logf("Operation: %s", pb.Operation)
			if pb.Operation == OperationResize {
				t.Logf("Resize: %v", pb.Resize)
			}
			if pb.Operation == OperationWrite {
				t.Logf("Write: %v", pb.Write)
			}
			if pb.Operation == OperationIccEmbed {
				t.Logf("ICC: %v", pb.ICCEmbedProfile)
			}
			if pb.Operation == OperationEncode {
				t.Logf("Encode: %v", pb.Encode)
			}
			if pb.Operation == OperationCrop {
				t.Logf("Crop: %v", pb.Crop)
			}

		}
	}
}

func TestProfileToYaml(t *testing.T) {

	config, err := LoadConfigFromFile("test_resources/test_full_conf.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Logf("Config: %s", config.ToYaml())
}

func TestDefaultConfigGeneration(t *testing.T) {

	// Note that this test is not deterministic.
	// It is only used to check if the function runs without error.

	default_config := GenerateDefaultConfig()
	t.Logf("Default config: %s", default_config.ToYaml())
}

func TestLoadConfigAndParse(t *testing.T) {

	// This is the test case for the config file `test_full_conf.yaml`.
	// Includes fully checking all fields.
	config, err := LoadConfigFromFile("test_resources/test_full_conf.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Profiles) == 0 {
		t.Fatalf("Expected at least one profile, got none.")
	}

	// Assign dummy input file.
	config.AssignInputFile("dummy.jpg")

	// Test first profile.

	pf := config.Profiles[0]
	pb := pf.PipelineBlocks

	if pf.ProfileName != "Profile1" {
		t.Fatalf("Expected profile name to be 'Profile1', got '%s'", config.Profiles[0].ProfileName)
	}

	// Process all pipeline blocks.

	// Block 1: Decode.
	if pb[0].Operation != OperationDecode {
		t.Fatalf("Expected operation to be 'Decode', got '%s'", pb[0].Operation)
	}

	// Block 2: Crop.
	if pb[1].Operation != OperationCrop {
		t.Fatalf("Expected operation to be 'Crop', got '%s'", pb[1].Operation)
	}

	if pb[1].Crop.Width != 50 {
		t.Fatalf("Expected width to be 50, got '%d'", pb[1].Crop.Width)
	}

	if pb[1].Crop.Height != 60 {
		t.Fatalf("Expected height to be 60, got '%d'", pb[1].Crop.Height)
	}

	if pb[1].Crop.Alignment != "center" {
		t.Fatalf("Expected alignment to be 'center', got '%s'", pb[1].Crop.Alignment)
	}

	// Block 3: Resize.
	if pb[2].Operation != OperationResize {
		t.Fatalf("Expected operation to be 'Resize', got '%s'", pb[2].Operation)
	}

	if pb[2].Resize.Width != 100 {
		t.Fatalf("Expected width to be 100, got '%d'", pb[2].Resize.Width)
	}

	if pb[2].Resize.Height != 200 {
		t.Fatalf("Expected height to be 200, got '%d'", pb[2].Resize.Height)
	}

	if pb[2].Resize.Factor != 0.9 {
		t.Fatalf("Expected factor to be 0.9, got '%f'", pb[2].Resize.Factor)
	}

	if pb[2].Resize.Algorithm != "catmullrom" {
		t.Fatalf("Expected algorithm to be 'catmullrom', got '%s'", pb[2].Resize.Algorithm)
	}

	// Block 4: ICC Embed.
	if pb[3].Operation != OperationIccEmbed {
		t.Fatalf("Expected operation to be 'IccEmbed', got '%s'", pb[3].Operation)
	}

	if pb[3].ICCEmbedProfile.ProfileName != "sRGB" {
		t.Fatalf("Expected profile name to be 'sRGB', got '%s'", pb[3].ICCEmbedProfile.ProfileName)
	}

	// Block 5: Encode.
	if pb[4].Operation != OperationEncode {
		t.Fatalf("Expected operation to be 'Encode', got '%s'", pb[4].Operation)
	}

	if pb[4].Encode.Format != "jpeg" {
		t.Fatalf("Expected format to be 'jpeg', got '%s'", pb[4].Encode.Format)
	}

	if pb[4].Encode.Options.Quality != 80 {
		t.Fatalf("Expected quality to be 80, got '%d'", pb[4].Encode.Options.Quality)
	}

	// Block 6: Write.
	if pb[5].Operation != OperationWrite {
		t.Fatalf("Expected operation to be 'Write', got '%s'", pb[5].Operation)
	}

	if pb[5].Write.NamePrefix != "prefix1_" {
		t.Fatalf("Expected prefix to be 'prefix1_', got '%s'", pb[5].Write.NamePrefix)
	}

	if pb[5].Write.NameSuffix != "_suffix1" {
		t.Fatalf("Expected suffix to be '_suffix1', got '%s'", pb[5].Write.NameSuffix)
	}

	if pb[5].Write.Format != "jpeg" {
		t.Fatalf("Expected format to be 'jpeg', got '%s'", pb[5].Write.Format)
	}
}
