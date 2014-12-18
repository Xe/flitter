# Contributing to Flitter

## Code Style

Idiomatic Go is a must. If you fail, fail loudly and with details.

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
<action>/<component>/<dash-separated-description-of-changes>
<action>/<component>-<subcomponent>/<dash-separated-description-of-changes>
```

Where action is one of:

| Action  | Meaning                          |
|:------- |:-------------------------------- |
| `feat`  | New features                     |
| `fix`   | Fixes                            |
| `doc`   | Documentation additions or fixes |
| `nfo`   | Nuking from orbit                |
| `idiom` | Idiomatic code changes           |

Please do close issues via issues and commits where needed.
