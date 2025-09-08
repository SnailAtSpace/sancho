package funcs

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

// RetroCompressionOptions controls the various effects to apply
type RetroCompressionOptions struct {
	JpegQuality       int     // 1-100, lower is more compressed
	CompressionCycles int     // How many times to re-compress
	DownscaleFactor   float64 // Factor to temporarily reduce size (0.5 = half size)
	ColorReduction    int     // Divide RGB values by this amount then multiply (higher = fewer colors)
	SaturationBoost   float64 // 1.0 = normal, 1.2 = 20% more saturated
	NoiseAmount       float64 // 0.0-1.0, amount of noise to add
}

// CompressGifToRetroJpeg applies mid-2000s JPEG-like compression effects to a GIF
func CompressGifToRetroJpeg(inputPath, outputPath string, options RetroCompressionOptions) error {
	// Seed the random number generator for noise
	//rand.Seed(time.Now().UnixNano())

	// Open the input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer inputFile.Close()

	// Decode the GIF
	img, err := gif.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("error decoding GIF: %w", err)
	}

	// Get image bounds
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	// Step 1: Resize to a smaller resolution
	if options.DownscaleFactor > 0 && options.DownscaleFactor < 1.0 {
		newWidth := uint(float64(width) * options.DownscaleFactor)
		newHeight := uint(float64(height) * options.DownscaleFactor)

		// Resize down
		img = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

		// Resize back up to original size for that pixelated look
		img = resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)
	}

	// Create a new RGBA image
	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)

	// Step 2: Color palette reduction and Step 3: Saturation increase
	if options.ColorReduction > 1 || options.SaturationBoost != 1.0 {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := rgbaImg.At(x, y).RGBA()

				// Convert from uint32 to uint8
				r8 := uint8(r >> 8)
				g8 := uint8(g >> 8)
				b8 := uint8(b >> 8)
				a8 := uint8(a >> 8)

				// Color reduction
				if options.ColorReduction > 1 {
					r8 = uint8((int(r8) / options.ColorReduction) * options.ColorReduction)
					g8 = uint8((int(g8) / options.ColorReduction) * options.ColorReduction)
					b8 = uint8((int(b8) / options.ColorReduction) * options.ColorReduction)
				}

				// Saturation boost
				if options.SaturationBoost != 1.0 {
					// Convert to HSL, boost saturation, convert back
					// Simple RGB saturation boost (not perfect but works for this effect)
					avg := (float64(r8) + float64(g8) + float64(b8)) / 3
					r8 = clamp(uint8(avg + (float64(r8)-avg)*options.SaturationBoost))
					g8 = clamp(uint8(avg + (float64(g8)-avg)*options.SaturationBoost))
					b8 = clamp(uint8(avg + (float64(b8)-avg)*options.SaturationBoost))
				}

				// Step 4: Add noise
				

				rgbaImg.Set(x, y, color.RGBA{r8, g8, b8, a8})
			}
		}
	}

	// Apply multiple JPEG compression cycles
	jpegImg := image.Image(rgbaImg)
	for i := 0; i < options.CompressionCycles; i++ {
		var jpegBuf bytes.Buffer
		err = jpeg.Encode(&jpegBuf, jpegImg, &jpeg.Options{
			Quality: options.JpegQuality,
		})
		if err != nil {
			return fmt.Errorf("error during JPEG compression cycle %d: %w", i, err)
		}

		jpegImg, err = jpeg.Decode(bytes.NewReader(jpegBuf.Bytes()))
		if err != nil {
			return fmt.Errorf("error during JPEG decompression cycle %d: %w", i, err)
		}
	}

	// Open output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	// Re-encode as GIF
	return gif.Encode(outputFile, jpegImg, nil)
}

// clamp ensures a uint8 value stays between 0-255
func clamp(v uint8) uint8 {
	return v
}

func main() {
	// Example usage with all effects enabled
	options := RetroCompressionOptions{
		JpegQuality:       15,  // Very low quality (1-100)
		CompressionCycles: 2,   // Apply compression twice
		DownscaleFactor:   0.5, // Resize to 50% and back
		ColorReduction:    16,  // Reduce color palette
		SaturationBoost:   1.2, // Increase saturation by 20%
		NoiseAmount:       0, // Add a moderate amount of noise
	}

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go input.gif output.gif")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := os.Args[2]

	fmt.Println("Applying retro compression to", inputPath)
	err := CompressGifToRetroJpeg(inputPath, outputPath, options)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	fmt.Println("Successfully created", outputPath)
}
