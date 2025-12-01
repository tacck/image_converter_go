package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileSystemManager はファイルシステム操作を提供します
type FileSystemManager struct{}

// NewFileSystemManager は新しいFileSystemManagerを作成します
func NewFileSystemManager() *FileSystemManager {
	return &FileSystemManager{}
}

// DirectoryExists はディレクトリが存在するかチェックします
func (fsm *FileSystemManager) DirectoryExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

// CreateDirectory は出力ディレクトリを作成します
func (fsm *FileSystemManager) CreateDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// ScanDirectory はディレクトリ内のすべてのファイルを走査します
func (fsm *FileSystemManager) ScanDirectory(path string) ([]string, error) {
	var files []string
	
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			fullPath := filepath.Join(path, entry.Name())
			files = append(files, fullPath)
		}
	}
	
	return files, nil
}

// IsImageFile は拡張子に基づいて画像ファイルかどうかを判定します
func (fsm *FileSystemManager) IsImageFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	
	imageExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}
	
	return imageExtensions[ext]
}

// ValidateInputDirectory は入力ディレクトリの存在と読み取り可能性を検証します
func (fsm *FileSystemManager) ValidateInputDirectory(path string) error {
	exists, err := fsm.DirectoryExists(path)
	if err != nil {
		return fmt.Errorf("failed to check input directory: %w", err)
	}
	
	if !exists {
		return fmt.Errorf("input directory does not exist: %s", path)
	}
	
	// 読み取り可能性をテスト
	_, err = os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("input directory is not readable: %w", err)
	}
	
	return nil
}

// EnsureOutputDirectory は出力ディレクトリが存在することを確認し、必要に応じて作成します
func (fsm *FileSystemManager) EnsureOutputDirectory(path string) error {
	exists, err := fsm.DirectoryExists(path)
	if err != nil {
		return fmt.Errorf("failed to check output directory: %w", err)
	}
	
	if !exists {
		if err := fsm.CreateDirectory(path); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}
	
	// 書き込み可能性をテスト
	testFile := filepath.Join(path, ".write_test")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("output directory is not writable: %w", err)
	}
	os.Remove(testFile)
	
	return nil
}
