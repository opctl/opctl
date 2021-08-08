package opgraph

import "time"

// LoadingSpinner has a string representation that indicates something's loading
type LoadingSpinner interface {
	String() string
}

// DotLoadingSpinner is a LoadingSpinner that uses brail dots that rotate
type DotLoadingSpinner struct {
	state       int
	lastChanged time.Time
}

var loadingRunes = []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}

func (l DotLoadingSpinner) String() string {
	now := time.Now()
	ms := now.UnixNano() / int64(time.Millisecond)
	r := loadingRunes[(ms/int64(100))%int64(len(loadingRunes))]
	return string(r)
}

// StaticLoadingSpinner is a LoadingSpinner that doesn't spin
type StaticLoadingSpinner struct{}

func (StaticLoadingSpinner) String() string {
	return "⋰"
}
