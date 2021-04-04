package models

import (
	"crypto/sha1"
	"crypto/tls"
	"errors"
	"github.com/EDDYCJY/go-gin-example/pkg/util"
	_ "github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"gopkg.in/mgo.v2/bson"
	"io"
	_ "log"
	_ "net/http"
	"regexp"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//type LoginUser struct {
//	Name     string `form:"name" bson:"name" json:"name" binding:"required"`
//	PassWord string `form:"password" bson:"password" json:"password,omitempty" binding:"required"`
//}

// CheckAuth checks if authentication information exists
func CheckAuth(username, password string) ([]User, error) {
	var users []User

	if e := accountDB.Find(bson.M{"name": username, "password": util.EncodeSha1(password)}).All(&users); e != nil {
		return nil, e
	}

	return users, nil
	//var auth Auth
	//err := db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error
	//if err != nil && err != gorm.ErrRecordNotFound {
	//	return false, err
	//}
	//
	//if auth.ID > 0 {
	//	return true, nil
	//}
	//
	//return false, nil
}
func Register(regUser RegUser) (err error) {
	//var regUser RegUser
	//if err := c.ShouldBind(&regUser); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"msg": "需要认证参数", "code": 400, "data": ""})
	//	return
	//}
	n, e := accountDB.Find(bson.M{"name": regUser.Name}).Count()
	if e != nil {
		return e
	}
	if n > 0 {
		//c.JSON(http.StatusBadRequest, gin.H{"msg": "用户已存在", "code": 400, "data": ""})
		return errors.New("用户已存在")
	}
	/**
	simple password encry
	*/
	hash := sha1.New()
	io.WriteString(hash, regUser.PassWord)
	regUser.PassWord = string(hash.Sum(nil))
	regUser.State = 0
	if e1 := accountDB.Insert(regUser); e1 != nil {
		return errors.New("注册失败,请重试")
	}

	if VerifyEmailFormat(regUser.EMail) {
		m := gomail.NewMessage()
		m.SetAddressHeader("From", "18736262687@163.com", "小书屋")
		m.SetHeader("To", regUser.EMail)
		m.SetHeader("Subject", "小书屋通知,不要回复")
		m.SetBody("text/html", " "+regUser.Name+",欢迎使用小书屋,这是开源项目非盈利...</br>如果觉得不错请点赞👉"+"<a href='https://github.com/leetomlee123/book'>项目地址</a><br/>")
		d := gomail.NewDialer("smtp.163.com", 465, "18736262687@163.com", "lx11427")
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		go func() {
			d.DialAndSend(m)
		}()
	}

	return nil
}
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
func ModifyPassword(user User) error {

	count, err2 := accountDB.Find(bson.M{"name": user.Name, "email": user.EMail}).Count()
	if err2 != nil {
		return err2
	}
	if count > 0 {
		hash := sha1.New()
		io.WriteString(hash, user.PassWord)
		user.PassWord = string(hash.Sum(nil))
		accountDB.Update(bson.M{"name": user.Name}, bson.M{"$set": bson.M{"password": user.PassWord}})
		return nil
	} else {
		return errors.New("用户信息有误")
	}
}
