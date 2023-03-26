package main

import "time"

//Generated thanks to Matt Holt's conversion tool - https://mholt.github.io/json-to-go/

type TadoConfig struct {
	Config struct {
		Version                                 string `json:"version"`
		Environment                             string `json:"environment"`
		DebugEnabled                            bool   `json:"debugEnabled"`
		LogEndpoint                             string `json:"logEndpoint"`
		BaseURL                                 string `json:"baseUrl"`
		TgaEndpoint                             string `json:"tgaEndpoint"`
		TgaRestAPIEndpoint                      string `json:"tgaRestApiEndpoint"`
		TgaRestAPIV2Endpoint                    string `json:"tgaRestApiV2Endpoint"`
		SusiAPIEndpoint                         string `json:"susiApiEndpoint"`
		HomeBackendBaseURL                      string `json:"homeBackendBaseUrl"`
		HvacAPIEndpoint                         string `json:"hvacApiEndpoint"`
		HvacIncludeInstallFlowsUnderDevelopment bool   `json:"hvacIncludeInstallFlowsUnderDevelopment"`
		GenieRestAPIV2Endpoint                  string `json:"genieRestApiV2Endpoint"`
		IvarRestAPIEndpoint                     string `json:"ivarRestApiEndpoint"`
		MinderRestAPIEndpoint                   string `json:"minderRestApiEndpoint"`
		GaTrackingID                            string `json:"gaTrackingId"`
		Oauth                                   struct {
			ClientAPIEndpoint string `json:"clientApiEndpoint"`
			APIEndpoint       string `json:"apiEndpoint"`
			ClientID          string `json:"clientId"`
			ClientSecret      string `json:"clientSecret"`
		} `json:"oauth"`
	} `json:"config"`
}

type TadoMe struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	ID       string `json:"id"`
	Homes    []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"homes"`
	Locale        string `json:"locale"`
	MobileDevices []struct {
		Name     string `json:"name"`
		ID       int    `json:"id"`
		Settings struct {
			GeoTrackingEnabled          bool `json:"geoTrackingEnabled"`
			SpecialOffersEnabled        bool `json:"specialOffersEnabled"`
			OnDemandLogRetrievalEnabled bool `json:"onDemandLogRetrievalEnabled"`
			PushNotifications           struct {
				LowBatteryReminder          bool `json:"lowBatteryReminder"`
				AwayModeReminder            bool `json:"awayModeReminder"`
				HomeModeReminder            bool `json:"homeModeReminder"`
				OpenWindowReminder          bool `json:"openWindowReminder"`
				EnergySavingsReportReminder bool `json:"energySavingsReportReminder"`
				IncidentDetection           bool `json:"incidentDetection"`
				EnergyIqReminder            bool `json:"energyIqReminder"`
			} `json:"pushNotifications"`
		} `json:"settings"`
		DeviceMetadata struct {
			Platform  string `json:"platform"`
			OsVersion string `json:"osVersion"`
			Model     string `json:"model"`
			Locale    string `json:"locale"`
		} `json:"deviceMetadata"`
	} `json:"mobileDevices"`
}

type TadoHome struct {
	ID                         int         `json:"id"`
	Name                       string      `json:"name"`
	DateTimeZone               string      `json:"dateTimeZone"`
	DateCreated                time.Time   `json:"dateCreated"`
	TemperatureUnit            string      `json:"temperatureUnit"`
	Partner                    interface{} `json:"partner"`
	SimpleSmartScheduleEnabled bool        `json:"simpleSmartScheduleEnabled"`
	AwayRadiusInMeters         float64     `json:"awayRadiusInMeters"`
	InstallationCompleted      bool        `json:"installationCompleted"`
	IncidentDetection          struct {
		Supported bool `json:"supported"`
		Enabled   bool `json:"enabled"`
	} `json:"incidentDetection"`
	Skills                  []string `json:"skills"`
	ChristmasModeEnabled    bool     `json:"christmasModeEnabled"`
	ShowAutoAssistReminders bool     `json:"showAutoAssistReminders"`
	ContactDetails          struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	} `json:"contactDetails"`
	Address struct {
		AddressLine1 string      `json:"addressLine1"`
		AddressLine2 interface{} `json:"addressLine2"`
		ZipCode      string      `json:"zipCode"`
		City         string      `json:"city"`
		State        interface{} `json:"state"`
		Country      string      `json:"country"`
	} `json:"address"`
	Geolocation struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"geolocation"`
	ConsentGrantSkippable bool     `json:"consentGrantSkippable"`
	EnabledFeatures       []string `json:"enabledFeatures"`
	IsAirComfortEligible  bool     `json:"isAirComfortEligible"`
	IsBalanceAcEligible   bool     `json:"isBalanceAcEligible"`
	IsEnergyIqEligible    bool     `json:"isEnergyIqEligible"`
}

