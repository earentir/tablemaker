package output

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// SystemFont represents a system font with fallback options
type SystemFont struct {
	Names []string // Font names to search for (in order of preference)
	Path  string   // Resolved path to the font file
}

// Common system fonts by category
var (
	SystemSansSerif = SystemFont{
		Names: []string{"DejaVu Sans", "Arial", "Helvetica", "Liberation Sans", "FreeSans", "Noto Sans", "Source Code Pro"},
	}
	SystemSansSerifBold = SystemFont{
		Names: []string{"DejaVu Sans Bold", "Arial Bold", "Helvetica Bold", "Liberation Sans Bold", "FreeSans Bold", "Noto Sans Bold", "Source Code Pro Bold"},
	}
	SystemMonospace = SystemFont{
		Names: []string{"DejaVu Sans Mono", "Courier New", "Liberation Mono", "FreeMono", "Noto Sans Mono", "Consolas", "Source Code Pro"},
	}
)

// getFontDirectories returns platform-specific font directories
func getFontDirectories() []string {
	var dirs []string

	switch runtime.GOOS {
	case "windows":
		dirs = []string{
			filepath.Join(os.Getenv("WINDIR"), "Fonts"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "Windows", "Fonts"),
		}
	case "darwin": // macOS
		dirs = []string{
			"/System/Library/Fonts",
			"/Library/Fonts",
			filepath.Join(os.Getenv("HOME"), "Library", "Fonts"),
		}
	case "linux":
		dirs = []string{
			"/usr/share/fonts",
			"/usr/local/share/fonts",
			filepath.Join(os.Getenv("HOME"), ".fonts"),
			filepath.Join(os.Getenv("HOME"), ".local/share/fonts"),
			"/var/lib/defoma/fontconfig.d/D", // Debian/Ubuntu
		}
	default:
		// Generic Unix-like system
		dirs = []string{
			"/usr/share/fonts",
			"/usr/local/share/fonts",
		}
		if home := os.Getenv("HOME"); home != "" {
			dirs = append(dirs, filepath.Join(home, ".fonts"))
			dirs = append(dirs, filepath.Join(home, ".local/share/fonts"))
		}
	}

	// Add current directory and local fonts directory
	dirs = append([]string{"./fonts", "."}, dirs...)

	return dirs
}

// findSystemFont searches for a font file in system directories
func findSystemFont(fontNames []string) (string, error) {
	dirs := getFontDirectories()

	// ONLY use TTF files - freetype library can't handle OTF properly
	extensions := []string{".ttf", ".TTF"}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}

		for _, fontName := range fontNames {
			// Try different naming conventions
			namingVariations := []string{
				fontName,
				strings.ReplaceAll(fontName, " ", ""),
				strings.ReplaceAll(fontName, " ", "-"),
				strings.ReplaceAll(fontName, " ", "_"),
				strings.ToLower(strings.ReplaceAll(fontName, " ", "")),
				strings.ToLower(strings.ReplaceAll(fontName, " ", "-")),
			}

			// Add Source Code Pro specific variations
			if strings.Contains(strings.ToLower(fontName), "source code pro") {
				if strings.Contains(strings.ToLower(fontName), "bold") {
					namingVariations = append(namingVariations, "SourceCodePro-Bold", "SourceCodePro-Semibold")
				} else {
					namingVariations = append(namingVariations, "SourceCodePro-Regular")
				}
			}

			for _, name := range namingVariations {
				for _, ext := range extensions {
					fontPath := filepath.Join(dir, name+ext)
					if _, err := os.Stat(fontPath); err == nil {
						return fontPath, nil
					}
				}
			}
		}

		// Also search recursively in subdirectories (limited depth) - TTF ONLY
		if foundPath := searchFontInSubdirs(dir, fontNames, 2); foundPath != "" {
			return foundPath, nil
		}
	}

	return "", fmt.Errorf("font not found: tried %v", fontNames)
}

// searchFontInSubdirs searches for fonts in subdirectories up to maxDepth
func searchFontInSubdirs(baseDir string, fontNames []string, maxDepth int) string {
	if maxDepth <= 0 {
		return ""
	}

	// ONLY TTF files
	extensions := []string{".ttf", ".TTF"}
	var foundPath string

	filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking even if there's an error
		}

		// Skip if we're too deep
		relPath, _ := filepath.Rel(baseDir, path)
		if strings.Count(relPath, string(filepath.Separator)) > maxDepth {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return nil
		}

		filename := info.Name()
		for _, ext := range extensions {
			if strings.HasSuffix(filename, ext) {
				baseFilename := strings.TrimSuffix(filename, ext)

				for _, fontName := range fontNames {
					namingVariations := []string{
						fontName,
						strings.ReplaceAll(fontName, " ", ""),
						strings.ReplaceAll(fontName, " ", "-"),
						strings.ReplaceAll(fontName, " ", "_"),
						strings.ToLower(strings.ReplaceAll(fontName, " ", "")),
						strings.ToLower(strings.ReplaceAll(fontName, " ", "-")),
					}

					// Add Source Code Pro specific variations for subdirectory search
					if strings.Contains(strings.ToLower(fontName), "source code pro") {
						if strings.Contains(strings.ToLower(fontName), "bold") {
							namingVariations = append(namingVariations, "SourceCodePro-Bold", "SourceCodePro-Semibold")
						} else {
							namingVariations = append(namingVariations, "SourceCodePro-Regular")
						}
					}

					for _, name := range namingVariations {
						if strings.EqualFold(baseFilename, name) {
							foundPath = path
							return filepath.SkipAll // Stop walking when found
						}
					}
				}
			}
		}

		return nil
	})

	return foundPath
}

// resolveSystemFont tries to resolve a font path, supporting both explicit paths and system font names
func resolveSystemFont(fontPath string) (string, error) {
	// If it's already a valid file path, use it
	if _, err := os.Stat(fontPath); err == nil {
		return fontPath, nil
	}

	// Check if it's a relative path that might exist
	if !filepath.IsAbs(fontPath) {
		if _, err := os.Stat(fontPath); err == nil {
			return fontPath, nil
		}
	}

	// Try to find it as a system font name
	return findSystemFont([]string{fontPath})
}

// getDefaultSystemFonts returns default system fonts with fallbacks
func getDefaultSystemFonts() (string, string, string, error) {
	// Try to find system fonts
	sansSerifPath, err1 := findSystemFont(SystemSansSerif.Names)
	boldPath, err2 := findSystemFont(SystemSansSerifBold.Names)
	monospacePath, err3 := findSystemFont(SystemMonospace.Names)

	// If we can't find system fonts, provide helpful error message
	if err1 != nil && err2 != nil && err3 != nil {
		return "", "", "", fmt.Errorf("no system fonts found. Please specify font paths in your JSON config or place TTF files in ./fonts/ directory. Searched for: %v, %v, %v",
			SystemSansSerif.Names, SystemSansSerifBold.Names, SystemMonospace.Names)
	}

	// Use fallbacks if some fonts are missing
	if err1 != nil {
		sansSerifPath = boldPath // Use bold as fallback for regular
	}
	if err2 != nil {
		boldPath = sansSerifPath // Use regular as fallback for bold
	}
	if err3 != nil {
		monospacePath = sansSerifPath // Use sans-serif as fallback for monospace
	}

	return sansSerifPath, boldPath, monospacePath, nil
}
