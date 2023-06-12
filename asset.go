package main

import (
	"context"
	cr "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/yyoshiki41/go-radiko"
	"github.com/yyoshiki41/radigo"
)

type Area struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Asset struct {
	AvailableStations []string
	AreaDevices       Devices
	Coordinates       Coordinates
	DefaultClient     *radiko.Client
	NextFetchTime     *time.Time
	OutputFormat      string
	Regions           Regions
	Rules             Rules
	Schedules         Schedules
	Stations          Stations
	Versions          Versions
}

// AddExtraStations appends stations to AvailableStations
func (a *Asset) AddExtraStations(es []string) {
	for _, s := range es {
		dup := false
		for _, as := range a.AvailableStations {
			if as == s {
				dup = true
				break
			}
		}
		if !dup {
			a.AvailableStations = append(a.AvailableStations, s)
		}
	}
}

// GenerateGPS returns the RadikoLocationHeader GPS string
// e.g., "35.689492,139.691701,gps"
func (a *Asset) GenerateGPSForAreaID(areaID string) string {
	if c, ok := a.Coordinates[areaID]; ok {
		negpos := []float64{-1.0, 1.0}
		// +/- 0 ~ 0.025 --> 0 ~ 1.5' -> +/- 0 ~ 2.77/2.13km
		rand.Shuffle(len(negpos), func(i, j int) {
			negpos[i], negpos[j] = negpos[j], negpos[i]
		})
		lat := c.Lat + (rand.Float64()/40.0)*negpos[0]
		rand.Shuffle(len(negpos), func(i, j int) {
			negpos[i], negpos[j] = negpos[j], negpos[i]
		})
		lng := c.Lng + (rand.Float64()/40.0)*negpos[0]
		return fmt.Sprintf("%.6f,%.6f,gps", lat, lng)
	}
	return ""
}

// GetAreaIDByStationID returns the first AreaID for the station
func (a *Asset) GetAreaIDByStationID(stationID string) string {
	if s, ok := a.Stations[stationID]; ok {
		return s.Areas[0]
	}
	return ""
}

// GetStationIDsByAreaID returns a slice of StationIDs
func (a *Asset) GetStationIDsByAreaID(areaID string) []string {
	sids := []string{}
	for sid, ss := range a.Stations {
		for _, aid := range ss.Areas {
			if aid == areaID {
				sids = append(sids, sid)
			}
		}
	}
	return sids
}

// LoadAvailableStations loads up the avaialable stations
func (a *Asset) LoadAvailableStations(areaID string) {
	// AvailableStations
	a.AvailableStations = a.GetStationIDsByAreaID(areaID)
}

// NewDevice returns a pointer to a new Device
func (a *Asset) NewDevice() *Device {
	// generate userID
	blob := make([]byte, 16)
	if _, err := cr.Read(blob); err != nil {
		return nil
	}
	userID := hex.EncodeToString(blob)

	// new Device placeholder
	device := &Device{
		AppName:    "aSmartPhone7a",
		Connection: "wifi",
		UserID:     userID,
	}

	// get an app version
	as := make([]string, len(a.Versions.Apps))
	copy(as, a.Versions.Apps)
	rand.Shuffle(len(as), func(i, j int) {
		as[i], as[j] = as[j], as[i]
	})
	device.AppVersion = as[0]

	// get an sdk version
	svs := reflect.ValueOf(a.Versions.SDKs).MapKeys()
	rand.Shuffle(len(svs), func(i, j int) {
		svs[i], svs[j] = svs[j], svs[i]
	})
	sdkVersion := svs[0]

	// get an sdk
	sdk := a.Versions.SDKs[sdkVersion.String()]

	// get a sdk build
	bs := make([]string, len(sdk.Builds))
	copy(bs, sdk.Builds)
	rand.Shuffle(len(bs), func(i, j int) {
		bs[i], bs[j] = bs[j], bs[i]
	})
	build := bs[0]

	// get a model
	ms := make([]string, len(a.Versions.Models))
	copy(ms, a.Versions.Models)
	rand.Shuffle(len(ms), func(i, j int) {
		ms[i], ms[j] = ms[j], ms[i]
	})
	model := ms[0]

	// assign the detail
	//Dalvik/2.1.0 (Linux; U; Android %SDK_VERSION%; %MODEL%/%BUILD%)
	device.UserAgent = fmt.Sprintf(
		"Dalvik/2.1.0 (Linux; U; Android %s; %s/%s)",
		sdkVersion,
		model,
		build,
	)
	//X-Radiko-Device: %SDK_ID%.%MODEL%
	device.Name = fmt.Sprintf("%s.%s", sdk.ID, model)

	return device
}

