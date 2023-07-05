package models

type Config struct {
	Jobs map[string]Job
}

type Job struct {
	Cmd      string
	Schedule string
}
