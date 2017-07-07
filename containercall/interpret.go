package containercall

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func (cc _ContainerCall) Interpret(
	scope map[string]*model.Value,
	scgContainerCall *model.SCGContainerCall,
	containerId string,
	rootOpId string,
	pkgPath string,
) (*model.DCGContainerCall, error) {

	dcgContainerCall := &model.DCGContainerCall{
		DCGBaseCall: &model.DCGBaseCall{
			RootOpId: rootOpId,
			PkgRef:   pkgPath,
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
		cc.rootFSPath,
		"dcg",
		rootOpId,
		"containers",
		containerId,
		"fs",
	)
	if err := cc.os.MkdirAll(scratchDirPath, 0700); nil != err {
		return nil, err
	}

	// construct cmd
	for _, cmdEntry := range scgContainerCall.Cmd {
		// interpolate each entry
		dcgContainerCall.Cmd = append(dcgContainerCall.Cmd, cc.interpolater.Interpolate(cmdEntry, scope))
	}

	// interpret dirs
	var err error
	dcgContainerCall.Dirs, err = cc.dirs.Interpret(pkgPath, scope, scgContainerCall.Dirs, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret envVars
	dcgContainerCall.EnvVars, err = cc.envVars.Interpret(scope, scgContainerCall.EnvVars)
	if nil != err {
		return nil, err
	}

	// interpret files
	dcgContainerCall.Files, err = cc.files.Interpret(pkgPath, scope, scgContainerCall.Files, scratchDirPath)
	if nil != err {
		return nil, err
	}

	// interpret image
	dcgContainerCall.Image, err = cc.image.Interpret(scope, scgContainerCall.Image)
	if nil != err {
		return nil, err
	}

	// interpret sockets
	dcgContainerCall.Sockets, err = cc.sockets.Interpret(scope, scgContainerCall.Sockets, scratchDirPath)

	return dcgContainerCall, err

}