type TadoZone []struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	DateCreated time.Time `json:"dateCreated"`
	DeviceTypes []string  `json:"deviceTypes"`
	Devices     []struct {
		DeviceType       string `json:"deviceType"`
		SerialNo         string `json:"serialNo"`
		ShortSerialNo    string `json:"shortSerialNo"`
		CurrentFwVersion string `json:"currentFwVersion"`
		ConnectionState  struct {
			Value     bool      `json:"value"`
			Timestamp time.Time `json:"timestamp"`
		} `json:"connectionState"`
		Characteristics struct {
			Capabilities []string `json:"capabilities"`
		} `json:"characteristics"`
		BatteryState string   `json:"batteryState"`
		Duties       []string `json:"duties"`
	} `json:"devices"`
	ReportAvailable   bool `json:"reportAvailable"`
	ShowScheduleSetup bool `json:"showScheduleSetup"`
	SupportsDazzle    bool `json:"supportsDazzle"`
	DazzleEnabled     bool `json:"dazzleEnabled"`
	DazzleMode        struct {
		Supported bool `json:"supported"`
		Enabled   bool `json:"enabled"`
	} `json:"dazzleMode"`
	OpenWindowDetection struct {
		Supported        bool `json:"supported"`
		Enabled          bool `json:"enabled"`
		TimeoutInSeconds int  `json:"timeoutInSeconds"`
	} `json:"openWindowDetection"`
}

type TadoZoneState struct {
	TadoMode                       string      `json:"tadoMode"`
	GeolocationOverride            bool        `json:"geolocationOverride"`
	GeolocationOverrideDisableTime interface{} `json:"geolocationOverrideDisableTime"`
	Preparation                    interface{} `json:"preparation"`
	Setting                        struct {
		Type        string `json:"type"`
		Power       string `json:"power"`
		Temperature struct {
			Celsius    float64 `json:"celsius"`
			Fahrenheit float64 `json:"fahrenheit"`
		} `json:"temperature"`
	} `json:"setting"`
	OverlayType        interface{} `json:"overlayType"`
	Overlay            interface{} `json:"overlay"`
	OpenWindow         interface{} `json:"openWindow"`
	NextScheduleChange struct {
		Start   time.Time `json:"start"`
		Setting struct {
			Type        string `json:"type"`
			Power       string `json:"power"`
			Temperature struct {
				Celsius    float64 `json:"celsius"`
				Fahrenheit float64 `json:"fahrenheit"`
			} `json:"temperature"`
		} `json:"setting"`
	} `json:"nextScheduleChange"`
	NextTimeBlock struct {
		Start time.Time `json:"start"`
	} `json:"nextTimeBlock"`
	Link struct {
		State string `json:"state"`
	} `json:"link"`
	ActivityDataPoints struct {
		HeatingPower struct {
			Type       string    `json:"type"`
			Percentage float64   `json:"percentage"`
			Timestamp  time.Time `json:"timestamp"`
		} `json:"heatingPower"`
	} `json:"activityDataPoints"`
	SensorDataPoints struct {
		InsideTemperature struct {
			Celsius    float64   `json:"celsius"`
			Fahrenheit float64   `json:"fahrenheit"`
			Timestamp  time.Time `json:"timestamp"`
			Type       string    `json:"type"`
			Precision  struct {
				Celsius    float64 `json:"celsius"`
				Fahrenheit float64 `json:"fahrenheit"`
			} `json:"precision"`
		} `json:"insideTemperature"`
		Humidity struct {
			Type       string    `json:"type"`
			Percentage float64   `json:"percentage"`
			Timestamp  time.Time `json:"timestamp"`
		} `json:"humidity"`
	} `json:"sensorDataPoints"`
}

