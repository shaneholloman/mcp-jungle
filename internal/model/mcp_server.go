package model

import (
	"encoding/json"
	"errors"
	"github.com/mcpjungle/mcpjungle/pkg/types"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type StreamableHTTPConfig struct {
	// URL must be a valid http/https URL.
	URL string `json:"url"`

	// TODO: Store the bearer token in a more secure way, e.g., encrypted instead of plaintext.
	// BearerToken is an optional token used for authenticating requests to the MCP server.
	// If present, it will be used to set the Authorization header in all requests to this MCP server.
	BearerToken string `json:"bearer_token,omitempty"`
}

type StdioConfig struct {
	// Command is the shell command to run the stdio mcp server.
	Command string `json:"command"`

	// Args contains a list of strings that are passed as arguments to the command
	Args []string `json:"args,omitempty"`

	// Env describes the environment variables to pass to the MCP server
	Env map[string]string `json:"env,omitempty"`
}

type McpServer struct {
	gorm.Model

	Name      string                   `json:"name" gorm:"uniqueIndex;not null"`
	Transport types.McpServerTransport `json:"transport" gorm:"type:varchar(30);not null"`

	Description string `json:"description"`

	// Config describes the transport-specific configuration for the MCP server.
	Config datatypes.JSON `json:"config" gorm:"type:jsonb;not null"`
}

func NewStreamableHTTPServer(name, description, url, bearerToken string) (*McpServer, error) {
	if url == "" {
		return nil, errors.New("url is required for streamable HTTP transport")
	}
	config := StreamableHTTPConfig{
		URL:         url,
		BearerToken: bearerToken,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	return &McpServer{
		Name:        name,
		Description: description,
		Transport:   types.TransportStreamableHTTP,
		Config:      configJSON,
	}, nil
}

func NewStdioServer(name, description, command string, args []string, env map[string]string) (*McpServer, error) {
	if command == "" {
		return nil, errors.New("command is required for stdio transport")
	}
	config := StdioConfig{
		Command: command,
		Args:    args,
		Env:     env,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return &McpServer{
		Name:        name,
		Description: description,
		Transport:   types.TransportStdio,
		Config:      datatypes.JSON(configJSON),
	}, nil
}

// GetStreamableHTTPConfig returns the configuration if this is a streamable HTTP server
func (s *McpServer) GetStreamableHTTPConfig() (*StreamableHTTPConfig, error) {
	if s.Transport != types.TransportStreamableHTTP {
		return nil, errors.New("server is not a streamable HTTP transport type")
	}
	var config StreamableHTTPConfig
	if err := json.Unmarshal(s.Config, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// GetStdioConfig returns the configuration if this is a stdio server
func (s *McpServer) GetStdioConfig() (*StdioConfig, error) {
	if s.Transport != types.TransportStdio {
		return nil, errors.New("server is not a stdio transport type")
	}
	var config StdioConfig
	if err := json.Unmarshal(s.Config, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
