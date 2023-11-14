package barbershop

import (
	"sync"
	"time"

	"github.com/fatih/color"
)

type BarberShop struct {
	WaitingRoomSize int
	ChWaitingList   chan string
	ChShop          chan bool
	Total_served    int
	sync.WaitGroup
	sync.Mutex
}

func (shop *BarberShop) AddBarber(barber string) {
	shop.Add(1)
	isSleep := false
	isFirst := true
	go func() {
		defer shop.Done()
		for {
			select {
			case client, isShopOpen := <-shop.ChWaitingList:
				if !isShopOpen {
					color.Cyan(">> %s is going home", barber)
					return
				}
				if isSleep {
					color.Yellow("@@ %s waking up %s", client, barber)
					isSleep = false
				}
				isFirst = false
				color.Green("** %s is cutting %s's hair", barber, client)
				<-time.NewTimer(1000 * time.Millisecond).C
				color.Green("*** %s done cutting %s's hair", barber, client)
				shop.Lock()
				shop.Total_served++
				shop.Unlock()
			case <-shop.ChShop:
			default:
				if !isSleep && !isFirst {
					color.Red("~~ %s is sleeping", barber)
					isSleep = true
				}
			}
		}
	}()
}
