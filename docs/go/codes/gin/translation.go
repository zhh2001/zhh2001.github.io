package main

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
)

type SignUpInfo struct {
	Username   string `json:"username" binding:"required,min=3,max=20"`
	Password   string `json:"password" binding:"required,min=8,max=20"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
	Email      string `json:"email" binding:"required,email"`
	Age        uint   `json:"age" binding:"lte=120"`
}

func initTranslator() (ut.Translator, error) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return nil, errors.New("unexpected validator engine")
	}
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	locale := zh.New()
	trans, ok := ut.New(locale, locale).GetTranslator("zh")
	if !ok {
		return nil, errors.New("Chinese translator not found")
	}
	if err := zhtranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}
	return trans, nil
}

func main() {
	trans, err := initTranslator()
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	router.POST("/signUp", signUp(trans))
	if err := router.Run("127.0.0.1:8000"); err != nil {
		log.Fatal(err)
	}
}

func signUp(trans ut.Translator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var info SignUpInfo
		if err := c.ShouldBindJSON(&info); err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				messages := make(map[string]string, len(validationErrors))
				for _, fieldError := range validationErrors {
					messages[fieldError.Field()] = fieldError.Translate(trans)
				}
				c.JSON(http.StatusBadRequest, gin.H{"error": messages})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "JSON 格式或字段类型不正确"})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "验证通过"})
	}
}
