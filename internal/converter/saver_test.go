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

// Feature: image-converter, Property 5: フォーマット変換の正確性
// Validates: Requirements 3.1, 3.2, 3.3
// 任意の画像と出力フォーマットに対して、保存後の画像ファイルは指定されたフォーマットでデコード可能でなければならない
func TestProperty_FormatConversionAccuracy(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "format_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	saver := NewImageSaver()
	loader := NewImageLoader()

	// サポートされているフォーマットのジェネレーター
	formatGen := gen.OneConstOf(
		types.FormatJPEG,
		types.FormatPNG,
		types.FormatWebP,
		types.FormatGIF,
		types.FormatBMP,
	)

	// 画像サイズのジェネレーター（小さめのサイズでテスト）
	sizeGen := gen.IntRange(10, 200)

	properties.Property("saved image can be decoded in specified format", prop.ForAll(
		func(width, height int, format types.ImageFormat) bool {
			// テスト用の画像を生成
			img := createTestImage(width, height)

			// 一時ファイルパスを生成
			tempFile := filepath.Join(tempDir, "test_image")
			
			// 画像を保存
			err := saver.Save(img, tempFile, format, 85)
			if err != nil {
				t.Logf("Failed to save image: %v", err)
				return false
			}

			// 保存した画像を読み込み
			loadedImg, err := loader.Load(tempFile)
			if err != nil {
				t.Logf("Failed to load saved image: %v", err)
				return false
			}

			// 画像が正しく読み込めたことを確認
			if loadedImg == nil {
				t.Logf("Loaded image is nil")
				return false
			}

			// サイズが保持されているか確認（GIFは色数制限があるため、サイズのみ確認）
			bounds := loadedImg.Bounds()
			if bounds.Dx() != width || bounds.Dy() != height {
				t.Logf("Image size mismatch: expected %dx%d, got %dx%d", 
					width, height, bounds.Dx(), bounds.Dy())
				return false
			}

			// ファイルをクリーンアップ
			os.Remove(tempFile)

			return true
		},
		sizeGen,
		sizeGen,
		formatGen,
	))

	properties.TestingRun(t)
}

// Feature: image-converter, Property 9: JPEG品質の影響
// Validates: Requirements 8.2
// 任意の画像と2つの異なるJPEG品質値（q1 < q2）に対して、
// 品質q2で保存されたファイルサイズは品質q1で保存されたファイルサイズ以上でなければならない
func TestProperty_JPEGQualityImpact(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSuccessfulTests = 100
	properties := gopter.NewProperties(parameters)

	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "jpeg_quality_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	saver := NewImageSaver()

	// 画像サイズのジェネレーター
	sizeGen := gen.IntRange(50, 300)
	
	// 品質値のジェネレーター（1-100の範囲）
	qualityGen := gen.IntRange(1, 100)

	properties.Property("higher JPEG quality results in larger or equal file size", prop.ForAll(
		func(width, height, q1, q2 int) bool {
			// q1 < q2 となるように調整
			if q1 >= q2 {
				q1, q2 = q2, q1
			}
			
			// 同じ品質値の場合はスキップ
			if q1 == q2 {
				return true
			}

			// テスト用の画像を生成
			img := createTestImage(width, height)

			// 低品質で保存
			lowQualityPath := filepath.Join(tempDir, "low_quality.jpg")
			err := saver.Save(img, lowQualityPath, types.FormatJPEG, q1)
			if err != nil {
				t.Logf("Failed to save low quality image: %v", err)
				return false
			}

			// 高品質で保存
			highQualityPath := filepath.Join(tempDir, "high_quality.jpg")
			err = saver.Save(img, highQualityPath, types.FormatJPEG, q2)
			if err != nil {
				t.Logf("Failed to save high quality image: %v", err)
				return false
			}

			// ファイルサイズを取得
			lowQualityInfo, err := os.Stat(lowQualityPath)
			if err != nil {
				t.Logf("Failed to stat low quality file: %v", err)
				return false
			}

			highQualityInfo, err := os.Stat(highQualityPath)
			if err != nil {
				t.Logf("Failed to stat high quality file: %v", err)
				return false
			}

			lowQualitySize := lowQualityInfo.Size()
			highQualitySize := highQualityInfo.Size()

			// 高品質のファイルサイズが低品質以上であることを確認
			if highQualitySize < lowQualitySize {
				t.Logf("Quality paradox: q1=%d (size=%d) vs q2=%d (size=%d)", 
					q1, lowQualitySize, q2, highQualitySize)
				return false
			}

			// クリーンアップ
			os.Remove(lowQualityPath)
			os.Remove(highQualityPath)

			return true
		},
		sizeGen,
		sizeGen,
		qualityGen,
		qualityGen,
	))

	properties.TestingRun(t)
}
