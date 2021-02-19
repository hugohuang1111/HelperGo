package utils

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Zip(fPath, zPath, dPath string) error {
	outFile, err := os.Create(dPath)
	if nil != err {
		return err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)

	addToZip(w, fPath, zPath)

	return w.Close()
}

func addFileToZip(w *zip.Writer, fPath, zPath, fName string) error {
	dat, err := ioutil.ReadFile(filepath.Join(fPath, fName))
	if nil != err {
		return err
	}

	f, err := w.Create(filepath.Join(zPath, fName))
	if nil != err {
		return err
	}
	_, err = f.Write(dat)
	return err
}

func addFolderToZip(w *zip.Writer, fPath, zPath string) error {
	fis, err := ioutil.ReadDir(fPath)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			err = addFileToZip(w, fPath, zPath, fi.Name())
		} else if fi.IsDir() {
			err = addFolderToZip(w, filepath.Join(fPath, fi.Name()), filepath.Join(zPath, fi.Name()))
		}
	}
	return err
}

func addToZip(w *zip.Writer, fPath, zPath string) {
	if isFile(fPath) {
		addFileToZip(w, filepath.Dir(fPath), zPath, filepath.Base(fPath))
	} else {
		addFolderToZip(w, fPath, zPath)
	}

}
