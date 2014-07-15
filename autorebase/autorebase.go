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
		fmt.Println("Processing", repo)
		os.Chdir(repo)
		// TODO: check that upstream is actually defined.
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
		up := fmt.Sprintf("git@github.com:/conformal/%v.git", file.Name())
		fmt.Println("Upstream:", up)
		// Do not check for errors: the upstream remote might already exist.
		exec.Command("git", "remote", "add", "upstream", up).Output()
		_, err = exec.Command("git", "fetch", "upstream").Output()
		if err != nil {
			fmt.Println("Error during upstream fetch:", err)
			continue
		}
		out, err := exec.Command("git", "rebase", "upstream/master").Output()
		out_str := string(out)
		if strings.Contains(out_str, "CONFLICT") {
			fmt.Println(file.Name(), "needs merging")
			continue
		}
		if err != nil {
			fmt.Println("Error during upstream rebase:", err)
			continue
		}
		out, err = exec.Command("git", "diff").Output()
		if err != nil {
			fmt.Println("Error during diff:", err)
			continue
		}
		if out == nil {
			fmt.Println("No changes, not updating monetas repo", file.Name())
		} else {
			fmt.Println("Updating monetas repo", file.Name())
			exec.Command("git", "push", "origin", "+HEAD").Output()
		}
	}
}
