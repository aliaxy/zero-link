// Package analyticsservicelogic contains link-rpc analytics service request logic.
package analyticsservicelogic

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

func hashIP(ip, salt string) string {
	mac := hmac.New(sha256.New, []byte(salt))
	mac.Write([]byte(ip))
	return hex.EncodeToString(mac.Sum(nil))
}

func detectDevice(ua string) string {
	if ua == "" {
		return "unknown"
	}
	lower := strings.ToLower(ua)
	if strings.Contains(lower, "bot") || strings.Contains(lower, "crawler") || strings.Contains(lower, "spider") {
		return "bot"
	}
	if strings.Contains(ua, "Mobile") || strings.Contains(ua, "Android") || strings.Contains(ua, "iPhone") {
		return "mobile"
	}
	return "desktop"
}

const (
	defaultRangeDays = 30
	maxRangeDays     = 90
)

var errInvalidDateRange = errors.New("invalid date range")

func parseDateRange(from, to string) (time.Time, time.Time, error) {
	now := time.Now().UTC().Truncate(24 * time.Hour)

	var fromT, toT time.Time
	var err error

	if from == "" {
		fromT = now.AddDate(0, 0, -defaultRangeDays)
	} else {
		fromT, err = time.Parse("2006-01-02", from)
		if err != nil {
			return time.Time{}, time.Time{}, errInvalidDateRange
		}
	}

	if to == "" {
		toT = now
	} else {
		toT, err = time.Parse("2006-01-02", to)
		if err != nil {
			return time.Time{}, time.Time{}, errInvalidDateRange
		}
	}

	if fromT.After(toT) {
		return time.Time{}, time.Time{}, errInvalidDateRange
	}
	if toT.Sub(fromT) > maxRangeDays*24*time.Hour {
		return time.Time{}, time.Time{}, errInvalidDateRange
	}

	return fromT, toT, nil
}
