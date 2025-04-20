package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Add middleware to recover from panics
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					stackTrace := string(debug.Stack())
					filename := fmt.Sprintf("panic_%d.txt", time.Now().Unix())
					if err := os.WriteFile(filename, []byte(stackTrace), 0644); err != nil {
						c.Logger().Errorf("Failed to write stack trace: %v", err)
					}
					c.Logger().Error("Panic occurred: ", r)
					c.String(http.StatusInternalServerError, "Internal Server Error")
				}
			}()
			return next(c)
		}
	})

	e.GET("/", func(c echo.Context) error {
		// Safe division with error handling
		x := 10
		y := 0
		
		if y == 0 {
			return c.String(http.StatusBadRequest, "Error: Division by zero is not allowed")
		}
		
		result := x / y
		return c.String(http.StatusOK, fmt.Sprintf("Result: %d", result))
	})

	e.Logger.Fatal(e.Start(":8080"))
}