package converter

import (
	"fmt"
	"runtime"
	"sync"

	"image-converter/internal/types"
)

// Converter は画像変換処理を統合します
type Converter struct {
	config          types.Config
	stats           types.ConversionStats
	statsMutex      sync.Mutex // 統計情報の更新を保護
	loader          *ImageLoader
	resizer         *ResizeCalculator
	saver           *ImageSaver
	formatDetector  *FormatDetector
}

// NewConverter は新しいConverterを作成します
func NewConverter(config types.Config) *Converter {
	return &Converter{
		config:         config,
		stats:          types.ConversionStats{},
		loader:         NewImageLoader(),
		resizer:        NewResizeCalculator(),
		saver:          NewImageSaver(),
		formatDetector: NewFormatDetector(),
	}
}

// ConvertImage は単一の画像ファイルを変換します
// 変換処理のフロー:
// 1. 画像の読み込み
// 2. リサイズ仕様の適用
// 3. 出力フォーマットの決定
// 4. 画像の保存
func (c *Converter) ConvertImage(sourcePath, outputDir string) types.ConversionResult {
	result := types.ConversionResult{
		SourcePath: sourcePath,
		Success:    false,
	}

	// 1. 画像の読み込み
	img, err := c.loader.Load(sourcePath)
	if err != nil {
		result.Error = fmt.Errorf("failed to load image: %w", err)
		return result
	}

	// 2. リサイズ仕様の作成と適用
	resizeSpec := types.ResizeSpec{
		Scale:  c.config.Scale,
		Width:  c.config.Width,
		Height: c.config.Height,
	}
	
	resizedImg := c.resizer.ResizeImage(img, resizeSpec)

	// 3. 出力フォーマットの決定
	var outputFormat types.ImageFormat
	if c.config.Format != "" {
		// ユーザーが指定したフォーマットを使用
		outputFormat = c.formatDetector.NormalizeFormat(c.config.Format)
	} else {
		// 元画像と同じフォーマットを使用
		detectedFormat, err := c.formatDetector.DetectFormat(sourcePath)
		if err != nil {
			result.Error = fmt.Errorf("failed to detect format: %w", err)
			return result
		}
		outputFormat = detectedFormat
	}

	// 4. 出力パスの生成
	outputPath := c.formatDetector.GenerateOutputPath(sourcePath, outputDir, outputFormat)
	result.OutputPath = outputPath

	// 5. 画像の保存
	quality := c.config.JPEGQuality
	if quality == 0 {
		quality = 85 // デフォルト品質
	}

	err = c.saver.Save(resizedImg, outputPath, outputFormat, quality)
	if err != nil {
		result.Error = fmt.Errorf("failed to save image: %w", err)
		return result
	}

	// 成功
	result.Success = true
	return result
}

// GetStats は現在の統計情報を返します
func (c *Converter) GetStats() types.ConversionStats {
	return c.stats
}

// UpdateStats は統計情報を更新します（スレッドセーフ）
func (c *Converter) UpdateStats(result types.ConversionResult) {
	c.statsMutex.Lock()
	defer c.statsMutex.Unlock()
	
	c.stats.Total++
	if result.Success {
		c.stats.Success++
	} else {
		c.stats.Failed++
	}
}

// IncrementSkipped はスキップされたファイル数を増やします（スレッドセーフ）
func (c *Converter) IncrementSkipped() {
	c.statsMutex.Lock()
	defer c.statsMutex.Unlock()
	
	c.stats.Total++
	c.stats.Skipped++
}

// FileSystemScanner はファイルシステム操作のインターフェースです
type FileSystemScanner interface {
	ScanDirectory(path string) ([]string, error)
	IsImageFile(path string) bool
}

// ProcessDirectory はディレクトリ内のすべてのファイルを並行処理します
// 画像ファイルと非画像ファイルを振り分け、統計情報を収集します
// エラーが発生しても処理を継続します
func (c *Converter) ProcessDirectory(inputDir, outputDir string, fsManager FileSystemScanner) error {
	// ディレクトリ内のファイルを走査
	files, err := fsManager.ScanDirectory(inputDir)
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	// 画像ファイル数をカウント
	imageFiles := []string{}
	for _, file := range files {
		if fsManager.IsImageFile(file) {
			imageFiles = append(imageFiles, file)
		} else {
			c.IncrementSkipped()
		}
	}

	// 要件6.1: 処理開始時の総ファイル数表示
	fmt.Printf("Processing %d images...\n", len(imageFiles))

	// 並行処理の設定
	numWorkers := runtime.NumCPU()
	fmt.Printf("Using %d workers (CPU count: %d)\n", numWorkers, numWorkers)
	sem := make(chan struct{}, numWorkers) // セマフォでCPU数に基づく並行数を制御
	var wg sync.WaitGroup
	var progressMutex sync.Mutex // 進行状況表示の保護
	processedCount := 0

	// 各画像ファイルを並行処理
	for _, file := range imageFiles {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			
			// セマフォを取得（並行数を制限）
			sem <- struct{}{}
			defer func() { <-sem }()

			// 進行状況の表示（スレッドセーフ）
			progressMutex.Lock()
			processedCount++
			currentIndex := processedCount
			fmt.Printf("[%d/%d] Converting %s... ", currentIndex, len(imageFiles), f)
			progressMutex.Unlock()

			// 画像の変換
			result := c.ConvertImage(f, outputDir)
			c.UpdateStats(result)

			// 結果の表示（スレッドセーフ）
			progressMutex.Lock()
			if result.Success {
				fmt.Printf("OK\n")
			} else {
				fmt.Printf("FAILED (%v)\n", result.Error)
			}
			progressMutex.Unlock()
		}(file)
	}

	// すべてのゴルーチンの完了を待機
	wg.Wait()

	// 要件6.5: 処理完了時の要約表示
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Total: %d\n", c.stats.Total)
	fmt.Printf("  Success: %d\n", c.stats.Success)
	fmt.Printf("  Failed: %d\n", c.stats.Failed)
	fmt.Printf("  Skipped: %d\n", c.stats.Skipped)

	return nil
}
