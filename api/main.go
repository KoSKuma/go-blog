package main

import (
	"net/http"

	"github.com/KoSKuma/go-blog/api/adapter"
	"github.com/KoSKuma/go-blog/api/entity"
	"github.com/KoSKuma/go-blog/api/repository"
	"github.com/KoSKuma/go-blog/api/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize
	r := gin.Default()
	mongoAdapter := adapter.MongoAdapter{}
	mongoAdapter.Setup("root", "root", "localhost", "27017", "blog", "posts") // Try using struct with all private fields -> need to call another function to setup
	databaseRepo := repository.Database{DBAdapter: &mongoAdapter}
	rabbitMQAdapter := adapter.RabbitMQAdapter{Username: "root", Password: "root", Host: "localhost", Port: "5672", Queue: "blog:logs"} // try using struct with all public fields
	logger := repository.Log{LogAdapter: &rabbitMQAdapter}
	usecase := usecase.PostUsecase{DBRepo: &databaseRepo, Logger: &logger}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/blog/posts", func(c *gin.Context) {
		posts, err := usecase.FindPosts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, posts)
	})

	r.POST("/blog/post", func(c *gin.Context) {
		var post entity.Post
		c.Bind(&post)
		id, err := usecase.CreatePost(post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"post_id": id,
		})
	})

	r.GET("/blog/post/:id", func(c *gin.Context) {
		id := c.Param("id")
		post, err := usecase.FindPost(id)
		if err != nil {
			errorCode := http.StatusInternalServerError
			if err.Error() == "Post not found" {
				errorCode = http.StatusNotFound
			}
			c.JSON(errorCode, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, post)
	})

	r.PATCH("/blog/post/:id", func(c *gin.Context) {
		var updatePost entity.PostUpdate
		id := c.Param("id")
		if err := c.BindJSON(&updatePost); err != nil {
			return
		}
		err := usecase.UpdatePost(id, updatePost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusAccepted, gin.H{
			"status": "success",
		})
	})

	r.DELETE("/blog/post/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := usecase.DeletePost(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
