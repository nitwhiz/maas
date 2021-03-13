package server

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"io/ioutil"
)

type ConfigMismatchError struct{}

func (e *ConfigMismatchError) Error() string {
	return "config mismatch"
}

type PathExistsError struct {
	Path string
}

func (e *PathExistsError) Error() string {
	return fmt.Sprintf("path %s already exists", e.Path)
}

type NoContainerFoundError struct{}

func (e *NoContainerFoundError) Error() string {
	return "no container found"
}

type Version struct {
	Type        string
	GameVersion string
}

type Settings struct {
	Version Version
}

type VMConfig struct {
	Image       string
	Environment []string
	ExposedPort int
}

type Server struct {
	ConfigPath string `json:"-"`
	VMConfig   VMConfig
	Settings   Settings
}

type ConfigOpts struct {
	IgnoreErrors bool
	NoDefaults   bool
}

// FromConfig reads a config file and populates the server struct accordingly
func FromConfig(filePath string, opts ConfigOpts) (Server, error) {
	srv := Server{}

	bs, err := ioutil.ReadFile(filePath)

	if err != nil {
		return srv, err
	}

	if err = json.Unmarshal(bs, &srv); err != nil {
		return srv, err
	}

	if !opts.IgnoreErrors {
		if err := srv.VerifyMandatoryFields(); err != nil {
			return srv, err
		}
	}

	if !opts.NoDefaults {
		srv.PopulateDefaults()
	}

	srv.ConfigPath = filePath

	return srv, nil
}

// FindContainer tries to find the docker container for this server
func (s *Server) FindContainer(docker *client.Client) (types.Container, error) {
	result := types.Container{}

	labelFilters := filters.NewArgs()

	labelFilters.Add("label", fmt.Sprintf("com.github.nitwhiz.maas.configPath=%s", s.ConfigPath))

	containers, err := docker.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Limit:   1,
		Filters: labelFilters,
	})

	if err != nil {
		return result, err
	}

	if len(containers) == 1 {
		result = containers[0]
	} else {
		return result, &NoContainerFoundError{}
	}

	return result, nil
}

// GetConfigHash generates the SHA-1 sum of the server config
func (s *Server) GetConfigHash() (string, error) {
	config, err := json.Marshal(s)

	if err != nil {
		return "", err
	}

	sha := sha1.New()
	sha.Write(config)
	hash := hex.EncodeToString(sha.Sum(nil))

	return hash, nil
}
