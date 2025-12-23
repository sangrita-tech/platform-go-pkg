package version

import (
	"time"
)

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)

type BuildInfo struct {
	Version string    `json:"version"`
	Commit  string    `json:"commit"`
	Date    time.Time `json:"date"`
}

func GetInfo() BuildInfo {
	b := BuildInfo{
		Version: Version,
		Commit:  Commit,
		Date:    time.Now(),
	}

	if Date != "" && Date != "unknown" {
		if t, err := time.Parse(time.RFC3339, Date); err == nil {
			b.Date = t
		}
	}

	return b
}
