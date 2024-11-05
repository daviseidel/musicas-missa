package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"embed"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
  "github.com/lithammer/fuzzysearch/fuzzy"
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
  
  // Serve os arquvis estáticos (js, css, e o json)
  app.Static("/static", "views/static") 

  app.Get("/", func(c *fiber.Ctx) error {
      return c.Render("views/index", fiber.Map{
          "var": ".",
      }, "views/layouts/main")
  })

  // /musicaCard?id
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
    
    if id < 0 || id >= 404 {
      return fiber.NewError(fiber.StatusServiceUnavailable, "Error: id out of range")
    } 
    
    return c.Render("views/musica", fiber.Map{
      "titulo": Musicas[id].Titulo,
      "letra": Musicas[id].Letra,
    }, "views/layouts/main")
  })

  app.Get("/musicas", func(c *fiber.Ctx) error {
    items := []string{}
    buf := new(bytes.Buffer) 
    
    for i:=0; i < len(Musicas); i++ {
      err := c.App().Config().Views.Render(buf ,"views/partials/item", fiber.Map{
        "titulo": Musicas[i].Titulo,
      }) 
      if err != nil {
        fmt.Println(err)
      }
      items = append(items, buf.String())
    }
    
    return c.SendString(strings.Join(items, ""))
    //return c.Render("views/layouts/blank", fiber.Map{
    //  "data": strings.Join(items, ""), 
    //}) 
  })

  app.Get("/search", func(c *fiber.Ctx) error {
    searchkey := c.Query("key")
    titulos := []string{}

    for i:=0; i < len(Musicas); i++ {
      titulos = append(titulos, Musicas[i].Titulo)
    }
    valid := fuzzy.Find(searchkey, titulos)
    buffer := new(bytes.Buffer) 
    fmt.Println(len(valid)) 
    validBuf := []string{}

    for i:=0; i < len(valid); i++ {
      err := c.App().Config().Views.Render(buffer ,"views/partials/item", fiber.Map{
        "titulo": valid[i],
      }) 
      if err != nil {
        fmt.Println(err)
      }
      fmt.Println(buffer.String())
      validBuf = append(valid, buffer.String())
    }
    return c.SendString(strings.Join(validBuf, ""))
  })

  readData()
  log.Fatal(app.Listen(":3000"))
}

// Abre o arquivo .json e leva pra memória. Deve rodar antes do app iniciar
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
