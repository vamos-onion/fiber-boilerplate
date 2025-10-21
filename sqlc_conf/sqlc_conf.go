package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the sqlc.yaml structure
type Config struct {
	Version string `yaml:"version"`
	SQL     []struct {
		Engine  string   `yaml:"engine"`
		Queries []string `yaml:"queries"`
		Schema  string   `yaml:"schema"`
		Rules   []string `yaml:"rules"`
		Gen     struct {
			Go struct {
				Package      string                   `yaml:"package"`
				Out          string                   `yaml:"out"`
				SqlPackage   string                   `yaml:"sql_package"`
				Overrides    []map[string]interface{} `yaml:"overrides,omitempty"`
				Rename       map[string]interface{}   `yaml:"rename,omitempty"`
				EmitJsonTags bool                     `yaml:"emit_json_tags"`
				EmitDbTags   bool                     `yaml:"emit_db_tags"`
			} `yaml:"go"`
		} `yaml:"gen"`
	} `yaml:"sql"`
	Overrides interface{} `yaml:"overrides,omitempty"`
}

func main() {
	// set Current Dir
	err := os.Chdir("sqlc_conf")
	if err != nil {
		log.Fatalf("Error changing directory: %v", err)
	}

	// Load sqlc.yaml
	sqlcData, err := os.ReadFile("sqlc.yaml")
	if err != nil {
		log.Fatalf("Error reading sqlc.yaml: %v", err)
	}
	var sqlcConfig Config
	if err := yaml.Unmarshal(sqlcData, &sqlcConfig); err != nil {
		log.Fatalf("Error unmarshalling sqlc.yaml: %v", err)
	}

	// Load overrides.yaml
	overridesData, err := os.ReadFile("overrides.yaml")
	if err != nil {
		log.Fatalf("Error reading overrides.yaml: %v", err)
	}

	// Unmarshal overrides.yaml
	var overrides []map[string]interface{} = make([]map[string]interface{}, 100)
	if err := yaml.Unmarshal(overridesData, &overrides); err != nil {
		log.Fatalf("Error unmarshalling overrides.yaml: %v", err)
	}

	// Apply overrides to sqlcConfig
	for _, override := range overrides {
		sqlcConfig.SQL[0].Gen.Go.Overrides = append(sqlcConfig.SQL[0].Gen.Go.Overrides, override)
	}

	// Output merged configuration (for demonstration)
	mergedData, err := yaml.Marshal(&sqlcConfig)
	if err != nil {
		log.Fatalf("Error marshalling merged config: %v", err)
	}

	// set Current Dir
	err = os.Chdir("../out")
	if err != nil {
		log.Fatalf("Error changing directory: %v", err)
	}

	// Create merged configuration to ../out/merged_sqlc.yaml
	f, err := os.Create("merged_sqlc.yaml")
	if err != nil {
		log.Fatalf("Error creating merged_sqlc.yaml: %v", err)
	}

	_, err = f.Write(mergedData)
	if err != nil {
		log.Fatalf("Error writing merged_sqlc.yaml: %v", err)
	}
}
