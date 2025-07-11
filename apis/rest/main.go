package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olegshulyakov/go-briefly-bot/briefly"
	"github.com/olegshulyakov/go-briefly-bot/briefly/summarization"
	"github.com/olegshulyakov/go-briefly-bot/briefly/transcript"
	"github.com/olegshulyakov/go-briefly-bot/briefly/transcript/youtube"
)

// Web server port
const port = 8080

func main() {
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

	// Setup locale endpoints
	localeRouter := router.Group("/locale")
	localeRouter.GET("/message", getLocaleMessage)

	// Start server
	slog.Info("Starting server", "port", port)
	_ = router.Run(fmt.Sprintf("localhost:%d", port))
}

func postSummarizeText(c *gin.Context) {
	type Request struct {
		Text         string `json:"text"`
		LanguageCode string `json:"languageCode"`
	}

	var request Request
	if err := c.BindJSON(&request); err != nil {
		return
	}

	summary, err := summarization.SummarizeText(request.Text, request.LanguageCode)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"summary":      summary,
		"languageCode": request.LanguageCode,
	})
}

func getVideoInfo(c *gin.Context) {
	url := c.Query("url")

	videoInfo, err := youtube.GetYoutubeVideoInfo(url)

	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, videoInfo)
}

func getVideoTranscript(c *gin.Context) {
	url := c.Query("url")

	videoTranscript, err := transcript.GetYoutubeVideoTranscript(url)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, videoTranscript)
}

func getVideoSummarize(c *gin.Context) {
	c.String(http.StatusNotImplemented, "")
}

func getLocaleMessage(c *gin.Context) {
	messageID := c.Query("messageId")
	languageCode := c.DefaultQuery("languageCode", "en")

	message, err := briefly.Localize(languageCode, messageID)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, message)
}
