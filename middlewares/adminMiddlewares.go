package middlewares

import (
	"C-web/database"
	"github.com/gin-gonic/gin"
)

var CookieDomain = "127.0.0.1"

// CheckAdminLogin 检查管理员是否已登录
func CheckAdminLogin(c *gin.Context) {
	adminname, err := c.Cookie("adminname")
	if err != nil {
		c.HTML(200, "adminLoginPage.html", gin.H{})
		c.Abort()
	} else {
		DB := database.DB
		admin := database.Admin{}
		DB.Where("adminname = ?", adminname).First(&admin)
		// 若浏览器中保存的cookie已经过期，数据库已找不到该管理员，则将其退出登录，删除cookie
		if admin.Password == "" {
			c.SetCookie("adminname", "", -1, "/admin", CookieDomain, false, false)
			c.HTML(200, "adminLoginPage.html", gin.H{})
			c.Abort()
		}
	}
}
