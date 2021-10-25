package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/meltemseyhan/inventoryservice/database"
)

func getProduct(id int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	row := database.DbConn.QueryRowContext(ctx, `SELECT productId, 
	manufacturer, 
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products WHERE productId = ?`, id)

	product := &Product{}
	err := row.Scan(&product.ProductID, &product.Manufacturer, &product.Sku, &product.Upc, &product.PricePerUnit, &product.QuantityOnHand, &product.ProductName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return product, nil
}

func removeProduct(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `DELETE FROM products WHERE productId=?`, id)
	if err != nil {
		return err
	}
	return nil
}

func getProductList() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT productId, 
	manufacturer, 
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products`)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	products := make([]Product, 0)
	for results.Next() {
		var nextProduct Product
		results.Scan(&nextProduct.ProductID, &nextProduct.Manufacturer, &nextProduct.Sku, &nextProduct.Upc, &nextProduct.PricePerUnit, &nextProduct.QuantityOnHand, &nextProduct.ProductName)
		products = append(products, nextProduct)
	}
	return products, nil
}

func getTopTenProducts() ([]Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	results, err := database.DbConn.QueryContext(ctx, `SELECT productId, 
	manufacturer, 
	sku,
	upc,
	pricePerUnit,
	quantityOnHand,
	productName
	FROM products ORDER BY quantityOnHand DESC LIMIT 10`)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	products := make([]Product, 0)
	for results.Next() {
		var nextProduct Product
		results.Scan(&nextProduct.ProductID, &nextProduct.Manufacturer, &nextProduct.Sku, &nextProduct.Upc, &nextProduct.PricePerUnit, &nextProduct.QuantityOnHand, &nextProduct.ProductName)
		products = append(products, nextProduct)
	}
	return products, nil
}

func updateProduct(product Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	_, err := database.DbConn.ExecContext(ctx, `UPDATE products SET 
	manufacturer=?, 
	sku=?,
	upc=?,
	pricePerUnit=?,
	quantityOnHand=?,
	productName=?
	WHERE productId=?`, product.Manufacturer, product.Sku, product.Upc, product.PricePerUnit, product.QuantityOnHand, product.ProductName, product.ProductID)
	if err != nil {
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	result, err := database.DbConn.ExecContext(ctx, `INSERT INTO products SET 
	manufacturer=?, 
	sku=?,
	upc=?,
	pricePerUnit=?,
	quantityOnHand=?,
	productName=?`, product.Manufacturer, product.Sku, product.Upc, product.PricePerUnit, product.QuantityOnHand, product.ProductName)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}
