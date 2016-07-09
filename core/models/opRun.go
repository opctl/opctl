package models

import (
  "time"
)

type OpRun interface {
  Id() string
  OpUrl() string
  ParentId() string
  RootId() string
  Timestamp() time.Time
}

func NewOpRun(
opUrl string,
id string,
parentId string,
rootId string,
timestamp time.Time,
) OpRun {

  return opRun{
    opUrl:opUrl,
    id:id,
    parentId:parentId,
    rootId:rootId,
    timestamp:timestamp,
  }

}

type opRun struct {
  id        string
  opUrl     string
  parentId  string
  rootId    string
  timestamp time.Time
}

func (this opRun) Id() string {
  return this.id
}

func (this opRun) OpUrl() string {
  return this.opUrl
}

func (this opRun) ParentId() string {
  return this.parentId
}

func (this opRun) RootId() string {
  return this.rootId
}

func (this opRun) Timestamp() time.Time {
  return this.timestamp
}
