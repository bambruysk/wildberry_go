package main

import (
	"fmt"
	"testing"
)

func TestGetArticlesFromCatalogPage (t *testing.T) {
	catalog := "https://www.wildberries.ru/catalog/elektronika/tehnika-dlya-kuhni/kuhonnye-vesy?page=2"

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