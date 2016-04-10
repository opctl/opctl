package logging

import "github.com/dev-op-spec/engine/core/models"

type Logger func(event models.LogEntryEmittedEvent)
