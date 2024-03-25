package utils

import (
	"bufio"
	"fmt"
	"os"
)

var (
	//GitoliteAdminPath string = "$HOME/gitolite-admin"
	GitoliteAdminPath string = "./utils"
)

func AddUserToRepo(username string, repo string, access string) error {
	errors := make(chan error)

	go addLineToFile(fmt.Sprintf("%s/conf/gitolite.conf", GitoliteAdminPath), fmt.Sprintf("repo %s", repo), fmt.Sprintf("    %s         =   %s", access, username), errors)

	return <-errors
}

func addLineToFile(filename string, line string, lineToAdd string, errors chan error) {
	file, err := os.Open(filename)
	if err != nil {
		errors <- err
	}

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		txt := scanner.Text()
		lines = append(lines, txt)
		if txt == line {
			lines = append(lines, lineToAdd)
		}
	}
	if err = scanner.Err(); err != nil {
		errors <- err
	}
	file.Close()

	file, err = os.Create(filename)
	if err != nil {
		errors <- err
	}

	writer := bufio.NewWriter(file)
	for _, l := range lines {
		_, err := writer.WriteString(l + "\n")
		if err != nil {
			errors <- err
		}
	}
	writer.Flush()
	errors <- nil
}
