package config

// Configuration structure
type Configuration struct {
	Root       string   `yaml:"root"`
	Exclusions []string `yaml:"exclusions"`
	Threshold  float64  `yaml:"threshold,omitempty"`
	Metric     string   `yaml:"thresholdType,omitempty"`
}
