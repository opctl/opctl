package model

import (
  "io"
)

type Arg struct {
  Dir *DirArg
  File *FileArg
  NetSocket *NetSocketArg
  String string
}

type DirArg struct {
  io.ReadCloser
}

type FileArg struct {
  io.ReadCloser
}

type NetSocketArg struct {
  Host string
  Port uint
}
