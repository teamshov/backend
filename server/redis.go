package main

import (
	"net/http"

	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
)

func InitRedisService(e *echo.Echo) {
	e.GET("/redis", redisTest)
}

func redisTest(c echo.Context) error {
	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer r.Close()

	n, _ := r.Do("PING")

	return c.String(http.StatusOK, n.(string))
}
