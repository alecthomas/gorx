package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/gorx"
)

func main() {
	gorx.
		FromTimeChannel(time.Tick(time.Second)).
		Take(5).
		Do(func(t time.Time) { fmt.Printf("%s\n", t) }).
		Wait()
}
