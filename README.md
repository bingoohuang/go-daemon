# go-daemon

Library for writing system daemons in Go. fork from [go-daemon](https://github.com/sevlyar/go-daemon)

Now supported only UNIX-based OS (Windows is not supported). But the library was tested only on Linux
and OSX, so that if you have an ability to test the library on other platforms, give me feedback, please (#26).

*Please, feel free to send me bug reports and fixes. Many thanks to all contributors.*

## Features

* Goroutine-safe daemonization;
* Out of box work with pid-files;
* Easy handling of system signals;
* The control of a daemon.

## Installation

	go get github.com/bingoohuang/godaemon@latest

## Examples

```go
package main

import (
	"github.com/bingoohuang/godaemon"
	flag "github.com/spf13/pflag"
)

func main() {
	daemon := flag.BoolP("daemon", "d", false, "run in daemonize mode")
	flag.Parse()
	godaemon.Daemonize(godaemon.WithDaemon(*daemon))

	// logic ...
}

```

* [Simple](cmd/gd-simple/)
* [Log rotation](cmd/gd-log-rotation/)
* [Signal handling](cmd/gd-signal-handling/)

```go
package main

import (
	"log"
	"github.com/bingoohuang/godaemon"
)

func main() {
	// Pre()

	context := new(godaemon.Context)
	child, _ := context.Reborn()

	if child != nil {
		// PostParent()
	} else {
		defer func() {
			if err := context.Release(); err != nil {
				log.Printf("Unable to release pid-file: %s", err.Error())
			}
		}()

		// PostChild()
	}
}
```

## How it works

We can not use `fork` syscall in Golang's runtime, because child process doesn't inherit
threads and goroutines in that case. The library uses a simple trick: it runs its own copy with
a mark - a predefined environment variable. Availability of the variable for the process means
an execution in the child's copy. So that if the mark is not set - the library executes
parent's operations and runs its own copy with mark, and if the mark is set - the library
executes child's operations:

![](img/idea.png)

## Sync logs

- Commit: 9d78c85, Commits on Mar 21, 2023，2023年04月24日
- Commit: f14d855, Commits on Feb 12, 2022，2022年09月07日
