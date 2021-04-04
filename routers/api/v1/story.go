package v1

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	"github.com/EDDYCJY/go-gin-example/service/story_service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func StoryInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	param := c.Param("id")
	data, err := models.StoryInfo(param)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}
func Rank(c *gin.Context) {
	appG := app.Gin{C: c}
	tpe, _ := strconv.Atoi(c.Param("type"))
	data, err := models.StoryRank(tpe)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}
func Process(c *gin.Context) {
	appG := app.Gin{C: c}
	account := c.Param("account")
	bookId := c.Param("bookId")
	process := c.Param("process")
	_ = models.PersonBookReadProcess(account, bookId, process)

	appG.Response(http.StatusOK, e.SUCCESS, "")

}
func GetProcess(c *gin.Context) {
	appG := app.Gin{C: c}
	account := c.Param("account")
	bookId := c.Param("bookId")
	process := models.GetReadProcess(account, bookId)

	appG.Response(http.StatusOK, e.SUCCESS, process)

}
func Shelf(c *gin.Context) {
	appG := app.Gin{C: c}
	name, exists := c.Get("username")
	if !exists {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	account, err2 := models.GetAccountByName(name.(string))
	if err2 != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	data, err := models.Shelf(account)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
func Hot(c *gin.Context) {
	appG := app.Gin{C: c}
	data, err := story_service.StoryRankService()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
func HealthCheck(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, "OK")
}
func ModifyShelf(c *gin.Context) {

	appG := app.Gin{C: c}
	name, exists := c.Get("username")
	if !exists {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	account, err2 := models.GetAccountByName(name.(string))
	if err2 != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	bookId := c.Param("bookId")
	action := c.Param("action")
	err := models.ModifyShelf(bookId, action, account)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, "")

}
func FreshToken(c *gin.Context) {

	appG := app.Gin{C: c}
	name, exists := c.Get("username")
	if !exists {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	token, err := util.GenerateToken(name.(string))
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})

}
func Category(c *gin.Context) {
	appG := app.Gin{C: c}
	data, err := story_service.CategoryService()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}
func GetStroysByCategory(c *gin.Context) {
	appG := app.Gin{C: c}
	page, e1 := strconv.Atoi(c.Param("page"))
	cate := c.Param("category")
	size, e2 := strconv.Atoi(c.Param("size"))

	if e1 != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	if e2 != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	data, err := models.GetStoryByCategoryWithPage(page, cate, size)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
func GetStoryChatersById(c *gin.Context) {
	id := c.Param("id")
	count, _ := strconv.Atoi(c.Param("count"))

	appG := app.Gin{C: c}
	data, err := models.GetStoryChapters(id, count)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
func DeleteChapte(c *gin.Context) {
	id := c.Param("id")

	appG := app.Gin{C: c}
	err := models.DeleteChapterById(id)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, "")
}
func GetChapterById(c *gin.Context) {
	id := c.Param("id")

	appG := app.Gin{C: c}
	data, err := story_service.GetChapterByIdService(id)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

func GetChapterByIdAsync(c *gin.Context) {
	id := c.Param("id")

	appG := app.Gin{C: c}
	data, err := story_service.GetChapterByIdServiceAsync(id)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)

}

func ReloadChapterById(c *gin.Context) {
	id := c.Param("id")

	appG := app.Gin{C: c}
	data, err := story_service.ReloadChapterByIdService(id)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
func GetBookByAuthorAndName(c *gin.Context) {
	name := c.Param("name")
	author := c.Param("author")

	appG := app.Gin{C: c}
	data, err := models.GetBookByAuthorAndName(author, name)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
func Search(c *gin.Context) {
	key := c.Query("key")
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	appG := app.Gin{C: c}
	data, err := models.Search(key, page, size)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, data)
}
func Register(c *gin.Context) {
	var regUser models.RegUser
	appG := app.Gin{C: c}
	if err := c.ShouldBind(&regUser); err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	regUser.Vip = 0
	if err := models.Register(regUser); err != nil {
		appG.Response(http.StatusOK, e.DUPLIT_USERNAME, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, "")

}
func ModifyPassword(c *gin.Context) {
	appG := app.Gin{C: c}
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		appG.Response(http.StatusOK, e.MODIFY_PASSWORD_FAILED, nil)
		return
	}

	err := models.ModifyPassword(user)
	if err != nil {
		appG.Response(http.StatusOK, e.MODIFY_PASSWORD_FAILED, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, "")
}
func Notice(c *gin.Context) {

	appG := app.Gin{C: c}
	err, infos := models.Notice()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, infos)
}
