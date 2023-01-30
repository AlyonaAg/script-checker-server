package model

type Script struct {
	ID            int64
	URL           string
	Script        string
	Result        bool
	DangerPercent float64
	VirusTotal    string
}

type Scripts []*Script

type ListScriptsFilter struct {
	Page  int64
	Limit int64
}
