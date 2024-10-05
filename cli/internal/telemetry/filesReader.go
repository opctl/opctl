package telemetry

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type filesReader struct {
	path         string
	unmarshaller *unmarshaller
	fileLock     *flock.Flock
}

func newFilesReader(path string, fileLock *flock.Flock) *filesReader {
	return &filesReader{
		path:         path,
		fileLock:     fileLock,
		unmarshaller: newUnmarshaller(),
	}
}

func (a *filesReader) Read() ([]pmetric.Metrics, error) {
	paths, err := a.listOldFiles()
	if err != nil {
		return nil, err
	}
	// append the current file
	paths = append(paths, a.path)

	locked, err := a.fileLock.TryLock()
	if err != nil {
		return nil, err
	}
	if locked {
		allMetrics := []pmetric.Metrics{}
		for _, path := range paths {
			fileMetrics, _ := a.readFile(path)
			allMetrics = append(allMetrics, fileMetrics...)
		}

		return accumulate(allMetrics), a.fileLock.Unlock()
	}

	return nil, fmt.Errorf("cannot lock the file: %w", err)
}

func (a *filesReader) Clear() error {
	paths, err := a.listOldFiles()
	if err != nil {
		return err
	}

	locked, err := a.fileLock.TryLock()
	if err != nil {
		return err
	}

	if locked {
		for _, path := range paths {
			if err := os.Remove(path); err != nil {
				a.fileLock.Unlock()
				return fmt.Errorf("cannot remove the file: %w", err)
			}
		}

		f, err := os.OpenFile(a.path, os.O_TRUNC|os.O_WRONLY, 0666)
		if err != nil {
			a.fileLock.Unlock()
			return fmt.Errorf("failed to clear file: %w", err)
		}
		err = f.Close()
		if err != nil {
			a.fileLock.Unlock()
			return fmt.Errorf("failed to close file: %w", err)
		}

		return a.fileLock.Unlock()
	}

	return nil
}

func (a *filesReader) readFile(path string) ([]pmetric.Metrics, error) {
	reader, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot open the file: %w", err)
	}
	defer reader.Close()

	allMetrics := []pmetric.Metrics{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		bytes := scanner.Bytes()
		metrics, err := a.unmarshaller.metrics(bytes)
		if err != nil {
			// Skip the line if it cannot be unmarshalled
			continue
		} else if metrics != nil {
			allMetrics = append(allMetrics, *metrics)
		}
	}
	if err := scanner.Err(); err != nil {
		return allMetrics, fmt.Errorf("cannot read the file: %w", err)
	}

	return allMetrics, nil
}

func (a *filesReader) listOldFiles() ([]string, error) {
	paths := []string{}

	locked, err := a.fileLock.TryLock()
	if err != nil {
		return paths, err
	}
	if locked {
		dir := filepath.Dir(a.path)
		allFiles, err := os.ReadDir(dir)
		if err != nil {
			a.fileLock.Unlock()
			return paths, fmt.Errorf("can't read log file directory: %s", err)
		}

		prefix, ext := a.prefixAndExt()
		filesWithTime := byFormatTime{}
		for _, f := range allFiles {
			if f.IsDir() {
				continue
			}

			if t, err := a.timeFromName(f.Name(), prefix, ext); err == nil {
				filesWithTime = append(filesWithTime, withFormatTime{t, f})
				continue
			}
		}

		sort.Sort(byFormatTime(filesWithTime))

		paths := []string{}
		for _, f := range filesWithTime {
			paths = append(paths, filepath.Join(dir, f.f.Name()))
		}
		return paths, a.fileLock.Unlock()
	}
	return paths, fmt.Errorf("cannot lock the file: %w", err)
}

func (a *filesReader) prefixAndExt() (prefix, ext string) {
	filename := filepath.Base(a.path)
	ext = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(ext)] + "-"
	return prefix, ext
}

// timeFromName extracts the formatted time from the filename by stripping off
// the filename's prefix and extension. This prevents someone's filename from
// confusing time.parse.
func (a *filesReader) timeFromName(filename, prefix, ext string) (time.Time, error) {
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, errors.New("mismatched extension")
	}
	ts := filename[len(prefix) : len(filename)-len(ext)]
	return time.Parse("2006-01-02T15-04-05.000", ts)
}
