package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

func execute(command string, args []string) ([]byte, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		fmt.Print(stderr.String())
	}

	return out.Bytes(), err
}

func initialize() *cli.App {
	var exclude string
	var include []string
	var git string

	app := cli.NewApp()
	app.Name = "git-delete-branches"
	app.Usage = "Delete local git branches"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "exclude, e",
			Value:       "",
			Usage:       "exclude deleting branches that match space separated string value. ie. \"feature/* release/*\"",
			Destination: &exclude,
		},
		cli.StringFlag{
			Name:        "git-alias, a",
			Value:       "git",
			Usage:       "alias for git command.",
			Destination: &git,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			include = []string{"*"}
		} else if c.NArg() == 1 {
			include = strings.Split(c.Args().Get(0), " ")
		} else {
			include = c.Args()
		}

		branches := filterBranches(getBranches(git), include, strings.Split(exclude, " "))

		if len(branches) > 0 {
			deleteBranches(branches, git)
		} else {
			fmt.Printf("no branches found with %v", strings.Join(include, ", "))
		}

		return nil
	}

	return app
}

func getBranches(alias string) []string {
	args := []string{"for-each-ref", "--format=%(refname:short)", "refs/heads"}
	out, _ := execute(alias, args)

	branches := strings.Split(string(out), "\n")

	// last element is an empty string
	return branches[:len(branches)-1]
}

func filterBranches(branches []string, includes []string, excludes []string) []string {
	filteredBranches := make([]string, 0)

	for _, branch := range branches {
		if matchBranch(branch, excludes) {
			continue
		}

		if matchBranch(branch, includes) {
			filteredBranches = append(filteredBranches, branch)
		}
	}

	return filteredBranches
}

func deleteBranches(branches []string, alias string) {
	args := []string{"branch", "-D"}
	args = append(args, branches...)

	out, _ := execute(alias, args)

	fmt.Printf("\n%s", string(out))
}

func matchBranch(branch string, matches []string) bool {
	replacer := strings.NewReplacer("*", "")

	for _, match := range matches {
		if len(match) == 0 {
			continue
		}

		if strings.Contains(match, "*") && strings.Contains(branch, replacer.Replace(match)) {
			return true
		}

		if !strings.Contains(match, "*") && branch == match {
			return true
		}
	}

	return false
}

func main() {
	app := initialize()
	err := app.Run(os.Args)

	if err != nil {
		fmt.Print(err)
	}
}
