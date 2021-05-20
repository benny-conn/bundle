package uploader

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bennycio/bundle/cli/term"
	"github.com/schollz/progressbar/v3"
)

type filewalker struct {
	target    string
	targetDir string
}

type walker interface {
	Walk() (*os.File, error)
}

func NewFileWalker(target string, targetDir string) walker {
	return &filewalker{target: target, targetDir: targetDir}
}

func (f *filewalker) Walk() (*os.File, error) {

	matches := make([]string, 0, 10)

	cur, err := os.Stat(f.target)
	if err == nil {
		matches = append(matches, cur.Name())
	}

	progress := progressbar.Default(-1, "Scanning for "+f.target+"...")

	err = filepath.Walk("../..", func(path string, info os.FileInfo, err error) error {
		progress.Add(1)
		if strings.Contains(strings.ToLower(path), strings.ToLower(f.targetDir)) {
			if info.Name() == f.target {
				matches = append(matches, path)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	progress.Finish()

	respath, err := f.matchesPrompt(matches)
	if err != nil {
		return nil, err
	}

	if respath == "" {
		return nil, errors.New("no readme to upload")
	}

	result, err := os.Open(respath)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (f *filewalker) matchesPrompt(matches []string) (string, error) {
	term.Print(fmt.Sprintf("Which %s file would you like to upload?\nType \"0\" for none and type a number to select the file that you wish to upload.\n", f.target))
	for i, v := range matches {
		fmt.Printf("%d. %s\n", i+1, v)
	}

	fmt.Print("\n")
	term.Print(fmt.Sprintf("Which %s file would you like to upload?\nType \"0\" for none and type a number to select the file that you wish to upload.\n", f.target))

	rd := bufio.NewReader(os.Stdin)
	matchString, err := rd.ReadString(byte('\n'))
	if err != nil {
		return "", err
	}
	matchString = strings.TrimSpace(strings.Trim(matchString, "\n"))

	match, err := strconv.Atoi(matchString)

	if err != nil {
		return "", err
	}

	if match < 1 {
		return "", nil
	}

	if len(matches) < match {
		return "", errors.New("invalid location")
	}
	return matches[match-1], nil
}
