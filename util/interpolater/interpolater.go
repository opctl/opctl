package interpolater

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
	"strings"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Interpolater

const (
	operator  = '$'
	refOpener = '('
	refCloser = ')'
)

type Interpolater interface {
	Interpolate(
		template string,
		scope map[string]*model.Value,
		pkgHandle model.PkgHandle,
	) (*model.Value, error)
}

func New() Interpolater {
	return _Interpolater{
		deReferencer: newDeReferencer(),
	}
}

type _Interpolater struct {
	deReferencer
}

func (itp _Interpolater) Interpolate(
	template string,
	scope map[string]*model.Value,
	pkgHandle model.PkgHandle,
) (*model.Value, error) {
	possibleRefCloserIndex := strings.Index(template, ")")
	var dir *model.Value
	var file *model.Value

	if strings.HasPrefix(template, "$(") && possibleRefCloserIndex > 0 {
		possibleRef := template[2:possibleRefCloserIndex]
		if dcgValue, ok := scope[possibleRef]; ok {
			if len(template)-1 == possibleRefCloserIndex {
				// this is an explicit ref; return value
				return dcgValue, nil
			}
			if nil != dcgValue.Dir {
				// this is a dir expansion
				dir = dcgValue
				// trim initial dir ref & interpolate remaining template
				template = template[possibleRefCloserIndex+1:]
			}
		}
	}

	// interpolate remaining template
	refBuffer := []byte{}
	i := 0
	for i < len(template) {
		switch {
		case operator == template[i]:
			result, consumed, err := itp.tryDeRef(template[i+1:], scope, pkgHandle)
			if nil != err {
				return nil, err
			}
			refBuffer = append(refBuffer, result...)
			i += consumed
		default:
			refBuffer = append(refBuffer, template[i])
		}

		// always increment loop counter
		i++
	}

	if nil != dir {
		expandedPath := filepath.Join(*dir.Dir, string(refBuffer))
		return &model.Value{Dir: &expandedPath}, nil
	}
	if nil != file {
		expandedPath := filepath.Join(*dir.File, string(refBuffer))
		return &model.Value{File: &expandedPath}, nil
	}

	valueAsString := string(refBuffer)
	return &model.Value{String: &valueAsString}, nil
}

// tryDeRef tries to de reference from possibleRef.
// It's assumed possibleRef doesn't contain the initial operator.
//
// returns the de referenced value (if any), number of bytes consumed, and any err
func (itp _Interpolater) tryDeRef(
	possibleRef string,
	scope map[string]*model.Value,
	pkgHandle model.PkgHandle,
) (string, int, error) {
	refBuffer := []byte{}
	i := 0
	for i < len(possibleRef) {
		switch {
		case refCloser == possibleRef[i]:
			if len(refBuffer) > 0 && refOpener == refBuffer[0] {
				result, ok, err := itp.deReferencer.DeReference(string(refBuffer[1:]), scope, pkgHandle)
				if nil != err {
					return "", 0, err
				}
				if ok {
					return result, i + 1, err
				}
			}
			refBuffer = append(refBuffer, possibleRef[i])
		case operator == possibleRef[i]:
			result, consumed, err := itp.tryDeRef(possibleRef[i+1:], scope, pkgHandle)
			if nil != err {
				return "", 0, err
			}
			refBuffer = append(refBuffer, result...)
			i += consumed
		default:
			refBuffer = append(refBuffer, possibleRef[i])
		}

		// always increment loop counter
		i++
	}

	return "$" + string(refBuffer), len(possibleRef), nil
}
