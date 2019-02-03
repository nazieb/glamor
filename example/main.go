package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/labstack/echo"
	"github.com/nazieb/glamor"
)

func main() {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Hello World",
		})
	})

	server := glamor.WrapServer(e)
	lambda.Start(server)
}
