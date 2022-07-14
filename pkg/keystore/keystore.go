package keystore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MarkusFreitag/keepassxc-go/internal"
	"github.com/kevinburke/nacl"
)

const FILENAME = ".keepassxc.keystore"

var (
	ErrEmptyKeystore  = errors.New("keystore does not contain any profiles")
	ErrToManyProfiles = errors.New("keystore has multiple profiles, please specify the one to use")
)

type Profile struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (p *Profile) NaclKey() nacl.Key {
	if p.Key == "" {
		return nil
	}
	return internal.B64ToNaclKey(p.Key)
}

type Keystore struct {
	Profiles []*Profile `json:"profiles"`
}

func Load() (*Keystore, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	storePath := filepath.Join(dir, FILENAME)
	if _, err := os.Stat(storePath); !os.IsNotExist(err) {
		content, err := os.ReadFile(filepath.Join(dir, FILENAME))
		if err != nil {
			return nil, err
		}

		store := new(Keystore)
		err = json.Unmarshal(content, store)
		if err != nil {
			return nil, err
		}

		return store, nil
	}

	return &Keystore{Profiles: make([]*Profile, 0)}, nil
}

func (k *Keystore) Add(prof *Profile) error {
	if p, err := k.Get(prof.Name); p != nil && err == nil {
		return fmt.Errorf("profile named '%s' already exists", prof.Name)
	}
	k.Profiles = append(k.Profiles, prof)
	return nil
}

func (k *Keystore) Get(name string) (*Profile, error) {
	switch len(k.Profiles) {
	case 0:
		return nil, ErrEmptyKeystore
	case 1:
		if profile := k.Profiles[0]; name == "" || profile.Name == name {
			return profile, nil
		}
	default:
		if name == "" {
			return nil, ErrToManyProfiles
		}
		for _, profile := range k.Profiles {
			if profile.Name == name {
				return profile, nil
			}
		}
	}
	return nil, fmt.Errorf("profile named '%s' not found", name)
}

func (k *Keystore) Save() error {
	content, err := json.Marshal(k)
	if err != nil {
		return err
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(dir, FILENAME), content, 0744)
}
