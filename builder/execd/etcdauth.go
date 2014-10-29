package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"

	"code.google.com/p/go.crypto/ssh"
	"github.com/coreos/go-etcd/etcd"
)

// handleAuth checks authentication against etcd using CanConnect and sets the needed
// environment variables for later parts of the builder to use. It takes in the SSH
// connection metadata, the public key of the user, and returns the SSH
// permissions of the connection and an error if they are not authorized.
func handleAuth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	if conn.User() != "git" {
		return nil, ErrUnauthorized
	}

	keydata := string(bytes.TrimSpace(ssh.MarshalAuthorizedKey(key)))

	etcd := etcd.NewClient([]string{*etcduplink})

	fp := getFingerprint(keydata)

	user, allowed := CanConnect(etcd, keydata)
	if allowed {
		log.Printf("User %s (%s) accepted with fingerprint %s", user, conn.RemoteAddr().String(), fp)
		return &ssh.Permissions{
			Extensions: map[string]string{
				"environ":     fmt.Sprintf("USER=%s\nKEY='%s'\nFINGERPRINT=%s\n", user, keydata, fp),
				"user":        user,
				"fingerprint": fp,
			},
		}, nil
	} else {
		log.Printf("Connection from %s rejected (bad key)", conn.RemoteAddr().String())
	}

	return nil, ErrUnauthorized
}

// CanConnect checks a given SSH key against the list of authorized users in etcd
// to check if a user with a given key is allowed to connect. It takes in an active
// etcd.Client struct pointer and the ssh key to test and returns a username
// and boolean representing if they are allowed to connect.
func CanConnect(e *etcd.Client, sshkey string) (user string, allowed bool) {
	reply, err := e.Get("/flitter/builder/users/", true, true)

	if err != nil {
		log.Printf("etcd: %s", err)
		return "", false
	}

	fp := getFingerprint(sshkey)

	for _, userdir := range reply.Node.Nodes {
		for _, fpnode := range userdir.Nodes {
			thisFpSplit := strings.Split(fpnode.Key, "/")
			thisFp := thisFpSplit[len(thisFpSplit)-1]

			if fp == thisFp {
				userpath := strings.Split(userdir.Key, "/")
				user := userpath[len(userpath)-1]
				return user, true
			}
		}
	}

	return
}

// getFingerprint takes an SSH key in and returns the fingerprint (an MD5 sum)
// and adds the needed colons to it.
func getFingerprint(key string) string {
	key = strings.Split(key, " ")[1]
	h := md5.New()

	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return key
	}

	io.WriteString(h, string(data))
	return addColons(fmt.Sprintf("%x", h.Sum(nil)))
}

// addColons adds colons every second character in a string to make the formatting
// of a public key fingerprint match the common standard:
//
//    dd3bb82e850406e9abffa80ac0046ed6
//
// becomes
//
//    dd:3b:b8:2e:85:04:06:e9:ab:ff:a8:0a:c0:04:6e:d6
func addColons(s string) (r string) {
	if len(s) == 0 {
		return ""
	}

	for i, c := range s {
		r = r + string(c)
		if i%2 == 1 && i != len(s)-1 { // Even number, add colon
			r = r + ":"
		}
	}

	return
}
