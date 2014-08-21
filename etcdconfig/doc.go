/*
Package etcdconfig provides a method for programs to scrape configuration from
etcd as well as continously update it. It is suggested that any program using
this package as a primary configuration source also have a "seed" option or
similar to ensure that etcd does not have any missing keys from the first
launch.

Currently etcdconfig supports string and boolean variables from etcd. More
will be added as time goes on. Usage is very simple:

    type Config struct {
    	Foo string `etcd:"/test/foo"`
    }

A string will always have its value copied verbatim but a boolean merely checks
for the presence or absense of a key.
*/
package etcdconfig
