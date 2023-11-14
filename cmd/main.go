package main

import (
	"fmt"
	"time"

	"github.com/andrizalchaidar/learn-go-concurrency/barbershop"
	"github.com/fatih/color"
)

func main() {

	barberShop := barbershop.BarberShop{
		WaitingRoomSize: 5,
		ChWaitingList:   make(chan string, 50),
		ChShop:          make(chan bool),
		Total_served:    0,
	}

	barberShop.AddBarber("foo")
	barberShop.AddBarber("bar")
	barberShop.AddBarber("baz")
	barberShop.AddBarber("qux")
	barberShop.AddBarber("corge")
	barberShop.AddBarber("garply")

	go func() {
		<-time.NewTimer(10 * time.Second).C
		close(barberShop.ChShop)
	}()
	go func() {
		i := 1
		for {
			select {
			case <-barberShop.ChShop:
				color.Red("!!! Waiting list order is closed")
				close(barberShop.ChWaitingList)
				return
			case <-time.NewTimer(100 * time.Millisecond).C:
				client := fmt.Sprintf("Client #%d", i)
				select {
				case barberShop.ChWaitingList <- client:
					color.Yellow("?? %s sit in waiting room", client)
				default:
					color.Red("!! Waiting room is full, so Client #%d leave.", i)
				}
				i++
			}
		}
	}()

	barberShop.Wait()
	color.Cyan("____________________________________")
	color.Cyan("All order done, everyone going home.")
	color.Cyan("%d", barberShop.Total_served)
}
