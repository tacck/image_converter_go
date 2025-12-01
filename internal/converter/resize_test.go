package converter

import (
	"image"
	"math"
	"testing"
	"testing/quick"

	"image-converter/internal/types"
)

// Feature: image-converter, Property 2: スケール倍率の正確性
// Validates: Requirements 2.2
func TestProperty_ScaleAccuracy(t *testing.T) {
	rc := NewResizeCalculator()

	f := func(srcWidth, srcHeight uint16, scale float32) bool {
		// 入力ドメインの制約
		if srcWidth == 0 || srcHeight == 0 {
			return true // スキップ
		}
		if scale <= 0 || scale > 100 {
			return true // スキップ（無効な倍率）
		}

		// 実際の倍率に変換
		actualScale := float64(scale)
		
		spec := types.ResizeSpec{
			Scale:  actualScale,
			Width:  0,
			Height: 0,
		}

		dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

		// 期待値を計算
		expectedWidth := int(math.Round(float64(srcWidth) * actualScale))
		expectedHeight := int(math.Round(float64(srcHeight) * actualScale))

		// 倍率を適用した結果が期待値と一致するか確認
		if dstWidth != expectedWidth {
			t.Logf("Width mismatch: got %d, expected %d (srcWidth=%d, scale=%f)", 
				dstWidth, expectedWidth, srcWidth, actualScale)
			return false
		}
		if dstHeight != expectedHeight {
			t.Logf("Height mismatch: got %d, expected %d (srcHeight=%d, scale=%f)", 
				dstHeight, expectedHeight, srcHeight, actualScale)
			return false
		}

		return true
	}

	config := &quick.Config{MaxCount: 100}
	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property violated: %v", err)
	}
}

// Feature: image-converter, Property 3: ピクセル指定の正確性
// Validates: Requirements 2.5, 2.6
func TestProperty_PixelSpecificationAccuracy(t *testing.T) {
	rc := NewResizeCalculator()

	t.Run("WidthOnly", func(t *testing.T) {
		f := func(srcWidth, srcHeight, targetWidth uint16) bool {
			// 入力ドメインの制約
			if srcWidth == 0 || srcHeight == 0 || targetWidth == 0 {
				return true // スキップ
			}

			spec := types.ResizeSpec{
				Scale:  0,
				Width:  int(targetWidth),
				Height: 0,
			}

			dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

			// 幅が指定された値と一致するか確認
			if dstWidth != int(targetWidth) {
				t.Logf("Width mismatch: got %d, expected %d", dstWidth, targetWidth)
				return false
			}

			// 高さが縦横比を維持しているか確認
			scale := float64(targetWidth) / float64(srcWidth)
			expectedHeight := int(math.Round(float64(srcHeight) * scale))
			if dstHeight != expectedHeight {
				t.Logf("Height mismatch: got %d, expected %d (scale=%f)", dstHeight, expectedHeight, scale)
				return false
			}

			return true
		}

		config := &quick.Config{MaxCount: 100}
		if err := quick.Check(f, config); err != nil {
			t.Errorf("Property violated: %v", err)
		}
	})

	t.Run("HeightOnly", func(t *testing.T) {
		f := func(srcWidth, srcHeight, targetHeight uint16) bool {
			// 入力ドメインの制約
			if srcWidth == 0 || srcHeight == 0 || targetHeight == 0 {
				return true // スキップ
			}

			spec := types.ResizeSpec{
				Scale:  0,
				Width:  0,
				Height: int(targetHeight),
			}

			dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

			// 高さが指定された値と一致するか確認
			if dstHeight != int(targetHeight) {
				t.Logf("Height mismatch: got %d, expected %d", dstHeight, targetHeight)
				return false
			}

			// 幅が縦横比を維持しているか確認
			scale := float64(targetHeight) / float64(srcHeight)
			expectedWidth := int(math.Round(float64(srcWidth) * scale))
			if dstWidth != expectedWidth {
				t.Logf("Width mismatch: got %d, expected %d (scale=%f)", dstWidth, expectedWidth, scale)
				return false
			}

			return true
		}

		config := &quick.Config{MaxCount: 100}
		if err := quick.Check(f, config); err != nil {
			t.Errorf("Property violated: %v", err)
		}
	})
}

