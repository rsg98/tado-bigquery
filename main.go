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
	"google.golang.org/protobuf/proto"

	"cloud.google.com/go/bigquery/storage/apiv1/storagepb"
	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"

	pb "tado-bigquery/tadodailyreport"
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

func writeDailyReportsToBq(ctx context.Context, data [][]byte) error {
	project := os.Getenv("BQ_PROJECT_ID")
	dataset := os.Getenv("BQ_DATASET")
	table := os.Getenv("BQ_TABLE")
	tableName := fmt.Sprintf("projects/%s/datasets/%s/tables/%s", project, dataset, table)

	client, err := managedwriter.NewClient(ctx, project)
	if err != nil {
		return fmt.Errorf("managedwriter.NewClient: %w", err)
	}
	defer client.Close()

	// get the protobuf descriptor
	m := &pb.Tadodailyreport{}
	descriptor, err := adapt.NormalizeDescriptor(m.ProtoReflect().Descriptor())
	if err != nil {
		return fmt.Errorf("NormalizeDescriptor: %w", err)
	}

	managedStream, err := client.NewManagedStream(ctx,
		managedwriter.WithDestinationTable(tableName),
		managedwriter.WithType(managedwriter.PendingStream),
		managedwriter.WithSchemaDescriptor(descriptor))
	if err != nil {
		return fmt.Errorf("NewManagedStream: %w", err)
	}

	result, err := managedStream.AppendRows(ctx, data, managedwriter.WithOffset(0))
	if err != nil {
		return fmt.Errorf("append %d returned error: %w", 0, err)
	}

	recvOffset, err := result.GetResult(ctx)
	if err != nil {
		return fmt.Errorf("append returned error: %w", err)
	}
	log.Printf("Successfully appended data at offset %d.\n", recvOffset)

	rowCount, err := managedStream.Finalize(ctx)
	if err != nil {
		return fmt.Errorf("error during Finalize: %w", err)
	}

	log.Printf("Stream %s finalized with %d rows.\n", managedStream.StreamName(), rowCount)

	// Run a batch commit.
	req := &storagepb.BatchCommitWriteStreamsRequest{
		Parent:       managedwriter.TableParentFromStreamName(managedStream.StreamName()),
		WriteStreams: []string{managedStream.StreamName()},
	}

	resp, err := client.BatchCommitWriteStreams(ctx, req)
	if err != nil {
		return fmt.Errorf("client.BatchCommit: %w", err)
	}
	if len(resp.GetStreamErrors()) > 0 {
		return fmt.Errorf("stream errors present: %v", resp.GetStreamErrors())
	}

	log.Printf("Table data committed at %s\n", resp.GetCommitTime().AsTime().Format(time.RFC3339Nano))

	return nil

}

func parseZoneReportToProto(zoneReports []TadoZoneDetailsReport) ([][]byte, error) {

	msgs := make([][]byte, len(zoneReports))

	for n, z := range zoneReports {

		// extract all the set temperatures from the Tado report
		settingsData := z.Report.Settings.DataIntervals
		settingsList := make([]*pb.Tadodailyreport_Settings, len(settingsData))
		for i := 0; i < len(settingsData); i++ {
			settingsList[i] = &pb.Tadodailyreport_Settings{
				From:        proto.Int64(settingsData[i].From.Unix()),
				To:          proto.Int64(settingsData[i].To.Unix()),
				Temperature: proto.Float32(float32(settingsData[i].Value.Temperature.Celsius)),
			}
		}

		// extract the weather data from the Tado report
		weatherData := z.Report.Weather.Condition.DataIntervals
		weatherList := make([]*pb.Tadodailyreport_Weather, len(weatherData))
		for i := 0; i < len(weatherData); i++ {
			weatherList[i] = &pb.Tadodailyreport_Weather{
				From:        proto.Int64(weatherData[i].From.Unix()),
				To:          proto.Int64(weatherData[i].To.Unix()),
				Temperature: proto.Float32(float32(weatherData[i].Value.Temperature.Celsius)),
				State:       proto.String(weatherData[i].Value.State),
			}
		}

		// extract when zones called for heat and at what rate
		callheatData := z.Report.CallForHeat.DataIntervals
		callheatList := make([]*pb.Tadodailyreport_Callforheat, len(callheatData))
		for i := 0; i < len(callheatData); i++ {
			callheatList[i] = &pb.Tadodailyreport_Callforheat{
				From:     proto.Int64(callheatData[i].From.Unix()),
				To:       proto.Int64(callheatData[i].To.Unix()),
				HeatRate: proto.String(callheatData[i].Value),
			}
		}

		humidityData := z.Report.MeasuredData.Humidity.DataPoints
		humidityList := make([]*pb.Tadodailyreport_Measureddata_Humidity, len(humidityData))
		for i := 0; i < len(humidityData); i++ {
			humidityList[i] = &pb.Tadodailyreport_Measureddata_Humidity{
				Timestamp: proto.Int64(humidityData[i].Timestamp.Unix()),
				Humidity:  proto.Float32(float32(humidityData[i].Value)),
			}
		}

		insidetempData := z.Report.MeasuredData.InsideTemperature.DataPoints
		insidetempList := make([]*pb.Tadodailyreport_Measureddata_Insidetemperature, len(insidetempData))
		for i := 0; i < len(insidetempData); i++ {
			insidetempList[i] = &pb.Tadodailyreport_Measureddata_Insidetemperature{
				Timestamp:   proto.Int64(insidetempData[i].Timestamp.Unix()),
				Temperature: proto.Float32(float32(insidetempData[i].Value.Celsius)),
			}
		}

		measuredList := make([]*pb.Tadodailyreport_Measureddata, 1)
		measuredList[0] = &pb.Tadodailyreport_Measureddata{
			InsideTemperature: insidetempList,
			Humidity:          humidityList,
		}

		reportInterval := &pb.Tadodailyreport_Interval{
			From: proto.Int64(z.Report.Interval.From.Unix()),
			To:   proto.Int64(z.Report.Interval.To.Unix()),
		}

		dayReportRow := &pb.Tadodailyreport{
			ZoneId:       proto.Int64(int64(z.ZoneID)),
			ZoneName:     proto.String(z.ZoneName),
			Interval:     reportInterval,
			MeasuredData: measuredList,
			Settings:     settingsList,
			CallForHeat:  callheatList,
			Weather:      weatherList,
		}

		b, err := proto.Marshal(dayReportRow)
		if err != nil {
			return nil, fmt.Errorf("error generating message %d: %w", n, err)
		}
		msgs[n] = b

	} // for range zoneReports

	return msgs, nil
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

	var zoneReports []TadoZoneDetailsReport
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

			zoneReport := TadoZoneDetailsReport{
				ZoneID:   z.ID,
				ZoneName: z.Name,
				Report:   currentZoneDayReport,
			}

			zoneReports = append(zoneReports, zoneReport)
		}
	}

	//Pass labelled zone reports to row builder
	reportsProtoData, err := parseZoneReportToProto(zoneReports)
	if err != nil {
		log.Fatal("parseZoneReportToProto: ", err)
	}

	err = writeDailyReportsToBq(ctx, reportsProtoData)
	if err != nil {
		log.Fatal("writeDailyReportsToBq: ", err)
	}
}
