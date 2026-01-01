//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"sync"
	"time"
)

//func runWithDone(process func(), done chan struct{}) {
//	process()
//	defer close(done)
//}
//
//// HandleRequest runs the processes requested by users. Returns false
//// if process had to be killed
//func HandleRequest(process func(), u *User) bool {
//	if u.IsPremium {
//		process()
//		return true
//	}
//
//	var (
//		done  = make(chan struct{})
//		timer = time.NewTimer(time.Second * 10)
//	)
//
//	go runWithDone(process, done)
//
//	select {
//	case <-timer.C:
//		return false
//	case <-done:
//		return true
//	}
//}

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	mu        sync.Mutex
	TimeUsed  time.Duration // in seconds
}

const freeLimit = 10 * time.Second

func runWithDone(process func(), done chan struct{}) {
	process()
	defer close(done)
}

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if u.IsPremium {
		process()
		return true
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	remaining := freeLimit - time.Second*u.TimeUsed
	if remaining <= 0 {
		return false
	}

	var (
		done  = make(chan struct{})
		timer = time.NewTimer(remaining)
		start = time.Now()
	)

	go runWithDone(process, done)

	select {
	case <-timer.C:
		u.TimeUsed = freeLimit
		return false
	case <-done:
		u.TimeUsed += time.Since(start)
		return true
	}
}

func main() {
	RunMockServer()
}
