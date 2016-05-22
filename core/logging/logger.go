package logging

import "github.com/opctl/engine/core/models"

type Logger func(event models.LogEntryEmittedEvent)
