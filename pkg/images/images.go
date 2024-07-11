package images

import (
	"context"
	"fmt"
	"github.com/containers/buildah/define"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/domain/entities"
	"github.com/containers/podman/v4/pkg/specgen"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	podmanArgsBase        = "--log-level=debug system service -t 0"
	podmanExec            = "/usr/bin/podman"
	podmanListenerLogFile = "podman-system-service.log"
	podmanSocketPath      = "/run/podman/podman.sock"
	podmanSocketURI       = "unix://%s"
	dockerfile            = "Dockerfile"
	podmanDirName         = "podman"
	podmanBuildLogFile    = "podman-image-build.log"
)

type Podman struct {
	context context.Context
	out     string
}

// New setups a podman listening service and returns a connected podman client.
//
// Parameters:
//   - out - location for podman to output any logs created as a result of podman commands
func New(out string) (*Podman, error) {
	if err := setupAPIListener(out); err != nil {
		return nil, fmt.Errorf("creating new podman instance: %w", err)
	}

	conn, err := bindings.NewConnection(context.Background(), fmt.Sprintf(podmanSocketURI, podmanSocketPath))
	if err != nil {
		return nil, fmt.Errorf("creating new podman connection: %w", err)
	}

	return &Podman{
		context: conn,
		out:     out,
	}, nil
}

// creates a listening service that answers API calls for Podman (https://docs.podman.io/en/v4.8.3/markdown/podman-system-service.1.html)
// only way to start the service from within a container - https://github.com/containers/podman/tree/v4.8.3/pkg/bindings#starting-the-service-manually
func setupAPIListener(out string) error {
	log.Println("Setting up Podman API listener...")

	logFile, err := os.Create(filepath.Join(out, podmanListenerLogFile))
	if err != nil {
		return fmt.Errorf("creating podman listener log file: %w", err)
	}

	defer logFile.Close()

	cmd := preparePodmanCommand(logFile)
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("error running podman system service: %w", err)
	}

	return waitForPodmanSock()
}

func preparePodmanCommand(out io.Writer) *exec.Cmd {
	args := strings.Split(podmanArgsBase, " ")
	cmd := exec.Command(podmanExec, args...)
	cmd.Stdout = out
	cmd.Stderr = out

	return cmd
}

func waitForPodmanSock() error {
	const (
		retries      = 5
		sleepSeconds = 3
	)

	for i := 0; i < retries; i++ {
		if _, err := os.Stat(podmanSocketPath); err == nil {
			return nil
		}
		time.Sleep(sleepSeconds * time.Second)
	}
	return fmt.Errorf("'%s' file was not created in the expected time", podmanSocketPath)
}
