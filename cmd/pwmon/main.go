package main

import (
	"fmt"

	"github.com/andrieee44/pwmon/pkg"
)

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var (
		infoChan <-chan *pwmon.Info
		errChan  <-chan error
		info     *pwmon.Info
		err      error
	)

	infoChan, errChan, err = pwmon.Monitor()
	panicIf(err)

	for {
		select {
		case info = <-infoChan:
			fmt.Printf("Volume: %d%%, Mute: %t\n", info.Volume, info.Mute)
		case err = <-errChan:
			panic(err)
		}
	}
}
