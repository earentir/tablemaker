package output

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// PNGConfig contains PNG rendering options
type PNGConfig struct {
	TitleFont   FontConfig `json:"title_font"`
	ContentFont FontConfig `json:"content_font"`
	ASCIIFont   FontConfig `json:"ascii_font"`
}

// FontConfig represents font configuration
type FontConfig struct {
	Path string  `json:"path"`
	Size float64 `json:"size"`
}

// TextType represents different types of text in the table
type TextType int

const (
	ASCIIText TextType = iota
	HeaderText
	ContentText
)

// TextSegment represents a segment of text with its type
type TextSegment struct {
	Text string
	Type TextType
}

// GeneratePNG creates a PNG image from the table text
func GeneratePNG(tableText string, config PNGConfig, outputPath string) error {
	// Load fonts with system font support
	fonts := make(map[TextType]*truetype.Font)

	// Use system defaults if paths are not specified or invalid
	var asciiPath, titlePath, contentPath string
	var err error

	if config.ASCIIFont.Path == "" || config.TitleFont.Path == "" || config.ContentFont.Path == "" {
		// Try to get default system fonts
		defaultContent, defaultTitle, defaultMono, sysErr := getDefaultSystemFonts()
		if sysErr != nil {
			return fmt.Errorf("failed to find system fonts and no paths specified: %v", sysErr)
		}

		if config.ASCIIFont.Path == "" {
			asciiPath = defaultMono
		} else {
			asciiPath = config.ASCIIFont.Path
		}

		if config.TitleFont.Path == "" {
			titlePath = defaultTitle
		} else {
			titlePath = config.TitleFont.Path
		}

		if config.ContentFont.Path == "" {
			contentPath = defaultContent
		} else {
			contentPath = config.ContentFont.Path
		}
	} else {
		asciiPath = config.ASCIIFont.Path
		titlePath = config.TitleFont.Path
		contentPath = config.ContentFont.Path
	}

	asciiFont, err := loadFont(asciiPath)
	if err != nil {
		return fmt.Errorf("failed to load ASCII font: %v", err)
	}
	fonts[ASCIIText] = asciiFont

	titleFont, err := loadFont(titlePath)
	if err != nil {
		return fmt.Errorf("failed to load title font: %v", err)
	}
	fonts[HeaderText] = titleFont

	contentFont, err := loadFont(contentPath)
	if err != nil {
		return fmt.Errorf("failed to load content font: %v", err)
	}
	fonts[ContentText] = contentFont

	// Parse the table text and create segments
	lines := strings.Split(tableText, "\n")
	var segments [][]TextSegment

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		segments = append(segments, parseLineSegments(line))
	}

	// Calculate image dimensions
	width, height := calculateImageDimensions(segments, fonts, config)

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with transparent background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.Point{}, draw.Src)

	// Render text
	if err := renderText(img, segments, fonts, config); err != nil {
		return fmt.Errorf("failed to render text: %v", err)
	}

	// Save PNG
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode PNG: %v", err)
	}

	return nil
}

// loadFont loads a TrueType font from file, with system font resolution
func loadFont(fontPath string) (*truetype.Font, error) {
	// Try to resolve system font if it's not a direct path
	resolvedPath, err := resolveSystemFont(fontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve font '%s': %v", fontPath, err)
	}

	fontBytes, err := ioutil.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file '%s': %v", resolvedPath, err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font file '%s': %v", resolvedPath, err)
	}

	return font, nil
}

