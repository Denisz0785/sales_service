package databasetest

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
)

// Container represents the container
type container struct {
	ID   string
	Host string // IP:PORT
}

// startContainer creates a new container.
func startContainer(t *testing.T) *container {
	t.Helper()

	cmd := exec.Command("docker", "run", "-P", "-d", "postgres:11.1-alpine")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not start container: %v", err)
	}

	id := out.String()[:12]
	t.Log("DB Container ID:", id)

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect container: %v", err)
	}

	var doc []struct {
		NetworkSettings struct {
			Port struct {
				TCP5432 []struct {
					HostIP   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				} `json:"5432/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json: %v", err)
	}

	network := doc[0].NetworkSettings.Port.TCP5432[0]

	c := container{
		ID:   id,
		Host: network.HostIP + ":" + network.HostPort,
	}
	t.Log("DB Container Host:", c.Host)

	return &c
}

// stopContainer stops the container.
func stopContainer(t *testing.T, c *container) {
	t.Helper()
	if err := exec.Command("docker", "stop", c.ID).Run(); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}

	t.Log("DB Container Stopped:", c.ID)

	if err := exec.Command("docker", "rm", c.ID, "-v").Run(); err != nil {
		t.Fatalf("could not remove container: %v", err)
	}

	t.Log("DB Container Removed:", c.ID)
}
