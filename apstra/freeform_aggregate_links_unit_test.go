package apstra_test

import (
	"encoding/json"
	"github.com/Juniper/apstra-go-sdk/apstra"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFreeformAggregateLinkRequest_MarshalJSON(t *testing.T) {
	type testcase struct {
		data     apstra.FreeformAggregateLinkData
		expected string
	}
	testcases := map[string]testcase{
		"foo": {
			data: apstra.FreeformAggregateLinkData{
				Endpoints: [2][]apstra.FreeformAggregateLinkMemberEndpoint{
					[]apstra.FreeformAggregateLinkMemberEndpoint{
						{
							SystemId:      "system1",
							PortChannelId: 1,
							LagMode:       apstra.RackLinkLagModeActive,
						},
						{
							SystemId:      "system2",
							PortChannelId: 2,
							LagMode:       apstra.RackLinkLagModeActive,
						},
					},
					[]apstra.FreeformAggregateLinkMemberEndpoint{
						{
							SystemId:      "system3",
							PortChannelId: 3,
							LagMode:       apstra.RackLinkLagModePassive,
						},
						{
							SystemId:      "system4",
							PortChannelId: 4,
							LagMode:       apstra.RackLinkLagModePassive,
						},
					},
				},
				MemberLinkIds: []apstra.ObjectId{"LinkId1", "LinkId2", "LinkId3"},
			},
			expected: `{
                "label": "",
                "endpoints": [
                  {
                    "system": {
                      "id": "system1"
                    },
                    "interface": {
                      "port_channel_id": 1,
                      "lag_mode": "lacp_active"
                    },
                    "endpoint_group": 0
                  },
                  {
                    "system": {
                      "id": "system2"
                    },
                    "interface": {
                      "port_channel_id": 2,
                      "lag_mode": "lacp_active"
                    },
                    "endpoint_group": 0
                  },
                  {
                    "system": {
                      "id": "system3"
                    },
                    "interface": {
                      "port_channel_id": 3,
                      "lag_mode": "lacp_passive"
                    },
                    "endpoint_group": 1
                  },
                  {
                    "system": {
                      "id": "system4"
                    },
                    "interface": {
                      "port_channel_id": 4,
                      "lag_mode": "lacp_passive"
                    },
                    "endpoint_group": 1
                  }
                ],
                "member_link_ids": [
                  "LinkId1",
                  "LinkId2",
                  "LinkId3"
                ]
              }`,
		},
	}

	for tName, tCase := range testcases {
		t.Run(tName, func(t *testing.T) {
			result, err := json.Marshal(tCase.data)
			require.NoError(t, err)
			require.JSONEq(t, tCase.expected, string(result))

			var clone apstra.FreeformAggregateLink
			err = json.Unmarshal(result, &clone)
			require.NoError(t, err)
			require.NotNil(t, *clone.Data)
			require.Equal(t, tCase.data, *clone.Data)
		})
	}
}
