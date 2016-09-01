package logging

import "github.com/opspec-io/engine/core/models"

type Logger func(event models.LogEntryEmittedEvent)
