package git

import (
	"errors"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"
)

type GitCollector struct {
	ContributionsMap map[int]*[366]int
	emails           []string
}

func NewGitCollector(emails []string) GitCollector {
	return GitCollector{ContributionsMap: make(map[int]*[366]int), emails: emails}
}

func (collector *GitCollector) CollectGitContributions(cwd string) error {
	log, err := exec.Command("git", "-C", cwd, "log", "--all", "--pretty=format:%ae %at").Output()
	if err != nil {
		return err
	}

	commits := strings.SplitSeq(string(log), "\n")

	for commit := range commits {
		commitParts := strings.Split(commit, " ")
		if len(commitParts) < 2 {
			return errors.New("malformed commit data")
		}

		if !slices.Contains(collector.emails, commitParts[0]) {
			continue
		}

		unixTimestamp, err := strconv.Atoi(commitParts[1])
		if err != nil {
			return err
		}

		datetime := time.Unix(int64(unixTimestamp), 0)
		year := datetime.Year()
		if collector.ContributionsMap[year] == nil {
			collector.ContributionsMap[year] = &[366]int{}
		}
		collector.ContributionsMap[year][datetime.YearDay()]++
	}

	return nil
}
