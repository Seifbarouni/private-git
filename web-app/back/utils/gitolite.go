package utils

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

var (
	GitoliteAdminPath string = "$HOME/gitolite-admin"
	wg                sync.WaitGroup
)

func AddUserToRepo(username string, pubKey string, repo string, access string) error {
	client, err := ConnectToServer()

	if err != nil {
		return err
	}

	defer client.Close()
	errChan := make(chan error, 2)
	wg = sync.WaitGroup{}
	wg.Add(2)

	// Add user to conf dir
	go addLineToFile(fmt.Sprintf("%s/conf/gitolite.conf", GitoliteAdminPath), fmt.Sprintf("repo %s", repo), fmt.Sprintf("            %s     =   %s", access, username), errChan, client)
	// Add user to keydir
	decodedKey, err := base64.StdEncoding.DecodeString(pubKey)
	if err != nil {
		return err
	}
	go addPubKey(username, string(decodedKey), errChan, client)

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// Add, commit and push changes
	err = addCommitPush(fmt.Sprintf("Add user %s to repo %s", username, repo), client)

	return err
}

func RemoveUserFromRepo(username string, repo string) error {
	client, err := ConnectToServer()

	if err != nil {
		return err
	}

	defer client.Close()

	err = removeLineFromFile(fmt.Sprintf("%s/conf/gitolite.conf", GitoliteAdminPath), repo, username, client)

	if err != nil {
		return err
	}

	err = addCommitPush(fmt.Sprintf("Remove user %s from repo %s", username, repo), client)

	return err
}

func addLineToFile(filename string, line string, lineToAdd string, errChan chan error, client *ssh.Client) {
	defer wg.Done()

	output, err := ExecuteCmd(client, fmt.Sprintf("cat %s", filename))

	if err != nil {
		errChan <- err
		return
	}

	tr := strings.TrimSpace(output.String())

	if tr == "" {
		errChan <- fmt.Errorf("file %s is empty", filename)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(tr))
	lines := []string{}
	f := false
	for scanner.Scan() {
		txt := scanner.Text()
		lines = append(lines, txt)
		if txt == line {
			lines = append(lines, lineToAdd)
			f = true
		}
	}
	if err = scanner.Err(); err != nil {
		errChan <- err
		return
	}

	if !f {
		lines = append(lines, line)
		lines = append(lines, lineToAdd)
	}

	newContent := strings.Join(lines, "\n")

	_, err = ExecuteCmd(client, fmt.Sprintf("echo '%s' > %s", newContent, filename))

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil
}

func addPubKey(username string, key string, errChan chan error, client *ssh.Client) {
	defer wg.Done()
	tk := strings.TrimSpace(key)
	_, err := ExecuteCmd(client, fmt.Sprintf("echo '%s' > %s/keydir/%s.pub", tk, GitoliteAdminPath, username))

	if err != nil {
		errChan <- err
		return
	}

	errChan <- nil
}

func addCommitPush(msg string, client *ssh.Client) error {
	if _, err := ExecuteCmd(client, fmt.Sprintf("cd %s && git add conf", GitoliteAdminPath)); err != nil {
		return err
	}

	if _, err := ExecuteCmd(client, fmt.Sprintf("cd %s && git add keydir", GitoliteAdminPath)); err != nil {
		return err
	}

	if _, err := ExecuteCmd(client, fmt.Sprintf("cd %s && git commit -m '%s'", GitoliteAdminPath, msg)); err != nil {
		return err
	}

	if _, err := ExecuteCmd(client, fmt.Sprintf("cd %s && git push", GitoliteAdminPath)); err != nil {
		return err
	}

	return nil
}

func removeLineFromFile(filename string, repo string, username string, client *ssh.Client) error {
	output, err := ExecuteCmd(client, fmt.Sprintf("cat %s", filename))
	if err != nil {
		return err
	}

	tr := strings.TrimSpace(output.String())

	if tr == "" {
		return fmt.Errorf("%s is empty", filename)
	}

	scanner := bufio.NewScanner(strings.NewReader(tr))
	lines := []string{}
	foundRepo := false
	for scanner.Scan() {
		txt := scanner.Text()

		if strings.Contains(txt, repo) {
			foundRepo = true
			lines = append(lines, txt)
			continue
		}

		if foundRepo && strings.Contains(txt, username) {
			foundRepo = false
			continue
		}
		lines = append(lines, txt)
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	newContent := strings.Join(lines, "\n")

	_, err = ExecuteCmd(client, fmt.Sprintf("echo '%s' > %s", newContent, filename))

	return err
}
