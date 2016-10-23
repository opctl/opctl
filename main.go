package main

import (
  "os"
  "fmt"
)

func main() {

  compositionRoot, err := newCompositionRoot(
  )
  if (nil != err) {
    fmt.Fprint(os.Stderr, err)
    os.Exit(1)
  }

  compositionRoot.TcpApi().Start()

}
