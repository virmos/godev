package letter
import (
	"fmt"
	"github.com/fatih/color"
)

// FreqMap records the frequency of each rune in a given text.
type FreqMap map[rune]int
var numberOfChar int

type Producer struct {
	data chan CharOrder
	quit chan chan error
}

type CharOrder struct {
	charNumber int
	message string
	success bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

// Frequency counts the frequency of each rune in a given text and returns this
// data as a FreqMap.
func Frequency(s string) FreqMap {
	m := FreqMap{}
	for _, r := range s {
		m[r]++
	}
	return m
}

func produce() *CharOrder {

}

func countCharJob(charMaker *Producer) {

}

func consume() {
	for j := range 
}

// ConcurrentFrequency counts the frequency of each rune in the given strings,
// by making use of concurrency.
func ConcurrentFrequency(l []string) FreqMap {
	numberOfChar = len(l)

	charJob := &Producer {
		data: make(chan CharOrder),
		quit: make(chan chan error)
	}
	go countCharJob(charJob)

	for i := range charJob.data {
		if i.charNumber <= numberOfChar {
			if i.success {
				
			} 
			else {
				
			}
		}
		else {
			color.Cyan("Done...")
			err := charJob.Close()
			if err != nil {
				color.Red("*** Error closing channel!", err)
			}
		}
	}
	consume()
}
