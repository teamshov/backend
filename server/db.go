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
	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo/v4"
)

func InitDBService(e *echo.Echo) {
	e.GET("/db/:db/:id", apidbGet)
	e.GET("/db/:db/:id/:attch", apidbGetAttch)
	e.PUT("/db/:db/:id", apidbPut)
	e.DELETE("/db/:db/:id", apidbDelete)
	e.GET("/db/all/:db", apidbAll)
}

func DBGet(dbname string, id string) (map[string]interface{}, error) {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		panic(err)
	}

	db, err := client.DB(context.TODO(), dbname)
	if err != nil {
		return nil, err
	}

	row, err := db.Get(context.TODO(), id)
	if err != nil {
		return nil, err
	}

	var doc map[string]interface{}
	if err := row.ScanDoc(&doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func apidbGet(c echo.Context) error {
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

func apidbGetAttch(c echo.Context) error {
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

func DBPut(dbparam string, idparam string, data map[string]interface{}) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return err
	}

	db, err := client.DB(context.TODO(), dbparam)
	if err != nil {
		return err
	}

	row, err := db.Get(context.TODO(), idparam)

	var doc map[string]interface{}
	if err != nil {
		if kivik.StatusCode(err) == kivik.StatusNotFound {
			doc2 := map[string]interface{}{
				"_id": idparam,
			}
			rev, err2 := db.Put(context.TODO(), idparam, doc2)
			if err2 != nil {
				return err2
			}
			doc2["_rev"] = rev
			doc = doc2
		} else {
			return err
		}
	} else {
		if err := row.ScanDoc(&doc); err != nil {
			return err
		}
	}

	for k, v := range data {
		doc[k] = v
	}

	rev, _ := db.Put(context.TODO(), idparam, doc)
	doc["_rev"] = rev

	return nil
}

func apidbPut(c echo.Context) error {
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
	row, err := db.Get(context.TODO(), idparam)

	var doc map[string]interface{}
	if err != nil {
		if kivik.StatusCode(err) == kivik.StatusNotFound {
			doc2 := map[string]interface{}{
				"_id": idparam,
			}
			rev, err2 := db.Put(context.TODO(), idparam, doc2)
			if err2 != nil {
				return c.String(http.StatusInternalServerError, err2.Error())
			}
			doc2["_rev"] = rev
			doc = doc2
		} else {
			return c.String(http.StatusInternalServerError, err.Error())
		}
	} else {
		if err := row.ScanDoc(&doc); err != nil {
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

func apidbDelete(c echo.Context) error {
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

func DBAll(dbparam string) ([]string, error) {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
	if err != nil {
		return nil, err
	}

	db, err := client.DB(context.TODO(), dbparam)
	if err != nil {
		return nil, err
	}

	rows, err := db.AllDocs(context.TODO())
	if err != nil {
		return nil, err
	}

	var ilist []string
	for rows.Next() {
		ilist = append(ilist, rows.ID())
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ilist, nil
}

func apidbAll(c echo.Context) error {
	client, err := kivik.New(context.TODO(), "couch", "http://admin:seniorshov@omaraa.ddns.net:5984/")
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

func RedisGetInt(key string) (int, error) {
	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return 0, err
	}
	defer r.Close()

	i, err := redis.Int(r.Do("GET", key))
	if err != nil {
		return 0, err
	}

	return i, nil
}

func RedisGetFloat64(key string) (float64, error) {
	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return 0, err
	}
	defer r.Close()

	f, err := redis.Float64(r.Do("GET", key))
	if err != nil {
		return f, err
	}

	return f, nil
}

func RedisGetString(key string) (string, error) {
	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return "", err
	}
	defer r.Close()

	f, err := redis.String(r.Do("GET", key))
	if err != nil {
		return f, err
	}

	return f, nil
}

func RedisSetInterface(key string, f map[string]interface{}) error {
	r, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer r.Close()

	_, err = r.Do("SET", key, f)
	if err != nil {
		return err
	}

	return nil
}