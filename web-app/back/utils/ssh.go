package utils

import (
	"bytes"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func ConnectToServer() (*ssh.Client, error) {
	publicKey, err := publicKeyFile(os.Getenv("SSH_KEY_PATH"))
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("SSH_USER"),
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", os.Getenv("SSH_SERVER"), sshConfig)
	return client, err
}

func ExecuteCmd(conn *ssh.Client, command string) (bytes.Buffer, error) {
	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err := session.Run(command)

	return stdoutBuf, err
}

func publicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
