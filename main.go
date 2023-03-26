package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/robertkrimen/otto"
	"golang.org/x/oauth2"

	storage "cloud.google.com/go/bigquery/storage/apiv1"
)

const tadoConfigEndpoint = "https://app.tado.com/env.js"

func getTadoConfig(tadoClient *http.Client) TadoConfig {
	req, err := http.NewRequest("GET", tadoConfigEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := tadoClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	defer res.Body.Close()
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	vm := otto.New()
	vm.Run(body)
	vm.Run(`
		var jsonConfig = JSON.stringify(TD);
	`)

	value, parseErr := vm.Get("jsonConfig")
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	jsonString, stringErr := value.ToString()
	if stringErr != nil {
		log.Fatal(stringErr)
	}

	var parsedConfig TadoConfig
	jsonErr := json.Unmarshal([]byte(jsonString), &parsedConfig)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return parsedConfig
}

func tadoApiUrl(config TadoConfig, format string, a ...interface{}) *url.URL {

	u, err := url.Parse(config.Config.TgaRestAPIV2Endpoint)
	if err != nil {
		panic(err)
	}

	apiUrl := u.JoinPath(fmt.Sprintf(format, a...))

	return apiUrl
}

func getTadoApi(c *http.Client, url string, target interface{}) error {
	res, err := c.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return fmt.Errorf("cannot unmarshal response from tado api: %w", err)
	}

	return nil
}

func writeBigqueryRow(ctx context.Context) {
	client, err := storage.NewBigQueryWriteClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
}

func main() {
	username := os.Getenv("TADO_USERNAME")
	password := os.Getenv("TADO_PASSWORD")
	homename := os.Getenv("TADO_HOME")

	ctx := context.Background()

	tadoClient := &http.Client{
		Timeout: time.Second * 5,
	}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, tadoClient)
	tadoConfig := getTadoConfig(tadoClient)
	tadoEndpoint := oauth2.Endpoint{
		TokenURL: tadoConfig.Config.Oauth.APIEndpoint + "/token",
		AuthURL:  tadoConfig.Config.Oauth.APIEndpoint + "/authorize",
	}

	tadoOAuthConfig := &oauth2.Config{
		ClientID:     tadoConfig.Config.Oauth.ClientID,
		ClientSecret: tadoConfig.Config.Oauth.ClientSecret,
		Endpoint:     tadoEndpoint,
	}

	token, err := tadoOAuthConfig.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		log.Fatal(err)
	}
	tadoOAuthClient := tadoOAuthConfig.Client(ctx, token)

	var tadoMe TadoMe
	if err := getTadoApi(tadoOAuthClient, tadoApiUrl(tadoConfig, "me").String(), &tadoMe); err != nil {
		log.Fatal(err)
	}

	var homeId int
	for _, h := range tadoMe.Homes {
		if h.Name == homename {
			homeId = h.ID
		}
	}

	var tadoZones TadoZone
	zonesUrl := tadoApiUrl(tadoConfig, "/homes/%d/zones", homeId)
	if err := getTadoApi(tadoOAuthClient, zonesUrl.String(), &tadoZones); err != nil {
		log.Fatal(err)
	}

	var dayReports []TadoZoneDayReport
	queryParams := url.Values{}
	yesterday := time.Now().AddDate(0, 0, -1)
	queryParams.Add("date", yesterday.Format("2006-01-02"))

	//For each HEATING zone, get a dayReport for yesterday
	for _, z := range tadoZones {
		if z.Type == "HEATING" {
			reportUrl := tadoApiUrl(tadoConfig, "/homes/%d/zones/%d/dayReport", homeId, z.ID)
			reportUrl.RawQuery = queryParams.Encode()

			var currentZoneDayReport TadoZoneDayReport

			if err := getTadoApi(tadoOAuthClient, reportUrl.String(), &currentZoneDayReport); err != nil {
				log.Fatal(err)
			}

			dayReports = append(dayReports, currentZoneDayReport)
		}
	}
	log.Println(dayReports)
}
