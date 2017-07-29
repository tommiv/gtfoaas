package main

import (
    "net/http"
    "html/template"
    "encoding/json"
    "io/ioutil"
    "fmt"
    "math/rand"
    "github.com/go-redis/redis"
)

type indexContext struct {
    FuckedCount     int64
    BackgroundIndex int
}

type apiResponseIncrement struct {
    FuckedCount     int64
}

type config struct {
    RedisServer   string
    RedisPassword string
    RedisDB       int
    Listen        string
}

var redisClient *redis.Client

func index(w http.ResponseWriter, r *http.Request) {
    counter, err := redisClient.Get("gtfo:counter").Int64()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println(err)
        return
    }

    ctx := indexContext {
        FuckedCount: counter,
        BackgroundIndex: rand.Intn(4),
    }
    indexTemplate, _ := template.ParseFiles("templates/index.html")
    indexTemplate.Execute(w, ctx)
}

func gtfo(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    counter := redisClient.Incr("gtfo:counter").Val()
    ctx := apiResponseIncrement {
        FuckedCount: counter,
    }

    bodyBuf, err := json.Marshal(ctx)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Println(err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(bodyBuf)
}

func main() {
    cfgBuf, err := ioutil.ReadFile("config.json")
    if err != nil {
        panic("Can't read config: " + err.Error())
    }

    var cfg config
    err = json.Unmarshal(cfgBuf, &cfg)
    if err != nil {
        panic("Can't parse config: " +  err.Error())
    }

    redisClient = redis.NewClient(&redis.Options{
    	Addr:     cfg.RedisServer,
    	Password: cfg.RedisPassword,
    	DB:       cfg.RedisDB,
    })

    http.HandleFunc("/", index)
    http.HandleFunc("/api/v1/gtfo", gtfo)

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.ListenAndServe(cfg.Listen, nil)
}
