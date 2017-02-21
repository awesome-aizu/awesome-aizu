package main

import (
	"io"
	"net/http"

	"log"

	"os"
	"strings"

	"flag"
	"net/url"

	"gopkg.in/xmlpath.v2"
	"gopkg.in/yaml.v2"
)

type Restaurant struct {
	Name                  string
	Time                  Time
	Tel                   string
	Category              string
	Location              Location
	IsAvailableCreditcard bool `yaml:"is_available_creditcard"`
}

type Time struct {
	Monday          []OpeningHours
	Tuesday         []OpeningHours
	Wednesday       []OpeningHours
	Thursday        []OpeningHours
	Friday          []OpeningHours
	Saturday        []OpeningHours
	Sunday          []OpeningHours
	NationalHoliday []OpeningHours `yaml:"national_holiday"`
}

type OpeningHours struct {
	Open  string
	Close string
}

type Location struct {
	Address   string
	URL       string
	Latitude  string
	Longitude string
}

var (
	nameXPath     *xmlpath.Path
	telXPath      *xmlpath.Path
	categoryXPath *xmlpath.Path
	addressXPath  *xmlpath.Path
	targetURL     string
)

func init() {
	nameXPath = xmlpath.MustCompile(`//*[@id="rstdata-wrap"]/table[1]/tbody/tr[1]/td/p/strong`)
	categoryXPath = xmlpath.MustCompile(`//*[@id="rstdata-wrap"]/table[1]/tbody/tr[2]/td/p`)
	telXPath = xmlpath.MustCompile(`//*[@id="rstdata-wrap"]/table[1]/tbody/tr[3]/td/div/div/p/strong`)
	addressXPath = xmlpath.MustCompile(`//*[@id="rstdata-wrap"]/table[1]/tbody/tr[4]/td/p`)
}

func getName(n *xmlpath.Node) string {
	if name, ok := nameXPath.String(n); ok {
		return sanitaize(name)
	}

	return ""
}

func getTel(n *xmlpath.Node) string {
	if tel, ok := telXPath.String(n); ok {
		return sanitaize(tel)
	}

	return ""
}

func getCateogry(n *xmlpath.Node) string {
	if cateogry, ok := categoryXPath.String(n); ok {
		return sanitaize(cateogry)
	}
	return ""
}

func getAddress(n *xmlpath.Node) string {
	if address, ok := addressXPath.String(n); ok {
		return sanitaize(address)
	}
	return ""
}

func getHTML(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func sanitaize(s string) string {
	return strings.TrimSpace(s)
}

func htmlToRestaurant(body io.Reader) (Restaurant, error) {
	n, err := xmlpath.ParseHTML(body)
	if err != nil {
		return Restaurant{}, err
	}

	return Restaurant{
		Name: getName(n),
		Time: Time{
			Monday:          []OpeningHours{OpeningHours{}},
			Tuesday:         []OpeningHours{OpeningHours{}},
			Wednesday:       []OpeningHours{OpeningHours{}},
			Thursday:        []OpeningHours{OpeningHours{}},
			Friday:          []OpeningHours{OpeningHours{}},
			Saturday:        []OpeningHours{OpeningHours{}},
			Sunday:          []OpeningHours{OpeningHours{}},
			NationalHoliday: []OpeningHours{OpeningHours{}},
		},
		Tel:      getTel(n),
		Category: getCateogry(n),
		Location: Location{
			Address: getAddress(n),
		},
	}, nil
}

func main() {
	flag.StringVar(&targetURL, "url", "", "tabelog url")
	flag.Parse()

	if targetURL == "" {
		log.Fatal("set a url")
	}

	if _, err := url.Parse(targetURL); err != nil {
		log.Fatal("set a valid url")
	}

	body, err := getHTML(targetURL)
	if err != nil {
		log.Fatal(err)
	}
	r, err := htmlToRestaurant(body)
	if err != nil {
		log.Fatal(err)
	}

	b, err := yaml.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stdout.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
