//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type safeCounter struct {
	mu      sync.Mutex
	counter int
}

func catchAndStop(proc *MockProcess, c *safeCounter, sigs chan os.Signal) {
	for {
		_, ok := <-sigs
		if ok {
			c.mu.Lock()

			c.counter++
			if c.counter == 1 {
				proc.Stop()
			} else {
				os.Exit(0)
			}

			c.mu.Unlock()
		}
	}
}

func main() {
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT)

	// Create a process
	proc := MockProcess{}
	counter := safeCounter{}

	go catchAndStop(&proc, &counter, sigs)

	// Run the process (blocking)
	proc.Run()
}
