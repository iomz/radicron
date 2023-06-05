package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/yyoshiki41/go-radiko"
)

// NewAuth retuns a new authorized device with token
func NewAuth(ctx context.Context, client *radiko.Client, asset *Asset, areaID string) *Device {
	// generate a new device
	device := asset.NewDevice()

	// auth1
	req, _ := http.NewRequest("GET", "https://radiko.jp/v2/api/auth1", nil)
	req = req.WithContext(ctx)
	headers := map[string]string{
		UserAgentHeader:        device.UserAgent,
		RadikoAppHeader:        device.AppName,
		RadikoAppVersionHeader: device.AppVersion,
		RadikoDeviceHeader:     device.Name,
		RadikoUserHeader:       device.UserID,
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}

	device.Token = resp.Header.Get(RadikoAuthTokenHeader)
	offset, _ := strconv.ParseInt(resp.Header.Get(RadikoKeyOffsetHeader), 10, 64)
	length, _ := strconv.ParseInt(resp.Header.Get(RadikoKeyLentghHeader), 10, 64)
	fullkey, _ := base64.StdEncoding.DecodeString(FullkeyB64)
	partial := base64.StdEncoding.EncodeToString([]byte(fullkey[offset : offset+length]))

	location := asset.GenerateGPSForAreaID(areaID)

	req, _ = http.NewRequest("GET", "https://radiko.jp/v2/api/auth2", nil)
	req = req.WithContext(ctx)
	headers = map[string]string{
		UserAgentHeader:        device.UserAgent,
		RadikoAppHeader:        device.AppName,
		RadikoAppVersionHeader: device.AppVersion,
		RadikoDeviceHeader:     device.Name,
		RadikoAuthTokenHeader:  device.Token,
		RadikoUserHeader:       device.UserID,
		RadikoLocationHeader:   location,
		RadikoConnectionHeader: device.Connection,
		RadikoPartialKeyHeader: partial,
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	return device
}
