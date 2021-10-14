package command

import (
	"os"
	"bufio"
	"path/filepath"
	"strings"
	"helper/log"
)

func LineCounter(root string) error {
	lineSum := 0
	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(fi.Name(), ".h") ||
			strings.HasSuffix(fi.Name(), ".hpp") ||
			strings.HasSuffix(fi.Name(), ".m") ||
			strings.HasSuffix(fi.Name(), ".mm") ||
			strings.HasSuffix(fi.Name(), ".c") ||
			strings.HasSuffix(fi.Name(), ".cpp") ||
			strings.HasSuffix(fi.Name(), ".cs") {
				lineSum += countLine(path)
		}

		return nil
	})

	log.Info("Total Lines: ", lineSum)
	return nil
}

func countLine(filePath string) int {
	lineSum := 0
    f, _ := os.Open(filePath)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
		lineSum ++
    }

    return lineSum
}
