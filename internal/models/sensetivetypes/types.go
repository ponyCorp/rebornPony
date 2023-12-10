package sensetivetypes

type SensetiveType string

func (s SensetiveType) String() string {
	return string(s)
}

const (
	Forbidden SensetiveType = "forbidden"
	Warn      SensetiveType = "warn"
)
