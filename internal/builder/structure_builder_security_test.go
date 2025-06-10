package builder_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KoHorizon/ForgeDir/internal/builder"
	"github.com/KoHorizon/ForgeDir/internal/config"
)

// trackingFS tracks what paths were accessed for security testing
type trackingFS struct {
	CreatedFolders []string
	WrittenFiles   []string
	ShouldFail     bool
}

func (t *trackingFS) CreateFolder(path string, perm os.FileMode) error {
	t.CreatedFolders = append(t.CreatedFolders, path)
	if t.ShouldFail {
		return os.ErrPermission
	}
	return nil
}

func (t *trackingFS) WriteFile(path string, content []byte, perm os.FileMode) error {
	t.WrittenFiles = append(t.WrittenFiles, path)
	if t.ShouldFail {
		return os.ErrPermission
	}
	return nil
}

func TestStructureBuilder_PathTraversalPrevention(t *testing.T) {
	tests := []struct {
		name          string
		structure     []config.StructureNode
		expectError   bool
		errorContains string
	}{
		{
			name: "valid simple structure",
			structure: []config.StructureNode{
				{Type: config.TypeDir, Name: "src"},
				{Type: config.TypeFile, Name: "main.go"},
			},
			expectError: false,
		},
		{
			name: "valid nested structure",
			structure: []config.StructureNode{
				{
					Type: config.TypeDir,
					Name: "src",
					Children: []config.StructureNode{
						{Type: config.TypeFile, Name: "main.go"},
					},
				},
			},
			expectError: false,
		},
		{
			name: "attempt parent directory traversal",
			structure: []config.StructureNode{
				{Type: config.TypeFile, Name: "../evil.txt"},
			},
			expectError:   true,
			errorContains: "path traversal sequences (..) are not allowed",
		},
		{
			name: "attempt nested parent directory traversal",
			structure: []config.StructureNode{
				{
					Type: config.TypeDir,
					Name: "src",
					Children: []config.StructureNode{
						{Type: config.TypeFile, Name: "../../etc/passwd"},
					},
				},
			},
			expectError:   true,
			errorContains: "path traversal sequences (..) are not allowed",
		},
		{
			name: "attempt absolute path",
			structure: []config.StructureNode{
				{Type: config.TypeFile, Name: "/etc/passwd"},
			},
			expectError:   true,
			errorContains: "absolute paths are not allowed",
		},
		{
			name: "attempt windows absolute path",
			structure: []config.StructureNode{
				{Type: config.TypeFile, Name: "C:\\Windows\\System32\\evil.exe"},
			},
			expectError:   true,
			errorContains: "absolute paths are not allowed",
		},
		{
			name: "path separator in name",
			structure: []config.StructureNode{
				{Type: config.TypeFile, Name: "src/main.go"},
			},
			expectError:   true,
			errorContains: "path separators are not allowed in names",
		},
		{
			name: "invalid characters",
			structure: []config.StructureNode{
				{Type: config.TypeFile, Name: "file<name.txt"},
			},
			expectError:   true,
			errorContains: "invalid character '<' in path name",
		},
		{
			name: "reserved windows filename",
			structure: []config.StructureNode{
				{Type: config.TypeFile, Name: "CON"},
			},
			expectError:   true,
			errorContains: "reserved filename not allowed",
		},
		{
			name: "empty name",
			structure: []config.StructureNode{
				{Type: config.TypeFile, Name: ""},
			},
			expectError:   true,
			errorContains: "path name cannot be empty",
		},
		{
			name: "whitespace only name",
			structure: []config.StructureNode{
				{Type: config.TypeDir, Name: "   "},
			},
			expectError:   true,
			errorContains: "path name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &trackingFS{}
			sb := builder.NewStructureBuilder(fs)

			cfg := &config.Config{Structure: tt.structure}
			root := t.TempDir() // Use a real temporary directory

			err := sb.Build(cfg, root)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tt.errorContains, err.Error())
				}
				// Ensure no files were created when validation fails
				if len(fs.CreatedFolders) > 0 || len(fs.WrittenFiles) > 0 {
					t.Error("no files should be created when validation fails")
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestStructureBuilder_PathContainment(t *testing.T) {
	// This test ensures that even if somehow malicious paths get through,
	// they're still contained within the project directory

	tempDir := t.TempDir()

	// Create a tracking filesystem that records all accessed paths
	fs := &trackingFS{}
	sb := builder.NewStructureBuilder(fs)

	// Valid structure that should pass validation
	cfg := &config.Config{
		Structure: []config.StructureNode{
			{
				Type: config.TypeDir,
				Name: "src",
				Children: []config.StructureNode{
					{Type: config.TypeFile, Name: "main.go"},
					{Type: config.TypeFile, Name: "config.json"},
				},
			},
			{Type: config.TypeFile, Name: "README.md"},
		},
	}

	err := sb.Build(cfg, tempDir)
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}

	// Verify all created paths are within the temp directory
	allPaths := append(fs.CreatedFolders, fs.WrittenFiles...)
	for _, path := range allPaths {
		// Convert to absolute path for comparison
		absPath, err := filepath.Abs(path)
		if err != nil {
			t.Fatalf("failed to get absolute path for %s: %v", path, err)
		}

		absTempDir, err := filepath.Abs(tempDir)
		if err != nil {
			t.Fatalf("failed to get absolute temp dir: %v", err)
		}

		// Check if the path is within the temp directory
		relPath, err := filepath.Rel(absTempDir, absPath)
		if err != nil {
			t.Fatalf("failed to get relative path: %v", err)
		}

		if strings.HasPrefix(relPath, "..") {
			t.Errorf("path %s is outside the project directory %s", absPath, absTempDir)
		}
	}

	// Verify expected structure was created
	expectedPaths := []string{
		filepath.Join(tempDir, "src"),                // folder
		filepath.Join(tempDir, "src", "main.go"),     // file
		filepath.Join(tempDir, "src", "config.json"), // file
		filepath.Join(tempDir, "README.md"),          // file
	}

	if len(fs.CreatedFolders) != 1 {
		t.Errorf("expected 1 folder, got %d: %v", len(fs.CreatedFolders), fs.CreatedFolders)
	}

	if len(fs.WrittenFiles) != 3 {
		t.Errorf("expected 3 files, got %d: %v", len(fs.WrittenFiles), fs.WrittenFiles)
	}

	// Check that all expected paths were created
	allCreated := append(fs.CreatedFolders, fs.WrittenFiles...)
	for _, expected := range expectedPaths {
		found := false
		for _, created := range allCreated {
			if created == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected path %s was not created", expected)
		}
	}
}
