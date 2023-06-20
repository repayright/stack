package fctl

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type SharedStore struct {
	data    interface{}
	profile *Profile
	config  *Config

	// Those data are not printed in the json output
	additionnalData map[string]interface{}
	//additionnalKeyType 	map[string]

}

func NewSharedStore() *SharedStore {
	return &SharedStore{
		additionnalData: make(map[string]interface{}),
	}
}

// GetSharedData returns the shared data store
func (s *SharedStore) GetData() interface{} {
	return s.data
}

func (s *SharedStore) GetProfile() *Profile {
	return s.profile
}

func (s *SharedStore) GetConfig() *Config {
	return s.config
}

func (s *SharedStore) SetConfig(c *Config) *SharedStore {
	s.config = c
	return s
}

func (s *SharedStore) SetData(data interface{}) *SharedStore {
	s.data = data
	return s
}

func (s *SharedStore) SetProfile(p *Profile) *SharedStore {
	s.profile = p
	return s
}

func (s *SharedStore) SetAdditionnalData(key string, value interface{}) {
	s.additionnalData[key] = value
}

func ShareStoreToJson() ([]byte, error) {
	if (sharedStore.data) == nil {
		errors.New("no data to marshal")
	}

	// Marshal to JSON then print to stdout
	return json.MarshalIndent(sharedStore.data, "", "  ")
}
