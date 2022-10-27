package files

import (
	"archive/zip"
	"github.com/4ra1n/swing-rce-inspector/log"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func unzipJar(path string, id string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Error("error jar path: %s", path)
		os.Exit(-1)
	}
	r, err := zip.OpenReader(absPath)
	if r == nil {
		log.Error("cannot read file: %s", absPath)
		os.Exit(-1)
	}
	for _, f := range r.File {
		tempPath := filepath.Join("temp", id, f.Name)
		if strings.HasSuffix(f.Name, "/") {
			_ = os.MkdirAll(tempPath, 0644)
		} else {
			if !strings.HasSuffix(f.Name, ".class") {
				continue
			}
			tempSplits := strings.Split(f.Name, "/")
			tempDir := strings.Join(tempSplits[0:len(tempSplits)-1], "/")
			tempDirPath := filepath.Join("temp", id, tempDir)
			_ = os.MkdirAll(tempDirPath, 0644)
			reader, _ := f.Open()
			data, _ := io.ReadAll(reader)
			_ = os.WriteFile(tempPath, data, 0644)
		}
	}
}

func UnzipJars(dir string) {
	dirPath, err := filepath.Abs(dir)
	if err != nil {
		log.Error("error dir path: %s", dir)
		os.Exit(-1)
	}
	fileList, _ := os.ReadDir(dirPath)
	for i, f := range fileList {
		if strings.HasSuffix(f.Name(), ".jar") {
			finalPath := filepath.Join(dirPath, f.Name())
			id := strconv.Itoa(i) + "_" + f.Name()
			unzipJar(finalPath, id)
		}
	}
}
