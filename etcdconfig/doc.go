/*
Package etcdconfig provides a method for programs to scrape configuration from
etcd as well as continously update it. It is suggested that any program using
this package as a primary configuration source also have a "seed" option or
similar to ensure that etcd does not have any missing keys from the first
launch.

Currently etcdconfig supports strings, string->string maps and boolean variables from etcd. More
will be added as time goes on. Usage is very simple:

    type Config struct {
    	Foo string            `etcd:"/test/foo"`
    	Bar bool              `etcd:"/test/bar"`
    	Baz map[string]string `etcd:"/test/baz"`
    }

    Demarshal(etcd, someStruct)

A string will always have its value copied verbatim.

A boolean merely checks for the presence or absense of a key.

With maps, the directory tree will be recursively walked and all subdirectories and
their keys will be added as-is.
*/
package etcdconfig
