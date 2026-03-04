package main

import "deploy-kit/internal/di"

func main() {
	app := di.NewApp()
	app.Run()
}
