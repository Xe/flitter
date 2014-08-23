package etcdconfig

import (
	"errors"
	"reflect"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

// Demarshal takes an etcd client and a an anonymous interface to
// seed with values from etcd. This will throw an error if there is an exceptional
// error in the etcd client or you are invoking this incorrectly with maps.
// Any missing keys in etcd will be filled in with blank strings.
func Demarshal(etcd *etcd.Client, target interface{}) (err error) {
	val := reflect.ValueOf(target).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		switch valueField.Kind() {
		case reflect.Bool:
			if _, notok := etcd.Get(tag.Get("etcd"), false, false); notok == nil {
				valueField.SetBool(true)
			} else {
				valueField.SetBool(false)
			}

		case reflect.String:
			etcdval, err := etcd.Get(tag.Get("etcd"), false, false)
			if err != nil {
				valueField.SetString("")
			}

			valueField.Set(reflect.ValueOf(etcdval.Node.Value))

		case reflect.Map:
			keyKind := typeField.Type.Key().Kind()

			if keyKind != reflect.String {
				return errors.New("Map must be string[string]")
			}

			resp, err := etcd.Get(tag.Get("etcd"), true, true)
			if err != nil {
				return err
			}

			if !resp.Node.Dir {
				return errors.New("maps must be pointed at an etcd directory")
			}

			SetMapOffDir(resp.Node, &valueField, "")
		}
	}

	return
}

// Subscribe continuously updates the target structure with data from etcd as it
// is changed. This is ideal for things such as configuration or a list of keys for
// authentication.
func Subscribe(e *etcd.Client, target interface{}, prefix string) (err error) {
	updates := make(chan *etcd.Response, 10)
	e.Watch(prefix, 0, true, updates, make(chan bool))

	go func() {
		for update := range updates {
			_ = update // TODO: replace me with a better method

			Demarshal(e, target)
		}
	}()

	return
}

// SetMapOffDir sets a map based off of an etcd directory and its contents.
func SetMapOffDir(parent *etcd.Node, target *reflect.Value, tack string) {
	for _, node := range parent.Nodes {
		key := strings.TrimPrefix(node.Key, parent.Key)
		key = tack + strings.TrimPrefix(key, "/")
		value := node.Value

		if node.Dir {
			SetMapOffDir(node, target, key+"/")
		} else {
			target.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
		}
	}
}
