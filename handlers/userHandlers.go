package handlers

import (
	"C-web/database"
	"github.com/gin-gonic/gin"
	"github.com/heejinzzz/logger"
	"strconv"
	"time"
)

var CookieDomain = "127.0.0.1"
var logFilename = "./webLog.log"

// CheckUserPassword 检查用户输入的用户名和密码，告知前端是否正确
func CheckUserPassword(c *gin.Context) {
	// 获取用户输入的用户名和密码
	username := c.Query("username")
	password := c.Query("password")

	// 验证用户输入的用户名与密码是否正确
	DB := database.DB
	user := database.User{}
	DB.Where("username = ?", username).First(&user)
	if user.Password == password {
		c.String(200, "correct")
	} else {
		c.String(200, "wrong")
	}
}

// GetUserLogin 设置并保存Cookie，实现用户登录
func GetUserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 再次验证密码，防止用户通过输入url直接登录
	DB := database.DB
	user := database.User{}
	DB.Where("username = ?", username).First(&user)
	if user.Password == password {
		c.SetCookie("username", username, 24*60*60, "/user", CookieDomain, false, false)

		// 记录用户登录日志
		userIP := c.ClientIP()
		logMsg := "{C-webLog}\tuser:" + username + " log in succeed.\tIP:" + userIP
		logger.RecordInfo(logFilename, logMsg)

		c.String(200, "success")
	} else {
		c.String(200, "fail")
	}
}

// UserHomePage 用户主页
func UserHomePage(c *gin.Context) {
	c.HTML(200, "userHomePage.html", gin.H{})
}

// GetUsername 向前端返回该用户的用户名
func GetUsername(c *gin.Context) {
	username, err := c.Cookie("username")
	if err == nil {
		c.String(200, username)
	} else {
		c.String(200, "")
	}
}

// UserLogout 用户退出当前账号并重新登录
func UserLogout(c *gin.Context) {
	username, _ := c.Cookie("username")
	// 使cookie失效
	c.SetCookie("username", "", -1, "/user", CookieDomain, false, false)

	// 记录日志
	userIP := c.ClientIP()
	logMsg := "{C-webLog}\tuser:" + username + " logout.  userIP:" + userIP
	logger.RecordInfo(logFilename, logMsg)

	c.String(200, "success")
}

// BlogList 博客园页面
func BlogList(c *gin.Context) {
	c.HTML(200, "blogList.html", gin.H{})
}

// GetBlogList 向前端返回博客列表
func GetBlogList(c *gin.Context) {
	blogs := []database.Article{}
	db := database.DB
	db.Select("id", "articlename", "author", "addtime").Find(&blogs)

	c.JSON(200, gin.H{"blogList": blogs})
}

// WatchBlog 阅读blog页面
func WatchBlog(c *gin.Context) {
	c.HTML(200, "watchBlog.html", gin.H{})
}

// GetBlog 根据文章id，向前端返回指定的博客文章
func GetBlog(c *gin.Context) {
	id := c.Query("id")
	article := database.Article{}
	database.DB.Where("id = ?", id).Find(&article)

	value := gin.H{
		"id":      id,
		"name":    article.Articlename,
		"author":  article.Author,
		"addtime": article.Addtime,
		"content": article.Content,
	}
	c.JSON(200, value)
}

// AddNewBlog 上传博客文章页面
func AddNewBlog(c *gin.Context) {
	c.HTML(200, "addNewBlog.html", gin.H{})
}

// SaveNewBlog 将前端传回的博客文章存入数据库
func SaveNewBlog(c *gin.Context) {
	author := c.PostForm("author")
	title := c.PostForm("title")
	content := c.PostForm("content")

	db := database.DB
	articles := []database.Article{}
	db.Order("id desc").Select("id").Find(&articles)

	// 获取当前时间作为上传时间
	now := time.Now().Format("2006-01-02 15:04:05")
	id := 1
	if len(articles) > 1 {
		id = articles[0].Id + 1
	}
	newBlog := database.Article{
		Id:          id,
		Articlename: title,
		Author:      author,
		Addtime:     now,
		Content:     content,
	}
	db.Create(newBlog)

	// 记录日志
	username, _ := c.Cookie("username")
	userIP := c.ClientIP()
	logMsg := "{C-webLog}\tuser:" + username + "upload a new blog: [id: " + strconv.Itoa(id) +
		", title:" + title + "] .  userIP:" + userIP
	logger.RecordInfo(logFilename, logMsg)

	c.String(200, "success")
}

// Cinema 影院页面
func Cinema(c *gin.Context) {
	c.HTML(200, "cinema.html", gin.H{})
}

// GetMovieList 向前端返回电影列表
func GetMovieList(c *gin.Context) {
	movies := []database.Movie{}
	db := database.DB
	db.Select("id", "moviename", "coverpath").Find(&movies)
	c.JSON(200, gin.H{"movieList": movies})
}

// WatchMovie 观影页面
func WatchMovie(c *gin.Context) {
	c.HTML(200, "watchMovie.html", gin.H{})
}

// GetMoviePath 根据电影id，向前端返回指定电影文件的存储路径和电影名
func GetMoviePath(c *gin.Context) {
	id := c.Query("id")
	movie := database.Movie{}
	db := database.DB
	db.Where("id = ?", id).Find(&movie)
	path := movie.Path[1:]
	c.JSON(200, gin.H{"name": movie.Moviename, "path": path})
}
