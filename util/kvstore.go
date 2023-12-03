package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type Store struct {
	path string
	data map[string]string
}

func KVStore(path string) (*Store, error) {
	s := &Store{
		path: path,
		data: map[string]string{},
	}
	f, err := os.Open(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&s.data)
	if err != nil {
		return nil, fmt.Errorf("error loading kvstore: %w", err)
	}
	return s, nil
}

func (s *Store) Set(k, v string) {
	s.data[k] = v
}

func (s Store) Get(k string) (string, bool) {
	v, ok := s.data[k]
	return v, ok
}

func (s Store) Save() error {
	f, err := os.Create(s.path)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", s.path, err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(s.data)
	if err != nil {
		return fmt.Errorf("error saving kvstore: %w", err)
	}
	return err
}