type TadoZoneDayReport struct {
	ZoneType string `json:"zoneType"`
	Interval struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
	} `json:"interval"`
	HoursInDay   int `json:"hoursInDay"`
	MeasuredData struct {
		MeasuringDeviceConnected struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			DataIntervals  []struct {
				From  time.Time `json:"from"`
				To    time.Time `json:"to"`
				Value bool      `json:"value"`
			} `json:"dataIntervals"`
		} `json:"measuringDeviceConnected"`
		InsideTemperature struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			Min            struct {
				Celsius    float64 `json:"celsius"`
				Fahrenheit float64 `json:"fahrenheit"`
			} `json:"min"`
			Max struct {
				Celsius    float64 `json:"celsius"`
				Fahrenheit float64 `json:"fahrenheit"`
			} `json:"max"`
			DataPoints []struct {
				Timestamp time.Time `json:"timestamp"`
				Value     struct {
					Celsius    float64 `json:"celsius"`
					Fahrenheit float64 `json:"fahrenheit"`
				} `json:"value"`
			} `json:"dataPoints"`
		} `json:"insideTemperature"`
		Humidity struct {
			TimeSeriesType string  `json:"timeSeriesType"`
			ValueType      string  `json:"valueType"`
			PercentageUnit string  `json:"percentageUnit"`
			Min            float64 `json:"min"`
			Max            float64 `json:"max"`
			DataPoints     []struct {
				Timestamp time.Time `json:"timestamp"`
				Value     float64   `json:"value"`
			} `json:"dataPoints"`
		} `json:"humidity"`
	} `json:"measuredData"`
	Stripes struct {
		TimeSeriesType string `json:"timeSeriesType"`
		ValueType      string `json:"valueType"`
		DataIntervals  []struct {
			From  time.Time `json:"from"`
			To    time.Time `json:"to"`
			Value struct {
				StripeType string `json:"stripeType"`
				Setting    struct {
					Type        string `json:"type"`
					Power       string `json:"power"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"setting"`
			} `json:"value"`
		} `json:"dataIntervals"`
	} `json:"stripes"`
	Settings struct {
		TimeSeriesType string `json:"timeSeriesType"`
		ValueType      string `json:"valueType"`
		DataIntervals  []struct {
			From  time.Time `json:"from"`
			To    time.Time `json:"to"`
			Value struct {
				Type        string `json:"type"`
				Power       string `json:"power"`
				Temperature struct {
					Celsius    float64 `json:"celsius"`
					Fahrenheit float64 `json:"fahrenheit"`
				} `json:"temperature"`
			} `json:"value"`
		} `json:"dataIntervals"`
	} `json:"settings"`
	CallForHeat struct {
		TimeSeriesType string `json:"timeSeriesType"`
		ValueType      string `json:"valueType"`
		DataIntervals  []struct {
			From  time.Time `json:"from"`
			To    time.Time `json:"to"`
			Value string    `json:"value"`
		} `json:"dataIntervals"`
	} `json:"callForHeat"`
	Weather struct {
		Condition struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			DataIntervals  []struct {
				From  time.Time `json:"from"`
				To    time.Time `json:"to"`
				Value struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"value"`
			} `json:"dataIntervals"`
		} `json:"condition"`
		Sunny struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			DataIntervals  []struct {
				From  time.Time `json:"from"`
				To    time.Time `json:"to"`
				Value bool      `json:"value"`
			} `json:"dataIntervals"`
		} `json:"sunny"`
		Slots struct {
			TimeSeriesType string `json:"timeSeriesType"`
			ValueType      string `json:"valueType"`
			Slots          struct {
				Zero400 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"04:00"`
				Zero800 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"08:00"`
				One200 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"12:00"`
				One600 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"16:00"`
				Two000 struct {
					State       string `json:"state"`
					Temperature struct {
						Celsius    float64 `json:"celsius"`
						Fahrenheit float64 `json:"fahrenheit"`
					} `json:"temperature"`
				} `json:"20:00"`
			} `json:"slots"`
		} `json:"slots"`
	} `json:"weather"`
}

type TadoZoneDetailsReport struct {
	ZoneID   int
	ZoneName string
	Report   TadoZoneDayReport
}
