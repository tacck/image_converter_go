package cli

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"

	"image-converter/internal/types"
)

// Feature: image-converter, Property 11: 無効な入力の拒否
// 任意の無効な設定（負の倍率、負のピクセル値、倍率とピクセルの同時指定、範囲外のJPEG品質）に対して、
// システムはエラーを返して処理を開始してはならない
// Validates: Requirements 2.8, 2.10, 2.11, 8.4

func TestProperty_InvalidInputRejection(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// プロパティ11a: 負の倍率は拒否される
	properties.Property("負の倍率は拒否される", prop.ForAll(
		func(scale float64) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				Scale:       scale,
				JPEGQuality: 85,
			}
			err := ValidateConfig(config)
			return err != nil // エラーが返されるべき
		},
		gen.Float64Range(-1000, -0.0001), // 負の倍率を生成
	))

	// プロパティ11b: 負の幅は拒否される
	properties.Property("負の幅は拒否される", prop.ForAll(
		func(width int) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				Width:       width,
				JPEGQuality: 85,
			}
			err := ValidateConfig(config)
			return err != nil // エラーが返されるべき
		},
		gen.IntRange(-1000, -1), // 負の幅を生成
	))

	// プロパティ11c: 負の高さは拒否される
	properties.Property("負の高さは拒否される", prop.ForAll(
		func(height int) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				Height:      height,
				JPEGQuality: 85,
			}
			err := ValidateConfig(config)
			return err != nil // エラーが返されるべき
		},
		gen.IntRange(-1000, -1), // 負の高さを生成
	))

	// プロパティ11d: 倍率とピクセル指定の同時指定は拒否される
	properties.Property("倍率とピクセル指定の同時指定は拒否される", prop.ForAll(
		func(scale float64, width int, height int) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				Scale:       scale,
				Width:       width,
				Height:      height,
				JPEGQuality: 85,
			}
			err := ValidateConfig(config)
			return err != nil // エラーが返されるべき
		},
		gen.Float64Range(0.1, 10.0),  // 正の倍率
		gen.IntRange(1, 5000),         // 正の幅
		gen.IntRange(0, 0),            // 高さは0でもOK
	))

	// プロパティ11e: 範囲外のJPEG品質は拒否される（下限）
	properties.Property("範囲外のJPEG品質は拒否される（下限）", prop.ForAll(
		func(quality int) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				JPEGQuality: quality,
			}
			err := ValidateConfig(config)
			return err != nil // エラーが返されるべき
		},
		gen.IntRange(-100, 0), // 1未満の品質を生成
	))

	// プロパティ11f: 範囲外のJPEG品質は拒否される（上限）
	properties.Property("範囲外のJPEG品質は拒否される（上限）", prop.ForAll(
		func(quality int) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				JPEGQuality: quality,
			}
			err := ValidateConfig(config)
			return err != nil // エラーが返されるべき
		},
		gen.IntRange(101, 1000), // 100超の品質を生成
	))

	// 各プロパティを最低100回テスト
	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// 有効な設定が受け入れられることを確認する補完的なプロパティテスト
func TestProperty_ValidInputAcceptance(t *testing.T) {
	properties := gopter.NewProperties(nil)

	// 有効な倍率指定は受け入れられる
	properties.Property("有効な倍率指定は受け入れられる", prop.ForAll(
		func(scale float64) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				Scale:       scale,
				JPEGQuality: 85,
			}
			err := ValidateConfig(config)
			return err == nil // エラーが返されないべき
		},
		gen.Float64Range(0.01, 100.0), // 正の倍率
	))

	// 有効なピクセル指定は受け入れられる
	properties.Property("有効なピクセル指定は受け入れられる", prop.ForAll(
		func(width int, height int) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				Width:       width,
				Height:      height,
				JPEGQuality: 85,
			}
			err := ValidateConfig(config)
			return err == nil // エラーが返されないべき
		},
		gen.IntRange(0, 10000), // 0以上の幅
		gen.IntRange(0, 10000), // 0以上の高さ
	))

	// 有効なJPEG品質は受け入れられる
	properties.Property("有効なJPEG品質は受け入れられる", prop.ForAll(
		func(quality int) bool {
			config := &types.Config{
				InputDir:    "/input",
				OutputDir:   "/output",
				JPEGQuality: quality,
			}
			err := ValidateConfig(config)
			return err == nil // エラーが返されないべき
		},
		gen.IntRange(1, 100), // 1-100の品質
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}
