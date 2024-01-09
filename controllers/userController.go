package controller

import (
	"context"
	"golang-restaurant-management/database"
	helper "golang-restaurant-management/helpers"
	"golang-restaurant-management/models"
	userServices "golang-restaurant-management/services"
	"golang-restaurant-management/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

// GetUsers godoc
// @Summary Get All Users
// @Description Get a list of users with pagination.
// @Tags User
// @Param recordPerPage query int false "Records per page (default is 10)"
// @Param page query int false "Page number (default is 1)"
// @Param startIndex query int false "Start index for slicing (default is 0)"
// @Success 200
// @Failure 500
// @Router /v1/users [GET]
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		limit, err := strconv.Atoi(c.Query("limit"))
		if err != nil || limit < 1 {
			limit = 10
		}

		page, err := strconv.Atoi(c.Query("page"))
		if err != nil || page < 1 {
			page = 1
		}

		result, err := userServices.GetUsers(ctx, limit, page)

		if err != nil {
			// Handle the error
			log.Println("Error during aggregation:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while aggregating documents"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

// GetUser godoc
// @Summary Get a User by ID
// @Description Get a user's details by their ID.
// @Tags User
// @Param user_id path string true "User ID"
// @Success 200
// @Failure 500
// @Router /v1/users/{user_id} [GET]
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		userId := c.Param("user_id")

		user, err := userServices.GetUserByID(ctx, userId)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
		}
		c.JSON(http.StatusOK, user)
	}
}

// SignUp godoc
// @Summary Sign Up a new User
// @Description Create a new user account.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.User true "User object to be created"
// @Success 200
// @Failure 400,500
// @Router /v1/users/signup [POST]
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		// convert the JSON data coming from request body to something that golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate the data based on user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// check the email and password has already been used by another user
		filter := bson.M{"email": user.Email, "phone": user.Phone}
		exists, err := userServices.UserExists(ctx, filter)
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		// hash password
		password := utils.HashPassword(*user.Password)
		user.Password = &password

		if exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or phone number already exsits"})
			return
		}

		// create some extra details for the user object
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		// insert this new user into the user collection
		createdUser, insertErr := userServices.CreateUser(ctx, user)
		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User item was not created"})
			return
		}
		defer cancel()

		// return status OK and send the result back
		c.JSON(http.StatusOK, createdUser)
	}
}

// Login godoc
// @Summary User Login
// @Description Authenticate user and generate access tokens.
// @Tags User
// @Accept  json
// @Produce  json
// @Param user body models.LoginRequest true "User object to be created"
// @Success 200
// @Failure 400,401,500
// @Router /v1/users/login [POST]
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.LoginRequest
		var foundUser models.User

		// convert the JSON data coming from request body to something that golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// find a user with that email
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found, login seems to be incorrect"})
			return
		}

		// verify the password
		passwordIsValid, msg := utils.VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		// generate tokens
		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		response := map[string]interface{}{
			"token":        token,
			"refreshToken": refreshToken,
		}

		for key, value := range utils.StructToMap(foundUser) {
			response[key] = value
		}

		// return statusOK with response
		c.JSON(http.StatusOK, response)
	}
}
