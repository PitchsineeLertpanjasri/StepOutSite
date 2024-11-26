package middlewares

import (
	"fmt"
	// "log"
	"net/http"
	"os"
	"stepoutsite/domain/entities"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func SetJWtHeaderHandler() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(os.Getenv("JWT_SECRET_KEY")),
			JWTAlg: jwtware.HS256,
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "Unauthorization Token."})
		},
	})
}

type TokenDetails struct {
	Token     *string `json:"token"`
	StudentID    string  `json:"student_id"`
	ExpiresIn *int64  `json:"exp"`
}

func DecodeJWTToken(ctx *fiber.Ctx) (*TokenDetails, error) {

	td := &TokenDetails{
		Token: new(string),
	}

	if ctx.Cookies("jwt") == "" {
		return nil, ctx.Status(http.StatusUnauthorized).SendString("Unauthorization Token.")
	}
	
	token, status := ctx.Locals("user").(*jwt.Token)

	if !status {
		return nil, ctx.Status(http.StatusUnauthorized).SendString("Unauthorization Token.")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ctx.Status(http.StatusUnauthorized).SendString("Unauthorization Token.")
	}

	for key, value := range claims {
		if key == "student_id" || key == "sub" {
			td.StudentID = value.(string)
		}
		if key == "uid" {
			td.StudentID = value.(string)
		}
	}
	*td.Token = token.Raw
	return td, nil
}

func GenerateJWTToken(studentID string) (*TokenDetails, error) {
	now := time.Now().UTC()

	td := &TokenDetails{
		ExpiresIn: new(int64),
		Token:     new(string),
	}
	*td.ExpiresIn = now.Add(time.Hour * 6).Unix()
	td.StudentID = studentID

	SigningKey := []byte(os.Getenv("JWT_SECRET_KEY"))

	atClaims := make(jwt.MapClaims)
	atClaims["student_id"] = studentID
	atClaims["exp"] = time.Now().Add(time.Hour * 6).Unix()
	atClaims["iat"] = time.Now().Unix()
	atClaims["nbf"] = time.Now().Unix()

	// log.Println("New claims: ", atClaims)

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString(SigningKey)
	if err != nil {
		return nil, fmt.Errorf("create: sign token: %w", err)
	}
	*td.Token = token
	
	return td, nil
}