package middleware

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

func ValidateAuth(c *gin.Context) {

	// Read the raw body first
    body, err := c.GetRawData()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Cannot read request body",
        })
        c.Abort()
        return
    }
    
    // Print out the raw body for debugging
    fmt.Println("Raw body:", string(body))

    var tokenBody struct {
        Token string `json:"token"`
    }
    if err := c.ShouldBindBodyWithJSON(&tokenBody); err != nil {
		fmt.Println("Binding error:", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid parameters",
			"detailed_error": err.Error(),
        })
        c.Abort() // Prevent further handlers from running.
        return
    }
    secret := os.Getenv("SECRET")

    token, err := jwt.Parse(tokenBody.Token, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(secret), nil
    })

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{})
        c.Abort() // Stop further handlers.
        return
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        if float64(time.Now().Unix()) > claims["exp"].(float64) {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Token expired",
            })
            c.Abort()
            return
        }

        encryptedData, ok := claims["ex"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": "Missing encrypted data in token",
            })
            c.Abort()
            return
        }

        decryptedData, err := decrypt(encryptedData, secret)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to decrypt token payload",
            })
            c.Abort()
            return
        }

        c.Set("data", decryptedData)
        c.Next() // Pass control to the next handler.
    } else {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "Invalid token claims",
        })
        c.Abort()
        return
    }
}


func decrypt(encryptedData, key string) (map[string]interface{}, error) {
	if len(key) < 32 {
		return nil, errors.New("key must be at least 32 characters long for AES-256 encryption")
	}

	data, err := hex.DecodeString(encryptedData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]
	keyBytes := []byte(key[:32])

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted, err = unpadPKCS7(decrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to unpad decrypted data: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(decrypted, &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}

func unpadPKCS7(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return nil, errors.New("invalid padding size")
	}

	for _, v := range data[len(data)-padding:] {
		if int(v) != padding {
			return nil, errors.New("invalid padding")
		}
	}

	return data[:len(data)-padding], nil
}
