package goapstra

import (
	"context"
	"encoding/json"
	"log"
	"testing"
)

func TestGetVersionsServer(t *testing.T) {
	clients, err := getCloudlabsTestClients()
	if err != nil {
		t.Fatal(err)
	}

	for _, client := range clients {
		aosdi, err := client.getVersionsAosdi(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		api, err := client.getVersionsApi(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		build, err := client.getVersionsBuild(context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		server, err := client.getVersionsServer(context.TODO())
		if err != nil {
			t.Fatal(err)
		}

		body, err := json.Marshal(&struct {
			Aosdi  *versionsAosdiResponse  `json:"aosdi"`
			Api    *versionsApiResponse    `json:"api"`
			Build  *versionsBuildResponse  `json:"build"`
			Server *versionsServerResponse `json:"server"`
		}{
			Aosdi:  aosdi,
			Api:    api,
			Build:  build,
			Server: server,
		})
		if err != nil {
			t.Fatal(err)
		}

		log.Println(string(body))
	}
}
