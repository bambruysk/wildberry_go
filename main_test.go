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
		product, err := ParseProductPage(a)
		if err !=  nil {
			t.Errorf("Error in parse products %v", err)
		}
		products = append(products, product)
	}
	fmt.Println(products)
}


func TestRetrieveAllArticlesFormCatalog( t * testing.T) {
	t.Log("Start testing")
	articles, err := RetrieveAllArticlesFormCatalog(catalog)
	if err != nil {
      t.Errorf(" Erroor parse %v", err)
	}
	fmt.Println("Articles found", articles, len(articles))
	t.Logf("Found %d articles", len(articles))
}


func TestParseCatalogPages( t * testing.T) {
	t.Log("Start testing")
	products, err := ParseCatalogPages(catalog)
	if err != nil {
      t.Errorf(" Erroor parse %v", err)
	}
	fmt.Println("products found", products, len(products))
	t.Logf("Found %d products", len(products))
}

func TestCheckDBConnect ( t * testing.T) {
	err := CheckDBConnect()
	if err != nil {
      t.Errorf(" Error db_connect %v", err)
	}
}