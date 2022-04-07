package test

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// Container ..
type Container struct {
	ID   string
	Host string
	Name string
}

type InspectLogs []struct {
	NetworkSettings struct {
		Ports struct {
			TCP3306 []struct {
				HostIP   string `json:"HostIp"`
				HostPort string `json:"HostPort"`
			} `json:"3306/tcp"`
		} `json:"Ports"`
	} `json:"NetworkSettings"`
}

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

const ContainerName = "mysql_test"

// StartContainer ..
func StartContainer(t *testing.T) *Container {
	var doc InspectLogs
	t.Helper()

	startContainerAndMigrate(t)

	cmd := exec.Command("docker", "ps", "-aqf", "name="+ContainerName)

	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &stderr
	_ = cmd.Run()

	ci := out.String()

	cmd = exec.Command("docker", "inspect", ContainerName)
	out.Reset()
	stderr.Reset()
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect container %s: %v", ci, stderr.String())

		return nil
	}

	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json: %v", err)
	}

	network := doc[0].NetworkSettings.Ports.TCP3306[0]

	c := Container{
		ID:   ci,
		Host: network.HostIP + ":" + network.HostPort,
		Name: ContainerName,
	}

	return &c
}

func startContainerAndMigrate(t *testing.T) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	t.Helper()

	tables := basePath + "/tables.sql"

	cmd := exec.Command("bash", basePath+"/run.sh", tables)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("could not build container %v", stderr.String())
	}

	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		var bb bytes.Buffer
		bb.WriteString("completed\n")

		if out.String() == bb.String() {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}
}

// StopContainer stops and removes the specified container.
func StopContainer(t *testing.T, c *Container) {
	t.Helper()

	t.Log("container id", c.ID)

	if err := exec.Command("docker", "kill", ContainerName).Run(); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}

	t.Log("Stopped:", c.ID)

	if err := exec.Command("docker", "container", "rm", "-f", ContainerName).Run(); err != nil {
		t.Fatalf("could not remove container: %v", err)
	}

	t.Log("Removed:", c.ID)
}

// DumpContainerLogs runs "docker logs" against the container and send it to t.Log
func DumpContainerLogs(t *testing.T, c *Container) {
	t.Helper()

	out, _ := exec.Command("docker", "logs", c.ID).CombinedOutput()

	t.Logf("Logs for %s\n%s:", c.ID, out)
}
