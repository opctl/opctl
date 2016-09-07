package models

type OpBundleManifest struct {
  BundleManifest
  Inputs []Param `yaml:"inputs,omitempty"`
  Run    *RunStatement `yaml:"run,omitempty"`
}
