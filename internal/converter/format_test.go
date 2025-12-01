package converter

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"image-converter/internal/types"
)

// Feature: image-converter, Property 6: ファイル名の一貫性
// Validates: Requirements 5.1, 5.2
// 任意の入力ファイルパスと出力フォーマットに対して、出力ファイル名のベース名は
// 入力ファイル名のベース名と等しく、拡張子は出力フォーマットに対応していなければならない
func TestProperty_FilenameConsistency(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	
	properties := gopter.NewProperties(parameters)
	
	fd := NewFormatDetector()
	
	// 入力ファイル名のジェネレーター
	genInputFilename := gen.OneConstOf(
		"image.jpg",
		"photo.jpeg",
		"picture.png",
		"graphic.webp",
		"animation.gif",
		"bitmap.bmp",
	).Map(func(base string) string {
		// ランダムなパスプレフィックスを追加
		return filepath.Join("input", "dir", base)
	})
	
	// 出力フォーマットのジェネレーター
	genOutputFormat := gen.OneConstOf(
		types.FormatJPEG,
		types.FormatPNG,
		types.FormatWebP,
		types.FormatGIF,
		types.FormatBMP,
	)
	
	properties.Property("output filename preserves base name and changes extension", 
		prop.ForAll(
			func(inputPath string, outputFormat types.ImageFormat) bool {
				// 出力ファイル名を生成
				outputFilename := fd.GenerateOutputFilename(inputPath, outputFormat)
				
				// 入力ファイルのベース名（拡張子なし）を取得
				inputBase := filepath.Base(inputPath)
				inputExt := filepath.Ext(inputBase)
				inputBaseName := strings.TrimSuffix(inputBase, inputExt)
				
				// 出力ファイルのベース名（拡張子なし）を取得
				outputExt := filepath.Ext(outputFilename)
				outputBaseName := strings.TrimSuffix(outputFilename, outputExt)
				
				// プロパティ1: ベース名が保持されている
				if inputBaseName != outputBaseName {
					t.Logf("Base name mismatch: input=%s, output=%s", inputBaseName, outputBaseName)
					return false
				}
				
				// プロパティ2: 拡張子が出力フォーマットに対応している
				expectedExt := getExpectedExtension(outputFormat)
				if outputExt != expectedExt {
					t.Logf("Extension mismatch: expected=%s, got=%s", expectedExt, outputExt)
					return false
				}
				
				return true
			},
			genInputFilename,
			genOutputFormat,
		))
	
	properties.TestingRun(t)
}

// getExpectedExtension は出力フォーマットに対応する拡張子を返します
func getExpectedExtension(format types.ImageFormat) string {
	switch format {
	case types.FormatJPEG:
		return ".jpg"
	case types.FormatPNG:
		return ".png"
	case types.FormatWebP:
		return ".webp"
	case types.FormatGIF:
		return ".gif"
	case types.FormatBMP:
		return ".bmp"
	default:
		return ".jpg"
	}
}

// ユニットテスト: フォーマット検出
func TestDetectFormat(t *testing.T) {
	fd := NewFormatDetector()
	
	tests := []struct {
		name     string
		path     string
		expected types.ImageFormat
		wantErr  bool
	}{
		{"JPEG with .jpg", "image.jpg", types.FormatJPEG, false},
		{"JPEG with .jpeg", "photo.jpeg", types.FormatJPEG, false},
		{"PNG", "picture.png", types.FormatPNG, false},
		{"WebP", "graphic.webp", types.FormatWebP, false},
		{"GIF", "animation.gif", types.FormatGIF, false},
		{"BMP", "bitmap.bmp", types.FormatBMP, false},
		{"Uppercase extension", "IMAGE.JPG", types.FormatJPEG, false},
		{"Unsupported format", "document.pdf", "", true},
		{"No extension", "image", "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format, err := fd.DetectFormat(tt.path)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("DetectFormat() expected error, got nil")
				}
				return
			}
			
			if err != nil {
				t.Errorf("DetectFormat() unexpected error: %v", err)
				return
			}
			
			if format != tt.expected {
				t.Errorf("DetectFormat() = %v, want %v", format, tt.expected)
			}
		})
	}
}

// ユニットテスト: フォーマットサポート確認
func TestIsFormatSupported(t *testing.T) {
	fd := NewFormatDetector()
	
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{"JPEG lowercase", "jpeg", true},
		{"JPEG uppercase", "JPEG", true},
		{"JPG", "jpg", true},
		{"PNG", "png", true},
		{"WebP", "webp", true},
		{"GIF", "gif", true},
		{"BMP", "bmp", true},
		{"Unsupported", "pdf", false},
		{"Empty", "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fd.IsFormatSupported(tt.format)
			if result != tt.expected {
				t.Errorf("IsFormatSupported(%s) = %v, want %v", tt.format, result, tt.expected)
			}
		})
	}
}

// ユニットテスト: フォーマット正規化
func TestNormalizeFormat(t *testing.T) {
	fd := NewFormatDetector()
	
	tests := []struct {
		name     string
		format   string
		expected types.ImageFormat
	}{
		{"JPG to JPEG", "jpg", types.FormatJPEG},
		{"JPG uppercase", "JPG", types.FormatJPEG},
		{"JPEG", "jpeg", types.FormatJPEG},
		{"PNG", "png", types.FormatPNG},
		{"WebP", "webp", types.FormatWebP},
		{"GIF", "gif", types.FormatGIF},
		{"BMP", "bmp", types.FormatBMP},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fd.NormalizeFormat(tt.format)
			if result != tt.expected {
				t.Errorf("NormalizeFormat(%s) = %v, want %v", tt.format, result, tt.expected)
			}
		})
	}
}

// ユニットテスト: 出力パス生成
func TestGenerateOutputPath(t *testing.T) {
	fd := NewFormatDetector()
	
	tests := []struct {
		name         string
		inputPath    string
		outputDir    string
		outputFormat types.ImageFormat
		expected     string
	}{
		{
			"JPEG to PNG",
			"/input/image.jpg",
			"/output",
			types.FormatPNG,
			"/output/image.png",
		},
		{
			"PNG to JPEG",
			"/input/photo.png",
			"/output",
			types.FormatJPEG,
			"/output/photo.jpg",
		},
		{
			"WebP to GIF",
			"/input/dir/graphic.webp",
			"/output/dir",
			types.FormatGIF,
			"/output/dir/graphic.gif",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fd.GenerateOutputPath(tt.inputPath, tt.outputDir, tt.outputFormat)
			if result != tt.expected {
				t.Errorf("GenerateOutputPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}
