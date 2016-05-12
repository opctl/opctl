package models

import "time"

type Event interface {
  Timestamp() time.Time
  CorrelationId() string
}

