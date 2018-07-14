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
	ids := []string{"dog", "cat", "fish"}
	sl.BatchLock(ids)
	fmt.Println("Acquired batch lock")
	for _, id := range ids {
		sl.Unlock(id)
	}
	fmt.Println("Goodbye.")
}
