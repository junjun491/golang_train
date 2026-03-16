package handler

import (
	"fmt"
	"log"
	"context"
	"net/http"

	"golang_train/backend-go/internal/db"
	"golang_train/backend-go/internal/model"

	"github.com/gin-gonic/gin"
)

var tasks = []model.Task{
	{ID: 1, Title: "study go"},
}

func GetTasks(c *gin.Context) {
	rows, err := db.DB.Query(context.Background(), `
		SELECT id, title, completed
		FROM tasks
		ORDER BY id ASC
	`)
	if err != nil {
		log.Println("GetTasks query error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	tasks := []model.Task{}

	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Completed); err != nil {
			log.Println("GetTasks scan error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		log.Println("GetTasks rows error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var newTask model.Task

	c.BindJSON(&newTask)

	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	c.JSON(201, newTask)
}

func UpdateTask(c *gin.Context) {
	var updatedTask model.Task

	c.BindJSON(&updatedTask)

	id := c.Param("id")

	for i, task := range tasks {
		if fmt.Sprint(task.ID) == id {
			tasks[i].Title = updatedTask.Title
			c.JSON(200, tasks[i])
			return
		}
	}

	c.JSON(404, gin.H{
		"error": "not found",
	})
}

func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	for i, task := range tasks {
		if fmt.Sprint(task.ID) == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(200, gin.H{
				"message": "deleted",
			})
			return
		}
	}

	c.JSON(404, gin.H{
		"error": "not found",
	})
}