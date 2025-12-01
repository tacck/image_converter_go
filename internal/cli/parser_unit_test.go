package cli

import (
	"testing"

	"image-converter/internal/types"
)

// ユニットテスト: 具体的なエッジケースと例をテスト

func TestValidateConfig_MissingInputDir(t *testing.T) {
	config := &types.Config{
		OutputDir:   "/output",
		JPEGQuality: 85,
	}
	err := ValidateConfig(config)
	if err == nil {
		t.Error("入力ディレクトリが指定されていない場合、エラーが返されるべき")
	}
}

func TestValidateConfig_MissingOutputDir(t *testing.T) {
	config := &types.Config{
		InputDir:    "/input",
		JPEGQuality: 85,
	}
	err := ValidateConfig(config)
	if err == nil {
		t.Error("出力ディレクトリが指定されていない場合、エラーが返されるべき")
	}
}

func TestValidateConfig_ScaleAndPixelsBoth(t *testing.T) {
	config := &types.Config{
		InputDir:    "/input",
		OutputDir:   "/output",
		Scale:       0.5,
		Width:       800,
		JPEGQuality: 85,
	}
	err := ValidateConfig(config)
	if err == nil {
		t.Error("倍率とピクセル指定を同時に使用した場合、エラーが返されるべき")
	}
}

func TestValidateConfig_ValidScale(t *testing.T) {
	config := &types.Config{
		InputDir:    "/input",
		OutputDir:   "/output",
		Scale:       0.5,
		JPEGQuality: 85,
	}
	err := ValidateConfig(config)
	if err != nil {
		t.Errorf("有効な倍率指定でエラーが返された: %v", err)
	}
}

func TestValidateConfig_ValidPixels(t *testing.T) {
	config := &types.Config{
		InputDir:    "/input",
		OutputDir:   "/output",
		Width:       800,
		Height:      600,
		JPEGQuality: 85,
	}
	err := ValidateConfig(config)
	if err != nil {
		t.Errorf("有効なピクセル指定でエラーが返された: %v", err)
	}
}

func TestValidateConfig_NoResize(t *testing.T) {
	config := &types.Config{
		InputDir:    "/input",
		OutputDir:   "/output",
		JPEGQuality: 85,
	}
	err := ValidateConfig(config)
	if err != nil {
		t.Errorf("リサイズ指定なしでエラーが返された: %v", err)
	}
}

func TestValidateConfig_UnsupportedFormat(t *testing.T) {
	config := &types.Config{
		InputDir:    "/input",
		OutputDir:   "/output",
		Format:      "tiff",
		JPEGQuality: 85,
	}
	err := ValidateConfig(config)
	if err == nil {
		t.Error("サポートされていないフォーマットでエラーが返されるべき")
	}
}

func TestValidateConfig_SupportedFormats(t *testing.T) {
	formats := []string{"jpeg", "jpg", "png", "webp", "gif", "bmp", "JPEG", "PNG"}
	for _, format := range formats {
		config := &types.Config{
			InputDir:    "/input",
			OutputDir:   "/output",
			Format:      format,
			JPEGQuality: 85,
		}
		err := ValidateConfig(config)
		if err != nil {
			t.Errorf("サポートされているフォーマット %s でエラーが返された: %v", format, err)
		}
	}
}

func TestValidateConfig_JPEGQualityBoundary(t *testing.T) {
	tests := []struct {
		quality int
		valid   bool
	}{
		{0, false},
		{1, true},
		{50, true},
		{100, true},
		{101, false},
	}

	for _, tt := range tests {
		config := &types.Config{
			InputDir:    "/input",
			OutputDir:   "/output",
			JPEGQuality: tt.quality,
		}
		err := ValidateConfig(config)
		if tt.valid && err != nil {
			t.Errorf("品質 %d は有効だがエラーが返された: %v", tt.quality, err)
		}
		if !tt.valid && err == nil {
			t.Errorf("品質 %d は無効だがエラーが返されなかった", tt.quality)
		}
	}
}
