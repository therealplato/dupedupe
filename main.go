package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
)

var hashPath map[string]string
var problemPaths []string
var total int
var hasher hash.Hash

func main() {
	var root = "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	hashPath = make(map[string]string)
	problemPaths = make([]string, 0)
	hasher = sha256.New() // Save some cpu cycles by not constructing this per file

	err := filepath.Walk(root, perFileOrDir)
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

	modeBits := info.Mode()
	// Are any bits set for directory, link, pipe, socket, device, etc?
	masked := modeBits & os.ModeType
	if masked != 0 {
		// something was set, skip
		return nil
	}

	total += 1
	bb := sum(path)
	if len(bb) == 0 {
		return errors.New("failed to hash " + path)
	}

	// make a hexadecimal string of the bytes:
	s := fmt.Sprintf("%x", bb)

	// check for collision:
	existed, ok := hashPath[s]
	if ok {
		// dupe was already there
		fmt.Printf("%q duplicates %q with hash %q\n", path, existed, s)
	} else {
		// store key,value = sum,path
		hashPath[s] = path
	}
	return nil
}

// sum returns the sha256 sum of the input path
func sum(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer f.Close()

	hasher.Reset()
	if _, err := io.Copy(hasher, f); err != nil {
		fmt.Println(err)
		return nil
	}

	return hasher.Sum(nil) // the nil just means don't add any more data to the copied data
}
