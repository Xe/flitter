package etcdconfig

import (
	"log"
	"testing"

	"github.com/coreos/go-etcd/etcd"
)

type ValidTestConfig struct {
	DefinedString string `etcd:"/test/definedstring"`
	DefinedBool   bool   `etcd:"/test/definedbool"`
	UndefinedBool bool   `etcd:"/test/undefinedbool"`
}

// Test basic confguration scraping from etcd
func TestBasicConfigScraping(t *testing.T) {
	cfg := &ValidTestConfig{}
	etcd := etcd.NewClient([]string{"http://127.0.0.1:4001"})

	etcd.CreateDir("/test", 0)
	etcd.Create("/test/definedstring", "bar", 0)
	etcd.Create("/test/definedbool", "this will be ignored", 0)

	err := Demarshall(etcd, cfg)
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("%v", cfg)

	if cfg.DefinedString != "bar" {
		t.Fatalf("DefinedString is %v, expected \"bar\"", cfg.DefinedString)
	}

	if !cfg.DefinedBool {
		t.Fatal("DefinedBool should be true, it is false.")
	}

	if cfg.UndefinedBool {
		t.Fatal("UndefinedBool should be false, it is true.")
	}
}
