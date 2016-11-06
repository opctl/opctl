package uniquestring

import (
  "github.com/nu7hatch/gouuid"
  "strings"
)

type UniqueStringFactory interface {
  Construct(
  ) (uniqueString string)
}

func NewUniqueStringFactory() UniqueStringFactory {
  return &uniqueStringFactory{}
}

type uniqueStringFactory struct{}

func (this uniqueStringFactory) Construct(
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
