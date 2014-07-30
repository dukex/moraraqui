package crawlers

import (
	"time"

	. "github.com/dukex/moraraqui/models"
)

type Bot interface {
	FirstRun(channel chan *Property, state, city, neighborhood string) int
	Get(channel chan *Property, page int, state, city, neighborhood string)
}

var bots []Bot

func Get(state, city, neighborhood string) (chan *Property, <-chan time.Time) {

	item, timeout := make(chan *Property), time.After(time.Second*20)

	go func() {
		for i := 0; i < len(bots); i++ {
			pages := bots[i].FirstRun(item, state, city, neighborhood)

			if pages > 1 {
				for page := 2; page <= pages; page++ {
					go bots[i].Get(item, page, state, city, neighborhood)
				}
			}
		}
	}()

	return item, timeout
}

func init() {
	bots = make([]Bot, 1)
	//bots[0] = new(ZimpBot)
	bots[0] = new(ImovelWebBot)
}
