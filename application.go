/* ThreadedEchoServer
 */
package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	go handleConnections()
	for {
		time.Sleep(40)
		tickGameworld()
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
