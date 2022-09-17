package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/DennisJensen95/golang-rest-api/config"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	fmt.Println("Wassup clean arch is nice")

	r := gin.Default()
	r.GET("/:number/add", adding_numbers)
	r.Run() // Listen and serve
}

func adding_numbers(c *gin.Context) {
	num := c.Param("number")
	intVar, err := strconv.Atoi(num)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Number not supported"})
		return
	}
	resulting_number := fmt.Sprintf("%d", intVar+intVar)
	c.JSON(http.StatusOK, gin.H{
		"result": resulting_number,
	})
}
