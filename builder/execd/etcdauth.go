package main

import (
	"bytes"
	"fmt"
	"log"

	"code.google.com/p/go.crypto/ssh"
	"github.com/Xe/flitter/lib/utils"
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

	fp := utils.GetFingerprint(keydata)

	user, allowed := utils.CanConnect(etcd, keydata)
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
