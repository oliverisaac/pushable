package types

import (
	errs "errors"
	"fmt"
	"net/mail"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/oliverisaac/goli"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Hostname          string
	AllowSignup       bool
	AllowSignupEmails []string
	CookeSecret       []byte
	DBPath            string
	VapidPublicKey    string
	VapidPrivateKey   string
}

func ConfigFromEnv() (Config, error) {
	ret := Config{}
	var retErr error
	var err error

	ret.AllowSignup, err = strconv.ParseBool(goli.DefaultEnv("PUSHABLE_ALLOW_SIGNUP", "false"))
	if err != nil {
		retErr = errs.Join(retErr, errors.Wrap(err, "parsing PUSHABLE_ALLOW_SIGNUP"))
	}

	allowedEmails := strings.Split(os.Getenv("PUSHABLE_ALLOW_SIGNUP_EMAILS"), ",")
	for _, e := range allowedEmails {
		if e == "" {
			continue
		}
		email, err := mail.ParseAddress(e)
		if err != nil {
			retErr = errs.Join(retErr, errors.Wrapf(err, "parsing email %q", e))
		} else {
			ret.AllowSignupEmails = append(ret.AllowSignupEmails, email.Address)
		}
	}
	logrus.Infof("Allowed signup emails: %v", ret.AllowSignupEmails)

	cookieSecret, ok := os.LookupEnv("PUSHABLE_COOKIE_STORE_SECRET")
	if !ok {
		retErr = errs.Join(retErr, fmt.Errorf("You must define env PUSHABLE_COOKIE_STORE_SECRET"))
	} else {
		ret.CookeSecret = []byte(cookieSecret)
	}

	ret.DBPath, ok = os.LookupEnv("PUSHABLE_DB_PATH")
	if !ok {
		retErr = errs.Join(retErr, fmt.Errorf("You must define env PUSHABLE_DB_PATH"))
	} else if _, err := os.Stat(path.Dir(ret.DBPath)); err != nil {
		retErr = errs.Join(retErr, errors.Wrap(err, "Directory for PUSHABLE_DB_PATH must exist"))
	}

	ret.VapidPrivateKey, ok = os.LookupEnv("VAPID_PRIVATE_KEY")
	if !ok {
		retErr = errs.Join(retErr, fmt.Errorf("You must define env VAPID_PRIVATE_KEY"))
	}

	ret.VapidPublicKey, ok = os.LookupEnv("VAPID_PUBLIC_KEY")
	if !ok {
		retErr = errs.Join(retErr, fmt.Errorf("You must define env VAPID_PUBLIC_KEY"))
	}

	ret.Hostname = goli.DefaultEnv("PUSHABLE_HOSTNAME", "localhost")

	return ret, retErr
}
