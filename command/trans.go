package command

import (
	"fmt"
	"helper/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	c "github.com/otiai10/copy"
)

func TransXCFramework(xcframework, output string) error {
	libName := utils.GetFileNameWithExt(xcframework)
	workDir, err := ioutil.TempDir("", "helper")
	if nil != err {
		panic("ERROR, create temp dir failed")
	}

	archFiles := make([]string, 4)
	templateFW := ""
	utils.Walk(xcframework, 3, func(fPath string, fi os.FileInfo, err error) error {
		if fi.Name() != libName {
			return nil
		}
		parentName := filepath.Base(filepath.Dir(fPath))
		if parentName != libName+".framework" {
			return nil
		}
		if "" == templateFW {
			templateFW = filepath.Dir(fPath)
		}
		archFiles = append(archFiles, splitLib(fPath, libName, workDir)...)

		return nil
	})

	fatLibPath := mergeLib(archFiles, libName, workDir)

	c.Copy(templateFW, filepath.Join(output, libName+".framework"))
	c.Copy(fatLibPath, filepath.Join(output, libName+".framework", libName))

	os.RemoveAll(workDir)

	return nil
}

func splitLib(libPath, libName, workDir string) []string {

	// check lib arch
	output := utils.RunCmd(fmt.Sprintf("lipo -archs %s", libPath), workDir, true)

	// split arch
	archFiles := make([]string, 4)
	archs := strings.Split(string(output), " ")
	if 0 == len(archs) {
		return []string{}
	}
	if 1 == len(archs) {
		return []string{libPath}
	}
	for _, arch := range archs {
		arch = strings.Trim(arch, "\n")
		arch = strings.Trim(arch, "\t")
		fileName := fmt.Sprintf("%s-%s", libName, arch)
		utils.RunCmd(fmt.Sprintf("lipo %s -thin %s -output %s", libPath, arch, fileName), workDir, true)
		archFiles = append(archFiles, fileName)
	}

	return archFiles
}

func mergeLib(archFiles []string, libName, workDir string) string {
	// merge arch
	utils.RunCmd(fmt.Sprintf("lipo -create %s -output %s", strings.Join(archFiles, " "), libName), workDir, true)

	return filepath.Join(workDir, libName)
}
