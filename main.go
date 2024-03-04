package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float32 `json:"price"`
}

var DB *gorm.DB

var cart = []Product{
	{Name: "Shirt", Qty: 2, Price: 200},
	{Name: "Pant", Qty: 3, Price: 400},
	{Name: "Toy", Qty: 5, Price: 500},
}

func getcart(c *gin.Context) {
	var products []Product
	if err := DB.Find(&products).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)

}

func getproductByName(name string) (*Product, error) {
	var product Product

	if err := DB.Where("name = ?", name).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Product not found")
		}
		return nil, err
	}

	return &product, nil
}

func productByName(c *gin.Context) {
	name := c.Param("name")
	var product Product

	if err := DB.Where("name = ?", name).First(&product).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, product)
}

func CheckoutProduct(c *gin.Context) {
	name := c.Query("name")
	var product Product

	if err := DB.Where("name = ?", name).First(&product).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	if product.Qty <= 0 || product.Qty > 8 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid quantity"})
		return
	}

	product.Qty--
	if err := DB.Save(&product).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to checkout product"})
		return
	}

	c.IndentedJSON(http.StatusOK, product)

}
func returnProduct(c *gin.Context) {
	name := c.Query("name")
	var product Product

	if err := DB.Where("name = ?", name).First(&product).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		return
	}

	product.Qty++
	if err := DB.Save(&product).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to return product"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Product quantity increased by 1"})
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

	dsn := "host=cornelius.db.elephantsql.com user=qzttriqe password=lcvShz27-sqGJUHBBn3WzuE1wTbPe3BH dbname=qzttriqe port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB.AutoMigrate(&Product{})

	router := gin.Default()
	router.GET("/cart", getcart)
	router.GET("/cart/:name", productByName)
	router.POST("/cart", addProduct)
	router.PATCH("/checkout", CheckoutProduct)
	router.PATCH("/return", returnProduct)
	router.Run("localhost:9090")

}
