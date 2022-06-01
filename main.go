package main

import (
	"C-web/handlers"
	"C-web/middlewares"
	"github.com/gin-gonic/gin"
	"html/template"
)

func main() {
	r := gin.Default()

	// 导入静态文件，其中第一个参数为路由，第二个参数为静态文件存放路径
	r.Static("/static", "./static")

	// 将自定义的函数注册进html模板，这一步要放在导入html模板之前
	r.SetFuncMap(template.FuncMap{})

	// 导入html模板，这一步要放在注册路由之前
	r.LoadHTMLGlob("./htmlFiles/**/*")

	// 注册全局中间件
	r.Use(
		middlewares.CrossDomain,
		middlewares.PingLog,
	)

	// 注册各组路由
	// basic路由分组
	basicRouters := r.Group("")
	{
		basicRouters.GET("/", handlers.HomePage)
	}

	// user路由分组
	userRouters := r.Group("/user")
	userRouters.Use(middlewares.CheckUserLogin) // 用中间件检查用户是否已登录
	{
		userRouters.GET("/", handlers.UserHomePage)
		userRouters.GET("/getUsername", handlers.GetUsername)
		userRouters.GET("/logout", handlers.UserLogout)
		userRouters.GET("/blogList", handlers.BlogList)
		userRouters.GET("/getBlogList", handlers.GetBlogList)
		userRouters.GET("/watchBlog", handlers.WatchBlog)
		userRouters.GET("/getBlog", handlers.GetBlog)
		userRouters.GET("/addNewBlog", handlers.AddNewBlog)
		userRouters.POST("/saveNewBlog", handlers.SaveNewBlog)
		userRouters.GET("/cinema", handlers.Cinema)
		userRouters.GET("/getMovieList", handlers.GetMovieList)
		userRouters.GET("/watchMovie", handlers.WatchMovie)
		userRouters.GET("getMoviePath", handlers.GetMoviePath)
	}
	r.GET("/checkUserPassword", handlers.CheckUserPassword) // 这一步操作不能用中间件，要放在括号外面，否则会陷入死循环
	r.GET("/getUserLogin", handlers.GetUserLogin)           // 这一步操作不能用中间件，要放在括号外面，原因同上

	// admin路由分组
	adminRouters := r.Group("/admin")
	adminRouters.Use(middlewares.CheckAdminLogin) // 用中间件检查管理员是否已登录
	{
		adminRouters.GET("/", handlers.AdminHomePage)
		adminRouters.GET("/getAdminname", handlers.GetAdminname)
		adminRouters.GET("/adminLogout", handlers.AdminLogout)
		adminRouters.GET("/webLog", handlers.WebLog)
		adminRouters.GET("/getUserList", handlers.UserList)
		adminRouters.POST("/deleteUser", handlers.DeleteUser)
		adminRouters.POST("/addNewUser", handlers.AddNewUser)
		adminRouters.GET("/getArticleList", handlers.ArticleList)
		adminRouters.GET("/deleteArticle", handlers.DeleteArticle)
		adminRouters.GET("/getMovieList", handlers.MovieList)
		adminRouters.GET("/deleteMovie", handlers.DeleteMovie)
		adminRouters.POST("/addNewMovie", handlers.AddNewMovie)
	}
	r.GET("/checkAdminPassword", handlers.CheckAdminPassword) // 这一步操作不能用中间件，要放在括号外面，原因同上
	r.GET("/getAdminLogin", handlers.GetAdminLogin)           // 这一步操作不能用中间件，要放在括号外面，原因同上

	r.Run(":9090")
}
