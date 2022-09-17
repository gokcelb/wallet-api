package main

import (
	"fmt"

	"github.com/gokcelb/wallet-api/config"
	"github.com/gokcelb/wallet-api/internal/wallet"
	"github.com/labstack/echo/v4"
)

func main() {
	conf, err := config.Read(fmt.Sprintf(".config/%s.json", config.GetEnvOrDefault()))
	if err != nil {
		panic(err)
	}

	e := echo.New()

	s := wallet.NewService(nil, conf)
	h := wallet.NewHandler(s)

	h.RegisterRoutes(e)

	if err := e.Start(":8000"); err != nil {
		panic(err)
	}
}
