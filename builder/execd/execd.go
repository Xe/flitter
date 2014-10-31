/*
Command execd is a simple SSH server that allows a user to run single commands on a remote server,
suitable for things like git deploys. This is a fork of execd by progrium. This version
of execd is far different.

    Usage: ./execd [options] <exec-handler>

      -debug=false: debug mode displays handler output
      -env-pass=false: pass environment to handlers
      -etcd-node="http://127.0.0.1:4001": etcd node to connect to
      -key="": pem file of private keys (read from SSH_PRIVATE_KEYS by default)
      -port="22": port to listen on

It is not suggested you run this outside of flitter as-is unless you know what you are doing.

Changes:

    - Authentication over etcd
      - Keys go in `/flitter/builder/users/$USERNAME/$FINGERPRINT` -> b64 encoded key
    - Only allow git pushes to be made
    - Command line flag for etcd endpoint added
    - Remove `auth-handler`
*/
package main

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"code.google.com/p/go.crypto/ssh"
	"github.com/flynn/go-shlex"
)

var (
	port       = flag.String("port", "22", "port to listen on")
	debug      = flag.Bool("debug", false, "debug mode displays handler output")
	env        = flag.Bool("env-pass", false, "pass environment to handlers")
	keys       = flag.String("key", "", "pem file of private keys (read from SSH_PRIVATE_KEYS by default)")
	etcduplink = flag.String("etcd-node", "http://10.1.42.1:4001", "etcd node to connect to")
)

var ErrUnauthorized = errors.New("execd: user is unauthorized")

// Struct exitStatusMsg represents the exit status of a command that execd runs.
type exitStatusMsg struct {
	Status uint32
}

// exitStatus gets the exit status of a command and returns an exitStatusMsg and an error if
// applicable.
func exitStatus(err error) (exitStatusMsg, error) {
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// There is no platform independent way to retrieve
			// the exit code, but the following will work on Unix
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return exitStatusMsg{uint32(status.ExitStatus())}, nil
			}
		}
		return exitStatusMsg{0}, err
	}
	return exitStatusMsg{0}, nil
}

// attachCmd attaches a given exec.Cmd pointer, standard output, error and input writers and runs
// the command with the given piping set. It returns a waitgroup representing the status of the command
// and an error if the command fails.
func attachCmd(cmd *exec.Cmd, stdout, stderr io.Writer, stdin io.Reader) (*sync.WaitGroup, error) {
	var wg sync.WaitGroup
	wg.Add(2)

	log.Printf("Running %s...", cmd.Args)

	if stdin != nil {
		stdinIn, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		go func() {
			io.Copy(stdinIn, stdin)
			stdinIn.Close()
		}()
	}

	stdoutOut, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		io.Copy(stdout, stdoutOut)
		wg.Done()
	}()

	stderrOut, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	go func() {
		io.Copy(stderr, stderrOut)
		wg.Done()
	}()

	return &wg, nil
}

// addKey parses an SSH private key for execd. It takes in the SSH server configuration and the key to add.
// It returns an error if the key is unsupported by execd.
func addKey(conf *ssh.ServerConfig, block *pem.Block) (err error) {
	var key interface{}

	switch block.Type {
	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(block.Bytes)
	case "DSA PRIVATE KEY":
		key, err = ssh.ParseDSAPrivateKey(block.Bytes)
	default:
		return fmt.Errorf("unsupported key type %q", block.Type)
	}
	if err != nil {
		return err
	}

	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		return err
	}

	conf.AddHostKey(signer)

	return nil
}

// parseKeys makes SSH private key structs out of the raw input from the disk.
func parseKeys(conf *ssh.ServerConfig, pemData []byte) error {
	var found bool
	for {
		var block *pem.Block
		block, pemData = pem.Decode(pemData)

		if block == nil {
			if !found {
				return errors.New("no private keys found")
			}
			return nil
		}

		if err := addKey(conf, block); err != nil {
			return err
		}

		found = true
	}
}

// init initializes the flag.Usage function.
func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] <exec-handler>\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

// main is the entry point for execd.
func main() {
	flag.Parse()

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return handleAuth(conn, key)
		},
		AuthLogCallback: func(conn ssh.ConnMetadata, method string, err error) {
		},
	}

	if keyEnv := os.Getenv("SSH_PRIVATE_KEYS"); keyEnv != "" {
		if err := parseKeys(config, []byte(keyEnv)); err != nil {
			log.Fatalln("Failed to parse private keys:", err)
		}
	} else {
		pemBytes, err := ioutil.ReadFile(*keys)
		if err != nil {
			log.Fatalln("Failed to load private keys:", err)
		}
		if err := parseKeys(config, pemBytes); err != nil {
			log.Fatalln("Failed to parse private keys:", err)
		}
	}

	if p := os.Getenv("PORT"); p != "" && *port == "22" {
		*port = p
	}

	listener, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalln("Failed to listen for connections:", err)
	}

	fmt.Println(logo)
	log.Printf("execd is now listening on port %s", *port)

	for {
		// SSH connections just house multiplexed connections
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Failed to accept incoming connection:", err)
			continue
		}
		go handleConn(conn, config)
	}
}

