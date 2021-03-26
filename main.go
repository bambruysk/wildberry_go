package main

import (
	"encoding/json"
	"errors"
	_ "io/ioutil"
	"strings"

	"fmt"
	"io"
	"log"
	"net/http"


	"github.com/PuerkitoBio/goquery"
)

func main() {
	product := "11119714"

	page, err := MakeRequest(CreateWBUrl(product))
	if err != nil {
		log.Println("Error in request creation")
	}

	//fil, err := io.ReadAll(page)

	//err = ioutil.WriteFile("dump_page.html", fil, 0644)

	ssrJSON, err := ExtractSsrModel(page)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(ssrJSON))

}

func MakeRequest(URL string) (io.ReadCloser, error) {

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


	return resp.Body, nil

}







func GetArticlesFromCatalogPage(URL string) ([]string, error) {

	
	
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

	// Load the HTML document
		doc, err := goquery.NewDocumentFromResponse(resp)
		if err != nil {
			log.Fatal(err)
		}
		



	// Find the review items
	doc.Find("div").Find(".dtList-inner").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		ref := s.Find("href")
		fmt.Print("Review",ref.Text(), ref)
	})



//	resp, err := soup.Get(url)



	//container = soup.select("div.dtList.i-dtList.j-card-item")
//	container := doc.FindAllStrict("div","dtList i-dtList j-card-item")
	// url_block = block.select_one(
    //         'a.ref_goods_n_p.j-open-full-product-card')

    //     url = ("https://www.wildberries.ru" +
    //            url_block.get('href')).replace("?targetUrl=GP", "")

	res := make([]string, 20)




	return  res , nil
}


func CreateWBUrl(article string) string {
	return "https://wildberries.ru/catalog/" + article + "/detail.aspx"
}

func ExtractSsrModel(body io.ReadCloser) (string, error) {
	bodyB, err := io.ReadAll(body)
	defer body.Close()
	if err != nil {
		return "", err
	}
	bodyT := string(bodyB)

	//fmt.Println(bodyT)

	start := strings.Index(bodyT, "ssrModel")
	if start < 0 {
		return "", errors.New("ssrModel not found")
	}
	text := bodyT[start:]
	parentFound := false
	for i, c := range text {
		if c == '{' {
			start += i
			parentFound = true
			break
		}
	}
	if parentFound == false {
		return "", errors.New("no { after ssrModel tag")
	}

	parenthesisCount := 0
	text = bodyT[start:]
	var end int
	for i, c := range text {
		switch c {
		case '{':
			parenthesisCount++
		case '}':
			parenthesisCount--
		}
		if parenthesisCount == 0 {
			end = start + i + 1
			break
		}
	}
	if end == 0 {

		return "", fmt.Errorf("No '}' after ssrMpodel. Parentehsis =  %d", parenthesisCount)
	}
	log.Println(bodyT[start:end])
	return bodyT[start:end], nil

}

type SuppliersInfo map [string] SupplierInfo


type SupplierInfo struct {
	supplierName string
	ogrn string
}

type Product struct {
	suppliersInfo SuppliersInfo
}

type Nomenclature struct {
}

type ProductCard struct {
	link                    int
	star                    int
	brandName               string
	brandID                 int
	brandDirectionID        int
	brandDirectionPicsCount int
	description             string
	goodsName               string
	nomenclatures           []Nomenclature
}

func parseProductInfoFromJSON(info []byte) (Product, error) {

	var f interface{}
	err := json.Unmarshal(info, &f)
	if err != nil {
		return Product{}, errors.New("Unable parse product")
	}

	m := f.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			switch k {
				case "suppliersInfo": {
					// for article, info := range vv {

					// }  
				}
			case "productCard": {
			//	card := ProductCard{}
			//	card.parse(vv)
			}
			}
			
		default:
			fmt.Println(k, "is of a type I don't know how to handle")
		}
	}

	return Product{}, nil
}


func (s * SupplierInfo) parse( m map[string] interface{}) error {
	for k,v :=  range m {
		switch vv := v.(type) {
		case string: {
			switch k {
			case "supplierName": 
				s.supplierName = vv
			case "ogrn":
				s.ogrn = vv
			}
		}
				
		}
	}
	return nil
}

// func (s * ProductCard) parse( m [] interface{}) error {
// 	for k,v :=  range (m.(map[string] interface{})) {
// 		switch vv := v.(type) {
// 		case string: {
// 			switch k {
// 			case "supplierName": 
// 				s.supplierName = vv
// 			case "ogrn":
// 				s.ogrn = vv
// 			}
// 		}
				
// 		}
// 	}
// 	return nil
// }