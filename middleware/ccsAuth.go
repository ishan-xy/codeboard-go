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

	utils "github.com/ItsMeSamey/go_utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateAuth(c *fiber.Ctx) error {
	var tokenBody struct {
		Token string `json:"token"`
	}

	// Parse JSON body
	if err := utils.WithStack(c.BodyParser(&tokenBody)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":          "Invalid parameters",
			"detailed_error": err.Error(),
		})
	}

	secret := os.Getenv("SECRET")

	// Debug print for token
	fmt.Println("Received Token:", tokenBody.Token)

	token, err := jwt.Parse(tokenBody.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err = utils.WithStack(err); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Debug print for claims
		fmt.Println("Token Claims:")
		for k, v := range claims {
			fmt.Printf("%s: %v\n", k, v)
		}

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token expired",
			})
		}

		encryptedData, ok := claims["ex"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing encrypted data in token",
			})
		}

		// Debug print for encrypted data
		fmt.Println("Encrypted Data:", encryptedData)

		decryptedData, err := decrypt(encryptedData, secret)
		if err = utils.WithStack(err); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to decrypt token payload",
			})
		}

		// Store decrypted data in Fiber's locals
		c.Locals("data", decryptedData)

		return c.Next()
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}
}

func decrypt(encryptedData, key string) (map[string]interface{}, error) {
	fmt.Println("Starting Decryption Process")
	fmt.Println("Encrypted Data Length:", len(encryptedData))
	fmt.Println("Key Length:", len(key))

	if len(key) < 32 {
		return nil, errors.New("key must be at least 32 characters long for AES-256 encryption")
	}

	data, err := hex.DecodeString(encryptedData)
	if err = utils.WithStack(err); err != nil {
		fmt.Println("Hex Decode Error:", err)
		return nil, fmt.Errorf("failed to decode encrypted data: %w", err)
	}

	fmt.Println("Decoded Data Length:", len(data))

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]
	keyBytes := []byte(key[:32])

	fmt.Println("IV:", hex.EncodeToString(iv))
	fmt.Println("Ciphertext Length:", len(ciphertext))
	fmt.Println("Key:", hex.EncodeToString(keyBytes))

	block, err := aes.NewCipher(keyBytes)
	if err = utils.WithStack(err); err != nil {
		fmt.Println("Cipher Creation Error:", err)
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	fmt.Println("Decrypted Bytes (before unpad):", hex.EncodeToString(decrypted))

	decrypted, err = unpadPKCS7(decrypted)
	if err = utils.WithStack(err); err != nil {
		fmt.Println("Unpad Error:", err)
		return nil, fmt.Errorf("failed to unpad decrypted data: %w", err)
	}

	fmt.Println("Decrypted Bytes (after unpad):", string(decrypted))

	var result map[string]interface{}
	if err := json.Unmarshal(decrypted, &result); err != nil {
		err = utils.WithStack(err)
		fmt.Println("JSON Unmarshal Error:", err)
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Println("Decrypted Data:")
	for k, v := range result {
		fmt.Printf("%s: %v\n", k, v)
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