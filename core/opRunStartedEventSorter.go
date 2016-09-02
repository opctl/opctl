package core

import "github.com/opspec-io/engine/core/models"

// A descending sorter for models.OpRunStartedEvent[] (satisfies sort.Interface)
type OpRunStartedEventDescSorter []models.OpRunStartedEvent

func (this OpRunStartedEventDescSorter) Len() int {
  return len(this)
}

func (this OpRunStartedEventDescSorter) Swap(i, j int) {
  this[i], this[j] = this[j], this[i]
}

func (this OpRunStartedEventDescSorter) Less(i, j int) bool {
  return this[i].Timestamp().After(this[j].Timestamp())
}
