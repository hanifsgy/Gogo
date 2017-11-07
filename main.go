package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
)

var db *gorm.DB

func main() {

	router := gin.Default()

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/",createTodo)
		v1.GET("/", fetchAll)
		v1.GET("/:id",fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
	router.Run()
}


func init() {
	// open DB Connection
	var err error
	db, err = gorm.Open("mysql", "root:root@/todos?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate schema
	db.AutoMigrate(&todoModel{})
}


type (
	// todoModel describes a todoModel type
	todoModel struct {
		gorm.Model
		Title 		string 	`json:"title"`
		Completed 	int 	`json:"completed"`
	}
	//transformed todo represents a formatted todo
	transformedTodo struct {
		ID			uint	`json:"id"`
		Title		string 	`json:"title"`
		Completed	bool	`json:"completed"`
	}
)

// func create Todo
func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}
// func fetch All
func fetchAll( c *gin.Context) {
	var todos [] todoModel
	var _todos [] transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	// transform the todos for building a good response
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID:item.ID, Title:item.Title, Completed:completed})
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data":_todos})
	}
}

// fetch single todo
func fetchSingleTodo( c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")

	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

// func updateTodo

func updateTodo ( c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}
}

func deleteTodo ( c *gin.Context) {
	var todo todoModel
	todoID := c.Param("id")
	db.First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK,gin.H{"status":http.StatusOK, "message": "Todo deleted created successfully"})
}