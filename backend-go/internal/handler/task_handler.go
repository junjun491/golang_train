package handler

import (
	"fmt"

	"golang_train/backend-go/internal/model"

	"github.com/gin-gonic/gin"
)

var tasks = []model.Task{
	{ID: 1, Title: "study go"},
}

func GetTasks(c *gin.Context) {
	c.JSON(200, tasks)
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