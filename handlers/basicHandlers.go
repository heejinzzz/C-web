package handlers

import (
	"github.com/gin-gonic/gin"
)

// HomePage 首页
func HomePage(c *gin.Context) {
	c.HTML(200, "homePage.html", gin.H{})
}
