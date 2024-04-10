//go:build gc && go1.22.2
// +build gc,go1.22.2

package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"sync"
)

var (
	gidLock sync.Mutex
)

type stack struct {
	lo uintptr
	hi uintptr
}

type gobuf struct {
	sp   uintptr
	pc   uintptr
	g    uintptr
	ctxt uintptr
	ret  uintptr
	lr   uintptr
	bp   uintptr
}

type g struct {
	stack       stack
	stackguard0 uintptr
	stackguard1 uintptr

	_panic       uintptr
	_defer       uintptr
	m            uintptr
	sched        gobuf
	syscallsp    uintptr
	syscallpc    uintptr
	stktopsp     uintptr
	param        uintptr
	atomicstatus uint32
	stackLock    uint32
	goid         int64 // Here it is!
}

// getSlow parses the goroutine ID from runtime.Stack() output. It's slower but works.
func Get() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	log.Printf("Stack: %+v", buf[:n])
	return ExtractGID(buf[:n])
}

func ExtractGID(s []byte) int64 {
	s = s[len("goroutine "):]
	s = s[:bytes.IndexByte(s, ' ')]
	gid, err := strconv.ParseInt(string(s), 10, 64)
	if err != nil {
		fmt.Println("Error extracting GID:", err)
		return -1
	}
	return gid
}

// getSlow parses the goroutine ID from runtime.Stack() output. It's slower but works.
func getSlow() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	gid := ExtractGID(buf[:n])
	return gid
}

// Example of tracking goroutine lifecycle
var goroutineMap sync.Map

func TrackGoroutineStart() int64 {
	gid := Get()
	goroutineMap.Store(gid, "running")
	return gid
}

func TrackGoroutineEnd(gid int64) {
	goroutineMap.Delete(gid)
}

// Use TrackGoroutineStart and TrackGoroutineEnd to wrap goroutine execution for tracking

func main() {
	fmt.Println("Main goroutine ID:", Get())

	done := make(chan bool)
	go func() {
		gid := TrackGoroutineStart()
		defer TrackGoroutineEnd(gid)

		fmt.Println("Another goroutine ID:", Get())
		done <- true
	}()

	// Wait for the goroutine to finish to avoid premature exit
	<-done // Replaces the select{} with a mechanism to prevent deadlock
}
