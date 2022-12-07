package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stephen10121/calenderapi/initializers"
	"github.com/stephen10121/calenderapi/models"
)

func Validate(c *gin.Context) {
	user2, _ := c.Get("user")
	user := user2.(models.User)

	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": gin.H{
			"userData": user,
		},
	})
}

func GoogleLogin(c *gin.Context) {
	var body struct {
		Token string `json:"token"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	if body.Token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/userinfo/v2/me", nil)
	req.Header.Add("Authorization", "Bearer "+body.Token)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Error authorizing user.",
		})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Error authorizing user.",
		})
		return
	}

	type GoogleRespond struct {
		Email          string `json:"email"`
		Family_name    string `json:"family_name"`
		Given_name     string `json:"given_name"`
		Id             string `json:"id"`
		Locale         string `json:"locale"`
		Name           string `json:"name"`
		Picture        string `json:"picture"`
		Verified_email bool   `json:"verified_email"`
	}

	defer resp.Body.Close()
	body2, err := io.ReadAll(resp.Body)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "Error authorizing user.",
		})
		return
	}

	var data GoogleRespond

	if err := json.Unmarshal(body2, &data); err != nil {
		fmt.Println("Error")
	}

	if data.Id == "" {
		c.JSON(http.StatusOK, gin.H{
			"error": "Invalid Token",
		})
		return
	}

	var user models.User
	initializers.DB.First(&user, "goog_id = ?", data.Id)

	if user.ID == 0 {
		// Users first time loggin in.

		user := models.User{Email: data.Email, GoogId: data.Id, FirstName: data.Given_name, LastName: data.Family_name, FullName: data.Name, Groups: "", PendingGroups: "", Locale: data.Locale, Picture: data.Picture, VerifiedEmail: data.Verified_email}
		result := initializers.DB.Create(&user)

		if result.Error != nil {
			fmt.Println(result.Error)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create user",
			})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create token",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"error": "",
			"data": gin.H{
				"userData": user,
				"token":    tokenString,
			},
		})
		return
	}

	var fullName string
	if user.FullName != data.Name {
		fullName = data.Name
		initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("full_name", fullName)
	}

	var firstName string
	if user.FirstName != data.Given_name {
		firstName = data.Given_name
		initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("first_name", firstName)
	}

	var lastName string
	if user.LastName != data.Family_name {
		lastName = data.Family_name
		initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("last_name", lastName)
	}

	var picture string
	if user.Picture != data.Picture {
		picture = data.Picture
		initializers.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("picture", picture)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": "",
		"data": gin.H{
			"userData": user,
			"token":    tokenString,
		},
	})
}
