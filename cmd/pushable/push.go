package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/labstack/echo/v4"
	"github.com/oliverisaac/pushable/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func removeSubscription(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := GetSessionUser(c)
		if !ok {
			return c.String(http.StatusUnauthorized, "unauthorized")
		}

		if err := db.Delete(user.PushSubscriptions).Error; err != nil {
			return errors.Wrap(err, "removing subscription")
		}

		return c.String(http.StatusOK, "subscription removed")
	}
}

func saveSubscription(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := GetSessionUser(c)
		if !ok {
			return c.String(http.StatusUnauthorized, "unauthorized")
		}

		var sub webpush.Subscription
		if err := c.Bind(&sub); err != nil {
			return errors.Wrap(err, "binding subscription")
		}

		keys, err := json.Marshal(sub.Keys)
		if err != nil {
			return errors.Wrap(err, "marshalling subscription keys")
		}

		pushSubscription := types.PushSubscription{
			UserID:   user.ID,
			Endpoint: sub.Endpoint,
			P256DH:   sub.Keys.P256dh,
			Auth:     sub.Keys.Auth,
			Keys:     string(keys),
		}

		if err := db.Create(&pushSubscription).Error; err != nil {
			return errors.Wrap(err, "saving subscription")
		}

		return c.String(http.StatusOK, "subscription saved")
	}
}

func pushNotification(cfg types.Config, db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		topic := c.FormValue("topic")
		title := c.FormValue("title")
		body := c.FormValue("body")
		icon := c.FormValue("icon")
		link := c.FormValue("link")

		var users []types.User
		if err := db.Preload("PushSubscriptions").Find(&users).Error; err != nil {
			return errors.Wrap(err, "finding users by topic")
		}

		for _, i := range []string{"fail", "success", "good", "bad", "neutral", "mid"} {
			if strings.HasPrefix(strings.ToLower(icon), i) {
				icon = fmt.Sprintf("https://%s/static/%s.png", cfg.Hostname, i)
			}
		}

		for _, user := range users {
			for _, subData := range user.PushSubscriptions {
				sub := &webpush.Subscription{
					Endpoint: subData.Endpoint,
					Keys: webpush.Keys{
						P256dh: subData.P256DH,
						Auth:   subData.Auth,
					},
				}

				pushPayload, err := json.Marshal(map[string]interface{}{
					"title": title,
					"body":  body,
					"icon":  icon,
					"data": map[string]string{
						"link": link,
					},
				})
				if err != nil {
					return errors.Wrap(err, "marshalling push payload")
				}

				resp, err := webpush.SendNotification(pushPayload, sub, &webpush.Options{
					Topic:           topic,
					VAPIDPublicKey:  cfg.VapidPublicKey,
					VAPIDPrivateKey: cfg.VapidPrivateKey,
					TTL:             3600,
					Urgency:         webpush.UrgencyNormal,
				})
				if err != nil {
					logrus.Error(errors.Wrap(err, "sending push notification"))
					continue
				}
				defer resp.Body.Close()

				if resp.StatusCode == 410 {
					if err := db.Delete(&subData).Error; err != nil {
						logrus.Error(errors.Wrap(err, "deleting subscription"))
					}
				}
			}
		}

		return c.String(http.StatusOK, "push notifications sent")
	}
}
