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
