package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/johnmikee/manifester/pkg/logger"
)

var log = logger.NewLogger(
	&logger.Config{
		ToFile:  false,
		Level:   logger.DEBUG,
		Service: "test",
		Env:     "dev",
	},
)

func top() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")

	// Capture the output of the command
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Convert the output to a string and trim any leading/trailing whitespace
	topLevelDir := strings.TrimSpace(string(output))

	return topLevelDir, nil
}

func TestCreateDeptManifest(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a Client instance for testing
	client := &Client{
		directory: tempDir + "/manifests",
		log:       &log,
	}

	// Test when department manifest does not exist
	dept := "TestDept"
	// copy the manifests to the temp directory
	top, _ := top()
	err := copyDir(top+"/munki_repo/manifests", client.directory)
	if err != nil {
		t.Errorf("Failed to copy manifests to temp directory: %v", err)
	}

	err = client.createDeptManifest(dept)
	if err != nil {
		t.Errorf("createDeptManifest returned an error: %v", err)
	}

	// Check if the department manifest file exists
	deptFile := fmt.Sprintf("%s/includes/%s", client.directory, dept)
	_, err = os.Stat(deptFile)
	if os.IsNotExist(err) {
		t.Errorf("Department manifest file does not exist")
	}

	// Test when department manifest already exists
	err = client.createDeptManifest(dept)
	if err != nil {
		t.Errorf("createDeptManifest returned an error for an existing department manifest: %v", err)
	}
}

func copyDir(src, dest string) error {
	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	// Get a list of files and subdirectories in the source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectories
			if err := copyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy files
			if err := copyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return err
	}

	return destFile.Sync()
}

func TestRemoveEntries(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create some mock files in the directory
	file1 := "file1.txt"
	file2 := "file2.txt"
	filePath1 := tempDir + "/" + file1
	filePath2 := tempDir + "/" + file2
	_, err := os.Create(filePath1)
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}
	_, err = os.Create(filePath2)
	if err != nil {
		t.Errorf("Failed to create test file: %v", err)
	}

	// Create a Client instance for testing
	client := &Client{
		directory:  tempDir,
		exclusions: []string{file1},
		log:        &log,
	}

	// Test the removeEntries function
	err = client.removeEntries()
	if err != nil {
		t.Errorf("removeEntries returned an error: %v", err)
	}

	// Check if the excluded file is still present
	_, err = os.Stat(filePath1)
	if os.IsNotExist(err) {
		t.Errorf("Excluded file was unexpectedly removed")
	}

	// Check if the non-excluded file was removed
	_, err = os.Stat(filePath2)
	if !os.IsNotExist(err) {
		t.Errorf("Non-excluded file was not removed")
	}
}
