// The autorebase tool automatically rebases all repositories from Conformal.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	path := os.Getenv("GOPATH")
	monetas := filepath.Join(path, "src", "github.com", "monetas")
	files, _ := ioutil.ReadDir(monetas)
	for _, file := range files {
		repo := filepath.Join(monetas, file.Name())
		os.Chdir(repo)

		_, err := exec.Command("git", "checkout", "master").Output()
		if err != nil {
			fmt.Println("Error during master checkout:", err)
			continue
		}
		_, err = exec.Command("git", "pull").Output()
		if err != nil {
			fmt.Println("Error during pull:", err)
			continue
		}

		// check to see if upstream exists
		out, err := exec.Command("git", "remote", "show", "upstream").Output()
		if err == nil {
			fmt.Println("------------")
			fmt.Println("Processing", repo)
		} else {
			up := fmt.Sprintf("git@github.com:/conformal/%v.git", file.Name())
			_, err = exec.Command("git", "remote", "add", "upstream", up).Output()
			if err != nil {
				// means this repo doesn't have a corresponding
				// one at conformal
				continue
			}
			fmt.Println("------------")
			fmt.Println("Processing", repo)
			fmt.Println("Adding upstream:", up)
		}

		// fetch and rebase
		_, err = exec.Command("git", "fetch", "upstream").Output()
		if err != nil {
			fmt.Println("Error during upstream fetch:", err)
			continue
		}
		out, err = exec.Command("git", "rebase", "upstream/master").Output()
		out_str := string(out)
		if strings.Contains(out_str, "CONFLICT") {
			fmt.Println(file.Name(), "needs merging")
			continue
		}
		if err != nil {
			fmt.Println("Error during upstream rebase:", err)
			continue
		}

		// check for any new imports
		out, err = exec.Command("grep", "-r", "-l", "github.com/conformal").Output()
		out_str = string(out)
		if strings.Contains(out_str, ".go") {
			fmt.Println(file.Name(), "has new imports")
			continue
		}

		// see if we actually did anything
		out, err = exec.Command("git", "diff").Output()
		if err != nil {
			fmt.Println("Error during diff:", err)
			continue
		}
		if len(out) == 0 {
			fmt.Println("No changes, not updating monetas repo", file.Name())
		} else {
			fmt.Println("Updating monetas repo", file.Name())
			exec.Command("git", "push", "origin", "+HEAD").Output()
		}
	}
}
