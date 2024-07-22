package controller

import (
        "cust-mobile-eu/handler"
        "strings"

        "github.com/gin-gonic/gin"
)

// ListProducts - Get all Products Information
//
//      @Summary        List all Products information from the server
//      @Description    Get method for listing all the products from server
//      @Tags           Products
//      @Produce        json
//      @Accept         json
//      @Success        200
//      @Failure        404
//      @Router         /products       [get]
func ListProducts(c *gin.Context) {
        
        p := handler.Product{}
        p.Build(nil)

        response, status := p.GetProducts()

        c.IndentedJSON(status, response)
        return
}

// GetProduct - Get information about particular Product
//
//      @Summary        Get particular Product's information from the server
//      @Description    Get method for particular product information from server
//      @Tags           Product
//      @Param          id      path    int     true    "Product ID"
//      @Produce        json
//      @Accept         json
//      @Success        200
//      @Failure        404
//      @Router         /product/{id}   [get]
func GetProduct(c *gin.Context) {
        id := strings.TrimSpace(c.Param("id"))
        
        p := handler.Product{}
        p.Build(nil)

        response, status := p.GetProduct(id)

        c.IndentedJSON(status, response)
        return
}
