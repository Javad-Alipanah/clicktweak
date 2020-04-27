package dispatcher

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"clicktweak/internal/pkg/db"
	exception "clicktweak/internal/pkg/error"
	"clicktweak/internal/pkg/model"
	"clicktweak/internal/pkg/util"
	"github.com/labstack/echo/v4"
	"github.com/mssola/user_agent"
)

func Redirect(url db.Url, logs chan<- model.Log) echo.HandlerFunc {
	return func(context echo.Context) error {
		id := context.Param("id")
		if !regexp.MustCompile(util.DefaultRegex).MatchString(id) {
			return context.JSON(http.StatusBadRequest, exception.ToJSON(exception.MalformedRequest))
		}

		result, err := url.GetByID(id)
		if err != nil {
			return context.JSON(http.StatusInternalServerError, exception.ToJSON(err))
		}

		if result == nil {
			return context.JSON(http.StatusNotFound, exception.ToJSON(exception.ResourceNotFound))
		}

		// send log elem to channel
		ua := user_agent.New(context.Request().UserAgent())
		browser, _ := ua.Browser()
		var device string
		if device = "mobile"; !ua.Mobile() {
			device = "desktop"
		}
		elem := model.Log{
			Id:         result.ID,
			Browser:    browser,
			Device:     device,
			RemoteAddr: strings.Split(context.Request().RemoteAddr, ":")[0],
			CreatedAt:  time.Now().Format(time.RFC3339),
		}
		logs <- elem

		return context.Redirect(http.StatusMovedPermanently, result.Url)
	}
}