// UnmarshalJSON loads up Coordinates with Regions
func (a *Asset) UnmarshalJSON(b []byte) error {
	var cs map[string][]float64
	if err := json.Unmarshal(b, &cs); err != nil {
		return err
	}

	a.Coordinates = Coordinates{}
	for areaName, latlng := range cs {
		for _, areas := range a.Regions {
			for _, area := range areas {
				if areaName == area.Name {
					a.Coordinates[area.ID] = &Coordinate{
						Lat: latlng[0],
						Lng: latlng[1],
					}
				}
			}
		}
	}
	return nil
}

type ContextKey string

type Coordinate struct {
	Lat float64
	Lng float64
}

type Coordinates map[string]*Coordinate

type Device struct {
	AppName    string
	AppVersion string
	Connection string
	Name       string
	Token      string
	UserAgent  string
	UserID     string
}

type Devices map[string]*Device

type Regions map[string][]Area

type SDK struct {
	ID     string   `json:"sdk"`
	Builds []string `json:"builds"`
}

type Station struct {
	Areas []string
	Name  string
	Ruby  string
}

type Stations map[string]*Station

type Versions struct {
	Apps   []string        `json:"apps"`
	Models []string        `json:"models"`
	SDKs   map[string]*SDK `json:"sdks"`
}

type XMLRegion struct {
	Region []XMLStations `xml:"stations"`
}

type XMLStations struct {
	Stations   []XMLStation `xml:"station"`
	RegionID   string       `xml:"region_id,attr"`
	RegionName string       `xml:"region_name,attr"`
}

type XMLStation struct {
	ID     string `xml:"id"`
	Name   string `xml:"name"`
	AreaID string `xml:"area_id"`
	Ruby   string `xml:"ruby"`
}

func FetchXMLRegion() (XMLRegion, error) {
	region := XMLRegion{}

	response, err := http.Get(RegionFullAPI)
	if err != nil {
		return region, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return region, err
	}
	defer response.Body.Close()

	err = xml.Unmarshal([]byte(string(body)), &region)
	if err != nil {
		return region, err
	}

	return region, nil
}

func GetAsset(ctx context.Context) *Asset {
	k := ContextKey("asset")
	asset, ok := ctx.Value(k).(*Asset)
	if !ok {
		return nil
	}
	return asset
}

func NewAsset(client *radiko.Client) (*Asset, error) {
	asset := &Asset{}
	// empty AreaDevices
	asset.AreaDevices = map[string]*Device{}
	// default client
	asset.DefaultClient = client
	// empty FileFormat
	asset.OutputFormat = radigo.AudioFormatAAC
	// nil *time.Time
	asset.NextFetchTime = nil
	// empty Schedules
	asset.Schedules = Schedules{}

	// Region
	regionsJSON, err := os.Open("assets/regions.json")
	if err != nil {
		return asset, err
	}
	defer regionsJSON.Close()
	blob, _ := io.ReadAll(regionsJSON)
	err = json.Unmarshal(blob, &asset.Regions)
	if err != nil {
		return asset, err
	}

	// Coordinate
	coordinatesJSON, err := os.Open("assets/coordinates.json")
	if err != nil {
		return asset, err
	}
	defer coordinatesJSON.Close()
	blob, _ = io.ReadAll(coordinatesJSON)
	err = json.Unmarshal(blob, asset)
	if err != nil {
		return asset, err
	}

	// Station
	xmlRegion, err := FetchXMLRegion()
	if err != nil {
		return asset, err
	}
	asset.Stations = Stations{}
	for _, xmlStations := range xmlRegion.Region {
		for _, xmlStation := range xmlStations.Stations {
			if station, ok := asset.Stations[xmlStation.ID]; ok {
				station.Areas = append(station.Areas, xmlStation.AreaID)
			} else {
				station := &Station{
					Areas: []string{xmlStation.AreaID},
					Name:  xmlStation.Name,
					Ruby:  xmlStation.Ruby,
				}
				asset.Stations[xmlStation.ID] = station
			}
		}
	}

	// Versions
	versionsJSON, err := os.Open("assets/versions.json")
	if err != nil {
		return asset, err
	}
	defer versionsJSON.Close()
	blob, _ = io.ReadAll(versionsJSON)
	err = json.Unmarshal(blob, &asset.Versions)
	if err != nil {
		return asset, err
	}

	return asset, nil
}
