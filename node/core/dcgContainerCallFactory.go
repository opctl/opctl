package core

//go:generate counterfeiter -o ./fakeDCGContainerCallFactory.go --fake-name fakeDCGContainerCallFactory ./ dcgContainerCallFactory

import (
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/dircopier"
	"github.com/golang-utils/filecopier"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type dcgContainerCallFactory interface {
	// Construct constructs a DCGContainerCall
	Construct(
		currentScope map[string]*model.Data,
		scgContainerCall *model.SCGContainerCall,
		containerId string,
		rootOpId string,
		pkgRef string,
	) (*model.DCGContainerCall, error)
}

func newDCGContainerCallFactory(
	dirCopier dircopier.DirCopier,
	fileCopier filecopier.FileCopier,
	interpolate interpolater.Interpolater,
	os ios.IOS,
	rootFSPath string,
) dcgContainerCallFactory {
	return _dcgContainerCallFactory{
		dirCopier:    dirCopier,
		fileCopier:   fileCopier,
		interpolater: interpolate,
		os:           os,
		rootFSPath:   rootFSPath,
	}
}

type _dcgContainerCallFactory struct {
	dirCopier    dircopier.DirCopier
	fileCopier   filecopier.FileCopier
	interpolater interpolater.Interpolater
	os           ios.IOS
	rootFSPath   string
}

func (df _dcgContainerCallFactory) Construct(
	currentScope map[string]*model.Data,
	scgContainerCall *model.SCGContainerCall,
	containerId string,
	rootOpId string,
	pkgRef string,
) (*model.DCGContainerCall, error) {

	dcgContainerCall := &model.DCGContainerCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId: rootOpId,
			PkgRef:   pkgRef,
		},
		Dirs:        map[string]string{},
		EnvVars:     map[string]string{},
		Files:       map[string]string{},
		Sockets:     map[string]string{},
		WorkDir:     scgContainerCall.WorkDir,
		ContainerId: containerId,
		Name:        scgContainerCall.Name,
		Ports:       scgContainerCall.Ports,
	}

	// construct dcg container path
	scratchDirPath := filepath.Join(
		df.rootFSPath,
		"dcg",
		rootOpId,
		"containers",
		containerId,
		"fs",
	)
	if err := df.os.MkdirAll(scratchDirPath, 0700); nil != err {
		return nil, err
	}

	// construct cmd
	for _, cmdEntry := range scgContainerCall.Cmd {
		// interpolate each entry
		dcgContainerCall.Cmd = append(dcgContainerCall.Cmd, df.interpolater.Interpolate(cmdEntry, currentScope))
	}

	// construct dirs
	for scgContainerDirPath, scgContainerDirBind := range scgContainerCall.Dirs {
		if "" == scgContainerDirBind {
			// implicitly bound
			scgContainerDirBind = scgContainerDirPath
		}

		if strings.HasPrefix(scgContainerDirBind, "/") {
			// is bound to pkg path
			dcgContainerCall.Dirs[scgContainerDirPath] = filepath.Join(scratchDirPath, scgContainerDirBind)
			if err := df.dirCopier.OS(
				filepath.Join(pkgRef, scgContainerDirBind),
				dcgContainerCall.Dirs[scgContainerDirPath],
			); nil != err {
				return nil, err
			}
		} else {
			// is bound to variable
			if varData, ok := currentScope[scgContainerDirBind]; ok {
				// bound to input
				dcgContainerCall.Dirs[scgContainerDirPath] = *varData.Dir
			} else {
				// bound to output
				// create placeholder dir on host so the output points to something
				dcgContainerCall.Dirs[scgContainerDirPath] = filepath.Join(scratchDirPath, scgContainerDirPath)
				if err := df.os.MkdirAll(
					dcgContainerCall.Dirs[scgContainerDirPath],
					0700,
				); nil != err {
					return nil, err
				}
			}
		}
	}

	// construct envVars
	for scgContainerEnvVarName, scgContainerEnvVar := range scgContainerCall.EnvVars {
		if "" == scgContainerCall.EnvVars[scgContainerEnvVarName] {
			// implicitly bound
			value := currentScope[scgContainerEnvVarName]
			switch {
			case nil != value.String:
				dcgContainerCall.EnvVars[scgContainerEnvVarName] = *value.String
			case nil != value.Number:
				dcgContainerCall.EnvVars[scgContainerEnvVarName] = strconv.FormatFloat(*value.Number, 'f', -1, 64)
			}
			continue
		}

		// otherwise interpolate value
		dcgContainerCall.EnvVars[scgContainerEnvVarName] = df.interpolater.Interpolate(scgContainerEnvVar, currentScope)

	}

	// construct files
	for scgContainerFilePath, scgContainerFileBind := range scgContainerCall.Files {
		if "" == scgContainerFileBind {
			// implicitly bound
			scgContainerFileBind = scgContainerFilePath
		}

		if strings.HasPrefix(scgContainerFileBind, "/") {
			// is bound to pkg path
			dcgContainerCall.Files[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFileBind)
			if err := df.fileCopier.OS(
				filepath.Join(pkgRef, scgContainerFileBind),
				dcgContainerCall.Files[scgContainerFilePath],
			); nil != err {
				return nil, err
			}
		} else {
			// is bound to variable
			if varData, ok := currentScope[scgContainerFileBind]; ok {
				// bound to input
				dcgContainerCall.Files[scgContainerFilePath] = *varData.File
			} else {
				// bound to output
				// create outputFile on host so the output points to something
				dcgContainerCall.Files[scgContainerFilePath] = filepath.Join(scratchDirPath, scgContainerFilePath)
				// create dir
				if err := df.os.MkdirAll(
					path.Dir(dcgContainerCall.Files[scgContainerFilePath]),
					0700,
				); nil != err {
					return nil, err
				}
				// create file
				var outputFile *os.File
				outputFile, err := df.os.Create(dcgContainerCall.Files[scgContainerFilePath])
				outputFile.Close()
				if nil != err {
					return nil, err
				}
			}
		}
	}

	// construct image
	if scgContainerCallImage := scgContainerCall.Image; scgContainerCallImage != nil {
		dcgContainerCall.Image = &model.DCGContainerCallImage{
			Ref: df.interpolater.Interpolate(scgContainerCall.Image.Ref, currentScope),
		}
		if "" != scgContainerCallImage.PullIdentity && "" != scgContainerCallImage.PullSecret {
			// fallback for deprecated cred format
			scgContainerCallImage.PullCreds = &model.SCGPullCreds{
				Username: scgContainerCallImage.PullIdentity,
				Password: scgContainerCallImage.PullSecret,
			}
		}

		if nil != scgContainerCallImage.PullCreds {
			dcgContainerCall.Image.PullCreds = &model.DCGPullCreds{
				Username: df.interpolater.Interpolate(scgContainerCall.Image.PullCreds.Username, currentScope),
				Password: df.interpolater.Interpolate(scgContainerCall.Image.PullCreds.Password, currentScope),
			}
		}
	}

	// construct sockets
	for scgContainerSocketAddress, scgContainerSocketBind := range scgContainerCall.Sockets {
		if boundArg, ok := currentScope[scgContainerSocketBind]; ok {
			// bound to var
			dcgContainerCall.Sockets[scgContainerSocketAddress] = *boundArg.Socket
		} else if isUnixSocketAddress(scgContainerSocketAddress) {
			// bound to output
			// create outputSocket on host so the output points to something
			dcgHostSocketAddress := filepath.Join(scratchDirPath, scgContainerSocketAddress)
			var outputSocket *os.File
			outputSocket, err := df.os.Create(dcgHostSocketAddress)
			outputSocket.Close()
			if nil != err {
				return nil, err
			}
			if err := os.Chmod(
				dcgHostSocketAddress,
				os.ModeSocket,
			); nil != err {
				return nil, err
			}
			dcgContainerCall.Sockets[scgContainerSocketAddress] = dcgHostSocketAddress
		}
	}

	return dcgContainerCall, nil

}
