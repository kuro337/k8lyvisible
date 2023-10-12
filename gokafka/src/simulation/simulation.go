package simulation

import (
	"fmt"
	"math/rand"
	"time"
)

func SimulateOperations() {
	for i := 0; i < 20; i++ {
		numObjects := rand.Intn(100)

		doSomethingWithObjects(numObjects)

		time.Sleep(time.Second)
	}
}

func doSomethingWithObjects(numObjects int) {
	objects := make([]int, numObjects)

	for i := 0; i < len(objects); i++ {
		objects[i] = i * i
	}

	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	fmt.Printf("Processed %d objects\n", len(objects))
}
