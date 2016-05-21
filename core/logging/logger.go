package logging

import "github.com/open-devops/engine/core/models"

type Logger func(event models.LogEntryEmittedEvent)
