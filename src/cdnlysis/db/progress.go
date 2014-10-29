package db

import tiedot "github.com/HouzuoGuo/tiedot/db"

var cachedDB tiedot.DB

const collectionName = "sync_progress"

var indexes = []string{"prefix", "marker"}

func getDB() *tiedot.DB {
	return &cachedDB
}

func LastMarker(Prefix string) string {
	syncDB := getDB()
	coll := syncDB.Use(collectionName)

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

func Update(Prefix string, Marker string) {
	syncDB := getDB()
	coll := syncDB.Use(collectionName)

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

	syncDB, err := tiedot.OpenDB(path)
	if err != nil {
		panic(err)
	}

	var hasCol bool
	for _, name := range syncDB.AllCols() {
		if name == collectionName {
			hasCol = true
			break
		}
	}

	if !hasCol {
		if err := syncDB.Create(collectionName); err != nil {
			panic(err)
		}
	}

	// ****** Index *********

	coll := syncDB.Use(collectionName)

	for _, indexPath := range indexes {
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

	cachedDB = *syncDB
}
