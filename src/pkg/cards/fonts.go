package cards

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/agladfield/postcart/pkg/shared/enum"
)

const cardsFontErrFmtStr = "cards fonts err: %w"

func ensureDirectory(dirPath string) error {
	// Check if directory exists
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		// Create directory with 0700 permissions (rwx for user)
		err = os.MkdirAll(dirPath, 0700)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
		}
		fmt.Printf("Created directory: %s\n", dirPath)
		return nil
	}
	if err != nil {
		return fmt.Errorf("error checking directory %s: %w", dirPath, err)
	}
	// fmt.Printf("Directory already exists: %s\n", dirPath)
	return nil
}

func fontExists(fontPath string) bool {
	_, err := os.Stat(fontPath)
	return !os.IsNotExist(err)
}

func installFonts() error {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf(cardsFontErrFmtStr, homeErr)
	}
	var fontsDir string
	switch runtime.GOOS {
	case "darwin":
		fontsDir = path.Join(home, "Library", "Fonts")
	case "linux":
		fontsDir = "/usr/share/fonts"
	case "windows":
		fontsDir = path.Join(home, "AppData", "Local", "Microsoft", "Windows", "Fonts")
	default:
		return fmt.Errorf(cardsFontErrFmtStr, errors.New("unsupported platform for installing fonts"))
	}

	// fontTempDir = tempDir
	dirErr := ensureDirectory(fontsDir)
	if dirErr != nil {
		return fmt.Errorf(cardsFontErrFmtStr, dirErr)
	}

	fontEntries, err := fontResources.ReadDir("res/fonts")
	if err != nil {
		return fmt.Errorf(cardsFontErrFmtStr, fmt.Errorf("failed to read embedded fonts directory: %w", err))
	}

	for _, entry := range fontEntries {
		if entry.IsDir() {
			continue // Skip directories
		}
		fontData, err := fontResources.ReadFile(filepath.Join("res/fonts", entry.Name()))
		if err != nil {
			return fmt.Errorf(cardsFontErrFmtStr, fmt.Errorf("failed to read font file %s: %w", entry.Name(), err))
		}
		fontPath := filepath.Join(fontsDir, entry.Name())
		if !fontExists(fontPath) {
			if err := os.WriteFile(fontPath, fontData, 0644); err != nil {
				return fmt.Errorf(cardsFontErrFmtStr, fmt.Errorf("failed to write font file %s: %w", entry.Name(), err))
			}
		}

	}

	return nil
}

func getTextFont(font enum.FontEnum) string {
	return fmt.Sprintf("%s 10", fontsMap[font])
}

var fontsMap = map[enum.FontEnum]string{
	enum.FontMarker:     "Fuzzy Bubbles",
	enum.FontPolite:     "Kavivanar",
	enum.FontTypewriter: "IM FELL English",
	enum.FontMidCentury: "Aoboshi One",
}
