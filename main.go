package main

import (
  dockercontainerengine "github.com/opspec-io/engine/pkg/containerengine/engines/docker"
  "github.com/opspec-io/engine/tcp"
  "github.com/opspec-io/engine/core"
)

func main() {

  containerEngine, err := dockercontainerengine.New()
  if (nil != err) {
    panic(err)
  }

  tcp.New(
    core.New(
      containerEngine,
    ),
  ).Start()

}
