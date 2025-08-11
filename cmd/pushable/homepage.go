package main

import (
	"github.com/labstack/echo/v4"
	"github.com/oliverisaac/pushable/types"
	"github.com/oliverisaac/pushable/views"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func homePageHandler(cfg types.Config, db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		pageData := types.HomePageData{Config: cfg}

		if user, ok := GetSessionUser(c); ok {
			logrus.Infof("Generating homepage for user %s", user.Email)
			pageData = pageData.WithUser(user)
		} else {
			logrus.Debug("Generating anonymous homepage")
		}

		return render(c, 200, views.Index(pageData))
	}
}
