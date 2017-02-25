package validate

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
)

// validates an value against a parameter
func (this validate) Param(
	value *model.Data,
	param *model.Param,
) (errs []error) {
	if nil == param {
		// panic as errs represents validation errors not execution errors
		panic("param required")
	}

	switch {
	case nil != param.String:
		errs = this.stringParam(value, param.String)
	case nil != param.Number:
		errs = this.numberParam(value, param.Number)
	case nil != param.Socket:
		errs = this.socketParam(value, param.Socket)
	}
	return
}
