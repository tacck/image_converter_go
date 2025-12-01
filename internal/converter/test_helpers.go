package converter

import (
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"golang.org/x/image/bmp"
	"image-converter/internal/types"
)

// createTestImage はテスト用の画像を生成します
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	
	// シンプルなグラデーションパターンを作成
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
	
	return img
}

// getExtension はフォーマットに対応する拡張子を返します
func getExtension(format types.ImageFormat) string {
	switch format {
	case types.FormatJPEG:
		return ".jpg"
	case types.FormatPNG:
		return ".png"
	case types.FormatGIF:
		return ".gif"
	case types.FormatBMP:
		return ".bmp"
	case types.FormatWebP:
		return ".webp"
	default:
		return ".jpg"
	}
}

// saveImageWithFormat は指定されたフォーマットで画像を保存します
func saveImageWithFormat(img image.Image, path string, format types.ImageFormat) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case types.FormatJPEG:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	case types.FormatPNG:
		return png.Encode(file, img)
	case types.FormatGIF:
		return gif.Encode(file, img, nil)
	case types.FormatBMP:
		return bmp.Encode(file, img)
	case types.FormatWebP:
		// WebPのエンコードは標準ライブラリにないため、
		// テストではPNGとして保存し、拡張子だけ変更
		// 実際のWebP保存は別のタスクで実装
		return png.Encode(file, img)
	default:
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 90})
	}
}