// handleConn takes in the socket connection of the user and the SSH server configuration.
// It then routes input and output to the right places.
func handleConn(conn net.Conn, conf *ssh.ServerConfig) {
	defer conn.Close()

	sshConn, chans, reqs, err := ssh.NewServerConn(conn, conf)
	if err != nil {
		log.Println("Failed to handshake:", err)
		return
	}

	go ssh.DiscardRequests(reqs)

	for ch := range chans {
		if ch.ChannelType() != "session" {
			ch.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		go handleChannel(sshConn, ch)
	}
}

// handleChannel runs the needed commands on a given SSH server connection against the user
// that opened the communication channel.
func handleChannel(conn *ssh.ServerConn, newChan ssh.NewChannel) {
	ch, reqs, err := newChan.Accept()
	if err != nil {
		log.Println("newChan.Accept failed:", err)
		return
	}

	defer ch.Close()

	for req := range reqs {
		switch req.Type {
		case "exec":
			assert := func(at string, err error) bool {
				if err != nil {
					log.Printf("%s failed: %s", at, err)
					ch.Stderr().Write([]byte("Internal error.\n"))
					return true
				}
				return false
			}

			defer func() {
				log.Printf("Connection lost from %s", conn.RemoteAddr().String())
			}()

			if req.WantReply {
				req.Reply(true, nil)
			}

			cmdline := string(req.Payload[4:])

			cmdargs, err := shlex.Split(cmdline)
			if assert("shlex.Split", err) {
				return
			}

			if len(cmdargs) != 2 {
				ch.Stderr().Write([]byte("Invalid arguments.\n"))
				return
			}

			if cmdargs[0] != "git-receive-pack" {
				ch.Stderr().Write([]byte("Only `git push` is supported.\n"))
				return
			}

			user := conn.Permissions.Extensions["user"]
			reponame := strings.TrimSuffix(strings.TrimPrefix(cmdargs[1], "/"), ".git")

			log.Printf("Push from %s at %s", user, reponame)

			if err := makeGitRepo(reponame); err != nil {
				ch.Stderr().Write([]byte("Error: " + err.Error()))
				return
			}

			log.Printf("Writing hooks...")

			err = ioutil.WriteFile(reponame+"/hooks/pre-receive", []byte(`#!/bin/bash
strip_remote_prefix() {
	sed -u "s/^/"$'\e[1G'"/"
}

set -eo pipefail; while read oldrev newrev refname; do
	/app/cloudchaser -etcd-machine `+*etcduplink+` pre $newrev 2>&1 | strip_remote_prefix
done`), 0755)
			if err != nil {
				return
			}

			err = ioutil.WriteFile(reponame+"/hooks/post-receive", []byte(`#!/bin/bash

strip_remote_prefix() {
	sed -u "s/^/"$'\e[1G'"/"
}

set -eo pipefail; while read oldrev newrev refname; do
	/app/builder -etcd-host `+*etcduplink+` $REPO ${refname##*/} $newrev 2>&1 | strip_remote_prefix
done`), 0755)
			if err != nil {
				return
			}

			log.Printf("Doing git receive...")

			receive := exec.Command("git-receive-pack", reponame)
			if conn.Permissions.Extensions["environ"] != "" {
				receive.Env = append(receive.Env, strings.Split(conn.Permissions.Extensions["environ"], "\n")...)
			}

			receive.Env = append(receive.Env, "USER="+conn.Permissions.Extensions["user"])
			receive.Env = append(receive.Env, "REMOTE_HOST="+conn.RemoteAddr().String())
			receive.Env = append(receive.Env, "REPO="+reponame)

			done, err := attachCmd(receive, ch, ch.Stderr(), ch)
			if err != nil {
				ch.Stderr().Write([]byte("Error: " + err.Error()))
				return
			}

			if assert("receive.Start", receive.Start()) {
				return
			}

			done.Wait()

			log.Printf("Receive done")

			status, rcvErr := exitStatus(receive.Wait())
			if rcvErr != nil {
				ch.Stderr().Write([]byte("Error: " + rcvErr.Error()))
				return
			}

			_, err = ch.SendRequest("exit-status", false, ssh.Marshal(&status))
			assert("sendExit", err)

			return

		case "env":
			if req.WantReply {
				req.Reply(true, nil)
			}
		default:
			return
		}
	}
}
