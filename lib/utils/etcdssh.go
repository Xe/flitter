package utils

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

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

	keybit := strings.Split(sshkey, " ")[1]
	fp := GetFingerprint(keybit)

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
func GetFingerprint(key string) string {
	h := md5.New()

	data, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "error"
	}

	log.Printf("key: %s", key)

	h.Write(data)

	ret := FPAddColons(fmt.Sprintf("%x", h.Sum(nil)))

	if key == ret {
		log.Fatal("Assertion failure")
	}

	return ret
}

// addColons adds colons every second character in a string to make the formatting
// of a public key fingerprint match the common standard:
//
//    dd3bb82e850406e9abffa80ac0046ed6
//
// becomes
//
//    dd:3b:b8:2e:85:04:06:e9:ab:ff:a8:0a:c0:04:6e:d6
func FPAddColons(s string) (r string) {
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
