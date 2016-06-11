package models

import (
  netUrl "net/url"
  "strings"
  "errors"
  "fmt"
  "encoding/json"
)

func NewUrl(
urlStr string,
) (url *Url, err error) {

  var parsedUrl *netUrl.URL
  parsedUrl, err = netUrl.Parse(urlStr)
  if (nil != err) {
    return
  }

  url = &Url{parsedUrl}

  if (parsedUrl.IsAbs() && strings.ToLower(parsedUrl.Scheme) != "file") {
    err = errors.New(
      fmt.Sprintf(
        "Urls referencing remote resources are currently unsupported. `%v` is not supported.",
        urlStr,
      ),
    )
  }

  return

}

type Url struct {
  *netUrl.URL
}

func (this *Url) MarshalJSON() ([]byte, error) {
  return json.Marshal((*this).String())
}

func (this *Url) UnmarshalJSON(b []byte) (err error) {

  urlString := new(string)
  err = json.Unmarshal(b, urlString)
  if (nil != err) {
    return
  }

  var constructedUrl *Url
  constructedUrl, err = NewUrl(*urlString)
  if (nil != err) {
    return
  }

  *this = *constructedUrl

  return

}