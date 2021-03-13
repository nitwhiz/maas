package server

import (
	"fmt"
)

type MandatoryFieldMissingError struct {
	FieldName string
}

func (e *MandatoryFieldMissingError) Error() string {
	return fmt.Sprintf("mandatory field %s missing", e.FieldName)
}

// VerifyMandatoryFields checks if all mandatory fields of this server are not empty
func (s *Server) VerifyMandatoryFields() error {
	if s.VMConfig.ExposedPort == 0 {
		return &MandatoryFieldMissingError{FieldName: "VMConfig.ExposedPort"}
	}

	if s.Settings.Version.Type == "" {
		return &MandatoryFieldMissingError{FieldName: "Settings.GameVersion.Type"}
	}

	if s.Settings.Version.GameVersion == "" {
		return &MandatoryFieldMissingError{FieldName: "Settings.GameVersion.GameVersion"}
	}

	return nil
}

// PopulateDefaults sets the default values to this server
func (s *Server) PopulateDefaults() {
	defaults := Server{
		VMConfig: VMConfig{
			Environment: []string{},
			Image:       "itzg/minecraft-server:java8",
		},
		Settings: Settings{},
	}

	if s.VMConfig.Image == "" {
		s.VMConfig.Image = defaults.VMConfig.Image
	}

	if s.VMConfig.Environment == nil {
		s.VMConfig.Environment = defaults.VMConfig.Environment
	}
}
