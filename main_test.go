package main

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var  catalog = "https://www.wildberries.ru/catalog/elektronika/tehnika-dlya-kuhni/kuhonnye-vesy"

func TestGetArticlesFromCatalogPage (t *testing.T) {
	//catalog := "https://www.wildberries.ru/catalog/elektronika/tehnika-dlya-kuhni/kuhonnye-vesy?page=2"

	articles, err := GetArticlesFromCatalogPage(catalog)

	if err !=  nil {
		fmt.Println(err)
		t.Errorf("Error in %v", err)
	}
	fmt.Println(articles)

}


func TestExtractArticleFromURL (t *  testing.T) {
	url  := "/catalog/19377339/detail.aspx?targetUrl=GP"

	article, err := extractArticleFromURL(url)
	
	if err != nil {
		t.Errorf("Error in %s %s", url, article)
	}

	if article  != "19377339"  {

		t.Errorf("Error in %s %s", url, article)
	}

}


func TestParseProductInfoFromJSON (t * testing.T) {
	article := "11119714"
	data, err := ioutil.ReadFile("./test.json")
    if err != nil {
    	t.Errorf(" Erroor read file %v", err)
    }
	product, err := parseProductInfoFromJSON(data, article)
	if err != nil {
      t.Errorf(" Erroor parse %v", err)
	}
	
	fmt.Println(product)
}


func TestParsePage ( t* testing.T) {
	articles, err := GetArticlesFromCatalogPage(catalog)
	if err != nil {
      t.Errorf(" Erroor parse %v", err)
	}
	products := make([]Product, 100)
	for _,a := range articles {
		product := ParseProductPage(a)
		products = append(products, product)
	}
	fmt.Println(products)
}