package controller

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/yingxv/flashcard-go/src/db"
	"github.com/yingxv/flashcard-go/src/middleware"
)

// Controller 控制器
type Controller struct {
	validate *validator.Validate
	trans    *ut.Translator
	auth     *middleware.Auth
	mongo    *db.MongoClient
}

// NewController 工厂方法
func NewController(
	validate *validator.Validate,
	trans *ut.Translator,
	auth *middleware.Auth,
	mongo *db.MongoClient,
) *Controller {

	return &Controller{
		validate,
		trans,
		auth,
		mongo,
	}
}
