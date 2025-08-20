package pushclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

func SendPush(endpointHostname string, push Push) error {
	formData := url.Values{}
	formData.Set("topic", push.Topic)
	formData.Set("title", push.Title)
	formData.Set("body", push.Body)
	formData.Set("icon", push.Icon)
	formData.Set("link", push.Link)

	endpoint := fmt.Sprintf("https://%s/push", endpointHostname)
	resp, err := http.Post(endpoint, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to send push: %s", string(respBody))
	}

	return nil
}
