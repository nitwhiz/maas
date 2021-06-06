package server

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Create writes the server config to a new data directory, if it does not exist already
func (s *Server) Create(name string) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	dataDir := path.Join(cwd, name)

	if _, err := os.Stat(dataDir); err == nil {
		return &PathExistsError{Path: dataDir}
	}

	err = os.MkdirAll(dataDir, 0775)

	if err != nil {
		return err
	}

	metaFile := path.Join(dataDir, "maas.json")

	bs, err := json.MarshalIndent(s, "", "  ")

	if err != nil {
		return err
	}

	return ioutil.WriteFile(metaFile, bs, 0664)
}
