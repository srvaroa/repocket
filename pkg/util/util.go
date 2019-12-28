package util

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type LogWriter struct{}

func EnsureDir(path string) {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err != nil {
			log.Fatalf("ERROR: expecting directory %s to exist", path)
		}
	}
	if !f.IsDir() {
		log.Fatalf("ERROR: expecting path %s to be a directory", path)
	}
}

func (writer LogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(string(bytes))
}

func DumpUrl(url string) ([]byte, error) {
	return exec.Command("w3m", "-dump", url).Output()
}

func EmptyDir(dirPath string) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatalf("Failed to list deletion queue: %s", err)
	}
	for _, f := range files {
		filePath := dirPath + "/" + f.Name()
		log.Printf("Deleting file %s\n", filePath)
		os.Remove(filePath)
	}
}
