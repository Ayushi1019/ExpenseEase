package main

import (
	"ExpenseEase/server/app"
)

func main() {
	a := app.App{}
	a.Initialize()
	a.Run(":8081")
}
