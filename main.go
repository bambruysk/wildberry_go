package wildberry_go

import (
	_ "github.com/anaskhan96/soup"
	"net/http"
	"log"
	"io/ioutil"
)

func main() {
	MakeRequest()
}

func MakeRequest() {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
}