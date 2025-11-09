package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func InitTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := zh.New()
		enT := en.New()
		uni := ut.New(zhT, zhT, enT)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return errors.New("translator not found")
		}
		switch locale {
		case "zh":
			if err := zh_translations.RegisterDefaultTranslations(v, trans); err != nil {
				return err
			}
		case "en":
			if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
				return err
			}
		}
	}
	return nil
}

var trans ut.Translator

func main() {
	if err := InitTrans("zh"); err != nil {
		fmt.Println(err)
	}
	router := gin.Default()
	router.POST("/signUp", signUp)
	router.Run("localhost:8000")
}

func signUp(context *gin.Context) {
	var info SignUpInfo
	if err := context.ShouldBindJSON(&info); err != nil {
		errs, _ := err.(validator.ValidationErrors)
		context.JSON(http.StatusBadRequest, gin.H{"error": errs.Translate(trans)})
		return
	}
	context.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
}
