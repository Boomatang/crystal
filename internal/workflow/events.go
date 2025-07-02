package workflow

type Event struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}
