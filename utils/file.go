package utils

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func rmFile(dir string) {
	dir = expendPath(dir)
	if !FileExist(dir) {
		return
	}
	if isFile(dir) {
		os.Remove(dir)
	} else {
		os.RemoveAll(dir)
	}
}

func expendPath(dir string) string {
	if !strings.HasPrefix(dir, "~") {
		return dir
	}
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	return strings.Replace(dir, "~", homeDir, 1)
}

func FileExist(fPath string) bool {
	if _, err := os.Stat(fPath); err == nil {
		// path/to/whatever exists
		return true
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		return false
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence

		return false
	}
}

func isFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}

	return !fi.IsDir()
}

func ReadLines(fPath string, cb func(line string) (skip bool)) {
	var lineBytes []byte
	prefix := false
	file, err := os.Open(fPath)
	if nil != err {
		return
	}
	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	for {
		lineBytes, prefix, err = reader.ReadLine()
		if nil != err {
			break
		}
		buffer.Write(lineBytes)
		if !prefix {
			if cb(buffer.String()) {
				break
			}
			buffer.Reset()
		}
	}
}

func WriteLines(lines []string, fPath string) {
	f, err := os.OpenFile(fPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if nil != err {
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	for _, s := range lines {
		w.WriteString(s)
		w.WriteString("\n")
	}

	w.Flush()
}

func AppendLines(lines []string, fPath string) {
	f, err := os.OpenFile(fPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if nil != err {
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	for _, s := range lines {
		w.WriteString(s)
		w.WriteString("\n")
	}

	w.Flush()
}

func Write2File(s string, fPath string) {
	f, err := os.OpenFile(fPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if nil != err {
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	w.WriteString(s)

	w.Flush()
}

func Walk(dirpath string, maxDepth int, f filepath.WalkFunc) {
	if isFile(dirpath) {
		fi, err := os.Stat(dirpath)
		f(dirpath, fi, err)
		return
	} else {
		walkDir(dirpath, 0, maxDepth, f)
	}
}

func walkDir(dirpath string, depth, maxDepth int, f filepath.WalkFunc) {
	if depth > maxDepth {
		return
	}
	fi, err := os.Stat(dirpath)
	if filepath.SkipDir == f(dirpath, fi, err) {
		return
	}
	fis, err := ioutil.ReadDir(dirpath)

	for _, fi := range fis {
		if !fi.IsDir() {
			if filepath.SkipDir == f(filepath.Join(dirpath, fi.Name()), fi, nil) {
				return
			}
		}
	}

	for _, fi := range fis {
		if fi.IsDir() {
			walkDir(filepath.Join(dirpath, fi.Name()), depth+1, maxDepth, f)
		}
	}
}

func FileNameWithoutSuffix(fileName string) string {
	base := filepath.Base(fileName)

	return base[:len(base)-len(filepath.Ext(base))]
}

func ListFiles(baseDir string, s string) []string {
	files := make([]string, 0)
	r, err := regexp.Compile(s)
	if nil != err {
		r = nil
	}
	filepath.Walk(baseDir, func(fPath string, fi os.FileInfo, err error) error {
		if nil != r {
			if r.Match([]byte(fPath)) {
				files = append(files, fPath)
			}
		} else if strings.HasSuffix(fPath, s) {
			files = append(files, fPath)
		}
		return nil
	})

	return files
}

func ExecPath() (string, error) {
	fPath, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(fPath)
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetInstallerHomeDir() string {
	home, err := homedir.Dir()
	if nil != err {
		return ""
	}

	return filepath.Join(home, ".sdkbox")
}

// MakeSureDirExist make sure directory exists
func MakeSureDirExist(path string) {
	dir := filepath.Dir(path)
	if !FileExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func IsFilePath(s string) bool {
	if strings.HasPrefix(s, "/") {
		return true
	}
	if strings.HasPrefix(s, "./") {
		return true
	}
	r, _ := regexp.Compile(`^[a-zA-Z]:\\`)
	if r.MatchString(s) {
		return true
	}

	return false
}

func GetFileNameWithExt(s string) string {
	fileName := filepath.Base(s)
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func ReadFileContent(fPath string) string {
	f, err := os.Open(fPath)
	if err != nil {
		return ""
	}

	defer f.Close()
	var chunk []byte
	buf := make([]byte, 1024)

	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return ""
		}
		if n == 0 {
			break
		}
		chunk = append(chunk, buf[:n]...)
	}

	return string(chunk)
}
