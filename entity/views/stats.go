package views

// StatsEntry is entry response object for a stats query
type StatsEntry struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// NewStatsEntry is a dummy constructor
func NewStatsEntry(name string, count int) *StatsEntry {
	return &StatsEntry{
		Name:  name,
		Count: count,
	}
}
