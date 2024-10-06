// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package telemetry // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/fileexporter"

import (
	"io"
	"sync"
	"time"
)

type fileWriter struct {
	path  string
	file  io.WriteCloser
	mutex sync.Mutex

	flushInterval time.Duration
	flushTicker   *time.Ticker
	stopTicker    chan struct{}
}

func (w *fileWriter) export(buf []byte) error {
	// Ensure only one write operation happens at a time.
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if _, err := w.file.Write(buf); err != nil {
		return err
	}
	if _, err := io.WriteString(w.file, "\n"); err != nil {
		return err
	}
	return nil
}

// startFlusher starts the flusher.
// It does not check the flushInterval
func (w *fileWriter) startFlusher() {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	ff, ok := w.file.(interface{ flush() error })
	if !ok {
		// Just in case.
		return
	}

	// Create the stop channel.
	w.stopTicker = make(chan struct{})
	// Start the ticker.
	w.flushTicker = time.NewTicker(w.flushInterval)
	go func() {
		for {
			select {
			case <-w.flushTicker.C:
				w.mutex.Lock()
				ff.flush()
				w.mutex.Unlock()
			case <-w.stopTicker:
				w.flushTicker.Stop()
				w.flushTicker = nil
				return
			}
		}
	}()
}

// Start starts the flush timer if set.
func (w *fileWriter) start() {
	if w.flushInterval > 0 {
		w.startFlusher()
	}
}

// Shutdown stops the exporter and is invoked during shutdown.
// It stops the flush ticker if set.
func (w *fileWriter) shutdown() error {
	// Flush and stop the flush ticker.
	if w.flushTicker != nil {
		// Stop the go routine.
		w.mutex.Lock()
		ff, ok := w.file.(interface{ flush() error })
		if !ok {
			// Just in case.
			return nil
		}
		ff.flush()
		close(w.stopTicker)
		w.mutex.Unlock()
	}
	return w.file.Close()
}
