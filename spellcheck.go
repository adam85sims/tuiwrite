package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/client9/gospell"
)

// Dictionary repository configuration
const dictRepoBase = "https://raw.githubusercontent.com/adam85sims/tuiwritedics/main"

// Available dictionaries
var availableDicts = map[string]struct {
	affFile string
	dicFile string
}{
	"uk": {"en_GB.aff", "en_GB.dic"},
	"gb": {"en_GB.aff", "en_GB.dic"},
	"us": {"en_US.aff", "en_US.dic"},
	"ca": {"en_CA.aff", "en_CA.dic"},
	"au": {"en_AU.aff", "en_AU.dic"},
	"es": {"es_ES.aff", "es_ES.dic"},
	"fr": {"fr_FR.aff", "fr_FR.dic"},
	"de": {"de_DE.aff", "de_DE.dic"},
	"it": {"it_IT.aff", "it_IT.dic"},
	"pt": {"pt_PT.aff", "pt_PT.dic"},
}

// SpellChecker manages spell-checking functionality
type SpellChecker struct {
	checker  *gospell.GoSpell
	enabled  bool
	language string
}

// newSpellChecker creates a new spell checker with the specified language
func newSpellChecker(language string) *SpellChecker {
	sc := &SpellChecker{
		enabled:  false,
		language: language,
	}
	return sc
}

// getDictPath returns the local dictionary directory (cross-platform)
func getDictPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var dictPath string

	// Use platform-appropriate config directory
	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\tuiwrite\dictionaries
		appData := os.Getenv("APPDATA")
		if appData == "" {
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		dictPath = filepath.Join(appData, "tuiwrite", "dictionaries")
	case "darwin":
		// macOS: ~/Library/Application Support/tuiwrite/dictionaries
		dictPath = filepath.Join(home, "Library", "Application Support", "tuiwrite", "dictionaries")
	default:
		// Linux/Unix: ~/.config/tuiwrite/dictionaries
		dictPath = filepath.Join(home, ".config", "tuiwrite", "dictionaries")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dictPath, 0755); err != nil {
		return "", err
	}

	return dictPath, nil
}

// hasDictionary checks if dictionary files exist locally
func (sc *SpellChecker) hasDictionary(lang string) bool {
	dictInfo, exists := availableDicts[lang]
	if !exists {
		return false
	}

	dictPath, err := getDictPath()
	if err != nil {
		return false
	}

	affPath := filepath.Join(dictPath, dictInfo.affFile)
	dicPath := filepath.Join(dictPath, dictInfo.dicFile)

	_, affErr := os.Stat(affPath)
	_, dicErr := os.Stat(dicPath)

	return affErr == nil && dicErr == nil
}

// downloadFile downloads a file from a URL
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		// Check if it's a network connectivity issue
		if strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "network is unreachable") {
			return fmt.Errorf("no internet connection - please connect to the internet and try again")
		}
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// checkInternetConnectivity checks if we can reach the dictionary repo
func checkInternetConnectivity() error {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Head(dictRepoBase + "/en_GB.aff")
	if err != nil {
		if strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "network is unreachable") ||
			strings.Contains(err.Error(), "timeout") {
			return fmt.Errorf("no internet connection detected - please connect to the internet to download dictionaries")
		}
		return fmt.Errorf("cannot reach dictionary repository: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("dictionary repository returned status: %s", resp.Status)
	}

	return nil
}

// downloadDictionary downloads a dictionary from the repo
func (sc *SpellChecker) downloadDictionary(lang string) error {
	dictInfo, exists := availableDicts[lang]
	if !exists {
		return fmt.Errorf("dictionary for language '%s' not available", lang)
	}

	// Check internet connectivity before attempting download
	if err := checkInternetConnectivity(); err != nil {
		return err
	}

	dictPath, err := getDictPath()
	if err != nil {
		return fmt.Errorf("failed to get dictionary path: %w", err)
	}

	// Download .aff file
	affURL := fmt.Sprintf("%s/%s", dictRepoBase, dictInfo.affFile)
	affPath := filepath.Join(dictPath, dictInfo.affFile)
	if err := downloadFile(affURL, affPath); err != nil {
		return fmt.Errorf("failed to download .aff file: %w", err)
	}

	// Download .dic file
	dicURL := fmt.Sprintf("%s/%s", dictRepoBase, dictInfo.dicFile)
	dicPath := filepath.Join(dictPath, dictInfo.dicFile)
	if err := downloadFile(dicURL, dicPath); err != nil {
		return fmt.Errorf("failed to download .dic file: %w", err)
	}

	return nil
}

// loadDictionary loads a dictionary into the spell checker
func (sc *SpellChecker) loadDictionary(lang string) error {
	dictInfo, exists := availableDicts[lang]
	if !exists {
		return fmt.Errorf("dictionary for language '%s' not available", lang)
	}

	dictPath, err := getDictPath()
	if err != nil {
		return err
	}

	affPath := filepath.Join(dictPath, dictInfo.affFile)
	dicPath := filepath.Join(dictPath, dictInfo.dicFile)

	// Load the dictionary using gospell
	checker, err := gospell.NewGoSpell(affPath, dicPath)
	if err != nil {
		return fmt.Errorf("failed to load dictionary: %w", err)
	}

	sc.checker = checker
	sc.language = lang
	sc.enabled = true

	return nil
}

// setLanguage changes the spell-check language
func (sc *SpellChecker) setLanguage(language string) error {
	lang := strings.ToLower(language)

	// Check if dictionary exists locally
	if !sc.hasDictionary(lang) {
		// Download it
		if err := sc.downloadDictionary(lang); err != nil {
			return fmt.Errorf("failed to download dictionary: %w", err)
		}
	}

	// Load the dictionary
	return sc.loadDictionary(lang)
}

// checkWord checks if a word is spelled correctly
func (sc *SpellChecker) checkWord(word string) bool {
	if !sc.enabled || sc.checker == nil {
		return true // All words correct when disabled
	}

	// Skip empty words or single characters
	if len(word) <= 1 {
		return true
	}

	// Skip words with numbers or special characters (likely proper nouns, URLs, etc.)
	hasLetter := false
	for _, r := range word {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
		} else if r != '\'' && r != '-' {
			return true // Contains special chars, skip
		}
	}

	if !hasLetter {
		return true
	}

	// Check against dictionary using gospell
	return sc.checker.Spell(word)
}

// getWordsInLine extracts words and their positions from a line
func getWordsInLine(line string) []wordPos {
	var words []wordPos
	var currentWord strings.Builder
	var startPos int
	inWord := false

	for i, r := range line {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '\'' || r == '-' {
			if !inWord {
				startPos = i
				inWord = true
			}
			currentWord.WriteRune(r)
		} else {
			if inWord {
				words = append(words, wordPos{
					word:  currentWord.String(),
					start: startPos,
					end:   i,
				})
				currentWord.Reset()
				inWord = false
			}
		}
	}

	// Add the last word if we ended in the middle of one
	if inWord {
		words = append(words, wordPos{
			word:  currentWord.String(),
			start: startPos,
			end:   len(line),
		})
	}

	return words
}

// wordPos represents a word and its position in a line
type wordPos struct {
	word  string
	start int
	end   int
}

// toggle enables or disables spell-checking
func (sc *SpellChecker) toggle() {
	sc.enabled = !sc.enabled
}
