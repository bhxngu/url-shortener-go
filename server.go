package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type RequestBody struct {
	URL string `json:"url"`
}

func main() {
	ginRouter := gin.Default()

	ginRouter.POST("/shorten", shorternerHandler)

	ginRouter.Run(":3001")
}

func shorternerHandler(ctx *gin.Context) {
	var body RequestBody
	binding_type := binding.JSON

	if error := ctx.ShouldBindWith(&body, binding_type); error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": body.URL,
	})
}
