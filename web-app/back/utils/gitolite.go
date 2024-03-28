package utils

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

var (
	GitoliteAdminPath string = "$HOME/gitolite-admin"
	wg                sync.WaitGroup
)

func AddUserToRepo(username string, pubKey string, repo string, access string) error {
	session, err := ConnectToServer()

	if err != nil {
		return err
	}

	defer session.Close()
	errChan := make(chan error, 2)
	wg = sync.WaitGroup{}
	wg.Add(2)

	// Add user to conf dir
	go addLineToFile(fmt.Sprintf("%s/conf/gitolite.conf", GitoliteAdminPath), fmt.Sprintf("repo %s", repo), fmt.Sprintf("    %s         =   %s", access, username), errChan, session)
	// Add user to keydir
	go addPubKey(username, pubKey, errChan, session)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// Add, commit and push changes
	err = addCommitPush(fmt.Sprintf("Add user %s to repo %s", username, repo), session)

	return err
}

func addLineToFile(filename string, line string, lineToAdd string, errChan chan error, session *ssh.Session) {
	defer wg.Done()

	text, err := session.Output(fmt.Sprintf("cat %s", filename))
	if err != nil {
		errChan <- err
		return
	}

	if string(text) == "" {
		errChan <- fmt.Errorf("file %s is empty", filename)
		return
	}

	// create a temp file to store the current content
	tempFilename := "tmp"

	file, err := os.Create(tempFilename)
	if err != nil {
		errChan <- err
		return
	}

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(string(text))
	if err != nil {
		errChan <- err
		return
	}
	writer.Flush()

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
		errChan <- err
		return
	}
	file.Close()

	file, err = os.Create(tempFilename)
	if err != nil {
		errChan <- err
		return
	}

	writer = bufio.NewWriter(file)
	for _, l := range lines {
		_, err := writer.WriteString(l + "\n")
		if err != nil {
			errChan <- err
			return
		}
	}
	writer.Flush()

	scanner = bufio.NewScanner(file)

	var newContent string
	for scanner.Scan() {
		newContent += scanner.Text() + "\n"
	}

	_, err = session.Output(fmt.Sprintf("echo %s > %s", newContent, filename))
	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil
}

func addPubKey(username string, key string, errChan chan error, session *ssh.Session) {
	defer wg.Done()

	_, err := session.Output(fmt.Sprintf("echo %s > %s/keydir/%s.pub", key, GitoliteAdminPath, username))

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil
}

func addCommitPush(msg string, session *ssh.Session) error {
	_, err := session.Output(fmt.Sprintf("cd %s", GitoliteAdminPath))

	if err != nil {
		return err
	}

	_, err = session.Output("git add conf")
	if err != nil {
		return err
	}

	_, err = session.Output("git add keydir")
	if err != nil {
		return err
	}

	_, err = session.Output(fmt.Sprintf("git commit -m %s", msg))
	if err != nil {
		return err
	}

	_, err = session.Output("git push")
	if err != nil {
		return err
	}

	return nil
}
