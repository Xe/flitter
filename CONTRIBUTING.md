# Contributing to Flitter

## Code Style

Idiomatic Go is a must. If you fail, fail loudly and with details. Avoid 
introducing dependencies where the maintainers do not strictly follow 
[semantic versioning](http://semver.org), and if you must please inform the 
Flitter authors so that a fork can be made.

Do not introduce pull requests for supporting tools such as Goop or Godep, they 
will be denied. `go get`-compatible code is a must.

## Documentation

All functions, public and private must be documented.

```go
// Command foobang processes the frozboz into flapnars for Flitter's asdf to
// function.
package main

import (
    "fmt"
)

// main is the entry point for Flitter's foobang
func main() {
    fmt.Println("Hi world!")
}
```

## Branching

Flitter uses [Github](http://github.com) for its version control and ticketing 
system. Take advantage of this.

When branching to fix a problem, always create a branch in one of the following 
formats:

```
<component>/<dash-separated-description-of-changes>
<component>/<subcomponent>/<dash-separated-description-of-changes>
```

Please do close issues via issues and commits where needed.
