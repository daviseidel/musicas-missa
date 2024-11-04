package main

import (
  "fmt"
  "encoding/json"
  "os"
  "log"
  "net/http"
  "html/template"

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

//go:embed views/static/*
var cssfs embed.FS



func main() {
  engine := html.NewFileSystem(http.FS(viewsfs), ".html")

  engine.AddFunc(
    "unescape", func(s string) template.HTML {
        return template.HTML(s)
    },
  )

  app := fiber.New(fiber.Config{
    Views: engine,
  })

  app.Static("/static", "views/static") 

  app.Get("/", func(c *fiber.Ctx) error {
      return c.Render("views/index", fiber.Map{
          "var": ".",
      }, "views/layouts/main")
  })

  app.Get("/musicaCard", func(c *fiber.Ctx) error { 
    id := c.QueryInt("id")
    if id < 0 || id >= 404 {
      return fiber.NewError(fiber.StatusServiceUnavailable, "Error: id out of range")
    } 

    return c.Render("views/partials/card", fiber.Map{
      "titulo": Musicas[id].Titulo,
    })
  })

  app.Get("/musica", func(c *fiber.Ctx) error {
    id := c.QueryInt("id")
    fmt.Println(id)
    if id < 0 || id >= 404 {
      return fiber.NewError(fiber.StatusServiceUnavailable, "Error: id out of range")
    } else {
      return c.Render("views/musica", fiber.Map{
        "titulo": Musicas[id].Titulo,
        "letra": Musicas[id].Letra,
      }, "views/layouts/main")
    }
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
