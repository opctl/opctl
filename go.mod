module github.com/opctl/opctl

go 1.15

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/appdataspec/sdk-golang v0.0.0-20170917062448-0c0ade7a92f7
	github.com/blang/semver v3.5.1+incompatible
	github.com/containerd/containerd v1.2.7 // indirect
	github.com/containers/image/v5 v5.0.0
	github.com/containers/storage v1.13.5 // indirect
	github.com/dgraph-io/badger/v3 v3.2103.2
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v17.12.0-ce-rc1.0.20200916142827-bd33bbf0497b+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/fatih/color v1.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-delve/delve v1.3.2
	github.com/go-git/go-git/v5 v5.2.0
	github.com/golang-interfaces/github.com-gorilla-websocket v0.0.0-20190604222234-f7c71d63fecb
	github.com/golang-interfaces/ihttp v0.0.0-20170731143308-228dd9eedf13
	github.com/golang-interfaces/iioutil v0.0.0-20170803194630-7d1c0886acdc // indirect
	github.com/golang-interfaces/ios v0.0.0-20170803194714-da59acb78efc // indirect
	github.com/golang-utils/dircopier v0.0.0-20170803194507-75bc9e581ed2
	github.com/golang-utils/filecopier v0.0.0-20170803193939-16f96e9dcff4
	github.com/golang-utils/lockfile v0.0.0-20170803195317-342df9650a96
	github.com/golang-utils/pscanary v0.0.0-20170803195345-167b86ee2e7e // indirect
	github.com/gorilla/handlers v1.4.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1
	github.com/jawher/mow.cli v1.1.0
	github.com/jfbus/httprs v0.0.0-20180614100821-7861a1150846
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.10 // indirect
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.3
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.1 // indirect
	github.com/mjibson/esc v0.2.0
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.15.0
	github.com/peterh/liner v1.1.0
	github.com/rakyll/statik v0.1.7-0.20191104211043-6b2f3ee522b6
	github.com/rhysd/go-github-selfupdate v1.2.3
	github.com/satori/go.uuid v0.0.0-20181028125025-b2ce2384e17b
	github.com/sethgrid/pester v0.0.0-20180430140037-03e26c9abbbf
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
	golang.org/x/net v0.0.0-20210805182204-aaa1db679c0d
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
	golang.org/x/sys v0.0.0-20210806184541-e5e7981a1069 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1
	gotest.tools v2.2.0+incompatible // indirect
	k8s.io/api v0.19.1
	k8s.io/apimachinery v0.19.1
	k8s.io/client-go v0.19.1
	k8s.io/utils v0.0.0-20201104234853-8146046b121e // indirect
)

replace github.com/opctl/opctl/sdks/go => ./sdks/go
