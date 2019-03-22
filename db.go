package main

import ( //"fmt"
	// Development version of Kivik
	//"github.com/go-kivik/kivik"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/flimzy/kivik"
	_ "github.com/go-kivik/couchdb" // The CouchDB driver
	"github.com/labstack/echo/v4"
)

func InitDBService(e *echo.Echo) {
	e.GET("/db/:db/:id", dbGet)
	e.GET("/db/:db/:id/:attch", dbGetAttch)
	e.PUT("/db/:db/:id", dbPut)
	e.DELETE("/db/:db/:id", dbDelete)
	e.GET("/db/all/:db", dbAll)
}

func DBGet(dbname string, id string) map[string]interface{} {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		panic(err)
	}

	db, err := client.DB(context.TODO(), dbname)
	if err != nil {
		panic(err)
	}

	row, _ := db.Get(context.TODO(), id)

	var doc map[string]interface{}
	if err = row.ScanDoc(&doc); err != nil {
		panic(err)
	}

	return doc
}

func dbGet(c echo.Context) error {
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

func dbGetAttch(c echo.Context) error {
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

func dbPut(c echo.Context) error {
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

func dbDelete(c echo.Context) error {
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

func dbAll(c echo.Context) error {
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
