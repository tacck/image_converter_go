# Image Converter CLI

A high-performance batch image conversion tool implemented in Go. It processes all image files in a specified directory, performing resizing and format conversion.

## Features

- 🚀 **Fast Batch Conversion**: Process all images in a directory at once
- 📐 **Flexible Resizing**: Free size adjustment by scale or pixel specification
- 🎨 **Multiple Format Support**: Supports JPEG, PNG, WebP, GIF, and BMP
- 📏 **Aspect Ratio Preservation**: Automatically maintains aspect ratio during resizing
- ⚙️ **Quality Control**: Fine-tune JPEG output quality
- 📊 **Detailed Progress Display**: Monitor processing status in real-time

## Installation

### Build from Source

```bash
git clone <repository-url>
cd image-converter
go build -o image-converter ./cmd/image-converter
```

### Install Binary

Place the built executable in a directory in your PATH:

```bash
# macOS/Linux
sudo mv image-converter /usr/local/bin/

# Or place in your home bin folder
mkdir -p ~/bin
mv image-converter ~/bin/
export PATH="$HOME/bin:$PATH"  # Add to .bashrc or .zshrc
```

## Usage

### Basic Usage

```bash
image-converter -input-dir <input-directory> -output-dir <output-directory> [options]
```

### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-input-dir` | Input directory path (required) | - |
| `-output-dir` | Output directory path (required) | - |
| `-scale` | Image scale factor (e.g., 0.5 for 50%, 2.0 for 200%) | - |
| `-width` | Output image width (pixels) | - |
| `-height` | Output image height (pixels) | - |
| `-format` | Output format (jpeg, png, webp, gif, bmp) | Original format |
| `-jpeg-quality` | JPEG quality (1-100) | 85 |

### Examples

#### 1. Reduce images to 50%

```bash
image-converter -input-dir ./photos -output-dir ./thumbnails -scale 0.5
```

#### 2. Resize to 800 pixels width (maintaining aspect ratio)

```bash
image-converter -input-dir ./photos -output-dir ./resized -width 800
```

#### 3. Resize to 600 pixels height

```bash
image-converter -input-dir ./photos -output-dir ./resized -height 600
```

#### 4. Fit within 800x600

```bash
image-converter -input-dir ./photos -output-dir ./resized -width 800 -height 600
```

When both width and height are specified, images are resized to fit within the specified dimensions while maintaining aspect ratio.

#### 5. Convert all images to JPEG

```bash
image-converter -input-dir ./photos -output-dir ./converted -format jpeg
```

#### 6. Convert to WebP and reduce to 50%

```bash
image-converter -input-dir ./photos -output-dir ./optimized -scale 0.5 -format webp
```

#### 7. Convert with specified JPEG quality

```bash
image-converter -input-dir ./photos -output-dir ./compressed -format jpeg -jpeg-quality 70
```

#### 8. Format conversion only (no resizing)

```bash
image-converter -input-dir ./photos -output-dir ./converted -format png
```

## Supported Formats

### Input Formats

- JPEG (.jpg, .jpeg)
- PNG (.png)
- WebP (.webp)
- GIF (.gif)
- BMP (.bmp)

### Output Formats

- JPEG (.jpeg)
- PNG (.png)
- WebP (.webp)
- GIF (.gif)
- BMP (.bmp)

## Resize Specifications

### Scale Factor (-scale)

Resize by multiplying the original image size by the scale factor.

- `0.5`: Reduce to 50%
- `1.0`: Maintain original size
- `2.0`: Enlarge to 200%

### Pixel Specification (-width, -height)

#### Width Only

Resize to the specified width and automatically calculate height to maintain aspect ratio.

#### Height Only

Resize to the specified height and automatically calculate width to maintain aspect ratio.

#### Both Width and Height

Resize to fit within both specified width and height while maintaining aspect ratio.

### Constraints

- Scale and pixel specifications cannot be used simultaneously
- Aspect ratio is maintained in all resize operations
- Original size is maintained if no size specification is provided

## Output Files

### File Names

Output files retain the base name of the original file, with only the extension changed.

Examples:
- Input: `photo.jpg`, Output format: `png` → Output: `photo.png`
- Input: `image.png`, No format specified → Output: `image.png`

### Overwriting Existing Files

If a file with the same name already exists in the output directory, it will be overwritten.

## Progress Display

During processing, progress is displayed as follows:

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

## Error Handling

### Configuration Errors

If invalid configuration is specified, an error message is displayed and the program exits immediately.

```bash
ERROR: Cannot use both scale and pixel specifications simultaneously
```

### File Errors

If an error occurs during individual file processing, the error is logged and processing continues with the next file.

```bash
[3/15] Converting corrupted.jpg... FAILED (failed to decode image)
```

### Exit Codes

- `0`: All images were successfully converted
- `1`: One or more images failed to convert, or a configuration error occurred

## Performance

- Image resizing using high-quality Catmull-Rom scaler
- Efficient batch processing with concurrent workers
- Memory-efficient image processing

## Troubleshooting

### Input Directory Not Found

```bash
ERROR: Input directory not found - /path/to/input
```

Verify that the specified path is correct.

### Failed to Create Output Directory

```bash
ERROR: Failed to create output directory - /path/to/output
```

Check that you have write permissions for the output path.

### Unsupported Format

```bash
[5/10] Converting image.xyz... SKIPPED (unsupported format)
```

Only supported formats (JPEG, PNG, WebP, GIF, BMP) are processed.

## License

This project is released under the MIT License.

## Contributing

Bug reports and feature requests are welcome via GitHub Issues.

## Technical Specifications

- **Language**: Go 1.21+
- **Image Processing**: `golang.org/x/image`
- **Resize Algorithm**: Catmull-Rom interpolation
- **Supported Platforms**: macOS, Linux, Windows
