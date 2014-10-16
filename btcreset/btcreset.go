// The btcreset tool resets all the monetas/btc* repositories
// so we don't have to "git reset --hard origin/master" over and over

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	path := os.Getenv("GOPATH")
	monetas := filepath.Join(path, "src", "github.com", "monetas")
	files, _ := ioutil.ReadDir(monetas)
	for _, file := range files {
		fmt.Printf("%v\n", file.Name())
		repo := filepath.Join(monetas, file.Name())
		os.Chdir(repo)

		_, err := exec.Command("git", "stash").Output()
		if err != nil {
			fmt.Println("Error during stash:", err)
			continue
		}

		_, err = exec.Command("git", "checkout", "master").Output()
		if err != nil {
			fmt.Println("Error during master checkout:", err)
			continue
		}

		// fetch and reset
		_, err = exec.Command("git", "fetch", "origin").Output()
		if err != nil {
			fmt.Println("Error during fetch:", err)
			continue
		}

		_, err = exec.Command("git", "reset", "--hard", "origin/master").Output()
		if err != nil {
			fmt.Println("Error during pull:", err)
			continue
		}
	}
}
