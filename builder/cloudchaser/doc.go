/*
Command cloudchaser is a minimal build sentry for Flitter. It will check permissions
based on environment variables and then return via exit code.

This is designed to be run from the context of a git push.

    2014/10/29 10:23:52 Usage:
      cloudchaser <revision> <sha>
*/
package main
