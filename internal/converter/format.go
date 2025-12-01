package converter

import (
	"fmt"
	"path/filepath"
	"strings"

	"image-converter/internal/types"
)

// FormatDetector は画像フォーマットの検出と変換を提供します
type FormatDetector struct{}

// NewFormatDetector は新しいFormatDetectorを作成します
func NewFormatDetector() *FormatDetector {
	return &FormatDetector{}
}

// DetectFormat はファイル拡張子から画像フォーマットを検出します
func (fd *FormatDetector) DetectFormat(path string) (types.ImageFormat, error) {
	ext := strings.ToLower(filepath.Ext(path))
	
	switch ext {
	case ".jpg", ".jpeg":
		return types.FormatJPEG, nil
	case ".png":
		return types.FormatPNG, nil
	case ".webp":
		return types.FormatWebP, nil
	case ".gif":
		return types.FormatGIF, nil
	case ".bmp":
		return types.FormatBMP, nil
	default:
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}
}

// IsFormatSupported はフォーマットがサポートされているかチェックします
func (fd *FormatDetector) IsFormatSupported(format string) bool {
	normalizedFormat := strings.ToLower(format)
	
	supportedFormats := map[string]bool{
		"jpeg": true,
		"jpg":  true,
		"png":  true,
		"webp": true,
		"gif":  true,
		"bmp":  true,
	}
	
	return supportedFormats[normalizedFormat]
}

// NormalizeFormat はフォーマット文字列を正規化します
func (fd *FormatDetector) NormalizeFormat(format string) types.ImageFormat {
	normalizedFormat := strings.ToLower(format)
	
	// jpgをjpegに正規化
	if normalizedFormat == "jpg" {
		return types.FormatJPEG
	}
	
	return types.ImageFormat(normalizedFormat)
}

// GenerateOutputFilename は出力ファイル名を生成します
// ベース名を保持し、拡張子を出力フォーマットに変更します
func (fd *FormatDetector) GenerateOutputFilename(inputPath string, outputFormat types.ImageFormat) string {
	// ファイル名からベース名（拡張子なし）を取得
	baseName := filepath.Base(inputPath)
	ext := filepath.Ext(baseName)
	baseNameWithoutExt := strings.TrimSuffix(baseName, ext)
	
	// 新しい拡張子を追加
	var newExt string
	switch outputFormat {
	case types.FormatJPEG:
		newExt = ".jpg"
	case types.FormatPNG:
		newExt = ".png"
	case types.FormatWebP:
		newExt = ".webp"
	case types.FormatGIF:
		newExt = ".gif"
	case types.FormatBMP:
		newExt = ".bmp"
	default:
		newExt = ".jpg" // デフォルト
	}
	
	return baseNameWithoutExt + newExt
}

// GenerateOutputPath は完全な出力パスを生成します
func (fd *FormatDetector) GenerateOutputPath(inputPath, outputDir string, outputFormat types.ImageFormat) string {
	outputFilename := fd.GenerateOutputFilename(inputPath, outputFormat)
	return filepath.Join(outputDir, outputFilename)
}
