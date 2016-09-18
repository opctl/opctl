package core

import "github.com/opspec-io/sdk-golang/models"

// A descending sorter for models.Event[] (satisfies sort.Interface)
type EventDescSorter []models.Event

func (this EventDescSorter) Len() int {
  return len(this)
}

func (this EventDescSorter) Swap(i, j int) {
  this[i], this[j] = this[j], this[i]
}

func (this EventDescSorter) Less(i, j int) bool {
  return this[i].Timestamp.After(this[j].Timestamp)
}
