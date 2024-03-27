package utils

import (
	"os"

	"golang.org/x/crypto/ssh"
)

func ConnectToServer() (*ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: os.Getenv("SSH_USER"),
		Auth: []ssh.AuthMethod{
			publicKeyFile(os.Getenv("SSH_KEY_PATH")),
		},
	}

	client, err := ssh.Dial("tcp", os.Getenv("SSH_SERVER"), sshConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	// TODO: pipe session output to stdout

	return session, nil
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}
