package main

import (
	"net/http"
	"sync"
	"strconv"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users  = make(map[int]User)
	nextID = 1
	mu     sync.Mutex
)

func main() {
	r := gin.Default()

	// Define routes
	r.POST("/users", createUser)
	r.GET("/users/:id", getUser)
	r.PUT("/users/:id", updateUser)
	r.DELETE("/users/:id", deleteUser)

	r.Run(":8080")
}

func createUser(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	var user User

	if err := c.ShouldBindBodyWithJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":"Invalid input"})
		return
	}

	user.ID = nextID
	nextID++

	users[user.ID] = user
	c.JSON(http.StatusCreated, user)
}

func getUser(c *gin.Context) {
	idstr := c.Param("id")

	id, err := strconv.Atoi(idstr)

	if err != nil || id < 1{
		c.JSON(http.StatusBadRequest,gin.H{"error":"Invalid user ID"})
	}

	mu.Lock()
	defer mu.Unlock()

	user, exists := users[id]

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error":"User not found"})
		return
	}

	c.JSON(http.StatusOK,user)
}

func updateUser(c *gin.Context) {
	idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    mu.Lock()
    defer mu.Unlock()

    _, exists := users[id]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    user.ID = id
    users[id] = user

    c.JSON(http.StatusOK, user)
}

func deleteUser(c *gin.Context) {
	 idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil || id < 1 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
        return
    }

    mu.Lock()
    defer mu.Unlock()

    _, exists := users[id]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    delete(users, id)
    c.Status(http.StatusNoContent) // 204 No Content
}
