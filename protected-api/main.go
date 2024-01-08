package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hashicorp/cap/jwt"
)

// Object represents the data structure of the stored object
type Object struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
}

var (
	objects          = make(map[string]Object)
	oidcIssuerURL    string
	oidcClientID     string
	oidcAudience     string
	tokenKeyset      jwt.KeySet
)

func main() {
	loadConfig()

	initializeOIDCVerifier()

	router := gin.Default()

	router.Use(authMiddleware)
	router.SetTrustedProxies(nil)

	router.POST("/object", createObject)
	router.GET("/object/:id", getObject)
	router.GET("/objects", getAllObjects)

	router.Run(":5010")
}

func loadConfig() {
	oidcIssuerURL = os.Getenv("OIDC_ISSUER_URL")
	if oidcIssuerURL == "" {
		panic("OIDC_ISSUER_URL environment variable is required")
	}

	oidcClientID = os.Getenv("OIDC_CLIENT_ID")
	if oidcClientID == "" {
		panic("OIDC_CLIENT_ID environment variable is required")
	}

	oidcAudience = os.Getenv("OIDC_AUDIENCE")
	if oidcAudience == "" {
		panic("OIDC_AUDIENCE environment variable is required")
	}
}

func initializeOIDCVerifier() {
	ctx := context.Background()

	keyset, err := jwt.NewOIDCDiscoveryKeySet(ctx, oidcIssuerURL, "")
	if err != nil {
		log.Fatal(err)
	}
	tokenKeyset = keyset
}

// authMiddleware is the middleware for JWT authentication
func authMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	tokenStrings := strings.Split(tokenString, " ")
	if len(tokenStrings) != 2 || tokenStrings[0] != "Bearer" {
	 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	 	c.Abort()
		return
	}
	token := tokenStrings[1]

	ctx := context.Background()
	claims, err := tokenKeyset.VerifySignature(ctx, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
	 	c.Abort()
	 	return
	}
	fmt.Printf("Claim: %+v\n", claims)

	// Additional checks, if needed
	// claims, ok := token.Claims.(jwt.MapClaims)
	// if !ok {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
	// 	c.Abort()
	// 	return
	// }

	// // Validate Audience
	// if aud, ok := claims["aud"].(string); !ok || aud != oidcAudience {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
	// 	c.Abort()
	// 	return
	// }

	// Attach user data to the request context if needed
	//c.Set("user", claims)

	// Continue with the next middleware or route handler
	c.Next()
}

// createObject handles the creation of objects
func createObject(c *gin.Context) {
	var obj Object
	if err := c.ShouldBindJSON(&obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique ID (UUID)
	id := uuid.New().String()

	// Save the object
	obj.ID = id
	objects[id] = obj

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// getObject handles the retrieval of objects by ID
func getObject(c *gin.Context) {
	id := c.Param("id")
	obj, exists := objects[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	c.JSON(http.StatusOK, obj)
}

// getObject handles the retrieval of objects by ID
func getAllObjects(c *gin.Context) {
	objIds := make([]string, len(objects))

	i := 0
	for k := range objects {
		objIds[i] = k
		i++
	}
	c.JSON(http.StatusOK, objIds)
}
