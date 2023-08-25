package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/khaosles/giz/xerror"
	"github.com/khaosles/go-contrib/api"
)

/*
   @File: recorvery.go
   @Author: khaosles
   @Time: 2023/8/23 11:50
   @Desc:
*/

func GlobalErrorHander(c context.Context, ctx *app.RequestContext) {

	defer func() {
		if err := recover(); err != nil {
			e := err.(error)
			var apiError xerror.IError
			if errors.As(e, &apiError) {
				ctx.JSON(http.StatusOK, api.NewNo(apiError.Code(), apiError.Error()))
			} else {
				ctx.JSON(http.StatusOK, api.NewNo(40000, e.Error()))
			}

		}
	}()
	ctx.Next(c)
}
