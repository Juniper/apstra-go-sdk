package apstra

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

const (
	apiUrlSwitchSystemLinks       = apiUrlBlueprintById + apiUrlPathDelim + "switch-system-links"
	apiUrlDeleteSwitchSystemLinks = apiUrlBlueprintById + apiUrlPathDelim + "delete-switch-system-links"
)

type CreateLinksWithNewServerRequest struct {
	Links  []CreateLinksWithNewServerRequestLink
	Server CreateLinksWithNewServerRequestServer
}

func (o *CreateLinksWithNewServerRequest) raw(ctx context.Context, client *Client) (*rawCreateLinksWithNewServerRequest, error) {
	rs, err := o.Server.raw(ctx, "server", client)
	if err != nil {
		return nil, err
	}

	links := make([]rawCreateLinksWithNewServerRequestLink, len(o.Links))
	for i, link := range o.Links {
		links[i] = *link.raw()
	}

	return &rawCreateLinksWithNewServerRequest{
		NewSystems: []rawCreateLinksWithNewServerRequestServer{*rs},
		Links:      links,
	}, nil
}

type rawCreateLinksWithNewServerRequest struct {
	NewSystems []rawCreateLinksWithNewServerRequestServer `json:"new_systems"`
	Links      []rawCreateLinksWithNewServerRequestLink   `json:"links"`
}

type CreateLinksWithNewServerRequestServer struct {
	Hostname         string
	Label            string
	LogicalDeviceId  ObjectId
	LogicalDevice    *LogicalDevice
	PortChannelIdMin int
	PortChannelIdMax int
	Tags             []string
}

func (o *CreateLinksWithNewServerRequestServer) raw(ctx context.Context, systemType string, client *Client) (*rawCreateLinksWithNewServerRequestServer, error) {
	if o.LogicalDeviceId != "" && o.LogicalDevice != nil {
		return nil, errors.New("both LogicalDevice (payload) and LogicalDeviceId (catalog ID) specified.")
	}

	var err error
	var rawLD *rawLogicalDevice

	if o.LogicalDeviceId != "" {
		rawLD, err = client.getLogicalDevice(ctx, o.LogicalDeviceId)
		if err != nil {
			return nil, fmt.Errorf("error fetching logical device %q - %w", o.LogicalDeviceId, err)
		}
		// wipe out the timestamps so we don't send 'em back to Apstra
		rawLD.CreatedAt = nil
		rawLD.LastModifiedAt = nil
	}

	if o.LogicalDevice != nil {
		rawLD = &rawLogicalDevice{
			DisplayName: "TF LD template",
			Id:          "tf-ld-template",
			Panels: []rawLogicalDevicePanel{{
				PanelLayout: LogicalDevicePanelLayout{RowCount: 1, ColumnCount: 1},
				PortIndexing: LogicalDevicePortIndexing{
					Order:      PortIndexingHorizontalFirst,
					StartIndex: 1,
					Schema:     PortIndexingSchemaAbsolute,
				},
				PortGroups: []rawLogicalDevicePortGroup{{
					Count: 1,
					Speed: rawLogicalDevicePortSpeed{Unit: "M", Value: 100},
				}},
			}},
		}
	}

	return &rawCreateLinksWithNewServerRequestServer{
		SystemType:       systemType,
		LogicalDevice:    *rawLD,
		PortChannelIdMin: o.PortChannelIdMin,
		PortChannelIdMax: o.PortChannelIdMax,
		Tags:             o.Tags,
		Label:            o.Label,
		Hostname:         o.Hostname,
	}, nil
}

type rawCreateLinksWithNewServerRequestServer struct {
	SystemType       string           `json:"system_type"`         // mandatory; only using "server"
	LogicalDevice    rawLogicalDevice `json:"logical_device"`      // mandatory
	PortChannelIdMin int              `json:"port_channel_id_min"` // mandatory; 0 is default
	PortChannelIdMax int              `json:"port_channel_id_max"` // mandatory; 0 is default
	Tags             []string         `json:"tags,omitempty"`
	Label            string           `json:"label,omitempty"`
	Hostname         string           `json:"hostname,omitempty"`
}

type CreateLinksWithNewServerRequestLink struct {
	Tags           []string
	SystemEndpoint SwitchLinkEndpoint
	SwitchEndpoint SwitchLinkEndpoint
	GroupLabel     string
	LagMode        RackLinkLagMode
}

