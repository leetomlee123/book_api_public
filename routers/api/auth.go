package api

import (
	"github.com/EDDYCJY/go-gin-example/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/EDDYCJY/go-gin-example/pkg/app"
	"github.com/EDDYCJY/go-gin-example/pkg/e"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// @Summary Get Auth
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}

	var user models.LoginUser
	if err := c.ShouldBind(&user); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"msg": "需要认证参数", "code": 400, "data": ""})
		appG.Response(http.StatusOK, 400, "")
		return
	}

	users, err := models.CheckAuth(user.Name, user.PassWord)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if len(users) == 0 {
		appG.Response(http.StatusOK, e.ERROR_LOGIN, nil)
		return
	}

	token, err := util.GenerateToken(user.Name)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS_LOGIN, map[string]string{
		"token": token,
		"email": users[0].EMail,
		"vip":   strconv.Itoa(int(users[0].Vip)),
	})
}
