package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"image-converter/internal/types"
)

// ParseArgs はコマンドライン引数を解析してConfig構造体を返します
func ParseArgs() (*types.Config, error) {
	config := &types.Config{}

	// フラグの定義
	flag.StringVar(&config.InputDir, "input-dir", "", "入力ディレクトリのパス（必須）")
	flag.StringVar(&config.OutputDir, "output-dir", "", "出力ディレクトリのパス（必須）")
	flag.Float64Var(&config.Scale, "scale", 0, "画像の倍率（例: 0.5で50%、2.0で200%）")
	flag.IntVar(&config.Width, "width", 0, "出力画像の幅（ピクセル）")
	flag.IntVar(&config.Height, "height", 0, "出力画像の高さ（ピクセル）")
	flag.StringVar(&config.Format, "format", "", "出力フォーマット（jpeg, png, webp, gif, bmp）")
	flag.IntVar(&config.JPEGQuality, "jpeg-quality", 85, "JPEG品質（1-100、デフォルト: 85）")

	flag.Parse()

	// 設定の検証
	if err := ValidateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfig は設定の妥当性を検証します
func ValidateConfig(config *types.Config) error {
	// 必須パラメータのチェック
	if config.InputDir == "" {
		return fmt.Errorf("入力ディレクトリが指定されていません")
	}

	if config.OutputDir == "" {
		return fmt.Errorf("出力ディレクトリが指定されていません")
	}

	// スケールとピクセル指定の排他チェック（要件 2.8）
	hasScale := config.Scale > 0
	hasPixels := config.Width > 0 || config.Height > 0

	if hasScale && hasPixels {
		return fmt.Errorf("倍率指定とピクセル指定を同時に使用できません")
	}

	// スケール値の検証（要件 2.1, 2.10）
	if config.Scale < 0 {
		return fmt.Errorf("倍率は0以上である必要があります")
	}

	// ピクセル値の検証（要件 2.3, 2.4, 2.11）
	if config.Width < 0 {
		return fmt.Errorf("幅は0以上である必要があります")
	}

	if config.Height < 0 {
		return fmt.Errorf("高さは0以上である必要があります")
	}

	// フォーマットの検証（要件 3.6）
	if config.Format != "" {
		format := strings.ToLower(config.Format)
		validFormats := map[string]bool{
			"jpeg": true,
			"jpg":  true,
			"png":  true,
			"webp": true,
			"gif":  true,
			"bmp":  true,
		}

		if !validFormats[format] {
			return fmt.Errorf("サポートされていないフォーマット: %s", config.Format)
		}

		// jpegとjpgを正規化
		if format == "jpg" {
			config.Format = "jpeg"
		} else {
			config.Format = format
		}
	}

	// JPEG品質の検証（要件 8.1, 8.4）
	if config.JPEGQuality < 1 || config.JPEGQuality > 100 {
		return fmt.Errorf("JPEG品質は1から100の範囲で指定してください")
	}

	return nil
}

// PrintUsage は使用方法を表示します
func PrintUsage() {
	fmt.Fprintf(os.Stderr, "Image Converter CLI - 画像一括変換ツール\n\n")
	fmt.Fprintf(os.Stderr, "使用方法:\n")
	fmt.Fprintf(os.Stderr, "  image-converter -input-dir <入力ディレクトリ> -output-dir <出力ディレクトリ> [オプション]\n\n")
	
	fmt.Fprintf(os.Stderr, "必須オプション:\n")
	fmt.Fprintf(os.Stderr, "  -input-dir string\n")
	fmt.Fprintf(os.Stderr, "        入力ディレクトリのパス（変換対象の画像が格納されているディレクトリ）\n")
	fmt.Fprintf(os.Stderr, "  -output-dir string\n")
	fmt.Fprintf(os.Stderr, "        出力ディレクトリのパス（変換後の画像を保存するディレクトリ）\n\n")
	
	fmt.Fprintf(os.Stderr, "リサイズオプション（いずれか1つを指定）:\n")
	fmt.Fprintf(os.Stderr, "  -scale float\n")
	fmt.Fprintf(os.Stderr, "        画像の倍率（例: 0.5で50%%、2.0で200%%）\n")
	fmt.Fprintf(os.Stderr, "  -width int\n")
	fmt.Fprintf(os.Stderr, "        出力画像の幅（ピクセル）。縦横比を維持して高さを自動計算\n")
	fmt.Fprintf(os.Stderr, "  -height int\n")
	fmt.Fprintf(os.Stderr, "        出力画像の高さ（ピクセル）。縦横比を維持して幅を自動計算\n")
	fmt.Fprintf(os.Stderr, "  -width と -height\n")
	fmt.Fprintf(os.Stderr, "        両方指定した場合、縦横比を維持しながら指定範囲内に収める\n\n")
	
	fmt.Fprintf(os.Stderr, "フォーマットオプション:\n")
	fmt.Fprintf(os.Stderr, "  -format string\n")
	fmt.Fprintf(os.Stderr, "        出力フォーマット: jpeg, png, webp, gif, bmp\n")
	fmt.Fprintf(os.Stderr, "        指定しない場合は元のフォーマットを維持\n")
	fmt.Fprintf(os.Stderr, "  -jpeg-quality int\n")
	fmt.Fprintf(os.Stderr, "        JPEG品質（1-100）（デフォルト: 85）\n\n")
	
	fmt.Fprintf(os.Stderr, "サポートされているフォーマット:\n")
	fmt.Fprintf(os.Stderr, "  入力: JPEG, PNG, WebP, GIF, BMP\n")
	fmt.Fprintf(os.Stderr, "  出力: JPEG, PNG, WebP, GIF, BMP\n\n")
	
	fmt.Fprintf(os.Stderr, "使用例:\n")
	fmt.Fprintf(os.Stderr, "  # 画像を50%%に縮小\n")
	fmt.Fprintf(os.Stderr, "  image-converter -input-dir ./photos -output-dir ./thumbnails -scale 0.5\n\n")
	fmt.Fprintf(os.Stderr, "  # 幅800ピクセルにリサイズ（縦横比を維持）\n")
	fmt.Fprintf(os.Stderr, "  image-converter -input-dir ./photos -output-dir ./resized -width 800\n\n")
	fmt.Fprintf(os.Stderr, "  # 800x600の範囲内に収める\n")
	fmt.Fprintf(os.Stderr, "  image-converter -input-dir ./photos -output-dir ./resized -width 800 -height 600\n\n")
	fmt.Fprintf(os.Stderr, "  # すべての画像をJPEGに変換\n")
	fmt.Fprintf(os.Stderr, "  image-converter -input-dir ./photos -output-dir ./converted -format jpeg\n\n")
	fmt.Fprintf(os.Stderr, "  # WebPに変換して50%%に縮小\n")
	fmt.Fprintf(os.Stderr, "  image-converter -input-dir ./photos -output-dir ./optimized -scale 0.5 -format webp\n\n")
	fmt.Fprintf(os.Stderr, "  # JPEG品質を指定して変換\n")
	fmt.Fprintf(os.Stderr, "  image-converter -input-dir ./photos -output-dir ./compressed -format jpeg -jpeg-quality 70\n\n")
	
	fmt.Fprintf(os.Stderr, "注意事項:\n")
	fmt.Fprintf(os.Stderr, "  - 倍率指定（-scale）とピクセル指定（-width/-height）は同時に使用できません\n")
	fmt.Fprintf(os.Stderr, "  - すべてのリサイズ操作で縦横比が維持されます\n")
	fmt.Fprintf(os.Stderr, "  - 出力ディレクトリに同名のファイルがある場合は上書きされます\n")
	fmt.Fprintf(os.Stderr, "  - サポートされていないフォーマットのファイルはスキップされます\n\n")
	
	fmt.Fprintf(os.Stderr, "詳細はREADME.mdを参照してください。\n")
}
