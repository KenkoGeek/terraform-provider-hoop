// Copyright (c) HashiCorp, Inc.

package hoop

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Agent struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Mode     string            `json:"mode"`
	Token    string            `json:"token"`
	Status   string            `json:"status"`
	Metadata map[string]string `json:"metadata"`
}

type AgentCreateRequest struct {
	Name string `json:"name"`
	Mode string `json:"mode"`
}

type AgentCreateResponse struct {
	Token string `json:"token"`
}

func (c *Client) GetAgent(nameOrID string) (*Agent, error) {
	req, err := c.newRequest("GET", "/agents/"+nameOrID, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.send(req, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return decodeAgent(resp.Body)
}

func (c *Client) CreateAgent(name, mode string) (*Agent, error) {
	body, err := json.Marshal(AgentCreateRequest{Name: name, Mode: mode})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal agent, reason=%v", err)
	}

	req, err := c.newRequest("POST", "/agents", body)
	if err != nil {
		return nil, err
	}
	resp, err := c.send(req, http.StatusCreated)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var created AgentCreateResponse
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, fmt.Errorf("failed decoding agent create response, reason=%v", err)
	}

	agent, err := c.GetAgent(name)
	if err != nil {
		return nil, err
	}
	agent.Token = created.Token
	return agent, nil
}

func (c *Client) DeleteAgent(nameOrID string) error {
	req, err := c.newRequest("DELETE", "/agents/"+nameOrID, nil)
	if err != nil {
		return err
	}
	return c.sendDiscard(req, http.StatusNoContent)
}

func decodeAgent(responseBody io.Reader) (*Agent, error) {
	var agent Agent
	if err := json.NewDecoder(responseBody).Decode(&agent); err != nil {
		return nil, fmt.Errorf("failed decoding agent resource, reason=%v", err)
	}
	return &agent, nil
}
