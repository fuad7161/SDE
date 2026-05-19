package channels

import "fmt"

func BufferChannel() {
	message := make(chan string, 2)

	message <- "buffered"
	message <- "channel"

	fmt.Println(<-message)
	fmt.Println(<-message)
}

func Solve() {
	message := make(chan string)

	go func() {
		message <- "ping"
	}()

	msg := <-message
	fmt.Println(msg)
}
