package etcdconfig

import (
	"log"
	"testing"

	"github.com/coreos/go-etcd/etcd"
)

type ValidTestConfig struct {
	DefinedString string            `etcd:"/test/definedstring"`
	DefinedBool   bool              `etcd:"/test/definedbool"`
	UndefinedBool bool              `etcd:"/test/undefinedbool"`
	TestMap       map[string]string `etcd:"/test/map"`
}

// Test basic confguration scraping from etcd
func TestBasicConfigScraping(t *testing.T) {
	cfg := &ValidTestConfig{
		TestMap: make(map[string]string),
	}
	etcd := etcd.NewClient([]string{"http://127.0.0.1:4001"})

	etcd.CreateDir("/test", 0)
	etcd.Create("/test/definedstring", "bar", 0)
	etcd.Create("/test/definedbool", "this will be ignored", 0)

	// Test map parsing
	etcd.CreateDir("/test/map", 0)
	etcd.Create("/test/map/foo", "foo", 0)
	etcd.Create("/test/map/bar", "bar", 0)

	// Test map subdirectory parsing
	etcd.CreateDir("/test/map/spam", 0)
	etcd.Create("/test/map/spam/eggs", "sausage", 0)

	err := Demarshal(etcd, cfg)
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
