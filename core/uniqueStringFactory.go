package core

import (
  "github.com/nu7hatch/gouuid"
  "strings"
)

type uniqueStringFactory interface {
  Construct(
  ) (uniqueString string)
}

func newUniqueStringFactory() uniqueStringFactory {
  return &_uniqueStringFactory{}
}

type _uniqueStringFactory struct{}

func (this _uniqueStringFactory) Construct(
) (uniqueString string) {
  v4Uuid, err := uuid.NewV4()
  if (nil != err) {
    panic(err)
  }

  uniqueString = strings.Replace(
    v4Uuid.String(),
    "-",
    "",
    -1,
  )

  return
}
