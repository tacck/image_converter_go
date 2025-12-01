package converter

import (
	"fmt"
	"image"
	_ "image/gif"  // GIFデコーダーを登録
	_ "image/jpeg" // JPEGデコーダーを登録
	_ "image/png"  // PNGデコーダーを登録
	"os"

	_ "golang.org/x/image/bmp"  // BMPデコーダーを登録
	_ "golang.org/x/image/webp" // WebPデコーダーを登録
)

// ImageLoader は画像ファイルの読み込みを提供します
type ImageLoader struct{}

// NewImageLoader は新しいImageLoaderを作成します
func NewImageLoader() *ImageLoader {
	return &ImageLoader{}
}

// Load は指定されたパスから画像を読み込みます
// サポートされているフォーマット: JPEG, PNG, GIF, WebP, BMP
func (il *ImageLoader) Load(path string) (image.Image, error) {
	// ファイルを開く
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 画像をデコード
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// デコード成功（フォーマット情報はログ用に保持可能）
	_ = format

	return img, nil
}
