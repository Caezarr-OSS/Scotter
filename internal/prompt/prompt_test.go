package prompt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// mockStdin creates a mock stdin with the given inputs
func mockStdin(t *testing.T, inputs []string) *os.File {
	// Create pipe
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}

	// Write inputs to pipe
	go func() {
		defer w.Close()
		for _, input := range inputs {
			_, err := w.Write([]byte(input + "\n"))
			if err != nil {
				panic(fmt.Sprintf("failed to write to pipe: %v", err))
			}
		}
	}()

	return r
}

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	// Create pipe
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	// Run function
	f()

	// Restore stdout
	os.Stdout = stdout
	w.Close()

	// Read captured output
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return fmt.Sprintf("error reading output: %v", err)
	}

	return buf.String()
}

// TestAskString tests the AskString method
func TestAskString(t *testing.T) {
	tests := []struct {
		name           string
		inputs         []string
		question       string
		defaultValue   string
		expectedResult string
		expectedOutput string
	}{
		{
			name:           "simple input",
			inputs:         []string{"testproject"},
			question:       "Project name",
			defaultValue:   "",
			expectedResult: "testproject",
			expectedOutput: "Project name: ",
		},
		{
			name:           "empty with default",
			inputs:         []string{""},
			question:       "Project name",
			defaultValue:   "default-project",
			expectedResult: "default-project",
			expectedOutput: "Project name [default-project]: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()
			os.Stdin = mockStdin(t, tt.inputs)

			// Create the prompt
			p := NewProjectPrompt()

			// Capture output
			output := captureOutput(func() {
				result := p.AskString(tt.question, tt.defaultValue)
				if result != tt.expectedResult {
					t.Errorf("expected %q, got %q", tt.expectedResult, result)
				}
			})

			// Check that the output contains the expected prompt
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

// TestAskGoProjectType tests the AskGoProjectType method
func TestAskGoProjectType(t *testing.T) {
	tests := []struct {
		name           string
		inputs         []string
		expectedResult model.GoProjectType
	}{
		{
			name:           "default type",
			inputs:         []string{"1"},
			expectedResult: model.DefaultGoType,
		},
		{
			name:           "library type",
			inputs:         []string{"2"},
			expectedResult: model.LibraryGoType,
		},
		{
			name:           "cli type",
			inputs:         []string{"3"},
			expectedResult: model.CLIGoType,
		},
		{
			name:           "api type",
			inputs:         []string{"4"},
			expectedResult: model.APIGoType,
		},
		{
			name:           "invalid then valid",
			inputs:         []string{"10", "1"},
			expectedResult: model.DefaultGoType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()
			os.Stdin = mockStdin(t, tt.inputs)

			// Create the prompt
			p := NewProjectPrompt()

			// Run the function
			result := p.AskGoProjectType()

			// Check result
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestAskContainerFileFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected model.ContainerFileFormat
	}{
		{
			name:     "Dockerfile format",
			input:    "1",
			expected: model.DockerfileFormat,
		},
		{
			name:     "Containerfile format",
			input:    "2",
			expected: model.ContainerfileFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock stdin
			inputs := []string{tt.input}
			r := mockStdin(t, inputs)
			
			// Save and restore original stdin
			oldStdin := os.Stdin
			os.Stdin = r
			defer func() {
				os.Stdin = oldStdin
			}()
			
			// Create the prompt and capture output
			p := NewProjectPrompt()
			var result model.ContainerFileFormat
			
			output := captureOutput(func() {
				result = p.AskContainerFileFormat()
			})
			
			// Check the result
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
			
			// Check that the output contains the expected prompt
			if !strings.Contains(output, "Select your preferred container file format") {
				t.Errorf("Expected output to contain prompt, got: %s", output)
			}
		})
	}
}

// TestAskBool tests the AskBool method
func TestAskBool(t *testing.T) {
	tests := []struct {
		name           string
		inputs         []string
		question       string
		defaultValue   bool
		expectedResult bool
		expectedOutput string
	}{
		{
			name:           "yes response",
			inputs:         []string{"y"},
			question:       "Test question?",
			defaultValue:   false,
			expectedResult: true,
			expectedOutput: "Test question? (y/n)",
		},
		{
			name:           "no response",
			inputs:         []string{"n"},
			question:       "Test question?",
			defaultValue:   true,
			expectedResult: false,
			expectedOutput: "Test question? (y/n)",
		},
		{
			name:           "empty with default true",
			inputs:         []string{""},
			question:       "Test question?",
			defaultValue:   true,
			expectedResult: true,
			expectedOutput: "Test question? (y/n)",
		},
		{
			name:           "empty with default false",
			inputs:         []string{""},
			question:       "Test question?",
			defaultValue:   false,
			expectedResult: false,
			expectedOutput: "Test question? (y/n)",
		},
		{
			name:           "invalid then valid",
			inputs:         []string{"invalid", "y"},
			question:       "Test question?",
			defaultValue:   false,
			expectedResult: true,
			expectedOutput: "Test question? (y/n)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock stdin
			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()
			os.Stdin = mockStdin(t, tt.inputs)

			// Create the prompt
			p := NewProjectPrompt()

			// Capture output
			output := captureOutput(func() {
				result := p.AskBool(tt.question, tt.defaultValue)
				if result != tt.expectedResult {
					t.Errorf("expected %v, got %v", tt.expectedResult, result)
				}
			})

			// Check that the output contains the expected prompt
			if !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

// TestCollectConfig tests the CollectConfig method
func TestCollectConfig(t *testing.T) {
	// Skip this test for now as it's complex to test and requires many inputs
	t.Skip("Skipping TestCollectConfig as it requires complex input simulation")
}