// Feature: image-converter, Property 4: 両次元指定時の範囲内収容
// Validates: Requirements 2.7
func TestProperty_BothDimensionsFitWithinBounds(t *testing.T) {
	rc := NewResizeCalculator()

	f := func(srcWidth, srcHeight, maxWidth, maxHeight uint16) bool {
		// 入力ドメインの制約
		if srcWidth == 0 || srcHeight == 0 || maxWidth == 0 || maxHeight == 0 {
			return true // スキップ
		}

		spec := types.ResizeSpec{
			Scale:  0,
			Width:  int(maxWidth),
			Height: int(maxHeight),
		}

		dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

		// 出力サイズが指定された範囲内に収まっているか確認
		if dstWidth > int(maxWidth) {
			t.Logf("Width exceeds bounds: got %d, max %d", dstWidth, maxWidth)
			return false
		}
		if dstHeight > int(maxHeight) {
			t.Logf("Height exceeds bounds: got %d, max %d", dstHeight, maxHeight)
			return false
		}

		// 少なくとも一方の次元が指定された値と等しいか確認
		if dstWidth != int(maxWidth) && dstHeight != int(maxHeight) {
			t.Logf("Neither dimension matches bound: width=%d (max %d), height=%d (max %d)", 
				dstWidth, maxWidth, dstHeight, maxHeight)
			return false
		}

		return true
	}

	config := &quick.Config{MaxCount: 100}
	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property violated: %v", err)
	}
}

// Feature: image-converter, Property 10: サイズ指定なしの恒等性
// Validates: Requirements 2.9
func TestProperty_NoResizeIdentity(t *testing.T) {
	rc := NewResizeCalculator()

	f := func(srcWidth, srcHeight uint16) bool {
		// 入力ドメインの制約
		if srcWidth == 0 || srcHeight == 0 {
			return true // スキップ
		}

		// サイズ指定なし
		spec := types.ResizeSpec{
			Scale:  0,
			Width:  0,
			Height: 0,
		}

		dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

		// 出力サイズが元のサイズと等しいか確認
		if dstWidth != int(srcWidth) {
			t.Logf("Width changed: got %d, expected %d", dstWidth, srcWidth)
			return false
		}
		if dstHeight != int(srcHeight) {
			t.Logf("Height changed: got %d, expected %d", dstHeight, srcHeight)
			return false
		}

		return true
	}

	config := &quick.Config{MaxCount: 100}
	if err := quick.Check(f, config); err != nil {
		t.Errorf("Property violated: %v", err)
	}
}

