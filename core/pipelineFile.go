package core

type pipelineFileStage struct {
  Name string `yaml:"name"`
  Type string `yaml:"type"`
}

type pipelineFile struct {
  Description string `yaml:"description"`
  Stages      []pipelineFileStage `yaml:"stages"`
}