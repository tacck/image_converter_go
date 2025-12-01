package converter

import (
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"image-converter/internal/types"
)

// saveTestImage はテスト用の画像をファイルに保存します
func saveTestImage(t *testing.T, path string, img image.Image) {
	t.Helper()
	
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}
}

func TestConverter_ConvertImage_Success(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")
	
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// テスト画像を作成
	testImg := createTestImage(100, 100)
	inputPath := filepath.Join(inputDir, "test.png")
	saveTestImage(t, inputPath, testImg)
	
	// Converterを作成
	config := types.Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Scale:       0.5,
		Width:       0,
		Height:      0,
		Format:      "jpeg",
		JPEGQuality: 85,
	}
	
	converter := NewConverter(config)
	
	// 画像を変換
	result := converter.ConvertImage(inputPath, outputDir)
	
	// 結果を検証
	if !result.Success {
		t.Errorf("Expected conversion to succeed, but got error: %v", result.Error)
	}
	
	if result.SourcePath != inputPath {
		t.Errorf("Expected source path %s, got %s", inputPath, result.SourcePath)
	}
	
	expectedOutputPath := filepath.Join(outputDir, "test.jpg")
	if result.OutputPath != expectedOutputPath {
		t.Errorf("Expected output path %s, got %s", expectedOutputPath, result.OutputPath)
	}
	
	// 出力ファイルが存在することを確認
	if _, err := os.Stat(result.OutputPath); os.IsNotExist(err) {
		t.Errorf("Output file does not exist: %s", result.OutputPath)
	}
}

func TestConverter_ConvertImage_InvalidInput(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "output")
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// 存在しないファイルパス
	nonExistentPath := filepath.Join(tempDir, "nonexistent.png")
	
	// Converterを作成
	config := types.Config{
		InputDir:    tempDir,
		OutputDir:   outputDir,
		Scale:       1.0,
		JPEGQuality: 85,
	}
	
	converter := NewConverter(config)
	
	// 画像を変換（失敗するはず）
	result := converter.ConvertImage(nonExistentPath, outputDir)
	
	// 結果を検証
	if result.Success {
		t.Error("Expected conversion to fail for non-existent file")
	}
	
	if result.Error == nil {
		t.Error("Expected error to be set")
	}
}

func TestConverter_ConvertImage_FormatDetection(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")
	
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}
	
	// テスト画像を作成
	testImg := createTestImage(50, 50)
	inputPath := filepath.Join(inputDir, "test.png")
	saveTestImage(t, inputPath, testImg)
	
	// フォーマット指定なしのConverterを作成
	config := types.Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Scale:       1.0,
		Format:      "", // フォーマット指定なし
		JPEGQuality: 85,
	}
	
	converter := NewConverter(config)
	
	// 画像を変換
	result := converter.ConvertImage(inputPath, outputDir)
	
	// 結果を検証
	if !result.Success {
		t.Errorf("Expected conversion to succeed, but got error: %v", result.Error)
	}
	
	// 元のフォーマット（PNG）が保持されることを確認
	expectedOutputPath := filepath.Join(outputDir, "test.png")
	if result.OutputPath != expectedOutputPath {
		t.Errorf("Expected output path %s, got %s", expectedOutputPath, result.OutputPath)
	}
}

func TestConverter_UpdateStats(t *testing.T) {
	config := types.Config{
		JPEGQuality: 85,
	}
	
	converter := NewConverter(config)
	
	// 初期状態を確認
	stats := converter.GetStats()
	if stats.Total != 0 || stats.Success != 0 || stats.Failed != 0 {
		t.Error("Expected initial stats to be zero")
	}
	
	// 成功結果を追加
	successResult := types.ConversionResult{
		SourcePath: "test1.png",
		OutputPath: "test1.jpg",
		Success:    true,
		Error:      nil,
	}
	converter.UpdateStats(successResult)
	
	stats = converter.GetStats()
	if stats.Total != 1 || stats.Success != 1 || stats.Failed != 0 {
		t.Errorf("Expected stats (1, 1, 0), got (%d, %d, %d)", stats.Total, stats.Success, stats.Failed)
	}
	
	// 失敗結果を追加
	failResult := types.ConversionResult{
		SourcePath: "test2.png",
		OutputPath: "",
		Success:    false,
		Error:      os.ErrNotExist,
	}
	converter.UpdateStats(failResult)
	
	stats = converter.GetStats()
	if stats.Total != 2 || stats.Success != 1 || stats.Failed != 1 {
		t.Errorf("Expected stats (2, 1, 1), got (%d, %d, %d)", stats.Total, stats.Success, stats.Failed)
	}
}

