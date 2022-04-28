package apstraTelemetry

type AosStreamingConfigResponse struct {
	Items []AosStreamingConfigItem `json:"items"`
}

type AosStreamingConfigItem struct {
	Status         AosStreamingConfigStatus `json:"status"`
	StreamingType  string                   `json:"streaming_type"`
	SequencingMode string                   `json:"sequencing_mode"`
	Protocol       string                   `json:"protocol"`
	Hostname       string                   `json:"hostname"`
	Id             string                   `json:"id"`
	Port           uint16                   `json:"Port"`
}

type AosStreamingConfigStatus struct {
	Status               AosStreamingConfigConnectionLog     `json:"status"`
	ConnectionTime       string                              `json:"connectionTime"`
	Epoch                string                              `json:"epoch"`
	ConnectionResetCount uint                                `json:"connnectionResetCount"`
	StreamingEndpoint    AosStreamingConfigStreamingEndpoint `json:"streamingEndpoint"`
	DnsLog               AosStreamingConfigDnsLog            `json:"dnsLog"`
	Connected            bool                                `json:"connected"`
	DisconnectionTime    string                              `json:"disconnectionTime"`
}

type AosStreamingConfigConnectionLog struct {
	Date    string `json:"date"'`
	Message string `json:"message"`
}

type AosStreamingConfigStreamingEndpoint struct {
	StreamingType  string `json:"streaming_type"`
	SequencingMode string `json:"sequencing_mode"`
	Protocol       string `json:"protocol"`
	Hostname       string `json:"Hostname"`
	Port           uint16 `json:"Port"`
}

type AosStreamingConfigDnsLog struct {
	Date    string `json:"date"`
	Message string `json:"message"`
}
