package utils

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

type FileList struct {
	FilePath string
	FileDir  string
}

func prepareDirs(dirs []string) ([]string, [][]os.FileInfo) {
	resultDir := make([]string, 0)
	resultDirFileInfo := make([][]os.FileInfo, 0)
	for _, dir := range dirs {
		if fi, err := os.Stat(dir); err != nil {
			if !os.IsNotExist(err) {
				continue
			}
			if err = os.MkdirAll(dir, 0700); err != nil {
				continue
			}
		} else if !fi.IsDir() {
			continue
		}
		if fis, err := ioutil.ReadDir(dir); err != nil {
		} else {
			resultDir = append(resultDir, dir)
			resultDirFileInfo = append(resultDirFileInfo, fis)
		}
	}
	return resultDir, resultDirFileInfo
}

func GetFileList(fileDir, suffix string) []*FileList {
	isRecursion := false
	if len(fileDir) != 0 && fileDir[len(fileDir)-1] == '*' {
		isRecursion = true
		fileDir = fileDir[:len(fileDir)-2]
	}
	arrFileLists := make([]*FileList, 0)
	suffixUp := strings.ToUpper(suffix)
	arrDirs, arrInfos := prepareDirs([]string{fileDir})
	for idx, dbDir := range arrDirs {
		for _, fi := range arrInfos[idx] {
			fileName := fi.Name()
			if fi.IsDir() && isRecursion {
				fileList := GetFileList(path.Join(dbDir, fileName), suffix)
				if len(fileList) != 0 {
					arrFileLists = append(arrFileLists, fileList...)
				}
				continue
			}
			// try match suffix and `ordinal_pubKey_bitLength.suffix`
			if !strings.HasSuffix(strings.ToUpper(fileName), suffixUp) {
				continue
			}
			filePath := filepath.Join(dbDir, fileName)
			arrFileLists = append(arrFileLists, &FileList{
				FilePath: filePath,
				FileDir:  dbDir,
			})
		}
	}
	return arrFileLists
}
