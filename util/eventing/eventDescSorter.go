package eventing

import "github.com/opspec-io/sdk-golang/pkg/models"

// A descending sorter for github.com/opspec-io/sdk-golang/pkg/models/models.Event[] (satisfies sort.Interface)
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
