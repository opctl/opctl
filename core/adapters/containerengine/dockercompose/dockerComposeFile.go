package dockercompose

type dockerComposeFileService struct {
  Image string `yaml:"image"`
}

type dockerComposeFile struct {
  Version  string `yaml:"version"`
  Services map[string]dockerComposeFileService `yaml:"services"`
}