// Feature: image-converter, Property 1: リサイズ後の縦横比保持
// Validates: Requirements 2.2, 2.5, 2.6, 2.7, 2.12
func TestProperty_AspectRatioPreservation(t *testing.T) {
	rc := NewResizeCalculator()

	// 縦横比を計算するヘルパー関数
	aspectRatio := func(width, height int) float64 {
		return float64(width) / float64(height)
	}

	// 浮動小数点の比較用の許容誤差
	// 整数への丸め誤差を考慮して、相対誤差で判定
	// 極端な縦横比や小さいサイズの場合、整数への丸めの影響が大きくなるため、
	// より大きな許容誤差を設定
	const relativeEpsilon = 0.15 // 15%の相対誤差を許容

	t.Run("ScaleResize", func(t *testing.T) {
		f := func(srcWidth, srcHeight uint16, scale float32) bool {
			// 入力ドメインの制約
			if srcWidth == 0 || srcHeight == 0 {
				return true // スキップ
			}
			if scale <= 0 || scale > 100 {
				return true // スキップ
			}

			spec := types.ResizeSpec{
				Scale:  float64(scale),
				Width:  0,
				Height: 0,
			}

			dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

			srcAspect := aspectRatio(int(srcWidth), int(srcHeight))
			dstAspect := aspectRatio(dstWidth, dstHeight)

			// 相対誤差で判定（整数への丸め誤差を考慮）
			relativeDiff := math.Abs(srcAspect-dstAspect) / srcAspect
			if relativeDiff > relativeEpsilon {
				t.Logf("Aspect ratio not preserved: src=%f, dst=%f, relativeDiff=%f", srcAspect, dstAspect, relativeDiff)
				return false
			}

			return true
		}

		config := &quick.Config{MaxCount: 100}
		if err := quick.Check(f, config); err != nil {
			t.Errorf("Property violated: %v", err)
		}
	})

	t.Run("WidthOnlyResize", func(t *testing.T) {
		f := func(srcWidth, srcHeight, targetWidth uint16) bool {
			// 入力ドメインの制約
			if srcWidth == 0 || srcHeight == 0 || targetWidth == 0 {
				return true // スキップ
			}

			spec := types.ResizeSpec{
				Scale:  0,
				Width:  int(targetWidth),
				Height: 0,
			}

			dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

			srcAspect := aspectRatio(int(srcWidth), int(srcHeight))
			dstAspect := aspectRatio(dstWidth, dstHeight)

			// 相対誤差で判定（整数への丸め誤差を考慮）
			relativeDiff := math.Abs(srcAspect-dstAspect) / srcAspect
			if relativeDiff > relativeEpsilon {
				t.Logf("Aspect ratio not preserved: src=%f, dst=%f, relativeDiff=%f", srcAspect, dstAspect, relativeDiff)
				return false
			}

			return true
		}

		config := &quick.Config{MaxCount: 100}
		if err := quick.Check(f, config); err != nil {
			t.Errorf("Property violated: %v", err)
		}
	})

	t.Run("HeightOnlyResize", func(t *testing.T) {
		f := func(srcWidth, srcHeight, targetHeight uint16) bool {
			// 入力ドメインの制約
			if srcWidth == 0 || srcHeight == 0 || targetHeight == 0 {
				return true // スキップ
			}

			spec := types.ResizeSpec{
				Scale:  0,
				Width:  0,
				Height: int(targetHeight),
			}

			dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

			srcAspect := aspectRatio(int(srcWidth), int(srcHeight))
			dstAspect := aspectRatio(dstWidth, dstHeight)

			// 相対誤差で判定（整数への丸め誤差を考慮）
			relativeDiff := math.Abs(srcAspect-dstAspect) / srcAspect
			if relativeDiff > relativeEpsilon {
				t.Logf("Aspect ratio not preserved: src=%f, dst=%f, relativeDiff=%f", srcAspect, dstAspect, relativeDiff)
				return false
			}

			return true
		}

		config := &quick.Config{MaxCount: 100}
		if err := quick.Check(f, config); err != nil {
			t.Errorf("Property violated: %v", err)
		}
	})

	t.Run("BothDimensionsResize", func(t *testing.T) {
		f := func(srcWidth, srcHeight, maxWidth, maxHeight uint16) bool {
			// 入力ドメインの制約
			if srcWidth == 0 || srcHeight == 0 || maxWidth == 0 || maxHeight == 0 {
				return true // スキップ
			}

			spec := types.ResizeSpec{
				Scale:  0,
				Width:  int(maxWidth),
				Height: int(maxHeight),
			}

			dstWidth, dstHeight := rc.CalculateOutputSize(int(srcWidth), int(srcHeight), spec)

			srcAspect := aspectRatio(int(srcWidth), int(srcHeight))
			dstAspect := aspectRatio(dstWidth, dstHeight)

			// 相対誤差で判定（整数への丸め誤差を考慮）
			relativeDiff := math.Abs(srcAspect-dstAspect) / srcAspect
			if relativeDiff > relativeEpsilon {
				t.Logf("Aspect ratio not preserved: src=%f, dst=%f, relativeDiff=%f", srcAspect, dstAspect, relativeDiff)
				return false
			}

			return true
		}

		config := &quick.Config{MaxCount: 100}
		if err := quick.Check(f, config); err != nil {
			t.Errorf("Property violated: %v", err)
		}
	})
}

