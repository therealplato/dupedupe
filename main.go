package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var hashPath map[string]string
var problemPaths []string
var total int

func main() {
	hashPath = make(map[string]string)
	problemPaths = make([]string, 0)
	err := filepath.Walk(".", perFileOrDir)
	if err != nil {
		// aborted early
		log.Fatalf("stopping early due to issue: %q", err)
	}
	fmt.Printf("hashed %v files, finding %v unique hashes\n", total, len(hashPath))
	if len(problemPaths) > 0 {
		fmt.Printf("the following paths could not be accessed, maybe permission issue")
		for _, path := range problemPaths {
			fmt.Println(path)
		}
	}
}

func perFileOrDir(path string, info os.FileInfo, err error) error {
	if err != nil {
		// couldnt access this path
		problemPaths = append(problemPaths, path)
		return nil
	}
	if info.IsDir() {
		return nil
	}

	total += 1
	s := sum(path)
	if s == "" {
		return errors.New("failed to hash " + path)
	}
	existed, ok := hashPath[s]
	if ok {
		// dupe was already there
		fmt.Printf("%q duplicates %q with hash %q\n", path, existed, sum)
	} else {
		// store key,value = sum,path
		hashPath[s] = path
	}
	return nil
}

// sum returns a string representation of the sha256 sum of the input path
func sum(path string) string {
	return ""
}
