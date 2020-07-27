package statx

// State definition. provide state management
type State struct {
	logTime string
	User
}

// New State instance
func New() *State {
	return &State{}
}
