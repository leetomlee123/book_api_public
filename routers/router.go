package routers

import (
	"github.com/EDDYCJY/go-gin-example/middleware/exception"
	"github.com/EDDYCJY/go-gin-example/middleware/logrus"
	"github.com/EDDYCJY/go-gin-example/routers/api"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/EDDYCJY/go-gin-example/docs"
	"github.com/EDDYCJY/go-gin-example/middleware/jwt"
	"github.com/EDDYCJY/go-gin-example/routers/api/v1"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	//r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(exception.ErrHandler())
	r.Use(logrus.LoggerToES())
	r.Use(func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,auth")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()

	})
	r.Any("/", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status":  "000000",
			"message": "service on line",
			"nick":    "hi,gays",
		})
	})
	apiv1 := r.Group("/v1")

	apiv1.GET("/health", v1.HealthCheck)
	apiv1.GET("/hot", v1.Hot)
	apiv1.GET("/info", v1.Notice)
	apiv1.POST("/login", api.GetAuth)
	apiv1.PATCH("/password", v1.ModifyPassword)
	apiv1.POST("/register", v1.Register)
	book := apiv1.Group("/book")
	book.GET("/detail/:id", v1.StoryInfo)
	book.PATCH("/process/:account/:bookId/:process", v1.Process)
	book.GET("/process/:account/:bookId", v1.GetProcess)
	book.GET("/rank/:type", v1.Rank)
	book.GET("/category", v1.Category)
	book.GET("/category/:category/:page/:size", v1.GetStroysByCategory)
	book.GET("/chapters/:id/:count", v1.GetStoryChatersById)
	book.DELETE("/chapter/:id", v1.DeleteChapte)
	book.GET("/chapter/:id", v1.GetChapterById)
	book.GET("/chapter/:id/async", v1.GetChapterByIdAsync)
	book.GET("/chapter/:id/reload", v1.ReloadChapterById)
	book.GET("/search", v1.Search)
	book.GET("/two/:name/:author", v1.GetBookByAuthorAndName)
	book.Use(jwt.JWT())
	{

		book.GET("/shelf", v1.Shelf)
		book.GET("/action/:bookId/:action", v1.ModifyShelf)
		book.GET("/freshToken", v1.FreshToken)

	}

	return r
}
