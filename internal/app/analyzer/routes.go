package analyzer

import (
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"

	"clicktweak/internal/pkg/db"
	"clicktweak/internal/pkg/model"

	exception "clicktweak/internal/pkg/error"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const day = time.Hour * 24

const (
	today     = "today"
	yesterday = "yesterday"
	lastWeek  = "lastweek"
	lastMonth = "lastmonth"
)

func GetStats(url db.Url, log db.Log) echo.HandlerFunc {
	return func(context echo.Context) error {
		user := context.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		temp := claims["id"].(float64)
		userID := uint(temp)

		fromStr := strings.ToLower(context.QueryParams().Get("from"))
		if !isFromValid(fromStr) {
			return context.JSON(http.StatusBadRequest, exception.ToJSON(exception.MalformedRequest))
		}

		until, from := calculateTimeRange(fromStr)

		// fetch all user urls
		urls, err := url.GetByUserID(userID)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
		}

		if urls == nil {
			return context.JSON(http.StatusNotFound, exception.ToJSON(exception.ResourceNotFound))
		}

		var result = make([]*model.Report, len(urls))
		for i, elem := range urls {
			result[i], err = getUrlStats(elem.ID, from, until, log)
			if err != nil {
				return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
			}
		}

		return context.JSONPretty(http.StatusOK, result, "  ")
	}
}

func GetUrlStats(url db.Url, log db.Log) echo.HandlerFunc {
	return func(context echo.Context) error {
		user := context.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		temp := claims["id"].(float64)
		userID := uint(temp)

		id := context.Param("id")

		fromStr := strings.ToLower(context.QueryParams().Get("from"))
		if !isFromValid(fromStr) {
			return context.JSON(http.StatusBadRequest, exception.ToJSON(exception.MalformedRequest))
		}

		until, from := calculateTimeRange(fromStr)

		// fetch url from database
		url, err := url.GetByID(id)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
		}

		if url == nil {
			return context.JSON(http.StatusNotFound, exception.ToJSON(exception.ResourceNotFound))
		}

		if url.UserID != userID {
			return context.JSON(http.StatusUnauthorized, exception.ToJSON(errors.New("you don't have enough permission to fetch this resource")))
		}

		result, err := getUrlStats(url.ID, from, until, log)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
		}
		return context.JSONPretty(http.StatusOK, result, "  ")
	}
}

func getUrlStats(id string, from, until string, db db.Log) (*model.Report, error) {
	report, err := db.GetStats(id, from[:len(from)-1], until[:len(until)-1])
	if err != nil {
		return nil, err
	}
	return report, nil
}

func calculateTimeRange(fromStr string) (until string, from string) {
	switch fromStr {
	case today:
		until = time.Now().Format(time.RFC3339)
		from = time.Now().Truncate(day).Format(time.RFC3339)
	case yesterday:
		until = time.Now().Truncate(day).Add(time.Minute * -1).Format(time.RFC3339)
		from = time.Now().Add(day * -1).Truncate(day).Format(time.RFC3339)
	case lastWeek:
		until = time.Now().Truncate(day).Add(time.Minute * -1).Format(time.RFC3339)
		from = time.Now().Truncate(day).Add(day * -7).Format(time.RFC3339)
	case lastMonth:
		until = time.Now().Truncate(day).Add(time.Minute * -1).Format(time.RFC3339)
		from = time.Now().Truncate(day).Add(day * -30).Format(time.RFC3339)
	}
	return
}

func isFromValid(from string) bool {
	return !(len(from) == 0 || from != today && from != yesterday && from != lastWeek && from != lastMonth)
}
