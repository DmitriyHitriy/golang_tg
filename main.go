package main

import (
	"fmt"
)

func main() {
	//generate_channel("Лежбище котиков", "Крутой канал про котиков", "bear.jpg")
	channels := search_contact("казино", 100)
	fmt.Println(channels)
}
