package main

import (
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "os"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "net/http"
  "fmt"
  )

type Reaction struct {
  Id int
  Title string
  Image string
}

func main() {
  db := NewDB()
  defer db.Close()

  m := martini.Classic()

  m.Use(Auth)
  m.Use(render.Renderer())

  m.Post("/reactions", func(rq *http.Request, r render.Render) {
    _, err := db.Exec(
      "INSERT INTO reactions (title, image) VALUES(?, ?)",
      rq.FormValue("title"),
      rq.FormValue("image"),
    )

    if err != nil {
      panic(err.Error())
    }

    r.JSON(201, nil)
  })

  m.Get("/randomreaction", func(r render.Render) {
    var reaction Reaction
    rows, err := db.Query("SELECT * FROM reactions ORDER BY RANDOM() LIMIT 1")

    if err != nil {
      panic(err.Error())
    }

    for rows.Next() {
      err := rows.Scan(&reaction.Id, &reaction.Title, &reaction.Image)

      if err != nil {
        panic(err.Error())
      }
    }
    fmt.Println(reaction)
    r.JSON(200, reaction)
  })

  m.Run()
}

func Auth(res http.ResponseWriter, req *http.Request) {
  if req.Header.Get("API-KEY") != os.Getenv("KEGIFTO_API_KEY") {
    http.Error(res, "Tu culo", 401)
  }
}

func NewDB() *sql.DB {
  db, err := sql.Open("sqlite3", "gogifs.sqlite")

  if err != nil {
    panic(err)
  }

  _, err = db.Exec(
    "create table if not exists reactions(id integer primary key autoincrement, title text, image text)",
  )

  if err != nil {
    panic(err)
  }

  return db
}
