package main

import (
	"strconv"
	"strings"

	opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"
)

type Options struct {
	Mode  string `json:"mode"`
	Retry int    `json:"retry"`
}
type Inputs struct {
	Env     opctl.String                 `opctl:"env,description='Target environment, dev or prod',default=dev,enum=dev|prod"`
	Count   opctl.Number                 `opctl:"count,default=2,minimum=1,maximum=3,format=integer"`
	Token   opctl.String                 `opctl:"token,secret,minLength=1"`
	Src     opctl.Dir                    `opctl:"src,default=$(./testdata/src)"`
	Config  opctl.File                   `opctl:"config,default=$(./testdata/config.txt)"`
	Prod    opctl.Boolean                `opctl:"wantProd"`
	Names   opctl.Array[string]          `opctl:"names,default='[\"a\",\"b\"]',minItems=1"`
	Options opctl.Object[Options]        `opctl:"options,default='{\"mode\":\"safe\",\"retry\":2}',required=mode|retry"`
	Labels  opctl.Object[map[string]int] `opctl:"labels,default='{\"limit\":3}'"`
}

func main() {
	opctl.Main(func(in Inputs) opctl.Call {
		// Hash the expected token so call events do not contain its literal value.
		args := []string{"sh", "-ce", `test "$#" -eq "$COUNT" && test "$(printf %s "$TOKEN" | sha256sum | cut -d ' ' -f 1)" = 1ec1c26b50d5d3c58d9583181af8076655fe00756bf7285940ba3670f99fcba0 && test -f /src/marker && grep -q hello /config && test "$PROD" = "$WANT_PROD" && test "$NAMES" = a,b && test "$MODE" = safe && test "$RETRY" = 2 && test "$LIMIT" = 3`, "sh"}
		for i := 0; i < int(in.Count.Value()); i++ {
			args = append(args, strconv.Itoa(i))
		}
		prod := "false"
		if in.Env.Value() == "prod" {
			prod = "true"
		}
		options := in.Options.Value()
		return opctl.Container("alpine:3").Cmd(args...).Env("PROD", prod).
			Env("NAMES", strings.Join(in.Names.Value(), ",")).Env("MODE", options.Mode).Env("RETRY", strconv.Itoa(options.Retry)).Env("LIMIT", strconv.Itoa(in.Labels.Value()["limit"])).
			EnvFrom("COUNT", in.Count).EnvFrom("TOKEN", in.Token).EnvFrom("WANT_PROD", in.Prod).Dir("/src", in.Src).File("/config", in.Config)
	})
}
