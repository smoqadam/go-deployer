package main

import (
	"encoding/json"
)

type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Commit struct {
	ID       string   `json:"id"`
	Message  string   `json:"message"`
	Authod   Author   `json:"author"`
	Added    []string `json:"added"`
	Removed  []string `json:"removed"`
	Modified []string `json:"modified"`
}

type Repository struct {
	Name string `json:"full_name"`
}

type Payload struct {
	Repo    Repository `json:"repository"`
	Commits []Commit   `json:"commits"`
}

func NewPayload(input []byte) (*Payload, error) {
	p := &Payload{}
	if err := json.Unmarshal(input, &p); err != nil {
		return p, err
	}
	return p, nil
}
