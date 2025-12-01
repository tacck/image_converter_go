package converter

import (
	"image"
	"math"

	"golang.org/x/image/draw"
	"image-converter/internal/types"
)

// ResizeCalculator はリサイズ計算を提供します
type ResizeCalculator struct{}

// NewResizeCalculator は新しいResizeCalculatorを作成します
func NewResizeCalculator() *ResizeCalculator {
	return &ResizeCalculator{}
}

// CalculateOutputSize はリサイズ仕様に基づいて出力サイズを計算します
func (rc *ResizeCalculator) CalculateOutputSize(srcWidth, srcHeight int, spec types.ResizeSpec) (dstWidth, dstHeight int) {
	// 倍率指定の場合
	if spec.Scale > 0 {
		dstWidth = int(math.Round(float64(srcWidth) * spec.Scale))
		dstHeight = int(math.Round(float64(srcHeight) * spec.Scale))
		return
	}

	// 幅と高さ両方指定の場合（範囲内に収める）
	if spec.Width > 0 && spec.Height > 0 {
		scaleW := float64(spec.Width) / float64(srcWidth)
		scaleH := float64(spec.Height) / float64(srcHeight)
		scale := math.Min(scaleW, scaleH)
		dstWidth = int(math.Round(float64(srcWidth) * scale))
		dstHeight = int(math.Round(float64(srcHeight) * scale))
		return
	}

	// 幅のみ指定の場合
	if spec.Width > 0 {
		scale := float64(spec.Width) / float64(srcWidth)
		dstWidth = spec.Width
		dstHeight = int(math.Round(float64(srcHeight) * scale))
		return
	}

	// 高さのみ指定の場合
	if spec.Height > 0 {
		scale := float64(spec.Height) / float64(srcHeight)
		dstWidth = int(math.Round(float64(srcWidth) * scale))
		dstHeight = spec.Height
		return
	}

	// サイズ指定なしの場合（元のサイズを維持）
	dstWidth = srcWidth
	dstHeight = srcHeight
	return
}

// CalculateOutputSizeFromImage は画像からリサイズ後のサイズを計算します
func (rc *ResizeCalculator) CalculateOutputSizeFromImage(img image.Image, spec types.ResizeSpec) (dstWidth, dstHeight int) {
	bounds := img.Bounds()
	srcWidth := bounds.Dx()
	srcHeight := bounds.Dy()
	return rc.CalculateOutputSize(srcWidth, srcHeight, spec)
}

// ResizeImage はCatmullRomスケーラーを使用して画像をリサイズします
func (rc *ResizeCalculator) ResizeImage(src image.Image, spec types.ResizeSpec) image.Image {
	// サイズ指定がない場合は元の画像をそのまま返す
	if spec.Scale == 0 && spec.Width == 0 && spec.Height == 0 {
		return src
	}

	// 出力サイズを計算
	dstWidth, dstHeight := rc.CalculateOutputSizeFromImage(src, spec)

	// 元のサイズと同じ場合は元の画像をそのまま返す
	bounds := src.Bounds()
	if dstWidth == bounds.Dx() && dstHeight == bounds.Dy() {
		return src
	}

	// 新しい画像を作成
	dst := image.NewRGBA(image.Rect(0, 0, dstWidth, dstHeight))

	// CatmullRomスケーラーを使用して高品質リサイズ
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return dst
}
