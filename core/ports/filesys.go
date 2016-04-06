package ports

type Filesys interface {
  CreateDir(
  pathToDir string,
  ) (err error)

  ListNamesOfChildDirs(
  pathToParentDir string,
  ) (namesOfChildDirs []string, err error)

  GetBytesOfFile(
  pathToFile string,
  ) (bytesOfFile []byte, err error)

  SaveFile(
  pathToFile string,
  bytesOfFile []byte,
  ) (err error)
}