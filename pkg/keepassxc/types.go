package keepassxc

import (
	"encoding/json"
	"errors"
	"strings"
)

// StringFields need this Prefix (incl. at least one space) to be returned by http api.
const StringFieldKeyPrefix = "KPH: "

var ErrInvalidResponse = errors.New("invalid response does not include entries")

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
		return nil, ErrInvalidResponse
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

// Example entry as returned from the api, if every field has some value (incl. tags, expire date and totp):
// [{"group":"foo","login":"myname","name":"bar","password":"myPa$$w0rd","stringFields":[{"KPH: bar":"barval"},{"KPH: foo":"fooval"}],
// "totp":"175413","uuid":"92bfee4f24614ef9ac6e1f440eff3292"}]
// If expired, the entry will actually not be found/returned by the api.
type Entry struct {
	Name     string   `json:"name"`
	Login    string   `json:"login"`
	Password Password `json:"password"`
	Totp     string   `json:"totp"`
	Group    string   `json:"group"`
	UUID     string   `json:"uuid"`
	Fields   Fields   `json:"stringFields"`
}

type Fields []map[string]Password

// ToMap converts the structure returned from the api (list of single entry maps)
// to a simple key value map (keys without the leading "KPH: ").
// The map uses Password values, since values could contain sensible data.
func (f Fields) ToMap() map[string]Password {
	fMap := make(map[string]Password, len(f))
	for _, e := range f {
		for k, v := range e {
			fMap[strings.TrimSpace(strings.TrimPrefix(k, StringFieldKeyPrefix))] = v
		}
	}
	return fMap
}

func (f Fields) String() string {
	v, _ := json.Marshal(f.ToMap())
	return string(v)
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
