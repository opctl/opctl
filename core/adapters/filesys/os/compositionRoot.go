package os

type compositionRoot interface {
  ListNamesOfChildDirsUseCase() listNamesOfChildDirsUseCase
  GetBytesOfFileUseCase() getBytesOfFileUseCase
  SaveFileUseCase() saveFileUseCase
  CreateDirUseCase() createDirUseCase
}

func newCompositionRoot(
) compositionRoot {

  return &_compositionRoot{
    listNamesOfChildDirsUseCase: newListNamesOfChildDirsUseCase(),
    getBytesOfFileUseCase:newGetBytesOfFileUseCase(),
    saveFileUseCase:newSaveFileUseCase(),
    createDirUseCase:newCreateDirUseCase(),
  }

}

type _compositionRoot struct {
  listNamesOfChildDirsUseCase listNamesOfChildDirsUseCase
  getBytesOfFileUseCase       getBytesOfFileUseCase
  saveFileUseCase             saveFileUseCase
  createDirUseCase            createDirUseCase
}

func (this _compositionRoot) ListNamesOfChildDirsUseCase(
) listNamesOfChildDirsUseCase {
  return this.listNamesOfChildDirsUseCase
}

func (this _compositionRoot) GetBytesOfFileUseCase(
) getBytesOfFileUseCase {
  return this.getBytesOfFileUseCase
}

func (this _compositionRoot) SaveFileUseCase(
) saveFileUseCase {
  return this.saveFileUseCase
}

func (this _compositionRoot) CreateDirUseCase(
) createDirUseCase {
  return this.createDirUseCase
}