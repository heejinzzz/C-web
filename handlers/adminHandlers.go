package handlers

import (
	"C-web/database"
	"github.com/gin-gonic/gin"
	"github.com/heejinzzz/logger"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// CheckAdminPassword 检查管理员输入的用户名和密码，告知前端是否正确
func CheckAdminPassword(c *gin.Context) {
	// 获取管理员输入的用户名和密码
	adminname := c.Query("adminname")
	password := c.Query("password")

	// 验证管理员输入的用户名与密码是否正确
	DB := database.DB
	admin := database.Admin{}
	DB.Where("adminname = ?", adminname).First(&admin)
	if admin.Password == password {
		c.String(200, "correct")
	} else {
		c.String(200, "wrong")
	}
}

// GetAdminLogin 设置并保存Cookie，实现管理员登录
func GetAdminLogin(c *gin.Context) {
	adminname := c.Query("adminname")
	password := c.Query("password")

	// 再次验证密码，防止用户通过输入url直接登录
	DB := database.DB
	admin := database.Admin{}
	DB.Where("adminname = ?", adminname).First(&admin)
	if admin.Password == password {
		c.SetCookie("adminname", adminname, 24*60*60,
			"/admin", CookieDomain, false, false)

		// 记录管理员登录日志
		adminIP := c.ClientIP()
		logMsg := "{C-webLog}\tadmin:" + adminname + " log in succeed.\tIP:" + adminIP
		logger.RecordInfo(logFilename, logMsg)

		c.String(200, "success")
	} else {
		c.String(200, "fail")
	}
}

// AdminHomePage 管理员主页
func AdminHomePage(c *gin.Context) {
	c.HTML(200, "adminHomePage.html", gin.H{})
}

// GetAdminname 向前端返回该管理员的用户名
func GetAdminname(c *gin.Context) {
	adminname, _ := c.Cookie("adminname")
	c.String(200, adminname)
}

// AdminLogout 管理员退出当前账号并重新登录
func AdminLogout(c *gin.Context) {
	adminname, _ := c.Cookie("adminname")
	// 使cookie过期
	c.SetCookie("adminname", "", -1, "/admin", CookieDomain, false, false)

	// 记录日志
	adminIP := c.ClientIP()
	logMsg := "{C-webLog}\tadmin:" + adminname + " logout.\tIP:" + adminIP
	logger.RecordInfo(logFilename, logMsg)

	c.String(200, "success")
}

// WebLog 向前端返回网站的日志
func WebLog(c *gin.Context) {
	file, err := ioutil.ReadFile(logFilename)
	if err != nil {
		// 日志文件读取失败
		c.String(200, "Read webLog file fail")

		// 记录日志
		adminname, _ := c.Cookie("adminname")
		logMsg := "{C-webLog}\tadmin:" + adminname + " read webLog failed."
		logger.RecordError(logFilename, logMsg)
	} else {
		c.String(200, string(file))
	}
}

// UserList 向前端返回用户列表
func UserList(c *gin.Context) {
	db := database.DB
	users := []database.User{}
	db.Find(&users)
	c.JSON(200, gin.H{"userList": users})
}

// DeleteUser 根据用户名删除指定用户
func DeleteUser(c *gin.Context) {
	db := database.DB
	username := c.PostForm("username")
	user := database.User{}
	db.Where("username = ?", username).Delete(&user)

	// 记录日志
	adminname, _ := c.Cookie("adminname")
	adminIP := c.ClientIP()
	logMsg := "{C-webLog}\tadmin:" + adminname + " delete user:" + username + " .  adminIP:" + adminIP
	logger.RecordWarning(logFilename, logMsg)

	c.String(200, "success")
}

// AddNewUser 根据用户名和密码添加新用户
func AddNewUser(c *gin.Context) {
	username, password := c.PostForm("username"), c.PostForm("password")
	user := database.User{Username: username, Password: password}
	database.DB.Create(&user)

	// 记录日志
	adminname, _ := c.Cookie("adminname")
	adminIP := c.ClientIP()
	logMsg := "{C-webLog}\tadmin:" + adminname + " add a new user: [username:" +
		username + ", password:" + password + "]  adminIP:" + adminIP
	logger.RecordInfo(logFilename, logMsg)

	c.String(200, "success")
}

// ArticleList 向前端返回文章列表
func ArticleList(c *gin.Context) {
	db := database.DB
	articles := []database.Article{}
	db.Select("id", "articlename", "author", "addtime").Find(&articles)
	c.JSON(200, gin.H{"articles": articles})
}

// DeleteArticle 根据文章id删除指定文章
func DeleteArticle(c *gin.Context) {
	id := c.Query("articleId")
	article := database.Article{}
	db := database.DB
	db.Where("id = ?", id).Select("articlename").Find(&article)
	db.Where("id = ?", id).Delete(&article)

	// 记录日志
	adminname, _ := c.Cookie("adminname")
	adminIP := c.ClientIP()
	logMsg := "{C-webLog}\tadmin:" + adminname + " delete article:[id:" + id +
		", articleName: " + article.Articlename + "] .  adminIP:" + adminIP
	logger.RecordWarning(logFilename, logMsg)

	c.String(200, "success")
}

// MovieList 向前端返回电影列表
func MovieList(c *gin.Context) {
	movies := []database.Movie{}
	db := database.DB
	db.Find(&movies)
	c.JSON(200, gin.H{"movieList": movies})
}

// DeleteMovie 根据电影id删除指定电影
func DeleteMovie(c *gin.Context) {
	id := c.Query("movieId")
	movie := database.Movie{}
	db := database.DB
	db.Where("id = ?", id).Find(&movie)
	path := movie.Path
	coverPath := movie.Coverpath
	name := movie.Moviename
	db.Where("id = ?", id).Delete(&movie)

	// 删除存储的电影文件和电影封面的图片文件
	os.Remove(path)
	os.Remove(coverPath)

	// 记录日志
	adminname, _ := c.Cookie("adminname")
	adminIP := c.ClientIP()
	logMsg := "{C-webLog}\tadmin:" + adminname + " delete movie:[id:" + id +
		", name:" + name + "] .  adminIP:" + adminIP
	logger.RecordWarning(logFilename, logMsg)

	c.String(200, "success")
}

// AddNewMovie 添加新电影
func AddNewMovie(c *gin.Context) {
	movieName := c.PostForm("name")
	movieCover, _ := c.FormFile("cover")
	movieFile, _ := c.FormFile("file")
	if movieCover != nil && movieFile != nil {
		now := time.Now().Format("2006-01-02 15:04:05")
		movies := []database.Movie{}
		db := database.DB
		db.Order("id desc").Select("id").Find(&movies)

		// 设置新添加的电影的id，将其id值作为其封面文件和电影文件的文件名
		newId := 1
		if len(movies) > 0 {
			newId = movies[0].Id + 1
		}
		coverFileName := strings.Split(movieCover.Filename, ".")
		coverPath := "./static/videos/covers/" + strconv.Itoa(newId) + "." + coverFileName[len(coverFileName)-1]
		movieFileName := strings.Split(movieFile.Filename, ".")
		path := "./static/videos/" + strconv.Itoa(newId) + "." + movieFileName[len(movieFileName)-1]

		// 将电影封面和电影文件保存到static文件夹下的videos文件夹中
		c.SaveUploadedFile(movieCover, coverPath)
		c.SaveUploadedFile(movieFile, path)

		// 将新添加的电影记录进数据库
		newMovie := database.Movie{
			Id:        newId,
			Moviename: movieName,
			Addtime:   now,
			Coverpath: coverPath,
			Path:      path,
		}
		db.Create(&newMovie)

		// 记录日志
		adminname, _ := c.Cookie("adminname")
		adminIP := c.ClientIP()
		logMsg := "{C-webLog}\tadmin:" + adminname + " add a new movie:[" + "id:" + strconv.Itoa(newId) +
			", name:" + movieName + "] .  adminIP:" + adminIP
		logger.RecordInfo(logFilename, logMsg)

		c.String(200, "success")
	} else {
		c.String(200, "fail")
	}
}
