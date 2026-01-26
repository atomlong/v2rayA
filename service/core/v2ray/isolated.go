package v2ray

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/v2rayA/v2rayA/core/v2ray/asset"
	"github.com/v2rayA/v2rayA/core/v2ray/where"
	"github.com/v2rayA/v2rayA/pkg/util/log"
)

// WriteTempConfig writes configuration to a temporary file and returns its path
// along with a cleanup function that removes the file when called.
func WriteTempConfig(content []byte) (path string, cleanup func(), err error) {
	tmpFile, err := ioutil.TempFile("", "v2ray-latency-test-*.json")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp config file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(content); err != nil {
		os.Remove(tmpFile.Name())
		return "", nil, fmt.Errorf("failed to write temp config file: %w", err)
	}

	cleanup = func() {
		// Remove the temporary file
		if err := os.Remove(tmpFile.Name()); err != nil && !os.IsNotExist(err) {
			log.Warn("failed to remove temp config file %v: %v", tmpFile.Name(), err)
		}
	}

	return tmpFile.Name(), cleanup, nil
}

// StartIsolatedProcess starts an isolated Xray-core process with the given configuration file.
// It returns the process handle and a cancel function that stops the process.
func StartIsolatedProcess(ctx context.Context, configPath string) (proc *os.Process, cancel func(), err error) {
	v2rayBinPath, err := where.GetV2rayBinPath()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get v2ray binary path: %w", err)
	}

	dir := filepath.Dir(v2rayBinPath)
	arguments := []string{
		v2rayBinPath,
		"run",
		"--config=" + configPath,
	}

	assetDir := asset.GetV2rayLocationAssetOverride()
	env := append(
		os.Environ(),
		"V2RAY_LOCATION_ASSET="+assetDir,
		"XRAY_LOCATION_ASSET="+assetDir,
	)

	cmd := exec.CommandContext(ctx, v2rayBinPath)
	cmd.Args = arguments
	cmd.Dir = dir
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("failed to start isolated process: %w", err)
	}

	// Give the process a moment to start and potentially fail
	time.Sleep(100 * time.Millisecond)

	// Check if the process has already exited
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}

	cancel = func() {
		if proc != nil {
			if err := proc.Kill(); err != nil && err.Error() != "os: process already finished" {
				log.Trace("failed to kill isolated process: %v", err)
			}
		}
	}

	return cmd.Process, cancel, nil
}