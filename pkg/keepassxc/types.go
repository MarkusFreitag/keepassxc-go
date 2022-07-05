package keepassxc

import (
	"encoding/json"
	"errors"
	"strings"
)

type Message map[string]interface{}

type Response map[string]interface{}

func (r Response) entries() (Entries, error) {
	var data []byte
	if msg, ok := r["message"]; ok {
		if v, ok := msg.(map[string]interface{})["entries"]; ok {
			var err error
			data, err = json.Marshal(v)
			if err != nil {
				return nil, err
			}
		}
	}
	if len(data) == 0 {
		return nil, errors.New("invalid response does not include entries")
	}

	var entries Entries
	err := json.Unmarshal(data, &entries)
	return entries, err
}

type Password string

func (p Password) String() string {
	return "*****"
}

func (p Password) Plaintext() string {
	return string(p)
}

type Entry struct {
	Name     string     `json:"name"`
	Login    string     `json:"login"`
	Password Password   `json:"password"`
	Group    string     `json:"group"`
	UUID     string     `json:"uuid"`
	Fields   Fields     `json:"stringFields"`
	Expired  BoolString `json:"expired"`
}

type Fields []string

func (f Fields) String() string {
	return strings.Join(f, ",")
}

type BoolString bool

func (bs *BoolString) UnmarshalJSON(b []byte) error {
	if strings.ToLower(strings.Trim(string(b), `"`)) == "true" {
		*bs = true
	}
	return nil
}

type Entries []*Entry

type DBGroup struct{}
