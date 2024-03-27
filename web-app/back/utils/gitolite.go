package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
	// stfu compiler
	_ = session

	defer wg.Done()
	file, err := os.Open(filename)
	if err != nil {
		errChan <- err
		return
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
		errChan <- err
		return
	}
	file.Close()

	file, err = os.Create(filename)
	if err != nil {
		errChan <- err
		return
	}

	writer := bufio.NewWriter(file)
	for _, l := range lines {
		_, err := writer.WriteString(l + "\n")
		if err != nil {
			errChan <- err
			return
		}
	}
	writer.Flush()
	errChan <- nil
}

func addPubKey(username string, key string, errChan chan error, session *ssh.Session) {
	// stfu compiler
	_ = session

	defer wg.Done()
	file, err := os.Create(fmt.Sprintf("%s/keydir/%s.pub", GitoliteAdminPath, username))
	if err != nil {
		errChan <- err
		return
	}

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(key)
	if err != nil {
		errChan <- err
		return
	}
	writer.Flush()
	errChan <- nil
}

func addCommitPush(msg string, session *ssh.Session) error {
	// stfu compiler
	_ = session

	// TODO: refactor this to use the ssh session to run the commands

	cmd := exec.Command("git", "--git-dir", GitoliteAdminPath, "--work-tree", GitoliteAdminPath, "add", fmt.Sprintf("%s/conf", GitoliteAdminPath))
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "--git-dir", GitoliteAdminPath, "--work-tree", GitoliteAdminPath, "add", fmt.Sprintf("%s/keydir", GitoliteAdminPath))
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "--git-dir", GitoliteAdminPath, "--work-tree", GitoliteAdminPath, "commit", "-m", msg)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("git", "--git-dir", GitoliteAdminPath, "--work-tree", GitoliteAdminPath, "push")
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
