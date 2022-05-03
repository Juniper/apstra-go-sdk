package telemetry

type TelemetryServiceMapping struct {
	ServiceName string `json:"service_name"`
	Devices     string `json:"devices"`
}

type GetTelemetryServiceMappingResult struct {
	Configured    []TelemetryServiceMapping `json:"configured"`
	LastRunError  []TelemetryServiceMapping `json:"last_run_error"`
	EnablingError []TelemetryServiceMapping `json:"enabling_error"`
}
