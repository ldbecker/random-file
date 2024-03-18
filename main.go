package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func listFiles(dir string, typesWanted map[string]bool) ([]string, error) {
	retFiles := make([]string, 0)
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return retFiles, fmt.Errorf("error listing dir '%v': %v", dir, err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.Type().IsDir() {
			subDir := path.Join(dir, dirEntry.Name())
			subEntries, err := listFiles(subDir, typesWanted)
			if err != nil {
				return []string{}, fmt.Errorf("error listing in dir '%v'", subDir)
			}
			retFiles = append(retFiles, subEntries...)
		} else if dirEntry.Type().IsRegular() {
			fileChunks := strings.Split(dirEntry.Name(), ".")
			filext := fileChunks[len(fileChunks)-1]
			if exists, ok := typesWanted[filext]; ok && exists {
				retFiles = append(retFiles, path.Join(dir, dirEntry.Name()))
			}
		}
	}
	return retFiles, nil
}

func main() {
	dirarg := os.Args[len(os.Args)-1]

	rr := rand.New(rand.NewSource(time.Now().Unix()))
	picTypes := make(map[string]bool)
	picTypes["jpg"], picTypes["jpeg"] = true, true

	datrstr := fmt.Sprintf("%v-%v-%v", time.Now().Month(), time.Now().Day(), time.Now().Year())
	fmt.Println(datrstr)

	fileList, err := listFiles(dirarg, picTypes)
	if err != nil {
		panic(err)
	}
	fc := len(fileList)
	fi := rr.Float64() * float64(fc)
	fn := fileList[int(fi)]

	exec.Command("open", "-a", "Preview", fn).Run()
}
