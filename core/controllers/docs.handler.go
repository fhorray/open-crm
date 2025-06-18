package controllers

import (
	"fmt"
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func ScalarHandler() fiber.Handler {
	// Cria o handler padrão net/http.HandlerFunc do Scalar
	scalarHTTPHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "http://localhost:8787/docs/swagger.json", // caminho relativo
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})

	// Converte para fasthttp.Handler
	fasthttpHandler := fasthttpadaptor.NewFastHTTPHandler(scalarHTTPHandler)

	// Retorna handler Fiber adaptado para fasthttp
	return func(c *fiber.Ctx) error {
		fasthttpHandler(c.Context())
		return nil
	}
}
