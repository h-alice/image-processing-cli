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
			if pb.Operation == "resize" {
				t.Logf("Resize: %v", pb.Resize)
			}
			if pb.Operation == "write" {
				t.Logf("Write: %v", pb.Write)
			}
			if pb.Operation == "icc_embed" {
				t.Logf("ICC: %v", pb.ICCEmbedProfile)
			}
			if pb.Operation == "encode" {
				t.Logf("Encode: %v", pb.Encode)
			}
			if pb.Operation == "crop" {
				t.Logf("Crop: %v", pb.Crop)
			}

		}
	}
}

/*
func TestLoadConfigAndParse(t *testing.T) {

	config, err := LoadConfigFromFile("test_resources/test1.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(config.Profiles) == 0 {
		t.Fatalf("Expected at least one profile, got none.")
	}

	// Test first profile.

	if config.Profiles[0].ProfileName != "Profile1" {
		t.Fatalf("Expected profile name to be 'Profile1', got '%s'", config.Profiles[0].ProfileName)
	}

	if config.Profiles[0].ICC != "ADOBE RGB" {
		t.Fatalf("Expected ICC to be 'ADOBE RGB', got '%s'", config.Profiles[0].ICC)
	}

	// Config 1: Resize config.
	if config.Profiles[0].Resize.Width != 100 {
		t.Fatalf("Expected width to be 100, got '%d'", config.Profiles[0].Resize.Width)
	}

	if config.Profiles[0].Resize.Height != 200 {
		t.Fatalf("Expected height to be 200, got '%d'", config.Profiles[0].Resize.Height)
	}

	if config.Profiles[0].Resize.Factor != 0.9 {
		t.Fatalf("Expected factor to be 0.9, got '%f'", config.Profiles[0].Resize.Factor)
	}

	if config.Profiles[0].Resize.Algorithm != "catmullrom" {
		t.Fatalf("Expected algorithm to be 'catmullrom', got '%s'", config.Profiles[0].Resize.Algorithm)
	}

	// Config 1: Output config.
	if config.Profiles[0].Output.Format != "jpeg" {
		t.Fatalf("Expected format to be 'jpeg', got '%s'", config.Profiles[0].Output.Format)
	}

	if config.Profiles[0].Output.NamePrefix != "prefix1_" {
		t.Fatalf("Expected prefix to be 'prefix1_', got '%s'", config.Profiles[0].Output.NamePrefix)
	}

	if config.Profiles[0].Output.NameSuffix != "_suffix1" {
		t.Fatalf("Expected suffix to be '_suffix1', got '%s'", config.Profiles[0].Output.NameSuffix)
	}

	if config.Profiles[0].Output.Options.Quality != 80 {
		t.Fatalf("Expected quality to be 80, got '%d'", config.Profiles[0].Output.Options.Quality)
	}

	// Test second profile.

	if config.Profiles[1].ProfileName != "Profile2" {
		t.Fatalf("Expected profile name to be 'Profile2', got '%s'", config.Profiles[1].ProfileName)
	}

	if config.Profiles[1].Output.Format != "png" {
		t.Fatalf("Expected format to be 'png', got '%s'", config.Profiles[1].Output.Format)
	}

	if config.Profiles[1].Output.Options != nil {
		t.Fatalf("Expected options to be nil, got '%v'", config.Profiles[1].Output.Options)
	}

	// Config 2: Resize config (omitted some fields)
	if config.Profiles[1].Resize.Width != 0 {
		t.Fatalf("Expected width to be 0, got '%d'", config.Profiles[1].Resize.Width)
	}

	if config.Profiles[1].Resize.Height != 0 {
		t.Fatalf("Expected height to be 0, got '%d'", config.Profiles[1].Resize.Height)
	}

	if config.Profiles[1].Resize.Factor != 0.0 {
		t.Fatalf("Expected factor to be 0.0, got '%f'", config.Profiles[1].Resize.Factor)
	}

	if config.Profiles[1].Resize.Algorithm != "nearestneighbor" {
		t.Fatalf("Expected algorithm to be 'nearestneighbor', got '%s'", config.Profiles[1].Resize.Algorithm)
	}

}

func TestPrettyPrint(t *testing.T) {

	// Note that this test is not deterministic.
	// It is only used to check if the function runs without error.

	config, err := LoadConfigFromFile("test_resources/test1.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Logf("Config: %s", config.PrettyPrint())
}

func TestProfileToYaml(t *testing.T) {

	config, err := LoadConfigFromFile("test_resources/test1.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Logf("Config: %s", config.ToYaml())
}
*/
