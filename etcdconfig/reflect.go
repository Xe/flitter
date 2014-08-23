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
