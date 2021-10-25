package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Fatal("file [%s] does not exist", fileName)
	}
	file, _ := ioutil.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	for _, nextProd := range productList {
		productMap.m[nextProd.ProductID] = nextProd
	}
}

func getProduct(id int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()
	if product, ok := productMap.m[id]; ok {
		return &product
	}
	return nil
}

func removeProduct(id int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, id)
}

func getProductList() []Product {
	productMap.RLock()
	defer productMap.RUnlock()
	products := make([]Product, 0, len(productMap.m))
	for _, nextProduct := range productMap.m {
		products = append(products, nextProduct)
	}
	return products
}

func getNextProductID() int {
	highestID := -1
	productList := getProductList()
	for _, product := range productList {
		if product.ProductID > highestID {
			highestID = product.ProductID
		}
	}
	return highestID + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		if oldProduct == nil {
			return 0, fmt.Errorf("product id [%d] does not exist", product.ProductID)
		}
	} else {
		product.ProductID = getNextProductID()
	}
	productMap.Lock()
	productMap.m[product.ProductID] = product
	productMap.Unlock()
	return product.ProductID, nil
}
