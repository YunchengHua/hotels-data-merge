package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DataSourcesRefreshIntervalSec int64 `yaml:"data_sources_refresh_interval_sec"`
	MergeDescriptions             bool  `yaml:"merge_description"`
}

const (
	defaultConfigPath = "./etc"
	defaultConfigFile = "local.yaml"
)

var (
	cfg = &Config{}
)

func GetConfig() Config {
	return *cfg
}

func MustInit() {
	fmt.Println("config start to init...")
	config := Config{}
	// Read file
	data, err := os.ReadFile(
		filepath.Join(defaultConfigPath, defaultConfigFile))
	if err != nil {
		fmt.Println("fail to read file", err)
		log.Fatalf("error: %v", err)
	}

	// Unmarshal YAML to struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("fail to unmarshal config", err)
		log.Fatalf("error: %v", err)
	}
	cfg = &config
	log.Printf("Config loaded: %v", config)
}
