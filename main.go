package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/is386/gcv/internal/git"
	"github.com/is386/gcv/internal/tui"
)

func main() {
	emailsFlag := flag.String("emails", "", "comma separated list of emails to pull git contributions for")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: [options] <project-dir>\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nPositional arguments:\n")
		fmt.Fprintf(os.Stderr, "  <project-dir>\t the directory where your projects live\n")
	}
	flag.Parse()

	dir := "."
	if len(flag.Args()) > 0 {
		dir = flag.Arg(0)
	}

	var emails []string
	if *emailsFlag == "" {
		email, err := git.GetGitEmail()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		emails = append(emails, email)
	} else {
		emails = strings.Split((*emailsFlag), ",")
	}

	projects, err := git.GetProjectsInDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	collector := git.NewGitCollector(emails)
	for _, project := range projects {
		err := collector.CollectGitContributions(filepath.Join(dir, project))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	if len(collector.ContributionsMap) == 0 {
		fmt.Fprintf(os.Stderr, "no contributions found for %v in %s\n", emails, dir)
		os.Exit(1)
	}

	m := tui.NewModel(collector.ContributionsMap)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
