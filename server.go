package main

import (
	"crypto/md5"
	"encoding/base64"
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

	hash := md5.Sum([]byte(body.URL))
	shortened := base64.URLEncoding.EncodeToString(hash[:])[:6]

	ctx.JSON(http.StatusOK, gin.H{
		"shortened": shortened,
	})
}

func main() {
	ginRouter := gin.Default()

	ginRouter.POST("/shorten", shorternerHandler)

	ginRouter.Run(":3001")
}
