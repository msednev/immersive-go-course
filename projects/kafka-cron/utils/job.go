package utils

type Job struct {
	Command string `json:"command,omitempty"`
	Spec    string `json:"spec,omitempty"`
}
