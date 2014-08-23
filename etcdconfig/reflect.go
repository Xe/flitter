package etcdconfig

import (
	"errors"
	"reflect"
	"strings"

	"github.com/coreos/go-etcd/etcd"
)

// Demarshall takes an etcd client and a an anonymous interface to
// seed with values from etcd. This will throw an error if there is an exceptional
// error in the etcd client.
func Demarshall(etcd *etcd.Client, target interface{}) (err error) {
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

			for _, node := range resp.Node.Nodes {
				if node.Dir {
					continue
				}

				key := strings.TrimPrefix(node.Key, resp.Node.Key)
				key = strings.TrimPrefix(key, "/")
				value := node.Value

				valueField.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
			}
		}
	}

	return
}
