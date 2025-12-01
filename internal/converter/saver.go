package converter

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/chai2010/webp"
	"golang.org/x/image/bmp"

	"image-converter/internal/types"
)

// ImageSaver は画像ファイルの保存を提供します
type ImageSaver struct{}

// NewImageSaver は新しいImageSaverを作成します
func NewImageSaver() *ImageSaver {
	return &ImageSaver{}
}

// Save は画像を指定されたパスとフォーマットで保存します
// formatはImageFormat型の文字列（jpeg, png, webp, gif, bmp）
// qualityはJPEG保存時の品質（1-100）、他のフォーマットでは無視されます
func (is *ImageSaver) Save(img image.Image, path string, format types.ImageFormat, quality int) error {
	// ファイルを作成
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// フォーマットに応じてエンコード
	switch format {
	case types.FormatJPEG:
		return is.saveJPEG(file, img, quality)
	case types.FormatPNG:
		return is.savePNG(file, img)
	case types.FormatWebP:
		return is.saveWebP(file, img, quality)
	case types.FormatGIF:
		return is.saveGIF(file, img)
	case types.FormatBMP:
		return is.saveBMP(file, img)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

// saveJPEG はJPEG形式で画像を保存します
func (is *ImageSaver) saveJPEG(file *os.File, img image.Image, quality int) error {
	options := &jpeg.Options{
		Quality: quality,
	}
	
	if err := jpeg.Encode(file, img, options); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}
	
	return nil
}

// savePNG はPNG形式で画像を保存します
func (is *ImageSaver) savePNG(file *os.File, img image.Image) error {
	encoder := &png.Encoder{
		CompressionLevel: png.DefaultCompression,
	}
	
	if err := encoder.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}
	
	return nil
}

// saveWebP はWebP形式で画像を保存します
func (is *ImageSaver) saveWebP(file *os.File, img image.Image, quality int) error {
	// WebPエンコーダーのオプション設定
	options := &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	}
	
	if err := webp.Encode(file, img, options); err != nil {
		return fmt.Errorf("failed to encode WebP: %w", err)
	}
	
	return nil
}

// saveGIF はGIF形式で画像を保存します
func (is *ImageSaver) saveGIF(file *os.File, img image.Image) error {
	options := &gif.Options{
		NumColors: 256,
	}
	
	if err := gif.Encode(file, img, options); err != nil {
		return fmt.Errorf("failed to encode GIF: %w", err)
	}
	
	return nil
}

// saveBMP はBMP形式で画像を保存します
func (is *ImageSaver) saveBMP(file *os.File, img image.Image) error {
	if err := bmp.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode BMP: %w", err)
	}
	
	return nil
}
