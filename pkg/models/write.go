package models

type ActionValue struct {
	Path  string      `json:"path,omitempty"`  // all components get added together
	Value interface{} `json:"value,omitempty"` // all components get added together
}

type ActionCommand struct {
	Action  []ActionValue `json:"action,omitempty"`  // Write everythign at once
	Comment string        `json:"comment,omitempty"` // Optional comment for the logs
	Page    string        `json:"page,omitempty"`    // The page the user was looking at
	Id      string        `json:"id,omitempty"`      // The page the user was looking at
}
