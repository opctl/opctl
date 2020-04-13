module github.com/opctl/opctl/cli

go 1.12

require (
	github.com/appdataspec/sdk-golang v0.0.0-20170917062448-0c0ade7a92f7
	github.com/equinox-io/equinox v1.2.0
	github.com/fatih/color v1.7.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-delve/delve v1.3.2
	github.com/golang-interfaces/ios v0.0.0-20170803194714-da59acb78efc
	github.com/golang-utils/lockfile v0.0.0-20170803195317-342df9650a96
	github.com/golang-utils/pscanary v0.0.0-20170803195345-167b86ee2e7e // indirect
	github.com/gorilla/handlers v1.4.1
	github.com/gorilla/mux v1.7.3
	github.com/jawher/mow.cli v1.1.0
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/maxbrunsfeld/counterfeiter/v6 v6.2.2
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/opctl/opctl/sdks/go v0.0.0-00010101000000-000000000000
	github.com/peterh/liner v1.1.0
	github.com/pkg/errors v0.8.1
	github.com/rakyll/statik v0.1.7-0.20191104211043-6b2f3ee522b6
)

replace github.com/opctl/opctl/sdks/go => ../sdks/go
