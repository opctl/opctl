package ports

type Filesys interface {
  CreateDevOpDir(
  devOpName string,
  ) (err error)

  CreatePipelineDir(
  pipelineName string,
  ) (err error)

  ListNamesOfDevOpDirs(
  ) (namesOfDevOpDirs []string, err error)

  ListNamesOfPipelineDirs(
  ) (namesOfPipelineDirs []string, err error)

  ReadDevOpFile(
  devOpName string,
  ) (file []byte, err error)

  ReadPipelineFile(
  pipelineName string,
  ) (file []byte, err error)

  SaveDevOpFile(
  devOpName string,
  data []byte,
  ) (err error)

  SavePipelineFile(
  pipelineName string,
  data []byte,
  ) (err error)
}