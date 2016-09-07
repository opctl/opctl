package models

type BundleManifest struct {
  Name        string `yaml:"name"`
  Description string `yaml:"description"`
  Version string `yaml:"version,omitempty"`
}
