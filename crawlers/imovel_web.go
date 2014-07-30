package crawlers

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	. "github.com/dukex/moraraqui/models"
)

const (
	imovelwebbaseurl = "http://www.imovelweb.com.br"
)

type ImovelWebBot struct {
}

func (i *ImovelWebBot) FirstRun(channel chan *Property, state, city, neighborhood string) int {
	url := i.urlFor("1", state, city, neighborhood)
	doc := i.parserPage(channel, url)

	lastPageS := doc.Find(".box-pagging .bt-pagging-num--p").Last().Text()
	lastPage, _ := strconv.Atoi(lastPageS)
	return lastPage
}

func (i *ImovelWebBot) Get(channel chan *Property, page int, state, city, neighborhood string) {
	pageS := strconv.Itoa(page)

	url := i.urlFor(pageS, state, city, neighborhood)
	i.parserPage(channel, url)

}

func (i *ImovelWebBot) parserPage(channel chan *Property, url string) *goquery.Document {
	log.Println(" Parsing", url, "...")

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return doc
	}

	doc.Find("ul[itemtype='http://www.schema.org/RealEstateAgent'] > li").Each(func(_ int, s *goquery.Selection) {

		pType := i.getType(s.Find(".busca-item-heading2").Text())

		if pType > 0 {
			url, _ := s.Find("a").Attr("href")
			var property Property

			if err := DB.Where("url = ?", imovelwebbaseurl+url).First(&property).Error; err != nil {
				property.Title = s.Find(".busca-item-heading1").Text()
				property.Address = s.Find(".busca-item-endereco").Text()
				property.Url = imovelwebbaseurl + url
				// property.Size = z.getSize(s.Find(".l-inline-list li").Eq(0).Text())
				property.Type = pType
				// property.Bedroom = z.getBedroom(s.Find(".l-inline-list li").Eq(1).Text())
				property.Value = i.getValue(s.Find(".busca-item-preco").Text())
				property.Neighborhood = strings.Split(property.Title, ",")[0]
				DB.Save(&property)
			}
			channel <- &property
		}
	})

	return doc
}

func (i *ImovelWebBot) getValue(text string) float64 {
	text = strings.Replace(text, "R$ ", "", -1)
	text = strings.Replace(text, ".", "", -1)
	value, _ := strconv.ParseFloat(text, 64)
	return value
}

func (i *ImovelWebBot) getType(text string) int {
	text = strings.TrimSpace(text)
	splited := strings.Split(text, " ")

	switch splited[0] {
	case "Casa":
		return House
	case "Apartamento":
		return Apartament
	}

	return Undefined
}

func (i *ImovelWebBot) urlFor(page, state, city, neighborhood string) string {
	base := "aluguel"
	query := "?pg=" + page
	return strings.Join([]string{imovelwebbaseurl, base, state, city, neighborhood, query}, "/")
}
