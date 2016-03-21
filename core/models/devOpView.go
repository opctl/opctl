package models

type DevOpView struct {
  Description string
  Name        string
}

func NewDevOpView(
description string,
name string,
) *DevOpView {

  return &DevOpView{
    Description:description,
    Name:name,
  }

}


