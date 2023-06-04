package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/yyoshiki41/go-radiko"
)

func getArea(stationID string) string {
	switch stationID {
	// hokkaido-tohoku
	case "HBC":
		return "JP1" // Hokkaido
	case "STV":
		return "JP1" // Hokkaido
	case "AIR-G":
		return "JP1" // Hokkaido
	case "NORTHWAVE":
		return "JP1" // Hokkaido
	case "RAB":
		return "JP1" // Hokkaido
	case "AFB":
		return "JP1" // Hokkaido
	case "IBC":
		return "JP1" // Hokkaido
	case "FMI":
		return "JP1" // Hokkaido
	case "TBC":
		return "JP1" // Hokkaido
	case "DATEFM":
		return "JP1" // Hokkaido
	case "ABS":
		return "JP1" // Hokkaido
	case "AFM":
		return "JP1" // Hokkaido
	case "YBC":
		return "JP1" // Hokkaido
	case "RFM":
		return "JP1" // Hokkaido
	case "RFC":
		return "JP1" // Hokkaido
	case "FMF":
		return "JP1" // Hokkaido
	case "JOIK":
		return "JP1" // Hokkaido
	case "JOHK":
		return "JP1" // Hokkaido
	// kanto
	case "TBS":
		return "JP13" // Tokyo
	case "QRR":
		return "JP13" // Tokyo
	case "LFR":
		return "JP13" // Tokyo
	case "INT":
		return "JP13" // Tokyo
	case "FMT":
		return "JP13" // Tokyo
	case "FMJ":
		return "JP13" // Tokyo
	case "JORF":
		return "JP13" // Tokyo
	case "BAYFM78":
		return "JP13" // Tokyo
	case "NACK5":
		return "JP13" // Tokyo
	case "YFM":
		return "JP13" // Tokyo
	case "IBS":
		return "JP13" // Tokyo
	case "CRT":
		return "JP13" // Tokyo
	case "RADIOBERRY":
		return "JP13" // Tokyo
	case "FMGUNMA":
		return "JP13" // Tokyo
	case "JOAK":
		return "JP13" // Tokyo
	// hokuriku-koushinetsu
	case "BSN":
		return "JP15" // Niigata
	case "FMNIIGATA":
		return "JP15" // Niigata
	case "KNB":
		return "JP15" // Niigata
	case "FMTOYAMA":
		return "JP15" // Niigata
	case "MRO":
		return "JP15" // Niigata
	case "HELLOFIVE":
		return "JP15" // Niigata
	case "FBC":
		return "JP15" // Niigata
	case "FMFUKUI":
		return "JP15" // Niigata
	case "YBS":
		return "JP15" // Niigata
	case "FM-FUJI":
		return "JP15" // Niigata
	case "SBC":
		return "JP15" // Niigata
	case "FMN":
		return "JP15" // Niigata
	case "JOCK":
		return "JP15" // Niigata
	// chubu
	case "CBC":
		return "JP23" // Aichi
	case "TOKAIRADIO":
		return "JP23" // Aichi
	case "GBS":
		return "JP23" // Aichi
	case "ZIP-FM":
		return "JP23" // Aichi
	case "FMAICHI":
		return "JP23" // Aichi
	case "FMGIFU":
		return "JP23" // Aichi
	case "SBS":
		return "JP23" // Aichi
	case "K-MIX":
		return "JP23" // Aichi
	case "FMMIE":
		return "JP23" // Aichi
	// kinki
	case "ABC":
		return "JP27" // Osaka
	case "MBS":
		return "JP27" // Osaka
	case "OBC":
		return "JP27" // Osaka
	case "CCL":
		return "JP27" // Osaka
	case "802":
		return "JP27" // Osaka
	case "FMO":
		return "JP27" // Osaka
	case "CRK":
		return "JP27" // Osaka
	case "KISSFMKOBE":
		return "JP27" // Osaka
	case "E-RADIO":
		return "JP27" // Osaka
	case "KBS":
		return "JP27" // Osaka
	case "ALPHA-STATION":
		return "JP27" // Osaka
	case "WBS":
		return "JP27" // Osaka
	case "JOBK":
		return "JP27" // Osaka
	// chugoku-shikoku
	case "BSS":
		return "JP33" // Okayama
	case "FM-SANIN":
		return "JP33" // Okayama
	case "RSK":
		return "JP33" // Okayama
	case "FM-OKAYAMA":
		return "JP33" // Okayama
	case "RCC":
		return "JP33" // Okayama
	case "HFM":
		return "JP33" // Okayama
	case "KRY":
		return "JP33" // Okayama
	case "FMY":
		return "JP33" // Okayama
	case "JRT":
		return "JP33" // Okayama
	case "FM807":
		return "JP33" // Okayama
	case "RNC":
		return "JP33" // Okayama
	case "FMKAGAWA":
		return "JP33" // Okayama
	case "RNB":
		return "JP33" // Okayama
	case "JOEU-FM":
		return "JP33" // Okayama
	case "RKC":
		return "JP33" // Okayama
	case "HI-SIX":
		return "JP33" // Okayama
	case "JOFK":
		return "JP33" // Okayama
	case "JOZK":
		return "JP33" // Okayama
	// kyushu
	case "RKB":
		return "JP40" // Fukuoka
	case "KBC":
		return "JP40" // Fukuoka
	case "LOVEFM":
		return "JP40" // Fukuoka
	case "CROSSFM":
		return "JP40" // Fukuoka
	case "FMFUKUOKA":
		return "JP40" // Fukuoka
	case "FMS":
		return "JP40" // Fukuoka
	case "NBC":
		return "JP40" // Fukuoka
	case "FMNAGASAKI":
		return "JP40" // Fukuoka
	case "RKK":
		return "JP40" // Fukuoka
	case "FMK":
		return "JP40" // Fukuoka
	case "OBS":
		return "JP40" // Fukuoka
	case "FM_OITA":
		return "JP40" // Fukuoka
	case "MRT":
		return "JP40" // Fukuoka
	case "JOYFM":
		return "JP40" // Fukuoka
	case "MBC":
		return "JP40" // Fukuoka
	case "MYUFM":
		return "JP40" // Fukuoka
	case "RBC":
		return "JP40" // Fukuoka
	case "ROK":
		return "JP40" // Fukuoka
	case "FM_OKINAWA":
		return "JP40" // Fukuoka
	case "JOLK":
		return "JP40" // Fukuoka
	// zenkoku
	case "RN1":
	case "RN2":
	case "HOUSOU-DAIGAKU":
	case "JOAK-FM":
	}
	return ""
}

