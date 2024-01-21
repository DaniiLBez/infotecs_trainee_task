package main

import "infotecs_trainee_task/internal/app"

const configPath = "./config/config.yaml"

func main() {
	app.Run(configPath)
}
