package main

import (
	_ "encoding/json"

	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/anaskhan96/soup"
)

func main() {
	product := "11119714"

	page, err := MakeRequest(CreateWBUrl(product))
	if err != nil {
		log.Println("Error in request creation")
	}
	pageText, err := io.ReadAll(page)
    if err != nil {
		print(err)

    }
	fmt.Print(string(pageText))

}

func MakeRequest(URL string) (io.ReadCloser, error)  {

	client := http.Client{}
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Println("Error in request creation")
		return nil, err
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.152 YaBrowser/21.2.2.101 Yowser/2.5 Safari/537.36")
	request.Header.Add("Acccept-Language", "ru")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	//var result map[string]interface{}
	//json.NewDecoder(resp.Body).Decode(&result)
	//log.Println(&resp.Body)
    body, err := io.ReadAll(resp.Body)
    if err != nil {
		print(err)
		return nil, err
    }
	fmt.Print(string(body))
	return  resp.Body, nil

}

func CreateWBUrl (article string) string {
	return "https://wildberries.ru/" + article + "/details.aspx" 
}



