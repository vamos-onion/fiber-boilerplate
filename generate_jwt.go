package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET not found in environment variables")
	}

	// Create claims
	claims := jwt.MapClaims{
		"uuid": "12956e54-503d-46f1-8b9b-7cf304fba601",
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days from now
		"iat":  time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Fatal("Error signing token:", err)
	}

	fmt.Println("\n=== JWT Token Generated ===")
	fmt.Println("\nSecret used:", jwtSecret)
	fmt.Println("\nClaims:")
	fmt.Printf("  uuid: %v\n", claims["uuid"])
	fmt.Printf("  exp: %v (expires at: %v)\n", claims["exp"], time.Unix(int64(claims["exp"].(int64)), 0))
	fmt.Printf("  iat: %v (issued at: %v)\n", claims["iat"], time.Unix(int64(claims["iat"].(int64)), 0))
	fmt.Println("\nToken:")
	fmt.Println(tokenString)
	fmt.Println("\nUse this in your Authorization header:")
	fmt.Printf("Authorization: Bearer %s\n", tokenString)
	fmt.Println()
}