// parseLineSegments parses a line into text segments with their types
func parseLineSegments(line string) []TextSegment {
	var segments []TextSegment

	// ASCII characters (box drawing) - expanded to include all table styles
	asciiChars := "┌┐└┘├┤┬┴┼─│╔╗╚╝╠╣╦╩╬═║"

	// Bold text pattern
	boldPattern := regexp.MustCompile(`\*\*([^*]+)\*\*`)

	// Find all bold matches
	boldMatches := boldPattern.FindAllStringSubmatchIndex(line, -1)

	lastIndex := 0
	for _, match := range boldMatches {
		start, end := match[0], match[1]

		// Add text before bold
		if start > lastIndex {
			beforeText := line[lastIndex:start]
			segments = append(segments, classifyTextSegments(beforeText, asciiChars)...)
		}

		// Add bold text as header
		boldText := line[match[2]:match[3]]
		segments = append(segments, TextSegment{Text: boldText, Type: HeaderText})

		lastIndex = end
	}

	// Add remaining text
	if lastIndex < len(line) {
		remainingText := line[lastIndex:]
		segments = append(segments, classifyTextSegments(remainingText, asciiChars)...)
	}

	return segments
}

// classifyTextSegments classifies text into ASCII or content segments
func classifyTextSegments(text string, asciiChars string) []TextSegment {
	var segments []TextSegment
	var currentSegment strings.Builder
	var currentType TextType = -1

	// Expand ASCII chars to include all common table drawing characters
	allASCIIChars := asciiChars + "┌┐└┘├┤┬┴┼─│╔╗╚╝╠╣╦╩╬═║"

	for _, char := range text {
		charType := ContentText
		if strings.ContainsRune(allASCIIChars, char) {
			charType = ASCIIText
		}

		if currentType == -1 {
			currentType = charType
		}

		if charType == currentType {
			currentSegment.WriteRune(char)
		} else {
			// Type changed, save current segment
			if currentSegment.Len() > 0 {
				segments = append(segments, TextSegment{
					Text: currentSegment.String(),
					Type: currentType,
				})
			}

			// Start new segment
			currentSegment.Reset()
			currentSegment.WriteRune(char)
			currentType = charType
		}
	}

	// Add final segment
	if currentSegment.Len() > 0 {
		segments = append(segments, TextSegment{
			Text: currentSegment.String(),
			Type: currentType,
		})
	}

	return segments
}

// renderText renders all text segments to the image
func renderText(img *image.RGBA, segments [][]TextSegment,
	fonts map[TextType]*truetype.Font, cfg PNGConfig) error {

	monoFace := newFace(fonts[ASCIIText], cfg.ASCIIFont.Size)
	lh := lineHeight(monoFace)
	ascent := monoFace.Metrics().Ascent.Ceil()

	y := 50 + ascent // top padding + ascent
	x := 50          // left padding is unchanged

	for _, lineSegs := range segments {
		var full strings.Builder
		for _, s := range lineSegs {
			full.WriteString(s.Text)
		}

		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(fonts[ASCIIText])
		c.SetFontSize(cfg.ASCIIFont.Size)
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(image.Black)

		// draw the whole line at baseline y
		if _, err := c.DrawString(full.String(), freetype.Pt(x, y)); err != nil {
			return err
		}
		y += lh
	}
	return nil
}

func newFace(f *truetype.Font, size float64) font.Face {
	return truetype.NewFace(f, &truetype.Options{Size: size, DPI: 72})
}

func lineHeight(face font.Face) int {
	m := face.Metrics()
	// Face metrics are 26.6-fixed-point; Ceil() converts to int.
	return (m.Ascent + m.Descent).Ceil()
}

// calculateImageDimensions calculates the required image size
func calculateImageDimensions(segments [][]TextSegment,
	fonts map[TextType]*truetype.Font, cfg PNGConfig) (int, int) {

	monoFace := newFace(fonts[ASCIIText], cfg.ASCIIFont.Size)
	lh := lineHeight(monoFace)

	maxWidth := 0
	for _, lineSegs := range segments {
		var full strings.Builder
		for _, s := range lineSegs {
			full.WriteString(s.Text)
		}
		w := font.MeasureString(monoFace, full.String())
		if int(w>>6) > maxWidth {
			maxWidth = int(w >> 6)
		}
	}

	totalHeight := len(segments) * lh
	// generous padding (unchanged)
	return maxWidth + 100, totalHeight + 100
}