func TestConverter_IncrementSkipped(t *testing.T) {
	config := types.Config{
		JPEGQuality: 85,
	}
	
	converter := NewConverter(config)
	
	// スキップを追加
	converter.IncrementSkipped()
	
	stats := converter.GetStats()
	if stats.Total != 1 || stats.Skipped != 1 {
		t.Errorf("Expected stats (1, 0, 0, 1), got (%d, %d, %d, %d)", 
			stats.Total, stats.Success, stats.Failed, stats.Skipped)
	}
	
	// もう一つスキップを追加
	converter.IncrementSkipped()
	
	stats = converter.GetStats()
	if stats.Total != 2 || stats.Skipped != 2 {
		t.Errorf("Expected stats (2, 0, 0, 2), got (%d, %d, %d, %d)", 
			stats.Total, stats.Success, stats.Failed, stats.Skipped)
	}
}

// Feature: image-converter, Property 7: バッチ処理の完全性
// Validates: Requirements 4.1, 4.2, 4.3, 4.5
//
// 任意のディレクトリ内の画像ファイルセットに対して、処理後の成功数、失敗数、スキップ数の合計は、
// ディレクトリ内の総ファイル数と等しくなければならない
func TestProperty_BatchProcessingCompleteness(t *testing.T) {
	// テストケース: 様々なファイルの組み合わせ
	testCases := []struct {
		name           string
		imageFiles     int // 有効な画像ファイル数
		nonImageFiles  int // 非画像ファイル数
		corruptedFiles int // 破損した画像ファイル数
	}{
		{"empty directory", 0, 0, 0},
		{"only images", 5, 0, 0},
		{"only non-images", 0, 5, 0},
		{"mixed files", 3, 2, 0},
		{"with corrupted", 2, 1, 1},
		{"large batch", 20, 5, 3},
		{"single file", 1, 0, 0},
		{"single non-image", 0, 1, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用の一時ディレクトリを作成
			tempDir := t.TempDir()
			inputDir := filepath.Join(tempDir, "input")
			outputDir := filepath.Join(tempDir, "output")

			if err := os.MkdirAll(inputDir, 0755); err != nil {
				t.Fatalf("Failed to create input directory: %v", err)
			}
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				t.Fatalf("Failed to create output directory: %v", err)
			}

			// 有効な画像ファイルを作成
			for i := 0; i < tc.imageFiles; i++ {
				testImg := createTestImage(100, 100)
				imagePath := filepath.Join(inputDir, filepath.Base(t.Name())+"-image-"+string(rune('a'+i))+".png")
				saveTestImage(t, imagePath, testImg)
			}

			// 非画像ファイルを作成
			for i := 0; i < tc.nonImageFiles; i++ {
				textPath := filepath.Join(inputDir, filepath.Base(t.Name())+"-text-"+string(rune('a'+i))+".txt")
				if err := os.WriteFile(textPath, []byte("not an image"), 0644); err != nil {
					t.Fatalf("Failed to create text file: %v", err)
				}
			}

			// 破損した画像ファイルを作成
			for i := 0; i < tc.corruptedFiles; i++ {
				corruptedPath := filepath.Join(inputDir, filepath.Base(t.Name())+"-corrupted-"+string(rune('a'+i))+".png")
				if err := os.WriteFile(corruptedPath, []byte("corrupted image data"), 0644); err != nil {
					t.Fatalf("Failed to create corrupted file: %v", err)
				}
			}

			// 総ファイル数を計算
			totalFiles := tc.imageFiles + tc.nonImageFiles + tc.corruptedFiles

			// Converterを作成
			config := types.Config{
				InputDir:    inputDir,
				OutputDir:   outputDir,
				Scale:       1.0,
				Format:      "jpeg",
				JPEGQuality: 85,
			}

			converter := NewConverter(config)

			// FileSystemManagerを作成
			fsManager := &mockFileSystemManager{
				scanFunc: func(path string) ([]string, error) {
					files, err := os.ReadDir(path)
					if err != nil {
						return nil, err
					}
					var result []string
					for _, f := range files {
						if !f.IsDir() {
							result = append(result, filepath.Join(path, f.Name()))
						}
					}
					return result, nil
				},
				isImageFunc: func(path string) bool {
					ext := filepath.Ext(path)
					return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".bmp" || ext == ".webp"
				},
			}

			// バッチ処理を実行
			err := converter.ProcessDirectory(inputDir, outputDir, fsManager)
			if err != nil {
				t.Fatalf("ProcessDirectory failed: %v", err)
			}

			// 統計情報を取得
			stats := converter.GetStats()

			// プロパティ7を検証: 成功数 + 失敗数 + スキップ数 = 総ファイル数
			actualTotal := stats.Success + stats.Failed + stats.Skipped
			if actualTotal != totalFiles {
				t.Errorf("Property 7 violated: expected total %d (success + failed + skipped), got %d (success=%d, failed=%d, skipped=%d)",
					totalFiles, actualTotal, stats.Success, stats.Failed, stats.Skipped)
			}

			// stats.Totalも総ファイル数と一致するべき
			if stats.Total != totalFiles {
				t.Errorf("stats.Total should equal total files: expected %d, got %d", totalFiles, stats.Total)
			}

			// 追加の検証: スキップ数は非画像ファイル数と一致するべき
			if stats.Skipped != tc.nonImageFiles {
				t.Errorf("Expected %d skipped files (non-images), got %d", tc.nonImageFiles, stats.Skipped)
			}

			// 追加の検証: 成功数は有効な画像ファイル数と一致するべき
			if stats.Success != tc.imageFiles {
				t.Errorf("Expected %d successful conversions, got %d", tc.imageFiles, stats.Success)
			}

			// 追加の検証: 失敗数は破損したファイル数と一致するべき
			if stats.Failed != tc.corruptedFiles {
				t.Errorf("Expected %d failed conversions (corrupted files), got %d", tc.corruptedFiles, stats.Failed)
			}
		})
	}
}

