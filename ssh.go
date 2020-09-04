package utils

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//utility for manage plugins ssh connection

//create a ssh client
func Getsshclient(host string, user string, key string) (*ssh.Client, err) {
	sshconf, err := getsshconfig(user, key)
	if err != nil {
		return nil, err
	}

	connection, err := ssh.Dial("tcp", host+":22", sshconf)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

//create a ssh client configuration with kay auth
func getsshconfig(user string, keyfile string) (*ssh.ClientConfig, err) {
	buffer, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}
	return sshConfig, nil
}

//run command on provided client and return the stdout
func Sshcmd(con *ssh.Client, cmd string) string {
	session, err := con.NewSession()
	defer session.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}
	out, err := session.Output(cmd)
	if err != nil {
		return ""
	}
	outs := string(out)
	return outs[0 : len(outs)-1]
}

func GetClientSFTP(con *ssh.Client) (*sftp.Client, err) {
	client, err := sftp.NewClient(con)
	if err != nil {
		return ""
	}
	return client, nil
}

func GetFileSFTP(con *sftp.Client, src string, dst string, perm os.FileMode) error {
	sfile, err := con.Open(src)
	if err != nil {
		return err
	}
	var bfile []byte
	_, err := sfile.Read(bfile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, bfile, perm)
	if err != nil {
		return err
	}
	return nil
}

func PutFileSFTP(con *sftp.Client, src string, dst string, perm os.FileMode) error {
	bfile, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	var sfile sftp.File
	sfile, err = con.Create(dst)
	if err != nil {
		sfile, err = con.Open(dst)
		if err != nil {
			return err
		}
	}
	err = sfile.Chmod(perm)
	if err != nil {
		return err
	}
	_, err := sfile.Write()
	if err != nil {
		return err
	}
	return nil
}
