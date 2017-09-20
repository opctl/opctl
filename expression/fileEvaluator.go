package expression

import (
  "fmt"
  "github.com/opspec-io/sdk-golang/data"
  "github.com/opspec-io/sdk-golang/expression/interpolater"
  "github.com/opspec-io/sdk-golang/model"
  "path/filepath"
  "strings"
)

type fileEvaluator interface {
  // EvalToFile evaluates an expression to a file value
  // expression must be a type supported by data.CoerceToFile
  // scratchDir will be used as the containing dir if file creation necessary
  //
  // Examples of valid file expressions:
  // scope ref: $(scope-ref)
  // scope ref w/ path expansion: $(scope-ref/file.txt)
  // scope ref w/ deprecated path expansion: $(scope-ref)/file.txt
  // pkg fs ref: $(/pkg-fs-ref)
  // pkg fs ref w/ path expansion: $(/pkg-fs-ref/file.txt)
  EvalToFile(
    scope map[string]*model.Value,
    expression interface{},
    pkgHandle model.PkgHandle,
    scratchDir string,
  ) (*model.Value, error)
}

func newFileEvaluator() fileEvaluator {
  return _fileEvaluator{
    data:         data.New(),
    interpolater: interpolater.New(),
  }
}

type _fileEvaluator struct {
  data         data.Data
  interpolater interpolater.Interpolater
}

func (etf _fileEvaluator) EvalToFile(
  scope map[string]*model.Value,
  expression interface{},
  pkgHandle model.PkgHandle,
  scratchDir string,
) (*model.Value, error) {
  switch expression := expression.(type) {
  case float64:
    return etf.data.CoerceToFile(&model.Value{Number: &expression}, scratchDir)
  case map[string]interface{}:
    return etf.data.CoerceToFile(&model.Value{Object: expression}, scratchDir)
  case string:
    if ref, ok := tryResolveExplicitRef(expression, scope); ok {
      // scope ref w/out expansion
      return etf.data.CoerceToFile(ref, scratchDir)
    } else if strings.HasPrefix(expression, "/") {
      // deprecated pkg path ref
      pkgFsRefPath, err := etf.interpolater.Interpolate(
        expression,
        scope,
        pkgHandle,
      )
      if nil != err {
        return nil, err
      }

      pkgFsRefPath = filepath.Join(pkgHandle.Ref(), pkgFsRefPath)
      return &model.Value{File: &pkgFsRefPath}, nil
    }

    possibleRefCloserIndex := strings.Index(expression, string(interpolater.RefCloser))
    var basePath string

    if strings.HasPrefix(expression, string(interpolater.Operator+interpolater.RefOpener)) && possibleRefCloserIndex > 0 {
      possibleRef := expression[2:possibleRefCloserIndex]
      interpolatedPossibleRef, err := etf.interpolater.Interpolate(possibleRef, scope, pkgHandle)
      if nil != err {
        return nil, err
      }

      if dcgValue, ok := scope[interpolatedPossibleRef]; ok && nil != dcgValue.Dir {
        // scope ref w/ deprecated path expansion
        basePath = *dcgValue.Dir
        // trim initial dir ref & interpolate remaining expression
        expression = expression[possibleRefCloserIndex+1:]
      }
    }

    result, err := etf.interpolater.Interpolate(expression, scope, pkgHandle)
    if nil != err {
      return nil, err
    }

    expandedPath := filepath.Join(basePath, result)

    return etf.data.CoerceToFile(&model.Value{File: &expandedPath}, scratchDir)
  }

  return nil, fmt.Errorf("unable to evaluate %+v to file; unsupported type", expression)

}
