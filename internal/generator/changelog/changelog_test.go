package changelog

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Caezarr-OSS/Scotter/internal/model"
)

// TestNewGenerator tests the constructor
func TestNewGenerator(t *testing.T) {
	cfg := &model.Config{
		ProjectName: "testproject",
	}
	generator := NewGenerator(cfg)
	
	if generator == nil {
		t.Fatal("expected generator to be created, got nil")
	}
	
	if generator.Config != cfg {
		t.Errorf("expected Config to be %v, got %v", cfg, generator.Config)
	}
}

// TestGenerateChangelog tests the changelog generation
func TestGenerateChangelog(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to the temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer os.Chdir(originalDir)
	
	// Create a config with changelog enabled
	cfg := &model.Config{
		ProjectName: "testproject",
		Features: model.Features{
			GitHub: model.GitHubFeatures{
				GenerateChangelog: true,
			},
		},
	}
	
	// Create and run the generator
	generator := NewGenerator(cfg)
	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	
	// Check that CHANGELOG.md was created
	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	if _, err := os.Stat(changelogPath); os.IsNotExist(err) {
		t.Error("expected CHANGELOG.md to exist")
	}
	
	// Check content
	content, err := os.ReadFile(changelogPath)
	if err != nil {
		t.Fatalf("failed to read CHANGELOG.md: %v", err)
	}
	
	// Verify key sections
	expectedSections := []string{
		"# Changelog",
		"The format is based on [Keep a Changelog]",
		"## [Unreleased]",
		"### Added",
		"### Changed",
		"## [0.1.0]",
	}
	
	for _, section := range expectedSections {
		if !strings.Contains(string(content), section) {
			t.Errorf("expected CHANGELOG.md to contain %q", section)
		}
	}
}

// TestGenerateCommitlintConfig tests the commitlint config generation
func TestGenerateCommitlintConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to the temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer os.Chdir(originalDir)
	
	// Create a config with commitlint enabled
	cfg := &model.Config{
		ProjectName: "testproject",
		Features: model.Features{
			GitHub: model.GitHubFeatures{
				UseCommitLint: true,
			},
		},
	}
	
	// Create and run the generator
	generator := NewGenerator(cfg)
	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	
	// Check that .commitlintrc.js was created
	commitlintPath := filepath.Join(tempDir, ".commitlintrc.js")
	if _, err := os.Stat(commitlintPath); os.IsNotExist(err) {
		t.Error("expected .commitlintrc.js to exist")
	}
	
	// Check content
	content, err := os.ReadFile(commitlintPath)
	if err != nil {
		t.Fatalf("failed to read .commitlintrc.js: %v", err)
	}
	
	// Verify key sections
	expectedSections := []string{
		"module.exports = {",
		"extends: ['@commitlint/config-conventional']",
		"'type-enum':",
		"'feat',",
		"'fix',",
	}
	
	for _, section := range expectedSections {
		if !strings.Contains(string(content), section) {
			t.Errorf("expected .commitlintrc.js to contain %q", section)
		}
	}
}

// TestGenerateCommitMsgHook tests the Git hook generation
func TestGenerateCommitMsgHook(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to the temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer os.Chdir(originalDir)
	
	// Create Git hooks directory
	if err := os.MkdirAll(filepath.Join(tempDir, ".git"), 0755); err != nil {
		t.Fatalf("failed to create .git directory: %v", err)
	}
	
	// Create a config with commitlint enabled
	cfg := &model.Config{
		ProjectName: "testproject",
		Features: model.Features{
			GitHub: model.GitHubFeatures{
				UseCommitLint: true,
			},
		},
	}
	
	// Create and run the generator
	generator := NewGenerator(cfg)
	if err := generator.GenerateCommitMsgHook(); err != nil {
		t.Fatalf("GenerateCommitMsgHook() failed: %v", err)
	}
	
	// Check that commit-msg hook was created
	hookPath := filepath.Join(tempDir, ".git", "hooks", "commit-msg")
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Error("expected commit-msg hook to exist")
	}
	
	// Check content
	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("failed to read commit-msg hook: %v", err)
	}
	
	// Verify key sections
	expectedSections := []string{
		"#!/bin/sh",
		"commitlint",
		"npx",
	}
	
	for _, section := range expectedSections {
		if !strings.Contains(string(content), section) {
			t.Errorf("expected commit-msg hook to contain %q", section)
		}
	}
	
	// Check executable permission for Unix systems
	// This test will be skipped on Windows
	if runtime.GOOS != "windows" {
		info, err := os.Stat(hookPath)
		if err != nil {
			t.Fatalf("failed to stat commit-msg hook: %v", err)
		}
		if info.Mode()&0100 == 0 {
			t.Error("expected commit-msg hook to be executable")
		}
	}
}

// TestGenerateWithBothDisabled tests that nothing is generated when features are disabled
func TestGenerateWithBothDisabled(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scotter-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// Change to the temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer os.Chdir(originalDir)
	
	// Create a config with both features disabled
	cfg := &model.Config{
		ProjectName: "testproject",
		Features: model.Features{
			GitHub: model.GitHubFeatures{
				GenerateChangelog: false,
				UseCommitLint:     false,
			},
		},
	}
	
	// Create and run the generator
	generator := NewGenerator(cfg)
	if err := generator.Generate(); err != nil {
		t.Fatalf("Generate() failed: %v", err)
	}
	
	// Check that neither file was created
	changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
	if _, err := os.Stat(changelogPath); !os.IsNotExist(err) {
		t.Error("expected CHANGELOG.md to not exist when disabled")
	}
	
	commitlintPath := filepath.Join(tempDir, ".commitlintrc.js")
	if _, err := os.Stat(commitlintPath); !os.IsNotExist(err) {
		t.Error("expected .commitlintrc.js to not exist when disabled")
	}
}
