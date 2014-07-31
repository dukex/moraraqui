package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/dukex/moraraqui/crawlers"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

// const (
//// According to Wikipedia, the Earth's radius is about 6,371km
//  EARTH_RADIUS = 6371
// )

// func inRadius(lat, lng, radius string) ([]Imovel, error) {
// 	// select_str := fmt.Sprintf("SELECT * FROM imoveis")
// 	lat1 := fmt.Sprintf("sin(radians(%s)) * sin(radians(imovels.lat))", lat)
// 	lng1 := fmt.Sprintf("cos(radians(%s)) * cos(radians(imovels.lat)) * cos(radians(imovels.lng) - radians(%s))", lat, lng)
// 	where_str := fmt.Sprintf("acos(%s + %s) * %f <= %s", lat1, lng1, float64(EARTH_RADIUS), radius)
// 	// query := fmt.Sprintf("%s %s", select_str, where_str)

// 	var imoveis []Imovel
// 	err := DB.Table("imovels").Where(where_str).Find(&imoveis).Error

// 	return imoveis, err
// }

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static("assets"))

	m.Get("/api/imoveis/:state/:city/:neighborhood", func(res http.ResponseWriter, req *http.Request, params martini.Params) {
		res.Header().Set("Content-Type", "text/event-stream")

		item, timeout := crawlers.Get(params["state"], params["city"], params["neighborhood"])

		for {
			select {
			case property, ok := <-item:
				if !ok {
					item = nil
				}
				b, _ := json.Marshal(property)
				b = append(b, []byte("\n")...)
				res.Write(b)
				res.(http.Flusher).Flush()
			case <-timeout:
				item = nil
				return
			}
		}
	})

	m.Get("/**", func(r render.Render) []byte {
		body, _ := ioutil.ReadFile("./assets/index.html")
		return body
	})

	m.Run()
}
