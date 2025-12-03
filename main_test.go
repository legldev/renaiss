package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// Integration-style test that feeds a temp file to main() and checks the printed counts.
func TestMainCountsWords(t *testing.T) {
	tmp, err := os.CreateTemp("", "wordfreq-*.txt")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(tmp.Name())

	content := "The cat and the hat.\nCat-cat cat! Dog."
	if _, err := tmp.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if err := tmp.Close(); err != nil {
		t.Fatalf("close temp file: %v", err)
	}

	// Simulate user input of the filename.
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	rIn, wIn, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdin: %v", err)
	}
	if _, err := wIn.WriteString(tmp.Name() + "\n"); err != nil {
		t.Fatalf("write stdin: %v", err)
	}
	wIn.Close()
	os.Stdin = rIn

	// Capture stdout.
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe stdout: %v", err)
	}
	os.Stdout = wOut

	main()

	// Close writer and read captured output.
	wOut.Close()
	var out bytes.Buffer
	if _, err := out.ReadFrom(rOut); err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	output := out.String()
	for _, expected := range []string{"cat: 4", "hat: 1", "dog: 1"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("output missing %q\n--- output ---\n%s", expected, output)
		}
	}

	if strings.Contains(output, "the:") {
		t.Fatalf("stopword 'the' should be filtered out, but output was:\n%s", output)
	}
}
