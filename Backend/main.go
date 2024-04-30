package main

import (
	analyzer_test "P1/Analyzer"
	utilities_test "P1/Utilities"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strings"
// )

//Para ejecutar el programa se debe correr el comando go run main.go

// func main() {
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		fmt.Print("Ingrese un comando: ")
// 		input, err := reader.ReadString('\n')
// 		if err != nil {
// 			fmt.Println("Error al leer la entrada:", err)
// 			continue
// 		}

// 		input = strings.TrimSpace(input) //quitamos el salto de linea
// 		//Para llamar una funcion desde otro archivo este debe ir en mayuscula al inicio
// 		analyzer_test.Command(input)
// 	}
// }

func main() {
	app := fiber.New()
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"status":  200,
			"message": "Funcion 1",
		})
	})
	app.Post("/comandos", func(c *fiber.Ctx) error {
		text := c.FormValue("texto")
		//fmt.Println(text)
		analyzer_test.Command(text)
		// Enviar respuesta de éxito al cliente
		return c.SendString(utilities_test.Resultados.String())
	})

	app.Post("/cargarArchivo", func(c *fiber.Ctx) error {
		// Recibe el contenido como texto
		body := c.FormValue("fileContent")

		// Divide el texto en líneas
		lines := strings.Split(body, "\n")

		// Procesa cada línea con la función analizar
		for _, line := range lines {
			analyzer_test.Command(line)
		}

		return c.SendString(utilities_test.Resultados.String())
	})

	log.Fatal(app.Listen(":3000"))

}
