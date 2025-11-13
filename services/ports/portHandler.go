package services

import (
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/bitmattz/nira_the_sniffer/models"
)

const (
	DefaultTimeout     = 500 * time.Millisecond
	DefaultConcurrency = 200
	MaxPort            = 65535
)

func ScanPort(protocol, hostname string, port int) models.PortScan {
	result := models.PortScan{
		Port: port,
	}
	address := net.JoinHostPort(hostname, strconv.Itoa(port))
	conn, err := net.DialTimeout(protocol, address, DefaultTimeout)
	log.Printf("Scan started for port %d", port)
	if err != nil {
		if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
			result.State = "filtered"
		} else {
			result.State = "closed"
		}
		log.Printf("Port %d -> %s", port, result.State)
		return result
	}
	defer conn.Close()
	result.State = "open"
	log.Printf("Port %d -> %s", port, result.State)
	return result
}

func ScanPorts(hostname string) []models.PortScan {
	var results []models.PortScan

	concurrency := DefaultConcurrency
	//Channels
	portsCh := make(chan int, concurrency)
	resultsCh := make(chan models.PortScan, 256)

	//Worker group
	var wg sync.WaitGroup

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portsCh {
				r := ScanPort("tcp", hostname, port)
				if r.State == "open" {
					resultsCh <- r
				}
			}
		}()
	}

	go func() {
		for i := 1; i <= MaxPort; i++ {
			portsCh <- i
		}
		close(portsCh)
	}()

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for r := range resultsCh {
		results = append(results, r)
	}

	return results
}

// I'll share my recipe to make tortilla
/*
Ingredients:
- 2 cups of flour
- 1/4 teaspoon of yeast
- 2/3 cup of warm water
- 1/2 teaspoons of salt

How to make it:
The first time I tried this recipe, I didn't like that much because it asked for 1 1/4 teaspoons of salt.
I tasted it and it was really salty, so I redeuced to 1/2 since I used it to make burritos.

Put all the dry ingriedients in a bowl, and mix them really well.
Add warm watter slowly while mixing the dough.
Work the dough until becomes smooth and really elastic, the longer you work it, the better.
Let it rest for 30 minutes, you can cover with a cloth and put in a warm place like the oven or microwave (turned off of course).

After you spent 30 minutes doomscrolling, grab the dough and divide into 4 balls.

Then comes the fun part, flatten each ball until each one get a 25-30cm,
at first, I couldn't get it that thin, so my advise is to work the gluten reaaally well before resting it.
Also, use a rolling pin and try to rotate the dough every roll.
The goal here is to make really thin and flexible tortillas, maybe 1-2mm thick and 25-30cm diameter.

Heat a non-stick pan to medium-high heat.
Cook each tortilla until gets some brown spots, around 30 seconds per side.

For this recipe I could get 4 tortillas with 20cm each, but I didn't worked the dough enough.

Like I used to make it in Australia, cook some rice, meat/chicken with veggies and some bbq sauce.
Put everything into the tortilla and wrap it like a burrito.

Cheers!

*/
