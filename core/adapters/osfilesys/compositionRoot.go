package osfilesys

type compositionRoot interface {
  ListNamesOfDevOpDirsUcExecuter() listNamesOfDevOpDirsUcExecuter
  ListNamesOfPipelineDirsUcExecuter() listNamesOfPipelineDirsUcExecuter
  ReadDevOpFileUcExecuter() readDevOpFileUcExecuter
  ReadPipelineFileUcExecuter() readPipelineFileUcExecuter
  SaveDevOpFileUcExecuter() saveDevOpFileUcExecuter
  SavePipelineFileUcExecuter() savePipelineFileUcExecuter
  CreateDevOpDirUcExecuter() createDevOpDirUcExecuter
  CreatePipelineDirUcExecuter() createPipelineDirUcExecuter
}

func newCompositionRoot(
) compositionRoot {

  relPathToDevOpDirFactory := newRelPathToDevOpDirFactory()
  relPathToDevOpFileFactory := newRelPathToDevOpFileFactory(relPathToDevOpDirFactory)
  relPathToPipelineDirFactory := newRelPathToPipelineDirFactory()
  relPathToPipelineFileFactory := newRelPathToPipelineFileFactory(relPathToPipelineDirFactory)

  return &_compositionRoot{
    listNamesOfDevOpDirsUcExecuter: newListNamesOfDevOpDirsUcExecuter(),
    listNamesOfPipelineDirsUcExecuter:newListNamesOfPipelineDirsUcExecuter(),
    readDevOpFileUcExecuter:newReadDevOpFileUcExecuter(relPathToDevOpFileFactory),
    readPipelineFileUcExecuter:newReadPipelineFileUcExecuter(relPathToPipelineFileFactory),
    saveDevOpFileUcExecuter:newSaveDevOpFileUcExecuter(relPathToDevOpFileFactory),
    savePipelineFileUcExecuter:newSavePipelineFileUcExecuter(relPathToPipelineFileFactory),
    createDevOpDirUcExecuter:newCreateDevOpDirUcExecuter(relPathToDevOpDirFactory),
    createPipelineDirUcExecuter:newCreatePipelineDirUcExecuter(relPathToPipelineDirFactory),
  }

}

type _compositionRoot struct {
  listNamesOfDevOpDirsUcExecuter    listNamesOfDevOpDirsUcExecuter
  listNamesOfPipelineDirsUcExecuter listNamesOfPipelineDirsUcExecuter
  readDevOpFileUcExecuter           readDevOpFileUcExecuter
  readPipelineFileUcExecuter        readPipelineFileUcExecuter
  saveDevOpFileUcExecuter           saveDevOpFileUcExecuter
  savePipelineFileUcExecuter        savePipelineFileUcExecuter
  createDevOpDirUcExecuter          createDevOpDirUcExecuter
  createPipelineDirUcExecuter       createPipelineDirUcExecuter
}

func (_compositionRoot _compositionRoot) ListNamesOfDevOpDirsUcExecuter(
) listNamesOfDevOpDirsUcExecuter {
  return _compositionRoot.listNamesOfDevOpDirsUcExecuter
}

func (_compositionRoot _compositionRoot) ListNamesOfPipelineDirsUcExecuter(
) listNamesOfPipelineDirsUcExecuter {
  return _compositionRoot.listNamesOfPipelineDirsUcExecuter
}

func (_compositionRoot _compositionRoot) ReadDevOpFileUcExecuter(
) readDevOpFileUcExecuter {
  return _compositionRoot.readDevOpFileUcExecuter
}

func (_compositionRoot _compositionRoot) ReadPipelineFileUcExecuter(
) readPipelineFileUcExecuter {
  return _compositionRoot.readPipelineFileUcExecuter
}

func (_compositionRoot _compositionRoot) SaveDevOpFileUcExecuter(
) saveDevOpFileUcExecuter {
  return _compositionRoot.saveDevOpFileUcExecuter
}

func (_compositionRoot _compositionRoot) SavePipelineFileUcExecuter(
) savePipelineFileUcExecuter {
  return _compositionRoot.savePipelineFileUcExecuter
}

func (_compositionRoot _compositionRoot) CreateDevOpDirUcExecuter(
) createDevOpDirUcExecuter {
  return _compositionRoot.createDevOpDirUcExecuter
}

func (_compositionRoot _compositionRoot) CreatePipelineDirUcExecuter(
) createPipelineDirUcExecuter {
  return _compositionRoot.createPipelineDirUcExecuter
}