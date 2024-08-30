package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql" // Import MySQL driver
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {
	// Initialize the database connection to MySQL
	var err error
	dsn := "****:****4@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local" // Replace with your MySQL connection string
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	DB.AutoMigrate(&User{})

	router := gin.Default()

	router.GET("/users", GetAllUser)
	router.GET("/users/:id", GetUser)
	router.POST("/users", SaveUser)
	router.PUT("/users/:id", UpdateUser)
	router.DELETE("/users/:id", DeleteUser)

	router.Run(":8080") // Adjust port number as needed
}

func GetAllUser(c *gin.Context) {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   users,
	})
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := DB.First(&user, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   user,
	})
}

func SaveUser(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := DB.Create(&user).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   user,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := DB.First(&user, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := DB.Save(&user).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user User

	if err := DB.First(&user, id).Error; err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if err := DB.Delete(&user).Error; err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   user,
	})
}

type User struct {
	Id        uint   `gorm:"primaryKey" json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}
