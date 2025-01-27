package GoSDK

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/websocket"
)

const (
	_REMOTE_SHELL_PREAMBLE = "/edge_shell"
)

type baseMessage struct {
	Type string `json:"type"`
}

type outputHandler interface {
	HandleOutput(data []byte)
	HandleError(err error)
}

func (d *DevClient) OpenRemoteShell(systemKey, edgeName string, handler outputHandler) (*RemoteShell, error) {
	u, err := url.Parse(d.getHttpAddr())
	if err != nil {
		return nil, err
	}

	server := fmt.Sprintf("wss://%s%s", u.Hostname(), _REMOTE_SHELL_PREAMBLE)
	config := new(websocket.Config)
	config.Header = http.Header(make(map[string][]string))
	config.Version = websocket.ProtocolVersionHybi13
	config.Location, err = url.ParseRequestURI(server)
	if err != nil {
		return nil, err
	}

	if d.DevToken == "" {
		return nil, fmt.Errorf("client is not authenticated")
	}

	config.Protocol = []string{"clearblade", systemKey, edgeName, d.DevToken}
	conn, err := websocket.DialConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	if err := openShell(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err := waitForShellToOpen(conn); err != nil {
		_ = conn.Close()
		return nil, err
	}

	go processShellOutput(conn, handler)
	return &RemoteShell{c: conn}, nil
}

func openShell(conn *websocket.Conn) error {
	data, err := json.Marshal(map[string]interface{}{
		"type": "open_shell",
		"options": map[string]interface{}{
			"shell": "/bin/bash",
		},
	})

	if err != nil {
		return fmt.Errorf("could not encode open shell message: %w", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("could not send open shell message: %w", err)
	}

	return nil
}

func waitForShellToOpen(conn *websocket.Conn) error {
	response := baseMessage{}
	if err := websocket.JSON.Receive(conn, &response); err != nil {
		return fmt.Errorf("could not receive shell open response: %w", err)
	}

	if response.Type != "shell_opened" {
		return fmt.Errorf("unexpected response type: %v", response.Type)
	}

	return nil
}

func processShellOutput(conn *websocket.Conn, handler outputHandler) {
	var rawData map[string]any
	if err := websocket.JSON.Receive(conn, &rawData); err != nil {
		handler.HandleError(fmt.Errorf("could not receive json message: %w", err))
		return
	}

	var baseMessage baseMessage
	if err := decodeMapToStruct(rawData, &baseMessage); err != nil {
		handler.HandleError(fmt.Errorf("could not decode websocket frame (%v) to base message: %w", rawData, err))
		return
	}

	switch baseMessage.Type {
	case "output":
		data, ok := rawData["data"].(string)
		if !ok {
			handler.HandleError(fmt.Errorf("could not convert data to string: %T", rawData["data"]))
			return
		}

		decoded := []byte{}
		if _, err := base64.StdEncoding.Decode(decoded, []byte(data)); err != nil {
			handler.HandleError(fmt.Errorf("could not decode base64 data: %w", err))
			return
		}

		handler.HandleOutput(decoded)
	case "error":
		err, ok := rawData["error"].(string)
		if !ok {
			handler.HandleError(fmt.Errorf("could not convert error to string: %T", rawData["error"]))
		} else {
			handler.HandleError(errors.New(err))
		}
	default:
		handler.HandleError(fmt.Errorf("unknown message type: %s", baseMessage.Type))
	}
}

type RemoteShell struct {
	c *websocket.Conn
}

func (r *RemoteShell) Close() error {
	return r.c.Close()
}

func (r *RemoteShell) Write(data []byte) error {
	msg := map[string]interface{}{
		"type":    "exec",
		"command": data,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = r.c.Write(data)
	return err
}
