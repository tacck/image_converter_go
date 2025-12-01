package filesystem

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirectoryExists(t *testing.T) {
	fsm := NewFileSystemManager()
	
	// 一時ディレクトリを作成
	tmpDir := t.TempDir()
	
	// 存在するディレクトリ
	exists, err := fsm.DirectoryExists(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected directory to exist")
	}
	
	// 存在しないディレクトリ
	nonExistent := filepath.Join(tmpDir, "nonexistent")
	exists, err = fsm.DirectoryExists(nonExistent)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected directory to not exist")
	}
	
	// ファイルを作成してテスト（ディレクトリではない）
	testFile := filepath.Join(tmpDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	
	exists, err = fsm.DirectoryExists(testFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if exists {
		t.Error("expected file to not be recognized as directory")
	}
}

func TestCreateDirectory(t *testing.T) {
	fsm := NewFileSystemManager()
	tmpDir := t.TempDir()
	
	// 新しいディレクトリを作成
	newDir := filepath.Join(tmpDir, "newdir")
	if err := fsm.CreateDirectory(newDir); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	
	// ディレクトリが作成されたことを確認
	exists, err := fsm.DirectoryExists(newDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected directory to be created")
	}
	
	// ネストされたディレクトリを作成
	nestedDir := filepath.Join(tmpDir, "parent", "child", "grandchild")
	if err := fsm.CreateDirectory(nestedDir); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}
	
	exists, err = fsm.DirectoryExists(nestedDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected nested directory to be created")
	}
}

func TestScanDirectory(t *testing.T) {
	fsm := NewFileSystemManager()
	tmpDir := t.TempDir()
	
	// テストファイルを作成
	testFiles := []string{"file1.txt", "file2.jpg", "file3.png"}
	for _, name := range testFiles {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}
	
	// サブディレクトリを作成（スキップされるべき）
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}
	
	// ディレクトリをスキャン
	files, err := fsm.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("failed to scan directory: %v", err)
	}
	
	// ファイル数を確認（サブディレクトリは含まれない）
	if len(files) != len(testFiles) {
		t.Errorf("expected %d files, got %d", len(testFiles), len(files))
	}
	
	// すべてのファイルが含まれていることを確認
	fileMap := make(map[string]bool)
	for _, f := range files {
		fileMap[filepath.Base(f)] = true
	}
	
	for _, name := range testFiles {
		if !fileMap[name] {
			t.Errorf("expected file %s to be in scan results", name)
		}
	}
}

func TestScanDirectory_NonExistent(t *testing.T) {
	fsm := NewFileSystemManager()
	
	// 存在しないディレクトリをスキャン
	_, err := fsm.ScanDirectory("/nonexistent/directory")
	if err == nil {
		t.Error("expected error when scanning non-existent directory")
	}
}

func TestIsImageFile(t *testing.T) {
	fsm := NewFileSystemManager()
	
	tests := []struct {
		path     string
		expected bool
	}{
		{"image.jpg", true},
		{"image.jpeg", true},
		{"image.JPG", true},
		{"image.JPEG", true},
		{"image.png", true},
		{"image.PNG", true},
		{"image.gif", true},
		{"image.bmp", true},
		{"image.webp", true},
		{"document.txt", false},
		{"document.pdf", false},
		{"archive.zip", false},
		{"noextension", false},
		{"image.tiff", false},
	}
	
	for _, tt := range tests {
		result := fsm.IsImageFile(tt.path)
		if result != tt.expected {
			t.Errorf("IsImageFile(%q) = %v, expected %v", tt.path, result, tt.expected)
		}
	}
}

func TestValidateInputDirectory(t *testing.T) {
	fsm := NewFileSystemManager()
	tmpDir := t.TempDir()
	
	// 存在するディレクトリ
	if err := fsm.ValidateInputDirectory(tmpDir); err != nil {
		t.Errorf("unexpected error for valid directory: %v", err)
	}
	
	// 存在しないディレクトリ
	nonExistent := filepath.Join(tmpDir, "nonexistent")
	if err := fsm.ValidateInputDirectory(nonExistent); err == nil {
		t.Error("expected error for non-existent directory")
	}
}

func TestEnsureOutputDirectory(t *testing.T) {
	fsm := NewFileSystemManager()
	tmpDir := t.TempDir()
	
	// 存在しないディレクトリを作成
	newDir := filepath.Join(tmpDir, "output")
	if err := fsm.EnsureOutputDirectory(newDir); err != nil {
		t.Fatalf("failed to ensure output directory: %v", err)
	}
	
	// ディレクトリが作成されたことを確認
	exists, err := fsm.DirectoryExists(newDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !exists {
		t.Error("expected output directory to be created")
	}
	
	// 既に存在するディレクトリ
	if err := fsm.EnsureOutputDirectory(newDir); err != nil {
		t.Errorf("unexpected error for existing directory: %v", err)
	}
}
