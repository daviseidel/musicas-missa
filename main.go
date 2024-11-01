package main

import (
  "fmt"
  "encoding/json"
  "os"
  "log"
  "net/http"
  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/template/html/v2"
  "embed"
)

type Musica struct {
  Titulo string `json:"titulo"`
  Letra  string `json:"letra"`
}

var Musicas []Musica

//go:embed views/*
var viewsfs embed.FS

//go:embed views/css/*
var cssfs embed.FS



func main() {
  engine := html.NewFileSystem(http.FS(viewsfs), ".html")

  app := fiber.New(fiber.Config{
    Views: engine,
  })
  app.Static("/static", "views/css") 
  app.Get("/", func(c *fiber.Ctx) error {
      return c.Render("views/index", fiber.Map{
          "var": "Não mascare a sua vingança com a sua justiça - Borges João",
      }, "views/layouts/main")
  })

  app.Get("/:id", func(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    fmt.Println(err)
    return c.JSON(Musicas[id])
  })

  readData()
  log.Fatal(app.Listen(":3000"))
}

func readData() {
  file, err := os.Open("songs.json")
  if err != nil {
    fmt.Println(err)
  }

  defer file.Close()

  decoder := json.NewDecoder(file)
  err = decoder.Decode(&Musicas)
  if err != nil {
    fmt.Println(err)
  }
}
