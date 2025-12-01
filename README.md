# Image Converter CLI

Go言語で実装された高性能な画像一括変換ツールです。指定されたディレクトリ内の画像ファイルを一括で変換し、サイズ変更とフォーマット変換を行います。

## 特徴

- 🚀 **高速な一括変換**: ディレクトリ内のすべての画像を一度に処理
- 📐 **柔軟なリサイズ**: 倍率またはピクセル指定による自由なサイズ変更
- 🎨 **多様なフォーマット対応**: JPEG、PNG、WebP、GIF、BMPをサポート
- 📏 **縦横比の維持**: リサイズ時に自動的に縦横比を保持
- ⚙️ **品質調整**: JPEG出力時の品質を細かく制御可能
- 📊 **詳細な進行状況表示**: 処理状況をリアルタイムで確認

## インストール

### ソースからビルド

```bash
git clone <repository-url>
cd image-converter
go build -o image-converter ./cmd/image-converter
```

### 実行ファイルの配置

ビルドした実行ファイルをPATHの通った場所に配置します：

```bash
# macOS/Linux
sudo mv image-converter /usr/local/bin/

# または、ホームディレクトリのbinフォルダに配置
mkdir -p ~/bin
mv image-converter ~/bin/
export PATH="$HOME/bin:$PATH"  # .bashrcや.zshrcに追加
```

## 使用方法

### 基本的な使い方

```bash
image-converter -input-dir <入力ディレクトリ> -output-dir <出力ディレクトリ> [オプション]
```

### オプション

| オプション | 説明 | デフォルト値 |
|-----------|------|-------------|
| `-input-dir` | 入力ディレクトリのパス（必須） | - |
| `-output-dir` | 出力ディレクトリのパス（必須） | - |
| `-scale` | 画像の倍率（例: 0.5で50%、2.0で200%） | - |
| `-width` | 出力画像の幅（ピクセル） | - |
| `-height` | 出力画像の高さ（ピクセル） | - |
| `-format` | 出力フォーマット（jpeg, png, webp, gif, bmp） | 元のフォーマット |
| `-jpeg-quality` | JPEG品質（1-100） | 85 |

### 使用例

#### 1. 画像を50%に縮小

```bash
image-converter -input-dir ./photos -output-dir ./thumbnails -scale 0.5
```

#### 2. 幅800ピクセルにリサイズ（縦横比を維持）

```bash
image-converter -input-dir ./photos -output-dir ./resized -width 800
```

#### 3. 高さ600ピクセルにリサイズ

```bash
image-converter -input-dir ./photos -output-dir ./resized -height 600
```

#### 4. 800x600の範囲内に収める

```bash
image-converter -input-dir ./photos -output-dir ./resized -width 800 -height 600
```

幅と高さの両方を指定した場合、縦横比を維持しながら指定された範囲内に収まるようにリサイズされます。

#### 5. すべての画像をJPEGに変換

```bash
image-converter -input-dir ./photos -output-dir ./converted -format jpeg
```

#### 6. WebPに変換して50%に縮小

```bash
image-converter -input-dir ./photos -output-dir ./optimized -scale 0.5 -format webp
```

#### 7. JPEG品質を指定して変換

```bash
image-converter -input-dir ./photos -output-dir ./compressed -format jpeg -jpeg-quality 70
```

#### 8. フォーマット変換のみ（サイズ変更なし）

```bash
image-converter -input-dir ./photos -output-dir ./converted -format png
```

## サポートされているフォーマット

### 入力フォーマット

- JPEG (.jpg, .jpeg)
- PNG (.png)
- WebP (.webp)
- GIF (.gif)
- BMP (.bmp)

### 出力フォーマット

- JPEG (.jpeg)
- PNG (.png)
- WebP (.webp)
- GIF (.gif)
- BMP (.bmp)

## リサイズの仕様

### 倍率指定（-scale）

元の画像サイズに倍率を掛けてリサイズします。

- `0.5`: 50%に縮小
- `1.0`: 元のサイズを維持
- `2.0`: 200%に拡大

### ピクセル指定（-width, -height）

#### 幅のみ指定

指定された幅にリサイズし、縦横比を維持して高さを自動計算します。

#### 高さのみ指定

指定された高さにリサイズし、縦横比を維持して幅を自動計算します。

#### 幅と高さの両方を指定

縦横比を維持しながら、指定された幅と高さの両方以下に収まるようにリサイズします。

### 制約

- 倍率指定とピクセル指定は同時に使用できません
- すべてのリサイズ操作で縦横比が維持されます
- サイズ指定がない場合、元のサイズが維持されます

## 出力ファイル

### ファイル名

出力ファイルは元のファイル名のベース名を保持し、拡張子のみが変更されます。

例：
- 入力: `photo.jpg`、出力フォーマット: `png` → 出力: `photo.png`
- 入力: `image.png`、フォーマット指定なし → 出力: `image.png`

### 既存ファイルの上書き

出力ディレクトリに同名のファイルが既に存在する場合、上書きされます。

## 進行状況の表示

処理中は以下のような進行状況が表示されます：

```
Processing 15 images...
Using 8 workers (CPU count: 8)
[1/15] Converting image1.jpg... OK
[2/15] Converting image2.png... OK
[3/15] Converting image3.gif... SKIPPED (unsupported format)
...
[15/15] Converting image15.webp... OK

Summary:
  Total: 15
  Success: 13
  Failed: 0
  Skipped: 2
```

## エラーハンドリング

### 設定エラー

無効な設定が指定された場合、エラーメッセージを表示して即座に終了します。

```bash
ERROR: 倍率指定とピクセル指定を同時に使用できません
```

### ファイルエラー

個別のファイル処理でエラーが発生した場合、エラーをログに記録して次のファイルの処理を継続します。

```bash
[3/15] Converting corrupted.jpg... FAILED (failed to decode image)
```

### 終了コード

- `0`: すべての画像が正常に変換された
- `1`: 1つ以上の画像の変換に失敗した、または設定エラーが発生した

## パフォーマンス

- 高品質なCatmullRomスケーラーを使用した画像リサイズ
- 効率的なバッチ処理
- メモリ効率の良い画像処理

## トラブルシューティング

### 入力ディレクトリが見つからない

```bash
ERROR: Input directory not found - /path/to/input
```

指定したパスが正しいか確認してください。

### 出力ディレクトリの作成に失敗

```bash
ERROR: Failed to create output directory - /path/to/output
```

出力先のパスに書き込み権限があるか確認してください。

### サポートされていないフォーマット

```bash
[5/10] Converting image.xyz... SKIPPED (unsupported format)
```

サポートされているフォーマット（JPEG、PNG、WebP、GIF、BMP）のみが処理されます。

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 貢献

バグ報告や機能リクエストは、GitHubのIssueでお願いします。

## 技術仕様

- **言語**: Go 1.21+
- **画像処理**: `golang.org/x/image`
- **リサイズアルゴリズム**: Catmull-Rom補間
- **対応プラットフォーム**: macOS、Linux、Windows
