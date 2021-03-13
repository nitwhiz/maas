package main

import (
	"fmt"
	"github.com/nitwhiz/maas/internal/server"
	"github.com/nitwhiz/maas/pkg/namesgenerator"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type CreateCmd struct {
	GameVersion string   `kong:"required,short='v',help='Minecraft game version'"`
	Port        int      `kong:"short='p',help='Port to listen on. Defaults to random.'"`
	Name        string   `kong:"short='n',help='Name for the data directory. Defaults to random.'"`
	Type        string   `kong:"short='t',default='vanilla',help='Server type. Defaults to vanilla'"`
	Image       string   `kong:"short='i',help='Container image to use for this server container.'"`
	Environment []string `kong:"short='e',help='Additional environment variables for the runtime.'"`
}

func (c *CreateCmd) Run() error {
	name := c.Name
	port := c.Port

	if name == "" {
		name = namesgenerator.GetRandomName()
	}

	if port == 0 {
		port = r.Intn(20000) + 25565
	}

	s := server.Server{
		VMConfig: server.VMConfig{
			Image:       c.Image,
			Environment: c.Environment,
			ExposedPort: port,
		},
		Settings: server.Settings{
			Version: server.Version{
				Type:        c.Type,
				GameVersion: c.GameVersion,
			},
		},
	}

	if err := s.VerifyMandatoryFields(); err != nil {
		return err
	}

	s.PopulateDefaults()

	if err := s.Create(name); err != nil {
		return err
	}

	fmt.Println(name)

	return nil
}
