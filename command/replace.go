package command

import (
	"helper/utils"
	"os"
	"path/filepath"
	"strings"

	c "github.com/otiai10/copy"
)

func Replace(src, dst, root string) error {

	filepath.Walk(root, func(path string, fi os.FileInfo, err error) error {
		if strings.HasSuffix(fi.Name(), ".h") ||
			strings.HasSuffix(fi.Name(), ".hpp") ||
			strings.HasSuffix(fi.Name(), ".m") ||
			strings.HasSuffix(fi.Name(), ".mm") ||
			strings.HasSuffix(fi.Name(), ".c") ||
			strings.HasSuffix(fi.Name(), ".cpp") {
			fileContent := utils.ReadFileContent(path)
			if strings.Contains(fileContent, src) {
				backupFile(root, path)
				fileContent = strings.ReplaceAll(fileContent, src, dst)
				utils.Write2File(fileContent, path)
			}
		}
		return nil
	})
	recoverFile(root)

	return nil
}

func backupFile(root, path string) {
	relPath, err := filepath.Rel(root, path)
	if nil != err {
		return
	}
	c.Copy(path, filepath.Join(root, ".replace_backup", relPath))
}

func recoverFile(root string) {
	backFolder := filepath.Join(root, ".replace_backup")
	c.Copy(backFolder, root)
	os.RemoveAll(backFolder)
}
