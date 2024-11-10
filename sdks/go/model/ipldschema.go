package model

import (
	_ "embed"
	"fmt"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/bindnode"
	"github.com/ipld/go-ipld-prime/schema"
	//"github.com/multiformats/go-multicodec"
)

var (
	//go:embed model.ipldsch
	schemaBytes                       []byte
	AuthAddedPrototype                schema.TypedPrototype
	AuthPrototype                     schema.TypedPrototype
	CallKillRequestedPrototype        schema.TypedPrototype
	ContainerCallImageSpecPrototype   schema.TypedPrototype
	ContainerCallPrototype            schema.TypedPrototype
	ContainerStdErrWrittenToPrototype schema.TypedPrototype
	ContainerStdOutWrittenToPrototype schema.TypedPrototype
	CredsPrototype                    schema.TypedPrototype
	LoopVarsPrototype                 schema.TypedPrototype
	LoopVarsSpecPrototype             schema.TypedPrototype
	SocketParamSpecPrototype          schema.TypedPrototype
)

// examples:
// - https://github.com/ipld/go-ipld-adl-hamt/pull/39/files
// - https://github.com/warptools/warpforge/blob/e3ec637a29aee4874de2ce2e70e9c9e85761ce22/wfapi/ipld.go#L8
func init() {
	ts, err := ipld.LoadSchemaBytes(schemaBytes)
	if err != nil {
		panic(fmt.Errorf("failed to load schema: %w", err))
	}

	AuthPrototype = bindnode.Prototype((*Auth)(nil), ts.TypeByName("Auth"))
	AuthAddedPrototype = bindnode.Prototype((*AuthAdded)(nil), ts.TypeByName("AuthAdded"))
	CallKillRequestedPrototype = bindnode.Prototype((*CallKillRequested)(nil), ts.TypeByName("CallKillRequested"))
	ContainerCallImageSpecPrototype = bindnode.Prototype((*ContainerCallImageSpec)(nil), ts.TypeByName("ContainerCallImageSpec"))
	ContainerCallPrototype = bindnode.Prototype((*ContainerCall)(nil), ts.TypeByName("ContainerCall"))
	ContainerStdErrWrittenToPrototype = bindnode.Prototype((*ContainerStdErrWrittenTo)(nil), ts.TypeByName("ContainerStdErrWrittenTo"))
	ContainerStdOutWrittenToPrototype = bindnode.Prototype((*ContainerStdOutWrittenTo)(nil), ts.TypeByName("ContainerStdOutWrittenTo"))
	CredsPrototype = bindnode.Prototype((*Creds)(nil), ts.TypeByName("Creds"))
	LoopVarsPrototype = bindnode.Prototype((*LoopVars)(nil), ts.TypeByName("LoopVars"))
	LoopVarsSpecPrototype = bindnode.Prototype((*LoopVarsSpec)(nil), ts.TypeByName("LoopVarsSpec"))
	SocketParamSpecPrototype = bindnode.Prototype((*SocketParamSpec)(nil), ts.TypeByName("SocketParamSpec"))
}
