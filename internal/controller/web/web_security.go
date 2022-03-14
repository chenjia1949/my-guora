package web

import (
	"net/http"
	"time"

	"my-guora/internal/constant"
	"my-guora/internal/h"
	"my-guora/internal/model"
	"my-guora/internal/service/authorization"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginInterface struct
type LoginInterface struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

// SecurityLogin func
func SecurityLogin(c *gin.Context) {

	var s LoginInterface
	var ra int64
	var err error

	if err = c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, h.Response{Status: 500, Message: err.Error()})
		return
	}

	// check password and return user model
	ra, user, err := s.login()
	if err != nil {
		c.JSON(200, h.Response{Status: 500, Message: err.Error()})
		return
	}

	// gen a ss(signed string)
	ss, err := authorization.Gen(user)
	if err != nil {
		c.JSON(200, h.Response{Status: 500, Message: err.Error()})
		return
	}

	// set the cookies
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    constant.SSKEY,
		Value:   ss,
		Path:    "/",
		Expires: time.Time.AddDate(time.Now(), 0, 0, 7),
	})

	c.JSON(200, h.Response{
		Status:  200,
		Message: ra,
	})

	return
}

// login func
func (s *LoginInterface) login() (ra int64, user model.User, err error) {

	var u model.User

	u.Mail = s.Mail

	user, err = u.Get()
	if err != nil {
		ra = -1
		return
	}

	// check the password
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(s.Password)); err != nil {
		ra = -2
		return
	}

	ra = 1
	return
}

// SignInterface struct
type SignInterface struct {
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

// SecuritySign func
func SecuritySign(c *gin.Context) {

	var s SignInterface
	var ra int64
	var err error

	if err = c.ShouldBindJSON(&s); err != nil {
		c.JSON(200, h.Response{Status: 500, Message: err.Error()})
		return
	}

	if ra, err = s.sign(); err != nil {
		c.JSON(200, h.Response{Status: 500, Message: err.Error()})
	} else {
		c.JSON(200, h.Response{
			Status:  200,
			Message: ra,
		})
	}

	return

}

func (s *SignInterface) sign() (ra int64, err error) {
	var p model.Profile
	var u model.User

	p.Name = s.Name
	p.Desc = "This is " + s.Name
	u.Mail = s.Mail
	u.Password = s.Password

	u.Profile = p

	ra, err = u.Create()

	return
}

// SecurityLogout func
func SecurityLogout(c *gin.Context) {

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     constant.SSKEY,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	c.JSON(200, h.Response{
		Status:  200,
		Message: 1,
	})

	return
}
