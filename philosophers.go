package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

const philosophers = 5
const meals = 3

var wg sync.WaitGroup

func main() {
    channelsLock := make([]chan bool, philosophers)
    channelsUnlock := make([]chan bool, philosophers)

    // Create forks
    for i := 0; i < philosophers; i++ {
        channelsLock[i] = make(chan bool)
        channelsUnlock[i] = make(chan bool)

        fork := Fork{i}
        go fork.listen(channelsLock[i], channelsUnlock[i])
    }

    // Create philosophers
    for i := 0; i < philosophers; i++ {
        wg.Add(1)
        philosopher := Philosopher{
            i,
            channelsLock[i],
            channelsUnlock[i],
            channelsLock[(i + 1) % philosophers],
            channelsUnlock[(i + 1) % philosophers],
        }
        go philosopher.dine()
    }

    wg.Wait()
    fmt.Println("Everyone has finished eating!")
}

// Fork logic
type Fork struct {
    id int
}

func (p *Fork) listen(lock chan bool, unlock chan bool) {
	// Fork will block the other incoming lock request, until the fork receives an unlock signal.
    for {
        <- lock
        <- unlock
    }
}

// Philosopher logic
type Philosopher struct {
    id int
    leftForkLock chan <- bool
    leftForkUnlock chan <- bool
    rightForkLock chan <- bool
    rightForkUnlock chan <- bool
}

func (p *Philosopher) dine() {
    for i := 0; i < meals; i++ {
		// Even philosophers will start by picking up the left fork, where the odd starts with the right.
		// This ensures that atleast one philosopher is always able to pickup the fork.
        if p.id % 2 == 0 {
            p.leftForkLock <- true
            p.rightForkLock <- true
        } else {
            p.rightForkLock <- true
            p.leftForkLock <- true
        }

        // Eating
        fmt.Println("Philosopher", p.id, "says: I am eating!", "(", i + 1, ")")
        time.Sleep(1 * time.Second)

        // Release the forks
        p.leftForkUnlock <- false
        p.rightForkUnlock <- false

        // Thinking
        fmt.Println("Philosopher", p.id, "says: I am thinking!")
        time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
    }
    defer wg.Done()
}