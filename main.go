package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Product struct {
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float32 `json:"price"`
}

var cart = []Product{
	{Name: "Shirt", Qty: 2, Price: 200},
	{Name: "Pant", Qty: 3, Price: 400},
	{Name: "Toy", Qty: 5, Price: 500},
}

func getcart(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, cart)
}

func getproductByName(name string) (*Product, error) {
	for i, p := range cart {
		if p.Name == name {
			return &cart[i], nil
		}
	}
	return nil, errors.New("Product not found")
}

func productByName(c *gin.Context) {
	name := c.Param("name")
	Product, err := getproductByName(name)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
	}
	c.IndentedJSON(http.StatusOK, Product)
}

func CheckoutProduct(c *gin.Context) {
	name, ok := c.GetQuery("name")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Product"})
		return
	}
	Product, err := getproductByName(name)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	if Product.Qty <= 0 || Product.Qty > 8 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Product"})
		return
	}

	Product.Qty -= 1
	c.IndentedJSON(http.StatusOK, Product)

}
func returnProduct(c *gin.Context) {

}
func addProduct(c *gin.Context) {
	var newproduct Product

	if err := c.BindJSON(&newproduct); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	if newproduct.Qty < 1 || newproduct.Qty > 8 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "qty mismatch"})
		return
	}

	cart = append(cart, newproduct)
	c.IndentedJSON(http.StatusCreated, newproduct)
}
func main() {
	router := gin.Default()
	router.GET("/cart", getcart)
	router.GET("/cart/:name", productByName)
	router.POST("/cart", addProduct)
	router.PATCH("/checkout", CheckoutProduct)
	router.PATCH("/return", returnProduct)
	router.Run("localhost:9090")

}