func getGPS(areaID string) string {
	switch areaID {
	case "JP1": // Hokkaido
		return "43.064612,141.346801,gps"
	case "JP13": // Tokyo
		return "35.689492,139.691701,gps"
	case "JP15": // Niigata
		return "37.902552,139.023091,gps"
	case "JP23": // Aichi
		return "35.180192,136.906561,gps"
	case "JP27": // Osaka
		return "34.686302,135.519661,gps"
	case "JP33":
		return "34.661752,133.934411,gps"
	case "JP40":
		return "33.606582,130.418301,gps"
	}
	return "35.689492,139.691701,gps"
}

func GetToken(ctx context.Context, client *radiko.Client, areaID string) string {
	var userID string
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		userID = "ba9619c7ac46d61a94aa35e343daba5e"
	}
	userID = hex.EncodeToString(bytes)

	// auth1
	req, _ := http.NewRequest("GET", "https://radiko.jp/v2/api/auth1", nil)
	req = req.WithContext(ctx)
	headers := map[string]string{
		UserAgentHeader:        DalvikAgent,
		RadikoAppHeader:        ASmartPhone7a,
		RadikoAppVersionHeader: ASmartPhone7aVersion,
		RadikoDeviceHeader:     ASmartPhone7aDevice,
		RadikoUserHeader:       userID,
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}

	token := resp.Header.Get(RadikoAuthTokenHeader)
	offset, _ := strconv.ParseInt(resp.Header.Get(RadikoKeyOffsetHeader), 10, 64)
	length, _ := strconv.ParseInt(resp.Header.Get(RadikoKeyLentghHeader), 10, 64)
	fullkey, _ := base64.StdEncoding.DecodeString(FullkeyB64)
	partial := base64.StdEncoding.EncodeToString([]byte(fullkey[offset : offset+length]))

	location := getGPS(areaID)

	req, _ = http.NewRequest("GET", "https://radiko.jp/v2/api/auth2", nil)
	req = req.WithContext(ctx)
	headers = map[string]string{
		UserAgentHeader:        DalvikAgent,
		RadikoAppHeader:        ASmartPhone7a,
		RadikoAppVersionHeader: ASmartPhone7aVersion,
		RadikoDeviceHeader:     ASmartPhone7aDevice,
		RadikoAuthTokenHeader:  token,
		RadikoUserHeader:       userID,
		RadikoLocationHeader:   location,
		RadikoConnectionHeader: WiFiConnection,
		RadikoPartialKeyHeader: partial,
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	return token
}
