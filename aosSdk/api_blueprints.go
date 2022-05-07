package aosSdk

const (
	apiUrlBlueprints = "/api/blueprints"
)

type BlueprintStatus string

type AdditionalProp struct {
	NumSucceeded int `json:"num_succeeded"`
	NumFailed    int `json:"num_failed"`
	NumPending   int `json:"num_pending"`
}

type DeploymentStatus struct {
	AdditionalProp1 AdditionalProp `json:"additionalProp1"`
	AdditionalProp2 AdditionalProp `json:"additionalProp2"`
	AdditionalProp3 AdditionalProp `json:"additionalProp3"`
}

type AnomalyCounts struct {
	AdditionalProp1 int `json:"additionalProp1"`
	AdditionalProp2 int `json:"additionalProp2"`
	AdditionalProp3 int `json:"additionalProp3"`
}

type Blueprint struct {
	Status           BlueprintStatus  `json:"status"`
	Version          int              `json:"version"`
	Design           string           `json:"design"`
	DeploymentStatus DeploymentStatus `json:"deployment_status"`
	AnomalyCounts    AnomalyCounts    `json:"anomaly_counts"`
	Id               string           `json:"id"`
	LastModifiedAt   string           `json:"last_modified_at"`
	Label            string           `json:"label"`
}

type GetBlueprintsResult struct {
	Items []Blueprint `json:"items"`
}

type BlueprintRelationships struct {
	AdditionalProp1 BlueprintRelationship `json:"additionalProp1"`
	AdditionalProp2 BlueprintRelationship `json:"additionalProp2"`
	AdditionalProp3 BlueprintRelationship `json:"additionalProp3"`
}

type BlueprintRelationship struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
	Type     string `json:"type"`
	Id       string `json:"id"`
}

type BlueprintNodes struct {
	AdditionalProp1 BlueprintNode `json:"additionalProp1"`
	AdditionalProp2 BlueprintNode `json:"additionalProp2"`
	AdditionalProp3 BlueprintNode `json:"additionalProp3"`
}

type BlueprintNode struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type BlueprintSourceVersions struct {
	AdditionalProp1 int `json:"additionalProp1"`
	AdditionalProp2 int `json:"additionalProp2"`
	AdditionalProp3 int `json:"additionalProp3"`
}

type BlueprintData struct {
	Relationships  BlueprintRelationships  `json:"relationships"`
	Version        int                     `json:"version"`
	Design         string                  `json:"design"`
	LastModifiedAt string                  `json:"last_modified_at"`
	Nodes          BlueprintNodes          `json:"nodes"`
	Id             string                  `json:"id"`
	SourceVersions BlueprintSourceVersions `json:"source_versions"`
}

// todo restore this function
//func (o Client) GetBlueprints() (*GetBlueprintsResult, error) {
//	var result GetBlueprintsResult
//	return &result, o.talkToAos(&talkToAosIn{
//		method:        httpMethodGet,
//		url:           apiUrlBlueprints,
//		fromServerPtr: &result,
//	})
//}