// mockFileSystemManager はテスト用のFileSystemManagerのモックです
type mockFileSystemManager struct {
	scanFunc    func(path string) ([]string, error)
	isImageFunc func(path string) bool
}

func (m *mockFileSystemManager) ScanDirectory(path string) ([]string, error) {
	if m.scanFunc != nil {
		return m.scanFunc(path)
	}
	return nil, nil
}

func (m *mockFileSystemManager) IsImageFile(path string) bool {
	if m.isImageFunc != nil {
		return m.isImageFunc(path)
	}
	return false
}

func (m *mockFileSystemManager) DirectoryExists(path string) (bool, error) {
	return true, nil
}

func (m *mockFileSystemManager) CreateDirectory(path string) error {
	return nil
}

func (m *mockFileSystemManager) ValidateInputDirectory(path string) error {
	return nil
}

func (m *mockFileSystemManager) EnsureOutputDirectory(path string) error {
	return nil
}

// TestConverter_ConcurrentProcessing_StatsAccuracy は並行処理時の統計情報の正確性をテストします
// 要件4.5: エラー発生時の処理継続
func TestConverter_ConcurrentProcessing_StatsAccuracy(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// 複数の画像ファイルを作成（並行処理をテストするため）
	numImages := 20
	for i := 0; i < numImages; i++ {
		testImg := createTestImage(100, 100)
		imagePath := filepath.Join(inputDir, filepath.Base(t.Name())+"-image-"+string(rune('a'+i%26))+string(rune('0'+i/26))+".png")
		saveTestImage(t, imagePath, testImg)
	}

	// 非画像ファイルも追加
	numNonImages := 5
	for i := 0; i < numNonImages; i++ {
		textPath := filepath.Join(inputDir, filepath.Base(t.Name())+"-text-"+string(rune('a'+i))+".txt")
		if err := os.WriteFile(textPath, []byte("not an image"), 0644); err != nil {
			t.Fatalf("Failed to create text file: %v", err)
		}
	}

	// Converterを作成
	config := types.Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Scale:       0.5,
		Format:      "jpeg",
		JPEGQuality: 85,
	}

	converter := NewConverter(config)

	// FileSystemManagerを作成
	fsManager := &mockFileSystemManager{
		scanFunc: func(path string) ([]string, error) {
			files, err := os.ReadDir(path)
			if err != nil {
				return nil, err
			}
			var result []string
			for _, f := range files {
				if !f.IsDir() {
					result = append(result, filepath.Join(path, f.Name()))
				}
			}
			return result, nil
		},
		isImageFunc: func(path string) bool {
			ext := filepath.Ext(path)
			return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".bmp" || ext == ".webp"
		},
	}

	// バッチ処理を実行（並行処理）
	err := converter.ProcessDirectory(inputDir, outputDir, fsManager)
	if err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// 統計情報を取得
	stats := converter.GetStats()

	// 並行処理でも統計情報が正確であることを検証
	expectedTotal := numImages + numNonImages
	if stats.Total != expectedTotal {
		t.Errorf("Expected total %d, got %d", expectedTotal, stats.Total)
	}

	if stats.Success != numImages {
		t.Errorf("Expected %d successful conversions, got %d", numImages, stats.Success)
	}

	if stats.Skipped != numNonImages {
		t.Errorf("Expected %d skipped files, got %d", numNonImages, stats.Skipped)
	}

	if stats.Failed != 0 {
		t.Errorf("Expected 0 failed conversions, got %d", stats.Failed)
	}

	// すべての出力ファイルが作成されていることを確認
	outputFiles, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(outputFiles) != numImages {
		t.Errorf("Expected %d output files, got %d", numImages, len(outputFiles))
	}
}

