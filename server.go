package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type RequestBody struct {
	URL string `json:"url"`
}

func shorternerHandler(ctx *gin.Context, db *sql.DB) {
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

	row := db.QueryRow("SELECT shortened FROM shortened_urls WHERE url = $1", body.URL)
	var shortened string
	err = row.Scan(&shortened)
	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"shortened": shortened,
		})
		return
	}

	hash := md5.Sum([]byte(body.URL))
	shortened = base64.URLEncoding.EncodeToString(hash[:])[:6]
	_, err = db.Exec("INSERT INTO shortened_urls (url, shortened) VALUES ($1, $2)", body.URL, shortened)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"shortened": shortened,
	})
}
func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	connStr := os.Getenv("CONNECTION_STRING")

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ginRouter := gin.Default()
	ginRouter.POST("/shorten", func(ctx *gin.Context) {
		shorternerHandler(ctx, db)
	})

	ginRouter.Run(":3001")
}
