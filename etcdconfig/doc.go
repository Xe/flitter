/*
Package etcdconfig provides a method for programs to scrape configuration from
etcd as well as continously update it. It is suggested that any program using
this package as a primary configuration source also have a "seed" option or
similar to ensure that etcd does not have any missing keys from the first
launch.
*/
package etcdconfig
