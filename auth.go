package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
)

// NewAuth retuns a new authorized device with token
func NewAuth(ctx context.Context, areaID string) error {
	asset := GetAsset(ctx)
	client := asset.DefaultClient
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
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	device.Token = resp.Header.Get(RadikoAuthTokenHeader)
	offset, _ := strconv.ParseInt(resp.Header.Get(RadikoKeyOffsetHeader), 10, 64)
	length, _ := strconv.ParseInt(resp.Header.Get(RadikoKeyLentghHeader), 10, 64)
	blob, _ := os.ReadFile("assets/base64-full.key")
	authKey, _ := base64.StdEncoding.DecodeString(string(blob))
	partialKey := base64.StdEncoding.EncodeToString([]byte(authKey[offset : offset+length]))

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
		RadikoPartialKeyHeader: partialKey,
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}

	asset.AreaDevices[areaID] = device
	return nil
}
