package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	//"fmt"

	"github.com/flimzy/kivik" // Development version of Kivik
	//"github.com/go-kivik/kivik"
	_ "github.com/go-kivik/couchdb" // The CouchDB driver
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORS())

	e.Static("/", "dist")
	e.File("/", "dist/index.html")

	e.GET("/beacons/:id", Getbeacon)
	e.PUT("/beacons/:id", Updatebeacon)
	e.DELETE("/beacons/:id", Deletebeacon)
	e.GET("/allbeacons", Getbeacons)

	e.GET("/db/:db/:id", DBGet)
	e.GET("/db/:db/:id/:attch", DBGetAttch)
	e.PUT("/db/:db/:id", DBPut)
	e.DELETE("/db/:db/:id", DBDelete)
	e.GET("/db/all/:db", DBAll)

	e.Logger.Fatal(e.Start(":62027"))
}

type Beacon struct {
	ID     string  `json:"_id"`
	Rev    string  `json:"_rev,omitempty"`
	Temp   string  `json:"temperature,omitempty"`
	XPos   float64 `json:"xpos,omitempty"`
	YPos   float64 `json:"ypos,omitempty"`
	Major  string  `json:"major,omitempty"`
	Minor  string  `json:"minor,omitempty"`
	Offset float64 `json:"offset,omitempty"`
	Shovid string  `json:"shovid,omitempty"`
}

func DBGet(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	dbparam := c.Param("db")

	db, err := client.DB(context.TODO(), dbparam)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	idparam := c.Param("id")

	row, _ := db.Get(context.TODO(), idparam)

	var doc map[string]interface{}
	if err = row.ScanDoc(&doc); err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, doc)
}

func DBGetAttch(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	dbparam := c.Param("db")

	db, err := client.DB(context.TODO(), dbparam)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	idparam := c.Param("id")
	attchparam := c.Param("attch")

	row, _ := db.Get(context.TODO(), idparam)

	var doc map[string]interface{}
	if err = row.ScanDoc(&doc); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	a, err := db.GetAttachment(context.TODO(), doc["_id"].(string), doc["_rev"].(string), attchparam)
	if err = row.ScanDoc(&doc); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	data, err := ioutil.ReadAll(a)
	if err = row.ScanDoc(&doc); err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.Blob(http.StatusOK, a.ContentType, data)
}

func DBPut(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	dbparam := c.Param("db")

	db, err := client.DB(context.TODO(), dbparam)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	idparam := c.Param("id")
	row, _ := db.Get(context.TODO(), idparam)

	var doc map[string]interface{}
	if err = row.ScanDoc(&doc); err != nil {
		if kivik.StatusCode(err) == kivik.StatusNotFound {
			doc2 := map[string]interface{}{
				"_id": idparam,
			}
			rev, err2 := db.Put(context.TODO(), idparam, doc2)
			if err2 != nil {
				return c.String(http.StatusInternalServerError, err2.Error())
			}
			doc["_rev"] = rev

		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	var data map[string]interface{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, &data)

	for k, v := range data {
		doc[k] = v
	}

	rev, _ := db.Put(context.TODO(), idparam, doc)
	doc["_rev"] = rev

	return c.JSON(http.StatusOK, doc)
}

func DBDelete(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	dbparam := c.Param("db")

	db, err := client.DB(context.TODO(), dbparam)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	idparam := c.Param("id")

	row, _ := db.Get(context.TODO(), idparam)

	var doc map[string]interface{}
	if err = row.ScanDoc(&doc); err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	deletedrev, err := db.Delete(context.TODO(), doc["_id"].(string), doc["_rev"].(string))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Deleted, rev: "+deletedrev)
}

func DBAll(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@localhost:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	dbparam := c.Param("db")

	db, err := client.DB(context.TODO(), dbparam)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	rows, err := db.AllDocs(context.TODO())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var ilist []string
	for rows.Next() {
		ilist = append(ilist, rows.ID())
	}

	if rows.Err() != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, ilist)
}

func Getbeacon(c echo.Context) error {

	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	beacons, err := client.DB(context.TODO(), "beacons")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	bid := c.Param("id")

	row, _ := beacons.Get(context.TODO(), bid)

	var beacon map[string]interface{}
	if err = row.ScanDoc(&beacon); err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, beacon)
}

func Updatebeacon(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	beacons, err := client.DB(context.TODO(), "beacons")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	bid := c.Param("id")
	row, _ := beacons.Get(context.TODO(), bid)

	var beacon map[string]interface{}
	if err = row.ScanDoc(&beacon); err != nil {
		if kivik.StatusCode(err) == kivik.StatusNotFound {
			doc := map[string]interface{}{
				"_id": bid,
			}
			rev, err2 := beacons.Put(context.TODO(), bid, doc)
			if err2 != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}
			beacon["_rev"] = rev

		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	temperature := c.FormValue("temperature")
	if temperature != "" {
		beacon["temperature"] = temperature
	}

	xpos := c.FormValue("xpos")
	if xpos != "" {
		beacon["xpos"], _ = strconv.ParseFloat(xpos, 64)
	}
	ypos := c.FormValue("ypos")
	if ypos != "" {
		beacon["ypos"], _ = strconv.ParseFloat(ypos, 64)
	}
	major := c.FormValue("major")
	if major != "" {
		beacon["major"] = major
	}
	minor := c.FormValue("minor")
	if minor != "" {
		beacon["minor"] = minor
	}
	offset := c.FormValue("offset")
	if offset != "" {
		beacon["offset"], _ = strconv.ParseFloat(offset, 64)
	}
	shovid := c.FormValue("shovid")
	if shovid != "" {
		beacon["shovid"] = shovid
	}

	rev, _ := beacons.Put(context.TODO(), bid, beacon)
	beacon["_rev"] = rev

	return c.JSON(http.StatusOK, beacon)
}

func Deletebeacon(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	beacons, err := client.DB(context.TODO(), "beacons")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	bid := c.Param("id")

	row, _ := beacons.Get(context.TODO(), bid)

	var beacon Beacon
	if err = row.ScanDoc(&beacon); err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	deletedrev, err := beacons.Delete(context.TODO(), beacon.ID, beacon.Rev)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Deleted, rev: "+deletedrev)
}

type BeaconList struct {
	Beacons []string `json:"beacons"`
}

func Getbeacons(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@localhost:5984/")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	beacons, err := client.DB(context.TODO(), "beacons")
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	rows, err := beacons.AllDocs(context.TODO())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	var bl BeaconList
	var sb strings.Builder
	sb.WriteString(`{"beacons":[`)
	for rows.Next() {
		bl.Beacons = append(bl.Beacons, rows.ID())
		sb.WriteString(rows.ID())
		sb.WriteString(",")
	}
	sb.WriteString("]}")

	if rows.Err() != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, bl)
}
