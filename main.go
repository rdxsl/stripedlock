package main

import "fmt"

func main() {
	fmt.Println("Hello")

	sl := NewStripedLock(256)

	// testing and hoping that we don't collide our hashcodes!
	sl.Lock("dog")
	fmt.Println("Acquired dog lock")
	sl.Lock("cat")
	fmt.Println("Acquired cat lock")
	sl.Unlock("dog")
	sl.Unlock("cat")
	fmt.Println("Goodbye.")
}
