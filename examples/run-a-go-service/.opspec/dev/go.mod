module dev

go 1.25

require github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0

// Outside this repository, delete this replace directive and require a released SDK version.
replace github.com/opctl/opctl/sdks/go/opspec/sdk => ../../../../sdks/go/opspec/sdk
