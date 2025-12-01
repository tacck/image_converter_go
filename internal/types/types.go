package types

import "image"

// Config はCLI設定を表します
type Config struct {
	InputDir    string
	OutputDir   string
	Scale       float64
	Width       int
	Height      int
	Format      string
	JPEGQuality int
}

// ResizeSpec は画像のリサイズ仕様を表します
type ResizeSpec struct {
	Scale  float64 // 倍率指定（0より大きい、0の場合は未指定）
	Width  int     // 幅のピクセル指定（0の場合は未指定）
	Height int     // 高さのピクセル指定（0の場合は未指定）
}

// ImageFormat はサポートされる画像フォーマットを表します
type ImageFormat string

const (
	FormatJPEG ImageFormat = "jpeg"
	FormatPNG  ImageFormat = "png"
	FormatWebP ImageFormat = "webp"
	FormatGIF  ImageFormat = "gif"
	FormatBMP  ImageFormat = "bmp"
)

// ConversionStats は変換処理の統計情報を表します
type ConversionStats struct {
	Total   int
	Success int
	Failed  int
	Skipped int
}

// ConversionResult は個別の変換結果を表します
type ConversionResult struct {
	SourcePath string
	OutputPath string
	Success    bool
	Error      error
}

// ImageProcessor は画像処理のインターフェースを定義します
type ImageProcessor interface {
	Load(path string) (image.Image, error)
	Resize(img image.Image, spec ResizeSpec) image.Image
	Save(img image.Image, path string, format string, quality int) error
}

// FileScanner はファイルシステム操作のインターフェースを定義します
type FileScanner interface {
	ScanDirectory(path string) ([]string, error)
	IsImageFile(path string) bool
}
