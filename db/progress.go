package db

import (
	"errors"
	"log"
	"os"

	tiedot "github.com/HouzuoGuo/tiedot/db"
)

var cachedDB tiedot.DB
var dbInit bool

const pathCollection = "path_progress"
const markerCollection = "marker_progress"

var markerIndices = []string{"prefix", "marker"}
var pathIndices = []string{"path"}

func getDB() *tiedot.DB {
	if !dbInit {
		err := errors.New("Database not initialized. Invoke InitDB to proceeed")
		panic(err)
	}

	return &cachedDB
}

func LastMarker(Prefix string) string {
	syncDB := getDB()
	coll := syncDB.Use(markerCollection)

	query := map[string]interface{}{
		"eq": Prefix,
		"in": []interface{}{"prefix"},
	}

	queryResult := make(map[int]struct{})
	if err := tiedot.EvalQuery(query, coll, &queryResult); err != nil {
		panic(err)
	}

	for id := range queryResult {
		readBack, err := coll.Read(id)
		if err != nil {
			panic(err)
		}
		return readBack["marker"].(string)
	}

	return ""
}

func HasVisited(path string) bool {
	syncDB := getDB()
	coll := syncDB.Use(pathCollection)

	query := map[string]interface{}{
		"eq": path,
		"in": []interface{}{"path"},
	}

	queryResult := make(map[int]struct{})
	if err := tiedot.EvalQuery(query, coll, &queryResult); err != nil {
		log.Println(err)
		return false
	}

	return len(queryResult) > 0
}

func SetVisited(path string) {
	syncDB := getDB()
	coll := syncDB.Use(pathCollection)
	coll.Insert(map[string]interface{}{
		"path":    path,
		"visited": true,
	})
}

func Update(Prefix string, Marker string) {
	syncDB := getDB()
	coll := syncDB.Use(markerCollection)

	query := map[string]interface{}{
		"eq": Prefix,
		"in": []interface{}{"prefix"},
	}

	queryResult := make(map[int]struct{})
	if err := tiedot.EvalQuery(query, coll, &queryResult); err != nil {
		panic(err)
	}

	for id := range queryResult {
		err := coll.Update(id, map[string]interface{}{
			"prefix": Prefix,
			"marker": Marker,
		})
		if err != nil {
			panic(err)
		}
		return
	}

	_, err := coll.Insert(map[string]interface{}{
		"prefix": Prefix,
		"marker": Marker,
	})

	if err != nil {
		panic(err)
	}
}

func InitDB(path string) {
	dbInit = true

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Println("DB Path does not exist. Initializing.")
		err := os.Mkdir(path, 0755)
		if err != nil {
			panic(err)
		}
	}

	syncDB, err := tiedot.OpenDB(path)
	if err != nil {
		panic(err)
	}

	var hasMarkerCol, hasPathCol bool

	for _, name := range syncDB.AllCols() {
		if name == markerCollection {
			hasMarkerCol = true
		} else if name == pathCollection {
			hasPathCol = true
		}
	}

	if !hasMarkerCol {
		if err := syncDB.Create(markerCollection); err != nil {
			panic(err)
		}
	}

	if !hasPathCol {
		if err := syncDB.Create(pathCollection); err != nil {
			panic(err)
		}
	}

	// ****** Index *********

	coll := syncDB.Use(markerCollection)
	for _, indexPath := range markerIndices {
		var indexFound bool
		for _, path := range coll.AllIndexes() {
			if path[0] == indexPath {
				indexFound = true
				break
			}
		}

		if !indexFound {
			if err := coll.Index([]string{indexPath}); err != nil {
				panic(err)
			}
		}
	}

	// ************** Path Indices **************

	path_coll := syncDB.Use(pathCollection)
	for _, indexPath := range pathIndices {
		var indexFound bool
		for _, path := range path_coll.AllIndexes() {
			if path[0] == indexPath {
				indexFound = true
				break
			}
		}

		if !indexFound {
			if err := path_coll.Index([]string{indexPath}); err != nil {
				panic(err)
			}
		}
	}

	cachedDB = *syncDB
}
