package cron

import (
	"github.com/tossp/tsgo/pkg/cron/core"
)

// Stop ALL active jobs from running at the next scheduled time
func Stop() {
	go MainCron.Stop()
}

// Remove a specific job from running
// Get EntryID from the list job entries jobrunner.Entries()
// If job is in the middle of running, once the process is finished it will be removed
func Remove(id core.EntryID) {
	MainCron.Remove(id)
}
