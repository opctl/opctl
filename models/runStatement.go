package models

type RunStatement struct {
  Op OpRunStatement `yaml:"op,omitempty"`
  Parallel *ParallelRunStatement `yaml:"parallel,omitempty"`
  Serial   *SerialRunStatement `yaml:"serial,omitempty"`
}
