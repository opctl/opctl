package core

import (
  "github.com/nu7hatch/gouuid"
  "strings"
)

type uniqueStringFactory interface {
  Construct(
  ) (uniqueString string, err error)
}

func newUniqueStringFactory() uniqueStringFactory {
  return &_uniqueStringFactory{}
}

type _uniqueStringFactory struct{}

func (this _uniqueStringFactory) Construct(
) (uniqueString string, err error) {
  v4Uuid, err := uuid.NewV4()
  if (nil != err) {
    return
  }

  uniqueStringValue := strings.Replace(
    v4Uuid.String(),
    "-",
    "",
    -1,
  )

  uniqueString = uniqueStringValue

  return
}