// TestConverter_ConcurrentProcessing_ErrorHandling は並行処理時のエラーハンドリングをテストします
// 要件4.5: エラー発生時の処理継続
func TestConverter_ConcurrentProcessing_ErrorHandling(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir := t.TempDir()
	inputDir := filepath.Join(tempDir, "input")
	outputDir := filepath.Join(tempDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// 有効な画像ファイルを作成
	numValidImages := 10
	for i := 0; i < numValidImages; i++ {
		testImg := createTestImage(100, 100)
		imagePath := filepath.Join(inputDir, filepath.Base(t.Name())+"-valid-"+string(rune('a'+i))+".png")
		saveTestImage(t, imagePath, testImg)
	}

	// 破損した画像ファイルを作成
	numCorruptedImages := 5
	for i := 0; i < numCorruptedImages; i++ {
		corruptedPath := filepath.Join(inputDir, filepath.Base(t.Name())+"-corrupted-"+string(rune('a'+i))+".png")
		if err := os.WriteFile(corruptedPath, []byte("corrupted image data"), 0644); err != nil {
			t.Fatalf("Failed to create corrupted file: %v", err)
		}
	}

	// Converterを作成
	config := types.Config{
		InputDir:    inputDir,
		OutputDir:   outputDir,
		Scale:       1.0,
		Format:      "jpeg",
		JPEGQuality: 85,
	}

	converter := NewConverter(config)

	// FileSystemManagerを作成
	fsManager := &mockFileSystemManager{
		scanFunc: func(path string) ([]string, error) {
			files, err := os.ReadDir(path)
			if err != nil {
				return nil, err
			}
			var result []string
			for _, f := range files {
				if !f.IsDir() {
					result = append(result, filepath.Join(path, f.Name()))
				}
			}
			return result, nil
		},
		isImageFunc: func(path string) bool {
			ext := filepath.Ext(path)
			return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".bmp" || ext == ".webp"
		},
	}

	// バッチ処理を実行（並行処理）
	err := converter.ProcessDirectory(inputDir, outputDir, fsManager)
	if err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// 統計情報を取得
	stats := converter.GetStats()

	// エラーが発生しても処理が継続されることを検証
	expectedTotal := numValidImages + numCorruptedImages
	if stats.Total != expectedTotal {
		t.Errorf("Expected total %d, got %d", expectedTotal, stats.Total)
	}

	// 有効な画像は成功するべき
	if stats.Success != numValidImages {
		t.Errorf("Expected %d successful conversions, got %d", numValidImages, stats.Success)
	}

	// 破損した画像は失敗するべき
	if stats.Failed != numCorruptedImages {
		t.Errorf("Expected %d failed conversions, got %d", numCorruptedImages, stats.Failed)
	}

	// 有効な画像の出力ファイルが作成されていることを確認
	outputFiles, err := os.ReadDir(outputDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(outputFiles) != numValidImages {
		t.Errorf("Expected %d output files, got %d", numValidImages, len(outputFiles))
	}
}

// TestConverter_ConcurrentProcessing_ThreadSafety は並行処理時のスレッドセーフ性をテストします
func TestConverter_ConcurrentProcessing_ThreadSafety(t *testing.T) {
	// このテストは複数回実行して競合状態を検出します
	for run := 0; run < 5; run++ {
		t.Run(filepath.Base(t.Name())+"-run-"+string(rune('0'+run)), func(t *testing.T) {
			// テスト用の一時ディレクトリを作成
			tempDir := t.TempDir()
			inputDir := filepath.Join(tempDir, "input")
			outputDir := filepath.Join(tempDir, "output")

			if err := os.MkdirAll(inputDir, 0755); err != nil {
				t.Fatalf("Failed to create input directory: %v", err)
			}
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				t.Fatalf("Failed to create output directory: %v", err)
			}

			// 多数の画像ファイルを作成（競合状態を引き起こしやすくする）
			numImages := 50
			for i := 0; i < numImages; i++ {
				testImg := createTestImage(50, 50)
				imagePath := filepath.Join(inputDir, filepath.Base(t.Name())+"-image-"+string(rune('a'+i%26))+string(rune('0'+i/26))+".png")
				saveTestImage(t, imagePath, testImg)
			}

			// Converterを作成
			config := types.Config{
				InputDir:    inputDir,
				OutputDir:   outputDir,
				Scale:       0.8,
				Format:      "jpeg",
				JPEGQuality: 85,
			}

			converter := NewConverter(config)

			// FileSystemManagerを作成
			fsManager := &mockFileSystemManager{
				scanFunc: func(path string) ([]string, error) {
					files, err := os.ReadDir(path)
					if err != nil {
						return nil, err
					}
					var result []string
					for _, f := range files {
						if !f.IsDir() {
							result = append(result, filepath.Join(path, f.Name()))
						}
					}
					return result, nil
				},
				isImageFunc: func(path string) bool {
					ext := filepath.Ext(path)
					return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" || ext == ".bmp" || ext == ".webp"
				},
			}

			// バッチ処理を実行（並行処理）
			err := converter.ProcessDirectory(inputDir, outputDir, fsManager)
			if err != nil {
				t.Fatalf("ProcessDirectory failed: %v", err)
			}

			// 統計情報を取得
			stats := converter.GetStats()

			// スレッドセーフ性を検証：統計情報が正確であることを確認
			if stats.Total != numImages {
				t.Errorf("Thread safety issue: expected total %d, got %d", numImages, stats.Total)
			}

			if stats.Success != numImages {
				t.Errorf("Thread safety issue: expected %d successful conversions, got %d", numImages, stats.Success)
			}

			// 成功数 + 失敗数 + スキップ数 = 総数
			actualTotal := stats.Success + stats.Failed + stats.Skipped
			if actualTotal != stats.Total {
				t.Errorf("Thread safety issue: success + failed + skipped (%d) != total (%d)", actualTotal, stats.Total)
			}
		})
	}
}
