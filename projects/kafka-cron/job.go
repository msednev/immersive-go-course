package main

type Job struct {
	Command string   `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
	Spec    string   `json:"spec,omitempty"`
}
