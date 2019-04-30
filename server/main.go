package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	buildings map[string]*Building
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())

	e.Static("/", "dist")
	e.File("/", "dist/index.html")
	e.File("/dangermap", "dist/danger.html")

	InitDBService(e)

	InitSensorService()
	defer DisconnectSensorService()

	buildings = make(map[string]*Building)
	buildings["eb2"] = initBuilding("eb2")

	e.PUT("/api/getpath/:building/:floor/:id", api_getPath)
	e.GET("/api/emergency/:building", api_EmergencyStatus)

	e.PUT("/getpath/:id", api_getPath_legacy)
	e.PUT("/getpath", api_getPath_legacy2)
	e.GET("/emergency", api_EmergencyStatus_legacy)

	e.Logger.Fatal(e.Start(":62027"))
}

func api_EmergencyStatus(c echo.Context) error {
	buildingID := c.Param("building")
	return c.JSON(http.StatusOK, buildings[buildingID].emergency)
}

func api_EmergencyStatus_legacy(c echo.Context) error {
	return c.JSON(http.StatusOK, true)
}

func api_getPath(c echo.Context) error {
	var input map[string]interface{}
	//var paths map[string]interface{}

	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, &input)

	buildingID := c.Param("building")
	floorID := c.Param("floor")
	id := c.Param("id")

	graph := buildings[buildingID].floors[floorID].graph

	RedisSetInterface(fmt.Sprintf("eb2:1:users:%s", id), input)

	x := input["x"].(float64)
	y := input["y"].(float64)

	return c.JSON(http.StatusOK, graph.getPathXY(x, y))
}

func api_getPath_legacy(c echo.Context) error {
	var input map[string]interface{}
	//var paths map[string]interface{}

	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, &input)

	//buildingID := "eb2"
	//floorID := "L1"
	id := c.Param("id")

	graph := g

	RedisSetInterface(fmt.Sprintf("eb2:1:users:%s", id), input)

	x := input["x"].(float64)
	y := input["y"].(float64)

	return c.JSON(http.StatusOK, graph.getPathXY(x, y))
}

func api_getPath_legacy2(c echo.Context) error {
	var input map[string]interface{}
	//var paths map[string]interface{}

	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, &input)

	buildingID := "eb2"
	floorID := "L1"
	id := "asdfasdfasdf"

	graph := buildings[buildingID].floors[floorID].graph

	RedisSetInterface(fmt.Sprintf("eb2:1:users:%s", id), input)

	x := input["x"].(float64)
	y := input["y"].(float64)

	return c.JSON(http.StatusOK, graph.getPathXY(x, y))
}
