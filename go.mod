module github.com/teamshov/backend

replace (
	github.com/teamshov/backend/build => ./build
	github.com/teamshov/backend/server => ./server
)

require (
	github.com/flimzy/diff v0.1.5 // indirect
	github.com/flimzy/kivik v1.8.1
	github.com/flimzy/testy v0.1.15 // indirect
	github.com/go-kivik/couchdb v1.8.1
	github.com/go-kivik/kivik v1.8.1 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/labstack/echo/v4 v4.0.0
	github.com/teamshov/backend/build v0.0.0-20190414123643-9ac1e98e4497 // indirect
	github.com/teamshov/backend/deploy v0.0.0-20190414120950-258508d4d406 // indirect
	github.com/teamshov/backend/server v0.0.0-20190414122341-f19d23914474 // indirect
	gopkg.in/resty.v1 v1.12.0
)
