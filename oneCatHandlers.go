package main

import "net/http"

func getCat(req *http.Request) (int, any) {
	catID := req.PathValue("catId")

	cat, found := catsDatabase[catID]
	if !found {
		Logger.Info("Cat not found")
		return http.StatusNotFound, "Cat not found"
	}

	return http.StatusOK, cat
}

func deleteCat(req *http.Request) (int, any) {
	catID := req.PathValue("catId")

	_, found := catsDatabase[catID]
	if !found {
		Logger.Info("Cat not found")
		return http.StatusNotFound, "Cat not found"
	}

	delete(catsDatabase, catID)
	return http.StatusNoContent, nil
}
