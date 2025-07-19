// Package main implements the Go Briefly API server, which provides services for
// text summarization, video information extraction, transcript generation, and
// localization support.
//
// The API exposes the following endpoints:
//   - /ping - Simple health check endpoint
//   - /summary/text - POST endpoint for summarizing text content
//   - /video/info - GET endpoint for retrieving video metadata
//   - /video/transcript - GET endpoint for retrieving video transcripts
//   - /video/summarize - GET endpoint for summarizing video content (not implemented)
//   - /locale/message - GET endpoint for localized message retrieval
//
// The server runs on port 8080 by default and uses Gin as the web framework.
package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/olegshulyakov/go-briefly-bot/lib"
	"github.com/olegshulyakov/go-briefly-bot/lib/loaders/video/youtube/ytdlp"
	"github.com/olegshulyakov/go-briefly-bot/lib/transformers/summarization"
)

// Web server port.
const port = 8080

// main initializes and starts the HTTP server with all configured routes.
// It sets up endpoints for ping checks, text summarization, video information,
// transcript retrieval, video summarization, and localized messages.
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

// postSummarizeText handles POST requests to summarize text content.
// It expects a JSON body with "text" and "languageCode" fields.
// Returns the summary and original language code in JSON format.
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

// getVideoInfo retrieves metadata for a YouTube video.
// Expects a "url" query parameter with the YouTube video URL.
// Returns video information in JSON format.
func getVideoInfo(c *gin.Context) {
	url := c.Query("url")

	videoInfo, err := ytdlp.New().VideoInfo(url)

	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, videoInfo)
}

// getVideoTranscript retrieves the transcript for a video.
// Expects a "url" query parameter with the video URL.
// Returns the transcript in JSON format.
func getVideoTranscript(c *gin.Context) {
	url := c.Query("url")

	videoTranscript, err := ytdlp.New().Transcript(url)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, videoTranscript)
}

// getVideoSummarize is a placeholder endpoint for video summarization.
// Currently returns 501 Not Implemented status.
func getVideoSummarize(c *gin.Context) {
	c.String(http.StatusNotImplemented, "")
}

// getLocaleMessage retrieves a localized message by ID.
// Expects "messageId" query parameter and optional "languageCode" (defaults to "en").
// Returns the localized message in JSON format.
func getLocaleMessage(c *gin.Context) {
	messageID := c.Query("messageId")
	languageCode := c.DefaultQuery("languageCode", "en")

	message, err := lib.Localize(languageCode, messageID)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.IndentedJSON(http.StatusOK, message)
}
