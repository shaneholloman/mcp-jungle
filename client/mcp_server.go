package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mcpjungle/mcpjungle/pkg/types"
	"io"
	"net/http"
)

// RegisterServer registers a new MCP server with the registry.
func (c *Client) RegisterServer(server *types.RegisterServerInput) (*types.McpServer, error) {
	u, _ := c.constructAPIEndpoint("/servers")
	body, err := json.Marshal(server)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize server data into JSON: %w", err)
	}

	req, err := c.newRequest(http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status: %d, message: %s", resp.StatusCode, body)
	}

	var registeredServer types.McpServer
	if err := json.NewDecoder(resp.Body).Decode(&registeredServer); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &registeredServer, nil
}

// ListServers fetches the list of registered servers.
func (c *Client) ListServers() ([]*types.McpServer, error) {
	u, _ := c.constructAPIEndpoint("/servers")
	req, err := c.newRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request failed with status: %d, message: %s", resp.StatusCode, body)
	}

	var servers []*types.McpServer
	if err := json.NewDecoder(resp.Body).Decode(&servers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return servers, nil
}

// DeregisterServer deletes a server by name.
func (c *Client) DeregisterServer(name string) error {
	u, _ := c.constructAPIEndpoint("/servers/" + name)
	req, _ := c.newRequest(http.MethodDelete, u, nil)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to %s: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status from server: %s, body: %s", resp.Status, body)
	}
	return nil
}
