package controllers

import (
	"net/http"
	"tusk/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func (u *UserController) Login(c *gin.Context) {
	user := models.User{}

	errBinJson := c.ShouldBindJSON(&user)

	if errBinJson != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBinJson.Error()})
		return

	}

	password := user.Password

	errDB := u.DB.Where("email=?", user.Email).Take(&user).Error

	if errDB != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username / Password salah"})
		return
	}

	errHash := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if errHash != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username / Password ppp"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *UserController) CreateAccount(c *gin.Context) {
	user := models.User{}

	errBinJson := c.ShouldBindJSON(&user)

	if errBinJson != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBinJson.Error()})
		return

	}

	emailExist := u.DB.Where("email=?", user.Email).First(&user).RowsAffected != 0

	if emailExist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email Sudah Ada"})
		return
	}
	hashedPasswordBytes, errHash := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

	if errHash != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errHash.Error()})
		return
	}

	user.Password = string(hashedPasswordBytes)
	user.Role = "employee"

	errDb := u.DB.Create(&user).Error

	if errDb != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDb.Error()})
		return

	}

	c.JSON(http.StatusOK, user)
}

func (u *UserController) Delete(c *gin.Context) {
	id := c.Param("id")

	// Check if user with the given ID exists
	var user models.User
	errDb := u.DB.First(&user, id).Error
	if errDb != nil {
		// If user is not found or there is a DB error, return appropriate response
		if errDb == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errDb.Error()})
		}
		return
	}

	// Delete the user if found
	errDb = u.DB.Delete(&user).Error
	if errDb != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDb.Error()})
		return
	}

	c.JSON(http.StatusOK, "Deleted")
}

func (u *UserController) GetEmployee(c *gin.Context) {

	users := []models.User{}

	errDb := u.DB.Select("id,name").Where("role=?", "employee").Find(&users).Error

	if errDb != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": errDb.Error()})
	}

	c.JSON(http.StatusOK, users	)
}
