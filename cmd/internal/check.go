package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"github.com/sergi/go-diff/diffmatchpatch"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Check(filePath string) error {
	extension := filepath.Ext(filePath)

	if extension != ".tf" && extension != ".tfvars" {
		return nil
	}

	original, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	formattedBytes, err := getFormattedContent(filePath)
	if err != nil {
		return err
	}

	return compare(original, formattedBytes)
}

func compare(original []byte, formatted []byte) error {
	originalHash, formattedHash := generateHash(original, formatted)

	if formattedHash == originalHash {
		return nil
	}

	return compareContent(original, formatted)
}

func generateHash(original []byte, formatted []byte) (string, string) {
	hasher := sha1.New()
	hasher.Write(original)

	originalHash := hex.EncodeToString(hasher.Sum(nil))

	hasher.Reset()

	hasher.Write(formatted)

	formattedHash := hex.EncodeToString(hasher.Sum(nil))
	return originalHash, formattedHash
}

func compareContent(original []byte, formatted []byte) error {
	dmp := diffmatchpatch.New()
	dmp.DiffTimeout = time.Hour
	src := string(original)
	dst := string(formatted)

	wSrc, wDst, warray := dmp.DiffLinesToRunes(src, dst)
	diffs := dmp.DiffMainRunes(wSrc, wDst, false)
	diffs = dmp.DiffCharsToLines(diffs, warray)

	var notEquals []diffmatchpatch.Diff
	for _, diff := range diffs {
		if diff.Type != diffmatchpatch.DiffEqual {
			notEquals = append(notEquals, diff)
		}
	}

	if notEquals == nil || len(notEquals) == 0 {
		return nil
	}

	var errorText strings.Builder
	errorText.WriteString("\n")
	errorText.WriteString(dmp.DiffPrettyText(diffs))
	return errors.New(errorText.String())
}
