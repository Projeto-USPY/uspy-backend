package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

var allowedTypes = []binding.Binding{
	binding.Query,
	binding.Form,
	binding.FormPost,
	binding.FormMultipart,
}

func Bind(name string, data interface{}, bindingType binding.Binding) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ok := false
		for _, b := range allowedTypes {
			if b == bindingType {
				ok = true
			}
		}
		if !ok {
			ctx.AbortWithError(
				http.StatusInternalServerError,
				fmt.Errorf("function Bind only allows the following binding types: %v", allowedTypes),
			)
		}
		if err := ctx.ShouldBindWith(data, bindingType); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
		ctx.Set(name, data)
	}
}
