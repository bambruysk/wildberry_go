package main

import (
	"fmt"
	"testing"
)

func TestGetArticlesFromCatalogPage (t *testing.T) {
	catalog := "https://www.wildberries.ru/catalog/elektronika/tehnika-dlya-kuhni/kuhonnye-vesy"

	articles, err := GetArticlesFromCatalogPage(catalog)

	if err !=  nil {
		fmt.Println(err)
		t.Errorf("Error in %v", err)
	}
	fmt.Println(articles)

}