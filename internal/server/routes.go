package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/bupd/goth/gothic"
)

var userTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>User Information</title>
</head>
<body>
    <h1>User Info</h1>
    <pre>{{.}}</pre>
</body>
</html>
`

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	// from blogs
	// start the google oauth
	// r.GET("/auth/google/start", BeginGoogleAuth)
	// handle the call back
	// r.GET("auth/google/callback", OAuthCallback)

	r.GET("/auth/:provider/callback", s.authCallbackHandler)
	r.GET("/logout/:provider", logoutHandler)
	r.GET("/auth/:provider", authHandler)

	return r
}

func (s *Server) authCallbackHandler(c *gin.Context) {
	// provider := c.Param("provider")

	res := c.Writer
	req := c.Request

	// req = req.WithContext(context.WithValue(req.Context(), "provider", provider))

	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Auth error: %v", err))
		return
	}

	t, err := template.New("foo").Parse(userTemplate)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Template error: %v", err))
		return
	}

	if err := t.Execute(res, user); err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Execution error: %v", err))
	}
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func getAuthCallbackFunction(c *gin.Context) {
	// Extract provider from URL parameters
	// provider := c.Param("provider")

	// Extract the http.ResponseWriter and *http.Request from Gin context
	res := c.Writer
	req := c.Request

	// Update request context with provider information
	// req = req.WithContext(context.WithValue(req.Context(), "provider", provider))

	// Complete user authentication using Gothic
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		// Log the error and respond with an error message
		c.String(http.StatusInternalServerError, fmt.Sprintf("Auth error: %v", err))
		return
	}

	// Log user information for debugging
	fmt.Println("Authenticated User:", user)

	// Send a success response (replace with your desired response)
	c.JSON(http.StatusOK, gin.H{
		"message": "Authentication successful",
		"user":    user,
	})

	// Redirect to your desired URL
	c.Redirect(http.StatusFound, "http://localhost:8080/")
}

func logoutHandler(c *gin.Context) {
	provider := c.Param("provider")

	res := c.Writer
	req := c.Request

	req = req.WithContext(context.WithValue(req.Context(), "provider", provider))

	gothic.Logout(res, req)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func authHandler(c *gin.Context) {
	provider := c.Param("provider")

	res := c.Writer
	req := c.Request

	req = req.WithContext(context.WithValue(req.Context(), "provider", provider))

	gothUser, err := gothic.CompleteUserAuth(res, req)
	if err == nil {
		t, err := template.New("foo").Parse(userTemplate)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Template error: %v", err))
			return
		}

		if err := t.Execute(res, gothUser); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Execution error: %v", err))
		}
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func BeginGoogleAuth(c *gin.Context) {
	log.Println("kumaaru beginning")
	q := c.Request.URL.Query()
	q.Add("provider", "google")
	c.Request.URL.RawQuery = q.Encode()
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func OAuthCallback(c *gin.Context) {
	log.Println("kumaaru this is callback daa")
	q := c.Request.URL.Query()

	q.Add("provider", "google")

	c.Request.URL.RawQuery = q.Encode()

	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	log.Println(user.Provider)
	log.Println(user.RefreshToken, "RefreshToken")

	res, err := json.Marshal(user)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	jsonString := string(res)

	c.JSON(http.StatusAccepted, jsonString)
}
