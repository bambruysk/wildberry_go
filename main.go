package main

import (
	_ "encoding/json"
	"errors"
	_ "io/ioutil"
	"strings"

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

	//fil, err := io.ReadAll(page)

	//err = ioutil.WriteFile("dump_page.html", fil, 0644)

	ssrJSON,err :=  ExtractSsrModel(page)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(ssrJSON))

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

	return  resp.Body, nil

}

func CreateWBUrl (article string) string {
	return "https://wildberries.ru/catalog/" + article + "/detail.aspx" 
}


func ExtractSsrModel ( body io.ReadCloser) (string, error) {
	bodyB, err := io.ReadAll (body)
	defer body.Close()
	if err != nil {
		return "", err
    }
	bodyT := string(bodyB)

	//fmt.Println(bodyT)

	start := strings.Index(bodyT, "ssrModel")
	if start < 0 {
		return "" , errors.New("ssrModel not found")
	}
	text := bodyT[start:]
	parentFound := false
	for i,c := range text {
		if c == '{' {
			start += i
			parentFound =  true
			break
		}
	}
	if parentFound == false {
		return "", errors.New("no { after ssrModel tag")
	}

	parenthesisCount :=  0
	text = bodyT[start:]
	var end int
	for i,c := range text {
		switch c {
		case '{': parenthesisCount++
		case '}': parenthesisCount--
		}
		if parenthesisCount == 0 {
			end = start + i +1 
			break
		}
	}
	if end == 0 {

		return "", fmt.Errorf("No '}' after ssrMpodel. Parentehsis =  %d", parenthesisCount)
	}
	log.Println(bodyT[start:end])
	return bodyT[start:end], nil

}

type SuppliersInfo struct {
	article string
}

type SupplierInfo struct {
	cod1S string

}


type Product struct {
	suppliersInfo SuppliersInfo

}

type Nomenclature struct {

}

type ProductCard struct {
	link int 
	star int 
	brandName string 
	brandId int 
	brandDirectionId int 
	brandDirectionPicsCount int  
	description string
	goodsName string
	nomenclatures [] Nomenclature
}

func parseProductInfoFromJSON( info string) ( Product, error) {



	return Product{}, nil
}