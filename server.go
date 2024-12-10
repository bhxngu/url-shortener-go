package main

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type RequestBody struct {
	URL string `json:"url"`
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

	u, err := url.Parse(body.URL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": body.URL,
	})
}

func main() {
	ginRouter := gin.Default()

	ginRouter.POST("/shorten", shorternerHandler)

	ginRouter.Run(":3001")
}
