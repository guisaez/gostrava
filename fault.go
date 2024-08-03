package gostrava

import "encoding/json"

type Fault struct {
	Errors  []Error `json:"errors"`  // The set of specific erros associated with this fault.
	Message string  `json:"message"` // The message of the fault.
}

type Error struct {
	Code     string `json:"code"`     // The code associated with this error.
	Field    string `json:"field"`    // The specific field or aspect of the resource associated witht this error.
	Resource string `json:"resource"` // The type of resource associated with his error.
}

func (f *Fault) Error() string {
	err, _ := json.Marshal(f)
	return string(err)
}
