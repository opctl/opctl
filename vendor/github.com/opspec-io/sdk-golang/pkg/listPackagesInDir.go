package pkg

//go:generate counterfeiter -o ./fakeListPackagesInDir.go --fake-name fakeListPackagesInDir ./ listPackagesInDir

import (
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

func (this pkg) ListPackagesInDir(
	dirPath string,
) (
	ops []*model.PackageView,
	err error,
) {

	childFileInfos, err := this.ioUtil.ReadDir(dirPath)
	if nil != err {
		return
	}

	for _, childFileInfo := range childFileInfos {
		packageView, err := this.viewFactory.Construct(
			path.Join(dirPath, childFileInfo.Name()),
		)
		if nil == err {
			ops = append(ops, packageView)
		}
	}

	return

}
