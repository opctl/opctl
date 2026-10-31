package node

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// scratch removes the dcg/<callID> scratch of a root call's tree once nothing will read it.
//
// A call's scratch holds the files and dirs its container mounts, and the values bound from them;
// parents read children's outputs from there, and API callers can read a root's outputs (a dir output
// is a path inside scratch) after it ends. So scratch is owned per root call, never per container:
// it's removed only when the root ended and either the caller asked for that (StartOpReq.CleanScratch,
// which the CLI sets since it reads nothing afterwards) or the node's TTL has passed since the root
// ended. Without either, scratch is kept, as before.
//
// Each root's call IDs are journaled to dcg/<rootID>.scratch as its calls start, so a tree is known
// even after the node restarts.
type scratch struct {
	dcg string
	ttl time.Duration // 0 keeps ended roots' scratch until the caller asks
	mu  sync.Mutex
}

const (
	journalSuffix = ".scratch"
	cleanMark     = "clean" // first line: remove at end
	keepMark      = "keep"
	endedMark     = "ended" // last line: the root ended; the file's mtime is when
)

// callIDs are uniquestring IDs; anything else (the events db) is never removed.
var callID = regexp.MustCompile(`^[0-9a-f]{32}$`)

func newScratch(dataDirPath string, ttl time.Duration) *scratch {
	return &scratch{dcg: filepath.Join(dataDirPath, "dcg"), ttl: ttl}
}

func (s *scratch) journal(rootID string) string {
	return filepath.Join(s.dcg, rootID+journalSuffix)
}

// begin starts rootID's journal.
func (s *scratch) begin(rootID string, clean bool) {
	if s == nil {
		return
	}

	mark := keepMark
	if clean {
		mark = cleanMark
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_ = os.WriteFile(s.journal(rootID), []byte(mark+"\n"), 0o600)
}

// track records callID as part of rootID's tree.
func (s *scratch) track(rootID, id string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.appendLine(rootID, id)
}

// end records that rootID ended, and removes its tree now if its caller asked.
func (s *scratch) end(rootID string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	mark, _, _ := s.read(rootID)
	if mark == cleanMark {
		s.remove(rootID)
		return
	}

	s.appendLine(rootID, endedMark)
}

// sweep runs on node start, before any call: a root that opted in but never recorded its end (the
// node died mid-run) is removed; any other unended root is marked ended now, so the TTL applies to it.
func (s *scratch) sweep(now time.Time) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, rootID := range s.roots() {
		switch mark, _, ended := s.read(rootID); {
		case !ended && mark == cleanMark:
			s.remove(rootID)
		case !ended:
			s.appendLine(rootID, endedMark)
		}
	}

	s.removeExpired(now)
}

// sweepEvery removes roots ended more than the TTL ago, every interval until done; without a TTL it
// returns at once.
func (s *scratch) sweepEvery(done <-chan struct{}, interval time.Duration) {
	if s == nil || s.ttl <= 0 {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case now := <-ticker.C:
			s.mu.Lock()
			s.removeExpired(now)
			s.mu.Unlock()
		}
	}
}

// removeExpired removes roots that ended more than the TTL ago; running roots are never touched.
func (s *scratch) removeExpired(now time.Time) {
	if s.ttl <= 0 {
		return
	}

	for _, rootID := range s.roots() {
		if _, _, ended := s.read(rootID); !ended {
			continue
		}

		if info, err := os.Stat(s.journal(rootID)); err == nil && now.Sub(info.ModTime()) > s.ttl {
			s.remove(rootID)
		}
	}
}

// roots lists the root calls with a journal.
func (s *scratch) roots() []string {
	journals, _ := filepath.Glob(filepath.Join(s.dcg, "*"+journalSuffix))

	ids := []string{}
	for _, j := range journals {
		if id := strings.TrimSuffix(filepath.Base(j), journalSuffix); callID.MatchString(id) {
			ids = append(ids, id)
		}
	}

	return ids
}

func (s *scratch) appendLine(rootID, line string) {
	f, err := os.OpenFile(s.journal(rootID), os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return // no journal: not a root this node started
	}
	defer f.Close()

	_, _ = f.WriteString(line + "\n")
}

// read returns a journal's mark, its call IDs, and whether its root ended.
func (s *scratch) read(rootID string) (mark string, ids []string, ended bool) {
	f, err := os.Open(s.journal(rootID))
	if err != nil {
		return "", nil, false
	}
	defer f.Close()

	lines := bufio.NewScanner(f)
	for lines.Scan() {
		line := strings.TrimSpace(lines.Text())

		switch {
		case mark == "":
			mark = line
		case line == endedMark:
			ended = true
		case callID.MatchString(line):
			ids = append(ids, line)
		}
	}

	return mark, ids, ended
}

// remove deletes rootID's tree, then its journal; a failed removal keeps the journal for the next sweep.
func (s *scratch) remove(rootID string) {
	_, ids, _ := s.read(rootID)

	for _, id := range append(ids, rootID) {
		if err := os.RemoveAll(filepath.Join(s.dcg, id)); err != nil {
			return
		}
	}

	_ = os.Remove(s.journal(rootID))
}
