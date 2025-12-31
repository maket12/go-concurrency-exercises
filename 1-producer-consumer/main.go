//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream, ch chan *Tweet) {
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			close(ch)
			return
		}

		select {
		case ch <- tweet:

		}
	}
}

func consumer(ch chan *Tweet) {
	for {
		v, ok := <-ch
		if !ok {
			return
		}

		if v.IsTalkingAboutGo() {
			fmt.Println(v.Username, "\ttweets about golang")
		} else {
			fmt.Println(v.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()

	var tweet = make(chan *Tweet)

	// Producer
	go producer(stream, tweet)
	// Consumer
	consumer(tweet)

	fmt.Printf("Process took %s\n", time.Since(start))
}