// TestResizeImage_NoResize はサイズ指定なしの場合に元の画像を返すことをテストします
func TestResizeImage_NoResize(t *testing.T) {
	rc := NewResizeCalculator()

	// テスト用の画像を作成
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))

	spec := types.ResizeSpec{
		Scale:  0,
		Width:  0,
		Height: 0,
	}

	result := rc.ResizeImage(src, spec)

	// 元の画像と同じオブジェクトが返されることを確認
	if result != src {
		t.Error("Expected same image object when no resize is specified")
	}
}

// TestResizeImage_Scale は倍率指定でリサイズが正しく行われることをテストします
func TestResizeImage_Scale(t *testing.T) {
	rc := NewResizeCalculator()

	// テスト用の画像を作成
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))

	spec := types.ResizeSpec{
		Scale:  0.5,
		Width:  0,
		Height: 0,
	}

	result := rc.ResizeImage(src, spec)

	bounds := result.Bounds()
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("Expected 50x50, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

// TestResizeImage_Width は幅指定でリサイズが正しく行われることをテストします
func TestResizeImage_Width(t *testing.T) {
	rc := NewResizeCalculator()

	// テスト用の画像を作成（200x100）
	src := image.NewRGBA(image.Rect(0, 0, 200, 100))

	spec := types.ResizeSpec{
		Scale:  0,
		Width:  100,
		Height: 0,
	}

	result := rc.ResizeImage(src, spec)

	bounds := result.Bounds()
	if bounds.Dx() != 100 {
		t.Errorf("Expected width 100, got %d", bounds.Dx())
	}
	// 縦横比を維持するので高さは50になるはず
	if bounds.Dy() != 50 {
		t.Errorf("Expected height 50, got %d", bounds.Dy())
	}
}

// TestResizeImage_Height は高さ指定でリサイズが正しく行われることをテストします
func TestResizeImage_Height(t *testing.T) {
	rc := NewResizeCalculator()

	// テスト用の画像を作成（100x200）
	src := image.NewRGBA(image.Rect(0, 0, 100, 200))

	spec := types.ResizeSpec{
		Scale:  0,
		Width:  0,
		Height: 100,
	}

	result := rc.ResizeImage(src, spec)

	bounds := result.Bounds()
	if bounds.Dy() != 100 {
		t.Errorf("Expected height 100, got %d", bounds.Dy())
	}
	// 縦横比を維持するので幅は50になるはず
	if bounds.Dx() != 50 {
		t.Errorf("Expected width 50, got %d", bounds.Dx())
	}
}

// TestResizeImage_BothDimensions は幅と高さ両方指定でリサイズが正しく行われることをテストします
func TestResizeImage_BothDimensions(t *testing.T) {
	rc := NewResizeCalculator()

	// テスト用の画像を作成（200x100）
	src := image.NewRGBA(image.Rect(0, 0, 200, 100))

	spec := types.ResizeSpec{
		Scale:  0,
		Width:  100,
		Height: 100,
	}

	result := rc.ResizeImage(src, spec)

	bounds := result.Bounds()
	// 縦横比を維持しながら100x100の範囲内に収めるので、100x50になるはず
	if bounds.Dx() != 100 {
		t.Errorf("Expected width 100, got %d", bounds.Dx())
	}
	if bounds.Dy() != 50 {
		t.Errorf("Expected height 50, got %d", bounds.Dy())
	}
}

// TestResizeImage_SameSize は元のサイズと同じサイズにリサイズする場合に元の画像を返すことをテストします
func TestResizeImage_SameSize(t *testing.T) {
	rc := NewResizeCalculator()

	// テスト用の画像を作成
	src := image.NewRGBA(image.Rect(0, 0, 100, 100))

	spec := types.ResizeSpec{
		Scale:  1.0,
		Width:  0,
		Height: 0,
	}

	result := rc.ResizeImage(src, spec)

	// 元の画像と同じオブジェクトが返されることを確認
	if result != src {
		t.Error("Expected same image object when resizing to same size")
	}
}
