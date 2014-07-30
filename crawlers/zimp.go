package crawlers

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	. "github.com/dukex/moraraqui/models"
)

const (
	zimpbaseurl = "http://zimp.infomoney.com.br"
)

type ZimpBot struct {
}

func (z *ZimpBot) Get(page int, state, city, neighborhood string) []*Property {
	pageS := strconv.Itoa(page)

	url := z.urlFor(pageS, state, city, neighborhood)
	properties, _ := z.parserPage(url)

	return properties
}

func (z *ZimpBot) FirstRun(state, city, neighborhood string) (int, []*Property) {
	url := z.urlFor("1", state, city, neighborhood)
	properties, doc := z.parserPage(url)

	lastPageS := doc.Find(".pagination-centered .pagination li:not(.arrow)").Last().Text()
	lastPage, _ := strconv.Atoi(lastPageS)

	return lastPage, properties
}

func (z *ZimpBot) parserPage(url string) ([]*Property, *goquery.Document) {
	log.Println(" Parsing", url, "...")
	properties := make([]*Property, 0)

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return properties, doc
	}

	doc.Find(".properties-list > li").Each(func(i int, s *goquery.Selection) {
		pType := z.getType(s.Find(".property-type").Text())

		if pType > 0 {
			var property Property
			link := s.Find("h2 a")
			url, _ := link.Attr("href")

			if err := DB.Where("url = ?", zimpbaseurl+url).First(&property).Error; err != nil {
				property.Title = link.Text()
				property.Address = link.Text()
				property.Url = zimpbaseurl + url
				property.Size = z.getSize(s.Find(".l-inline-list li").Eq(0).Text())
				property.Type = pType
				property.Bedroom = z.getBedroom(s.Find(".l-inline-list li").Eq(1).Text())
				property.Value = z.getValue(s.Find(".property-list-price.show-for-large").Text())
				DB.Save(&property)
			}

			properties = append(properties, &property)
		}
	})

	return properties, doc
}

func (z *ZimpBot) getSize(htmlsize string) int {
	slited := strings.Split(htmlsize, " ")
	size, _ := strconv.Atoi(slited[0])
	return size
}

func (z *ZimpBot) getType(text string) int {
	splited := strings.Split(text, " ")
	switch splited[0] {
	case "Casa":
		return House
	case "Apartamento":
		return Apartament
	}

	return Undefined
}

func (z *ZimpBot) getBedroom(text string) int {
	onlyNumber := strings.Replace(text, "quartos: ", "", -1)
	number, _ := strconv.Atoi(onlyNumber)
	return number
}

func (z *ZimpBot) getValue(text string) float64 {
	text = strings.Replace(text, "R$ ", "", -1)
	text = strings.Replace(text, ".", "", -1)
	text = strings.Replace(text, ",", ".", -1)
	value, _ := strconv.ParseFloat(text, 64)
	return value
}

func (z *ZimpBot) urlFor(page, state, city, neighborhood string) string {
	query := "?p=" + page + "&psize=30"
	base := "imoveis/busca/aluguel"

	return strings.Join([]string{zimpbaseurl, base, state, city, neighborhood, query}, "/")
}
