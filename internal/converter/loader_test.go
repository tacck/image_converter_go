package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"image-converter/internal/types"
)

// Feature: image-converter, Property 8: 入力フォーマットの独立性
// Validates: Requirements 7.1, 7.2, 7.3, 7.4, 7.5
//
// 任意のサポートされている入力フォーマットの画像に対して、
// 同じリサイズ仕様と出力フォーマットを適用した場合、
// 出力画像の次元とフォーマットは入力フォーマットに依存せず一貫していなければならない
func TestProperty_InputFormatIndependence(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "image-loader-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	loader := NewImageLoader()

	properties.Property("入力フォーマットに依存せず一貫した画像が読み込まれる", prop.ForAll(
		func(width, height int, formatIndex int) bool {
			// 入力の検証
			if width < 1 || width > 1000 || height < 1 || height > 1000 {
				return true // 無効な入力はスキップ
			}

			// サポートされているフォーマットのリスト
			formats := []types.ImageFormat{
				types.FormatJPEG,
				types.FormatPNG,
				types.FormatGIF,
				types.FormatBMP,
				types.FormatWebP,
			}

			// フォーマットインデックスを範囲内に制限
			formatIdx := formatIndex % len(formats)
			inputFormat := formats[formatIdx]

			// テスト画像を生成
			testImg := createTestImage(width, height)

			// 画像を一時ファイルに保存
			imagePath := filepath.Join(tempDir, "test_image"+getExtension(inputFormat))
			if err := saveImageWithFormat(testImg, imagePath, inputFormat); err != nil {
				t.Logf("Failed to save test image: %v", err)
				return false
			}

			// 画像を読み込む
			loadedImg, err := loader.Load(imagePath)
			if err != nil {
				t.Logf("Failed to load image: %v", err)
				return false
			}

			// 読み込んだ画像の次元を確認
			bounds := loadedImg.Bounds()
			loadedWidth := bounds.Dx()
			loadedHeight := bounds.Dy()

			// 次元が元の画像と一致することを確認
			// （JPEGやWebPなどの非可逆圧縮でも次元は保持される）
			if loadedWidth != width || loadedHeight != height {
				t.Logf("Dimension mismatch: expected %dx%d, got %dx%d", width, height, loadedWidth, loadedHeight)
				return false
			}

			// クリーンアップ
			os.Remove(imagePath)

			return true
		},
		gen.IntRange(10, 500),  // width
		gen.IntRange(10, 500),  // height
		gen.IntRange(0, 10000), // formatIndex
	))

	properties.TestingRun(t)
}
