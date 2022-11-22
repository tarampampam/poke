package run

import (
	"fmt"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/tarampampam/poke/internal/js/events"
)

type scriptRunningStat struct {
	events   []events.Event
	duration time.Duration
	err      error
}

type OverallRunningStats struct {
	mu sync.Mutex
	m  map[string]*scriptRunningStat

	summaryDuration time.Duration
}

func NewOverallRunningStats() *OverallRunningStats {
	return &OverallRunningStats{m: make(map[string]*scriptRunningStat)}
}

func (r *OverallRunningStats) SetEvents(scriptName string, events []events.Event) {
	r.mu.Lock()

	if v, ok := r.m[scriptName]; ok {
		v.events = events
	} else {
		r.m[scriptName] = &scriptRunningStat{events: events}
	}

	r.mu.Unlock()
}

func (r *OverallRunningStats) SetDuration(scriptName string, d time.Duration) {
	r.mu.Lock()

	if v, ok := r.m[scriptName]; ok {
		v.duration = d
	} else {
		r.m[scriptName] = &scriptRunningStat{duration: d}
	}

	r.mu.Unlock()
}
func (r *OverallRunningStats) SetError(scriptName string, err error) {
	r.mu.Lock()

	if v, ok := r.m[scriptName]; ok {
		v.err = err
	} else {
		r.m[scriptName] = &scriptRunningStat{err: err}
	}

	r.mu.Unlock()
}

func (r *OverallRunningStats) SetSummaryDuration(d time.Duration) {
	r.mu.Lock()

	r.summaryDuration = d

	r.mu.Unlock()
}

func (r *OverallRunningStats) ToConsole() string { // TODO make this printer great again!
	tbl := table.NewWriter()
	tbl.SetStyle(table.StyleLight)
	tbl.AppendHeader(table.Row{"File", "Events count", "Success", "Duration"})

	r.mu.Lock()

	for name, stat := range r.m {
		var success string

		if stat.err == nil {
			success = "yes"
		} else {
			success = stat.err.Error()
		}

		tbl.AppendRow(table.Row{
			name,
			len(stat.events),
			success,
			stat.duration.Round(time.Millisecond).String(),
		})
	}

	tbl.AppendFooter(table.Row{
		fmt.Sprintf("Total files: %d", len(r.m)),
		"",
		"",
		fmt.Sprintf("Elapsed time: %s", r.summaryDuration.Round(time.Millisecond)),
	})

	r.mu.Unlock()

	return tbl.Render()
}
