package cron

import (
	"time"

	"github.com/tossp/tsgo/pkg/cron/core"
)

type StatusData struct {
	Id        core.EntryID
	JobRunner *Job
	Next      time.Time
	Prev      time.Time
}

// Return detailed list of currently running recurring jobs
// to remove an entry, first retrieve the ID of entry
func Entries() []core.Entry {
	return MainCron.Entries()
}

func StatusPage() []StatusData {

	ents := MainCron.Entries()

	Statuses := make([]StatusData, len(ents))
	for k, v := range ents {
		Statuses[k].Id = v.ID
		Statuses[k].JobRunner = AddJob(v.Job)
		Statuses[k].Next = v.Next
		Statuses[k].Prev = v.Prev

	}

	// t := template.New("status_page")

	// var data bytes.Buffer
	// t, _ = t.ParseFiles("views/Status.html")

	// t.ExecuteTemplate(&data, "status_page", Statuses())
	return Statuses
}

func StatusJson() map[string]interface{} {

	return map[string]interface{}{
		"jobrunner": StatusPage(),
	}

}

func AddJob(job core.Job) *Job {
	return job.(*Job)
}