func (o *CreateLinksWithNewServerRequestLink) raw() *rawCreateLinksWithNewServerRequestLink {
	return &rawCreateLinksWithNewServerRequestLink{
		Tags:           o.Tags,
		SystemEndpoint: o.SystemEndpoint.raw(),
		SwitchEndpoint: o.SwitchEndpoint.raw(),
		GroupLabel:     o.GroupLabel,
		LagMode:        rackLinkLagMode(o.LagMode.String()),
	}
}

type rawCreateLinksWithNewServerRequestLink struct {
	Tags           []string              `json:"tags,omitempty"`
	SystemEndpoint rawSwitchLinkEndpoint `json:"system"`
	SwitchEndpoint rawSwitchLinkEndpoint `json:"switch"`
	GroupLabel     string                `json:"link_group_label,omitempty"`
	LagMode        rackLinkLagMode       `json:"lag_mode,omitempty"`
}

type SwitchLinkEndpoint struct {
	TransformationId int
	SystemId         ObjectId
	IfName           string
	//LagMode          RackLinkLagMode
}

func (o *SwitchLinkEndpoint) raw() rawSwitchLinkEndpoint {
	var systemIdPtr *string
	if s := o.SystemId.String(); s != "" {
		systemIdPtr = &s
	}

	return rawSwitchLinkEndpoint{
		TransformationId: o.TransformationId,
		SystemId:         systemIdPtr,
		IfName:           o.IfName,
		//LagMode:          o.LagMode.String(),
	}
}

type rawSwitchLinkEndpoint struct {
	TransformationId int     `json:"transformation_id,omitempty"`
	SystemId         *string `json:"system_id"` // must send `null` when creating a new system, so no `omitempty`
	IfName           string  `json:"if_name,omitempty"`
	LagMode          string  `json:"lag_mode,omitempty"`
}

func (o *TwoStageL3ClosClient) CreateLinksWithNewServer(ctx context.Context, req *CreateLinksWithNewServerRequest) ([]ObjectId, error) {
	apiInput, err := req.raw(ctx, o.Client())
	if err != nil {
		return nil, fmt.Errorf("error processing CreateLinksWithNewServerRequest, - %w", err)
	}

	apiResponse := struct {
		IDs []ObjectId `json:"ids"`
	}{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlSwitchSystemLinks, o.blueprintId),
		apiInput:    apiInput,
		apiResponse: &apiResponse,
	})

	return apiResponse.IDs, convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) DeleteLinks(ctx context.Context, ids []ObjectId) error {
	apiInput := struct {
		LinkIds []ObjectId `json:"link_ids"`
	}{ids}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPost,
		urlStr:   fmt.Sprintf(apiUrlDeleteSwitchSystemLinks, o.blueprintId),
		apiInput: &apiInput,
	})

	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) DeleteGenericSystem(ctx context.Context, id ObjectId) error {
	response := struct {
		Items []struct {
			Link struct {
				ID ObjectId `json:"id"`
			} `json:"n_link"`
		} `json:"items"`
	}{}
	query := new(PathQuery).
		SetBlueprintId(o.blueprintId).
		SetBlueprintType(BlueprintTypeStaging).
		SetClient(o.client).
		Node([]QEEAttribute{NodeTypeSystem.QEEAttribute(),
			{Key: "role", Value: QEStringVal("generic")},
			{Key: "id", Value: QEStringVal(id)},
		}).
		Out([]QEEAttribute{RelationshipTypeHostedInterfaces.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeInterface.QEEAttribute(),
			{Key: "if_type", Value: QEStringVal("ethernet")},
		}).
		Out([]QEEAttribute{RelationshipTypeLink.QEEAttribute()}).
		Node([]QEEAttribute{NodeTypeLink.QEEAttribute(),
			{Key: "name", Value: QEStringVal("n_link")},
		})

	err := query.Do(ctx, &response)
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	if len(response.Items) == 0 {
		return ApstraClientErr{
			errType: ErrNotfound,
			err:     fmt.Errorf("query %q returned no results", query),
		}
	}

	linkIds := make([]ObjectId, len(response.Items))
	for i, item := range response.Items {
		linkIds[i] = item.Link.ID
	}

	return o.DeleteLinks(ctx, linkIds)
}
