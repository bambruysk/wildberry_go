package main

import (
	"encoding/json"
	"errors"
	_ "io/ioutil"
	"strconv"
	"strings"
	"sync"
	"time"

	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/sqlite"
  	"gorm.io/gorm"
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

	res := make([]string, 0)

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
	return bodyT[start:end], nil

}

type SuppliersInfo map[string]SupplierInfo

type SupplierInfo struct {
	supplierName string
	ogrn         string
}

type Product struct {
	gorm.Model 
	SupplierName string
	BrandName    string
	Article      string `gorm:"primaryKey"` 
	URL          string
	Price        int
	OrderCount   int
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
		fmt.Println("Trouble in parse", article, info)
		return Product{}, errors.New("Unable parse product")
	}

	m := f.(map[string]interface{})
	product := Product{Article: article}
	for k, v := range m {
		switch k {
		case "suppliersInfo":
			{
				suppliersInfo := v.(map[string]interface{})[article]
				supplierInfo := suppliersInfo.(map[string]interface{})["supplierName"]
				product.SupplierName = supplierInfo.(string)
			}
		case "productCard":
			{
				// TODO: add checking
				productCard := v.(map[string]interface{})
				product.BrandName = productCard["brandName"].(string)
				n, ok := productCard["nomenclatures"]
				var nomenclatures map[string]interface{}
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

func ParseProductPage(article string) (Product, error) {
	if len(article) == 0 {
		return Product{}, errors.New("Article must be not empty")
	}
	page, err := GetProductPage(article)
	if err != nil {
		log.Printf("Found err at %v \n", err)
	}
	ssr, err := ExtractSsrModel(page)

	product, err := parseProductInfoFromJSON([]byte(ssr), article)

	if err != nil {
		log.Printf("Parse  err at %v \n", err)
	}

	return product, nil
}

func RetrieveAllArticlesFormCatalog(URL string) ([]string, error) {

	var articles []string
	//article_count := 0
	for page := 1; ; page++ {
		artcls, err := GetArticlesFromCatalogPage(URL + "?page=" + strconv.Itoa(page))
		log.Println(artcls)
		if err != nil {
			return articles, fmt.Errorf("Article get error, %v ", err)
		}
		if len(artcls) == 0 {
			return articles, nil
		}
		articles = append(articles, artcls...)
	}

}

func ParseCatalogPages(URL string) ([]Product, error) {
	/// retrinve all articles
	var products []Product
	aritcles, err := RetrieveAllArticlesFormCatalog(URL)

	if err != nil {
		return nil, err
	}

	artChan := make(chan string,8)
	productChan := make(chan Product,8)
	numsOfWorkers := 8

	// for _, article :=  range aritcles {
	// 	product, err :=  ParseProductPage(article)
	// 	if err != nil {
	// 		return products, err
	// 	}
	// 	products = append(products, product)
	// }
	var wg sync.WaitGroup
	for i := 0; i < numsOfWorkers; i++ {
		go func(c chan string, out chan Product, wg * sync.WaitGroup ) {
			wg.Add(1)
			for article := range c {
				product, err := ParseProductPage(article)
				if err != nil {
					// Check it
					panic(err)
				}
				out <- product
				time.Sleep(1*time.Millisecond)
			}
			wg.Done()
		}(artChan, productChan, &wg)
	}
	productsCh := make(chan []Product)

	// retreive results
	go func(in chan Product, out chan []Product) {
		var res []Product
		db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		  // Migrate the schema
		db.AutoMigrate(&Product{})
		  
		for p := range in {
			db.Create(&p)
			res = append(res, p)
			//fmt.Println(p)
		}

		productsCh <- res
	}(productChan, productsCh)

	for _, a := range aritcles {
		artChan <- a
	}

	close(artChan)

	wg.Wait()

	close(productChan)

	products = <-productsCh

	//close(productsCh)

	return products, nil
}

func CheckDBConnect () (error) {

	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	fmt.Println(db)
	return nil
}


func GetDBConnect () ( error) {

	// github.com/mattn/go-sqlite3
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	fmt.Println(db)
	return nil
}
