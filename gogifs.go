package main

import (
  "os"
  "math/rand"
  "time"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "net/http"
  )

type Reaction struct {
  Title string
  Image string
}

func Auth(res http.ResponseWriter, req *http.Request) {
  if req.Header.Get("API-KEY") != os.Getenv("KEGIFTO_API_KEY") {
    http.Error(res, "Tu culo", 401)
  }
}

func DB() martini.Handler {
  dburl := os.Getenv("MONGOLAB_URI")

  session, err := mgo.Dial(dburl)

  if err!= nil {
    panic(err)
  }

  return func(c martini.Context) {
    s := session.Clone()
    c.Map(s.DB(os.Getenv("MONGOLAB_DB")))
    defer s.Close()
    c.Next()
  }
}

func RandomNumber(max int) int {
  rand.Seed(time.Now().Unix())

  if max == 1 {
    return 0
  }

  return rand.Intn(max)
}

func main() {
  m := martini.Classic()

  m.Use(Auth)
  m.Use(render.Renderer())
  m.Use(DB())

  m.Post("/reactions", func(res http.ResponseWriter, rq *http.Request, r render.Render, db *mgo.Database) {
    title, image := rq.FormValue("title"), rq.FormValue("image")

    if title == "" || image == "" {
      http.Error(res, "Fill title and image", 422)
      return
    }

    reaction := Reaction{title, image}

    db.C("reactions").Insert(reaction)

    r.JSON(201, "Created")
  })

  m.Get("/reactions", func(r render.Render, rq *http.Request, db *mgo.Database) {
    query := rq.FormValue("q")
    var reaction Reaction

    db.C("reactions").Find(bson.M{"title": query}).One(&reaction)

    r.JSON(200, reaction)
  })

  m.Get("/randomreaction", func(r render.Render, db *mgo.Database) {
    var reaction Reaction

    total, err := db.C("reactions").Count()
    if err != nil{
      panic(err.Error())
    }

    skip := RandomNumber(total)

    db.C("reactions").Find(nil).Limit(-1).Skip(skip).One(&reaction)

    r.JSON(200, reaction)
  })

  m.Run()
}
