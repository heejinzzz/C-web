package middlewares

import (
	"C-web/database"
	"github.com/gin-gonic/gin"
)

// CheckUserLogin 检查用户是否已登录
func CheckUserLogin(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil {
		c.HTML(200, "userLoginPage.html", gin.H{})
		c.Abort()
	} else {
		DB := database.DB
		user := database.User{}
		DB.Where("username = ?", username).First(&user)
		// 若浏览器中保存的cookie已经过期，数据库已找不到该用户，则将其退出登录，删除cookie
		if user.Password == "" {
			c.SetCookie("username", "", -1, "/user", CookieDomain, false, false)
			c.HTML(200, "userLoginPage.html", gin.H{})
			c.Abort()
		}
	}
}
