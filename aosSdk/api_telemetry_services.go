package aosSdk

const (
	apiUrlTelemetryServices = "/api/telemetry/services"
)

type TelemetryServiceMapping struct {
	ServiceName string   `json:"service_name"`
	Devices     []string `json:"devices"`
}

type GetTelemetryServiceMappingResult struct {
	Configured    []TelemetryServiceMapping `json:"configured"`
	LastRunError  []TelemetryServiceMapping `json:"last_run_error"`
	EnablingError []TelemetryServiceMapping `json:"enabling_error"`
}

func (o Client) GetTelemetryServicesDeviceMapping() (*GetTelemetryServiceMappingResult, error) {
	var result GetTelemetryServiceMappingResult
	_, err := o.talkToAos(&talkToAosIn{
		method:        httpMethodGet,
		url:           apiUrlTelemetryServices,
		fromServerPtr: &result,
	})
	return &result, err
}
