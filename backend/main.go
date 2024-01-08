package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var tasks []Task

// Custom middleware function
func customMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Perform tasks before the route handler
		fmt.Println("Executing custom middleware")

		// Call the next middleware or the route handler
		return next(c)
	}
}

func createTask(c echo.Context) error {
	task := new(Task)
	if err := c.Bind(task); err != nil {
		return err
	}

	task.ID = len(tasks) + 1
	tasks = append(tasks, *task)

	return c.JSON(http.StatusCreated, task)
}

func getTasks(c echo.Context) error {
	queryParams := c.QueryParams()

	// id := c.QueryParam("id")

	// Log the id
	// fmt.Printf("ID:-", id)

	// Log all query parameters
	for key, values := range queryParams {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	return c.JSON(http.StatusOK, tasks)
}

func getTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	for _, task := range tasks {
		if task.ID == id {
			return c.JSON(http.StatusOK, task)
		}
	}

	return c.String(http.StatusNotFound, "Task not found")
}

func updateTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	for i, task := range tasks {
		if task.ID == id {
			updatedTask := new(Task)
			if err := c.Bind(updatedTask); err != nil {
				return err
			}

			updatedTask.ID = id
			tasks[i] = *updatedTask
			return c.JSON(http.StatusOK, updatedTask)
		}
	}

	return c.String(http.StatusNotFound, "Task not found")
}

func deleteTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return c.String(http.StatusNoContent, "Task deleted")
		}
	}

	return c.String(http.StatusNotFound, "Task not found")
}

func main() {
	e := echo.New()

	// Root level middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		fmt.Println("log here")
		return c.String(http.StatusOK, "Hello, World!")
	})

	// e.Use(customMiddleware)

	// Routes
	e.POST("/tasks", createTask)
	e.GET("/tasks", getTasks, customMiddleware)
	e.GET("/tasks/:id", getTask)
	e.PUT("/tasks/:id", updateTask)
	e.DELETE("/tasks/:id", deleteTask)

	e.Logger.Fatal(e.Start(":1323"))
}
