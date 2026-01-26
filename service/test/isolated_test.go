package test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/v2rayA/v2rayA/core/v2ray"
)

func TestWriteTempConfig(t *testing.T) {
	// Test successful write and cleanup
	content := []byte(`{"test": "config"}`)
	path, cleanup, err := v2ray.WriteTempConfig(content)
	if err != nil {
		t.Fatalf("WriteTempConfig failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Temp config file was not created: %s", path)
	}

	// Verify content
	readContent, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read temp config file: %v", err)
	}
	if string(readContent) != string(content) {
		t.Fatalf("Content mismatch: expected %s, got %s", content, readContent)
	}

	// Cleanup
	cleanup()

	// Verify file was removed
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Temp config file was not cleaned up: %s", path)
	}
}

func TestWriteTempConfig_FailOnWrite(t *testing.T) {
	// Test cleanup is called on write failure
	// This is hard to test without mocking, so we skip it
	t.Skip("Cannot reliably test write failure without mocking")
}

func TestStartIsolatedProcess_InvalidConfig(t *testing.T) {
	// Test with invalid config - process should start but fail quickly
	invalidConfig := []byte(`{"invalid": "config"}`)
	path, cleanup, err := v2ray.WriteTempConfig(invalidConfig)
	if err != nil {
		t.Fatalf("Failed to create temp config: %v", err)
	}
	defer cleanup()

	// Start process with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	proc, procCancel, err := v2ray.StartIsolatedProcess(ctx, path)
	if err != nil {
		// This is expected if v2ray binary is not available
		t.Skipf("v2ray binary not available: %v", err)
	}
	defer procCancel()

	if proc == nil {
		t.Fatal("Process handle should not be nil")
	}

	// Process may have already exited due to invalid config
	// Wait a bit to see if it exits
	time.Sleep(500 * time.Millisecond)

	// Cancel should be safe to call even if process already exited
	procCancel()
}

func TestStartIsolatedProcess_ContextCancel(t *testing.T) {
	// Test context cancellation
	validConfig := []byte(`{
		"log": {"loglevel": "warning"},
		"inbounds": [],
		"outbounds": [{"protocol": "freedom"}]
	}`)
	path, cleanup, err := v2ray.WriteTempConfig(validConfig)
	if err != nil {
		t.Fatalf("Failed to create temp config: %v", err)
	}
	defer cleanup()

	// Start process with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	proc, procCancel, err := v2ray.StartIsolatedProcess(ctx, path)
	if err != nil {
		// This is expected if v2ray binary is not available
		t.Skipf("v2ray binary not available: %v", err)
	}

	if proc != nil {
		procCancel()
	}

	// Context should be cancelled
	select {
	case <-ctx.Done():
		// Expected
	case <-time.After(200 * time.Millisecond):
		t.Error("Context should have been cancelled")
	}
}

func TestWriteTempConfig_FileNaming(t *testing.T) {
	// Test that temp files have correct naming pattern
	content := []byte(`{"test": "config"}`)
	path, cleanup, err := v2ray.WriteTempConfig(content)
	if err != nil {
		t.Fatalf("WriteTempConfig failed: %v", err)
	}
	defer cleanup()

	base := filepath.Base(path)
	if !strings.HasPrefix(base, "v2ray-latency-test-") {
		t.Fatalf("Temp config file has incorrect prefix: %s", base)
	}
	if !strings.HasSuffix(base, ".json") {
		t.Fatalf("Temp config file has incorrect suffix: %s", base)
	}
}

func TestWriteTempConfig_MultipleCalls(t *testing.T) {
	// Test that multiple calls create different files
	content := []byte(`{"test": "config"}`)
	path1, cleanup1, err := v2ray.WriteTempConfig(content)
	if err != nil {
		t.Fatalf("WriteTempConfig failed: %v", err)
	}
	defer cleanup1()

	path2, cleanup2, err := v2ray.WriteTempConfig(content)
	if err != nil {
		cleanup1()
		t.Fatalf("WriteTempConfig failed: %v", err)
	}
	defer cleanup2()

	if path1 == path2 {
		t.Fatalf("Multiple calls should create different files: %s == %s", path1, path2)
	}
}

func TestStartIsolatedProcess_Cleanup(t *testing.T) {
	// Test that cleanup function works correctly
	validConfig := []byte(`{
		"log": {"loglevel": "warning"},
		"inbounds": [],
		"outbounds": [{"protocol": "freedom"}]
	}`)
	path, cleanup, err := v2ray.WriteTempConfig(validConfig)
	if err != nil {
		t.Fatalf("Failed to create temp config: %v", err)
	}
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	proc, procCancel, err := v2ray.StartIsolatedProcess(ctx, path)
	if err != nil {
		t.Skipf("v2ray binary not available: %v", err)
	}

	// Process should be running
	if proc == nil {
		t.Fatal("Process should be running")
	}

	// Call cleanup - should not panic
	procCancel()

	// Calling cleanup again should be safe
	procCancel()
}