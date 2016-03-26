package osfilesys

type compositionRoot interface {
  ListNamesOfDevOpDirsUseCase() listNamesOfDevOpDirsUseCase
  ListNamesOfPipelineDirsUseCase() listNamesOfPipelineDirsUseCase
  ReadDevOpFileUseCase() readDevOpFileUseCase
  ReadPipelineFileUseCase() readPipelineFileUseCase
  SaveDevOpFileUseCase() saveDevOpFileUseCase
  SavePipelineFileUseCase() savePipelineFileUseCase
  CreateDevOpDirUseCase() createDevOpDirUseCase
  CreatePipelineDirUseCase() createPipelineDirUseCase
}

func newCompositionRoot(
) compositionRoot {

  relPathToDevOpDirFactory := newRelPathToDevOpDirFactory()
  relPathToDevOpFileFactory := newRelPathToDevOpFileFactory(relPathToDevOpDirFactory)
  relPathToPipelineDirFactory := newRelPathToPipelineDirFactory()
  relPathToPipelineFileFactory := newRelPathToPipelineFileFactory(relPathToPipelineDirFactory)

  return &_compositionRoot{
    listNamesOfDevOpDirsUseCase: newListNamesOfDevOpDirsUseCase(),
    listNamesOfPipelineDirsUseCase:newListNamesOfPipelineDirsUseCase(),
    readDevOpFileUseCase:newReadDevOpFileUseCase(relPathToDevOpFileFactory),
    readPipelineFileUseCase:newReadPipelineFileUseCase(relPathToPipelineFileFactory),
    saveDevOpFileUseCase:newSaveDevOpFileUseCase(relPathToDevOpFileFactory),
    savePipelineFileUseCase:newSavePipelineFileUseCase(relPathToPipelineFileFactory),
    createDevOpDirUseCase:newCreateDevOpDirUseCase(relPathToDevOpDirFactory),
    createPipelineDirUseCase:newCreatePipelineDirUseCase(relPathToPipelineDirFactory),
  }

}

type _compositionRoot struct {
  listNamesOfDevOpDirsUseCase    listNamesOfDevOpDirsUseCase
  listNamesOfPipelineDirsUseCase listNamesOfPipelineDirsUseCase
  readDevOpFileUseCase           readDevOpFileUseCase
  readPipelineFileUseCase        readPipelineFileUseCase
  saveDevOpFileUseCase           saveDevOpFileUseCase
  savePipelineFileUseCase        savePipelineFileUseCase
  createDevOpDirUseCase          createDevOpDirUseCase
  createPipelineDirUseCase       createPipelineDirUseCase
}

func (this _compositionRoot) ListNamesOfDevOpDirsUseCase(
) listNamesOfDevOpDirsUseCase {
  return this.listNamesOfDevOpDirsUseCase
}

func (this _compositionRoot) ListNamesOfPipelineDirsUseCase(
) listNamesOfPipelineDirsUseCase {
  return this.listNamesOfPipelineDirsUseCase
}

func (this _compositionRoot) ReadDevOpFileUseCase(
) readDevOpFileUseCase {
  return this.readDevOpFileUseCase
}

func (this _compositionRoot) ReadPipelineFileUseCase(
) readPipelineFileUseCase {
  return this.readPipelineFileUseCase
}

func (this _compositionRoot) SaveDevOpFileUseCase(
) saveDevOpFileUseCase {
  return this.saveDevOpFileUseCase
}

func (this _compositionRoot) SavePipelineFileUseCase(
) savePipelineFileUseCase {
  return this.savePipelineFileUseCase
}

func (this _compositionRoot) CreateDevOpDirUseCase(
) createDevOpDirUseCase {
  return this.createDevOpDirUseCase
}

func (this _compositionRoot) CreatePipelineDirUseCase(
) createPipelineDirUseCase {
  return this.createPipelineDirUseCase
}