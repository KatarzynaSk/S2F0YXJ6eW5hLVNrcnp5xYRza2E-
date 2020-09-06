package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type correctResp struct {
	ID int `json:"id"`
}

type errorResp struct {
	Message string `json:"message"`
}

func addRequestHandler(db *sqlx.DB, f *fetcher) echo.HandlerFunc {
	return func(c echo.Context) error {
		var r request
		err := json.NewDecoder(c.Request().Body).Decode(&r)
		if err != nil {
			resp := errorResp{Message: "invalid json"}
			return c.JSON(http.StatusBadRequest, resp)
		}

		id, err := insertRequest(db, r)
		if err != nil {
			log.Print(err)
			return err
		}

		r.ID = id
		f.StartJob(r)

		return c.JSON(http.StatusOK, correctResp{ID: id})
	}
}

func deleteRequestHandler(db *sqlx.DB, f *fetcher) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			resp := errorResp{Message: "invalid id param"}
			return c.JSON(http.StatusBadRequest, resp)
		}

		err = deactivateRequest(db, id)
		if err != nil {
			log.Print(err)
			return err
		}

		f.StopJob(id)

		return c.JSON(http.StatusOK, correctResp{ID: id})
	}
}

func selectRequestsHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		all, err := selectAllActiveRequest(db)
		if err != nil {
			log.Print(err)
			return err
		}

		return c.JSON(http.StatusOK, all)
	}
}

func selectHistoryHandler(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			resp := errorResp{Message: "invalid id param"}
			return c.JSON(http.StatusBadRequest, resp)
		}

		r, err := selectActiveRequest(db, id)
		if err != nil {
			log.Print(err)
			return err
		}

		if r == nil {
			resp := errorResp{Message: "request does not exist"}
			return c.JSON(http.StatusNotFound, resp)
		}

		all, err := selectAllRequestResult(db, id)
		if err != nil {
			log.Print(err)
			return err
		}

		return c.JSON(http.StatusOK, all)
	}
}
