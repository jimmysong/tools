/*

Autorebase is a tool to automatically rebase all repositories from Conformal

*/

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
		// check that upstream is actually defined
		exec.Command("git", "checkout", "master").Output()
		exec.Command("git", "pull").Output()
		up := fmt.Sprintf("git@github.com:/conformal/%v.git", file.Name())
		fmt.Println("upstream", up)
		exec.Command("git", "remote", "add", "upstream", up).Output()
		exec.Command("git", "fetch", "upstream").Output()

		out, _ := exec.Command("git", "rebase", "upstream/master").Output()
		s := string(out)
		if strings.Contains(s, "CONFLICT") {
			fmt.Println(file.Name(), "needs merging")
		} else {
			fmt.Println("updating monetas repo", file.Name())
			exec.Command("git", "push", "origin", "HEAD", "-f")
		}
	}
}
