package main

import (
	functions "golang_tg/cmd"
)

func main() {
	//generate_channel("Лежбище котиков", "Крутой канал про котиков", "bear.jpg")
	//TGToolsGenerateChannel("Лютые заносы", "Канал про выигрыши в казино. Лучшие моменты стримов и ссылки на игры которые заносят.", "casino.jpg")
	//channels := search_contact("казино", 100)
	//fmt.Println(channels)
	//check7()
	//search_messages("занос", 1000)
	functions.Get_rows_in_file("input/channel_list")

}
