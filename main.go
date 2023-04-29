package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/yongxin/lanxin-gpt/lxgpt"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	client := lxgpt.New(config)

	if err := client.Start(); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() (*lxgpt.ClientConfig, error) {
	godotenv.Load(".env", "../.env")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	return &lxgpt.ClientConfig{
		LxAPIUrl:   os.Getenv("LX_API_URL"),
		AppID:      os.Getenv("APP_ID"),
		AppSecret:  os.Getenv("APP_SECRET"),
		OrgID:      os.Getenv("ORG_ID"),
		HookToken:  os.Getenv("HOOK_TOKEN"),
		HookSecret: os.Getenv("HOOK_SECRET"),
		ServerPort: port,
	}, nil
}
