package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/olegshulyakov/go-briefly-bot/briefly"
)

func main() {
	// Load configuration
	_, err := briefly.LoadConfiguration()
	if err != nil {
		briefly.Error("Failed to load config: %v", "error", err)
		os.Exit(1)
	}

	// Set port
	port := 8080

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Go Briefly API")
	})

	// Setup summary endpoints
	summaryRouter := router.Group("/summary")
	summaryRouter.POST("/text", postSummarizeText)

	// Setup video endpoints
	videoRouter := router.Group("/video")
	videoRouter.GET("/info", getVideoInfo)
	videoRouter.GET("/transcript", getVideoTranscript)
	videoRouter.GET("/summarize", getVideoSummarize)

	// Start server
	briefly.Info("Starting server", "port", port)
	router.Run(fmt.Sprintf("localhost:%d", port))
}

func postSummarizeText(c *gin.Context) {
	text := c.GetString("text")
	languageCode := c.GetString("languageCode")

	summary, err := briefly.SummarizeText(text, languageCode)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"summary":      summary,
		"languageCode": languageCode,
	})
}

func getVideoInfo(c *gin.Context) {
	url := c.Query("url")

	videoInfo, err := briefly.GetYoutubeVideoInfo(url)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"videoId":  videoInfo.ID,
		"language": videoInfo.Language,
		"uploader": videoInfo.Uploader,
		"title":    videoInfo.Title,
		"thumbnai": videoInfo.Thumbnail,
	})
}

func getVideoTranscript(c *gin.Context) {
	url := c.Query("url")
	languageCode := c.DefaultQuery("languageCode", "en")

	transcript, err := briefly.GetYoutubeTranscript(url, languageCode)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"videoId":    url,
		"language":   languageCode,
		"transcript": transcript,
	})
}

func getVideoSummarize(c *gin.Context) {
	c.String(http.StatusNotImplemented, "")
}
