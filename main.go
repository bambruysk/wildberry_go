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

	res := make([]string, 20)

	err = nil

	// Find the review items
	doc.Find("div").Find(".dtList-inner").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		linkTag := s.Find("a")
		ref, _ := linkTag.Attr("href")
		article, e := extractArticleFromURL(ref)
		if e != nil {
			err = e
		}

		res = append(res, article)
		log.Println("Found ref", ref)
	})

	//	resp, err := soup.Get(url)

	//container = soup.select("div.dtList.i-dtList.j-card-item")
	//	container := doc.FindAllStrict("div","dtList i-dtList j-card-item")
	// url_block = block.select_one(
	//         'a.ref_goods_n_p.j-open-full-product-card')

	//     url = ("https://www.wildberries.ru" +
	//            url_block.get('href')).replace("?targetUrl=GP", "")

	return res, err
}

//  /catalog/19377339/detail.aspx?targetUrl=GP
func extractArticleFromURL(URL string) (string, error) {
	surl := strings.Split(URL, "/")
	if len(surl) < 3 {
		return "", errors.New("Error in url format")
	}
	return surl[2], nil
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

type SuppliersInfo map[string]SupplierInfo

type SupplierInfo struct {
	supplierName string
	ogrn         string
}

type Product struct {
	SupplierName string
	BrandName    string
	Article      string
	URL 		string
	Price		int
	OrderCount 	int
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
	supplierName            string
	description             string
	goodsName               string
	nomenclatures           []Nomenclature
}

func parseProductInfoFromJSON(info []byte, article string) (Product, error) {

	var f interface{}
	err := json.Unmarshal(info, &f)
	if err != nil {
		return Product{}, errors.New("Unable parse product")
	}

	m := f.(map[string]interface{})
	product := Product{Article: article}
	for k, v := range m {
		switch k {
		case "suppliersInfo":
			{
				fmt.Println(v)
				suppliersInfo := v.(map[string]interface{})[article]
				supplierInfo := suppliersInfo.(map[string]interface{})["supplierName"]
				product.SupplierName = supplierInfo.(string)
				log.Println("Supplier Name is ", product.SupplierName)
			}
		case "productCard":
			{
				// TODO: add checking
				productCard := v.(map[string]interface{})
				product.BrandName = productCard["brandName"].(string)
				n, ok := productCard["nomenclatures"]
				var  nomenclatures map[string]interface{}
				if ok {
					nomenclatures = n.(map[string]interface{})
				} else {
					log.Println("Nomenclatures not found")
				}
				nomenclature := nomenclatures[article].(map[string]interface{})
				product.OrderCount = int(nomenclature["ordersCount"].(float64))
				sizes := nomenclature["sizes"].([]interface{})
				size := sizes[0].(map[string]interface{})
				product.Price = int(size["price"].(float64))

				//product.OrderCount = nomenclatures["orderCount"]


			}

		default:
			{
			}
		}
	}


	return product, nil
}

func (s *SupplierInfo) parse(m map[string]interface{}) error {
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			{
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

func GetProductPage(article string) (io.ReadCloser, error) {
	addr := CreateWBUrl(article)
	return MakeRequest(addr)
}

func ParseProductPage(article string) (product Product) {
	page, err := GetProductPage(article)
	if err != nil {
		log.Printf("Found err at %v \n", err)
	}
	ssr, err := ExtractSsrModel(page)

	product, err = parseProductInfoFromJSON([]byte(ssr), article)

	if err != nil {
		log.Printf("Parse  err at %v \n", err)
	}

	return product
}
