package goapstra

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	apiUrlDesignLogicalDevices       = apiUrlDesignPrefix + "logical-devices"
	apiUrlDesignLogicalDevicesPrefix = apiUrlDesignLogicalDevices + apiUrlPathDelim
	apiUrlDesignLogicalDeviceById    = apiUrlDesignLogicalDevicesPrefix + "%s"

	PortIndexingVerticalFirst   = "L-R, T-B"
	PortIndexingHorizontalFirst = "T-B, L-R"
	PortIndexingSchemaAbsolute  = "absolute"
)

type LogicalDevicePortRoleFlags uint16
type logicalDevicePortRole string
type logicalDevicePortRoles []string

const (
	LogicalDevicePortRoleAccess = LogicalDevicePortRoleFlags(1 << iota)
	LogicalDevicePortRoleGeneric
	LogicalDevicePortRoleL3Server
	LogicalDevicePortRoleLeaf
	LogicalDevicePortRolePeer
	LogicalDevicePortRoleServer
	LogicalDevicePortRoleSpine
	LogicalDevicePortRoleSuperspine
	LogicalDevicePortRoleUnused
	LogicalDevicePortRoleUnknown = "unknown logical device port role '%s'"

	logicalDevicePortRoleAccess     = logicalDevicePortRole("access")
	logicalDevicePortRoleGeneric    = logicalDevicePortRole("generic")
	logicalDevicePortRoleL3Server   = logicalDevicePortRole("l3_server")
	logicalDevicePortRoleLeaf       = logicalDevicePortRole("leaf")
	logicalDevicePortRolePeer       = logicalDevicePortRole("peer")
	logicalDevicePortRoleServer     = logicalDevicePortRole("server")
	logicalDevicePortRoleSpine      = logicalDevicePortRole("spine")
	logicalDevicePortRoleSuperspine = logicalDevicePortRole("superspine")
	logicalDevicePortRoleUnused     = logicalDevicePortRole("unused")
)

func (o LogicalDevicePortRoleFlags) raw() []logicalDevicePortRole {
	var result []logicalDevicePortRole
	if o&LogicalDevicePortRoleAccess != 0 {
		result = append(result, logicalDevicePortRoleAccess)
	}
	if o&LogicalDevicePortRoleGeneric != 0 {
		result = append(result, logicalDevicePortRoleGeneric)
	}
	if o&LogicalDevicePortRoleL3Server != 0 {
		result = append(result, logicalDevicePortRoleL3Server)
	}
	if o&LogicalDevicePortRoleLeaf != 0 {
		result = append(result, logicalDevicePortRoleLeaf)
	}
	if o&LogicalDevicePortRolePeer != 0 {
		result = append(result, logicalDevicePortRolePeer)
	}
	if o&LogicalDevicePortRoleServer != 0 {
		result = append(result, logicalDevicePortRoleServer)
	}
	if o&LogicalDevicePortRoleSpine != 0 {
		result = append(result, logicalDevicePortRoleSpine)
	}
	if o&LogicalDevicePortRoleSuperspine != 0 {
		result = append(result, logicalDevicePortRoleSuperspine)
	}
	if o&LogicalDevicePortRoleUnused != 0 {
		result = append(result, logicalDevicePortRoleUnused)
	}
	return result
}

func (o logicalDevicePortRole) parse() (LogicalDevicePortRoleFlags, error) {
	switch o {
	case logicalDevicePortRoleAccess:
		return LogicalDevicePortRoleAccess, nil
	case logicalDevicePortRoleGeneric:
		return LogicalDevicePortRoleGeneric, nil
	case logicalDevicePortRoleL3Server:
		return LogicalDevicePortRoleL3Server, nil
	case logicalDevicePortRoleLeaf:
		return LogicalDevicePortRoleLeaf, nil
	case logicalDevicePortRolePeer:
		return LogicalDevicePortRolePeer, nil
	case logicalDevicePortRoleServer:
		return LogicalDevicePortRoleServer, nil
	case logicalDevicePortRoleSpine:
		return LogicalDevicePortRoleSpine, nil
	case logicalDevicePortRoleSuperspine:
		return LogicalDevicePortRoleSuperspine, nil
	case logicalDevicePortRoleUnused:
		return LogicalDevicePortRoleUnused, nil
	default:
		return 0, fmt.Errorf(LogicalDevicePortRoleUnknown, o)
	}
}

type optionsLogicalDevicesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type LogicalDevicePanelLayout struct {
	RowCount    int `json:"row_count"`
	ColumnCount int `json:"column_count"`
}

type LogicalDevicePortIndexing struct {
	Order      string `json:"order"`
	StartIndex int    `json:"start_index"`
	Schema     string `json:"schema"` // Valid choices: absolute
}

type LogicalDevicePortGroup struct {
	Count int                        `json:"count"`
	Speed LogicalDevicePortSpeed     `json:"speed"`
	Roles LogicalDevicePortRoleFlags `json:"roles"`
}

func (o LogicalDevicePortGroup) raw() *rawLogicalDevicePortGroup {
	return &rawLogicalDevicePortGroup{
		Count: o.Count,
		Speed: o.Speed,
		Roles: o.Roles.raw(),
	}
}

type rawLogicalDevicePortGroup struct {
	Count int                     `json:"count"`
	Speed LogicalDevicePortSpeed  `json:"speed"`
	Roles []logicalDevicePortRole `json:"roles"`
}

func (o *rawLogicalDevicePortGroup) parse() (*LogicalDevicePortGroup, error) {
	var roles LogicalDevicePortRoleFlags
	for _, role := range o.Roles {
		parsed, err := role.parse()
		if err != nil {
			return nil, err
		}
		roles = roles & parsed
	}
	return &LogicalDevicePortGroup{
		Count: o.Count,
		Speed: o.Speed,
		Roles: roles,
	}, nil
}

type LogicalDevicePortSpeed struct {
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

type LogicalDevicePanel struct {
	PanelLayout  LogicalDevicePanelLayout  `json:"panel_layout"`
	PortIndexing LogicalDevicePortIndexing `json:"port_indexing"`
	PortGroups   []LogicalDevicePortGroup  `json:"port_groups"`
}

func (o LogicalDevicePanel) raw() *rawLogicalDevicePanel {
	var portGroups []rawLogicalDevicePortGroup
	for _, pg := range o.PortGroups {
		portGroups = append(portGroups, *pg.raw())
	}

	return &rawLogicalDevicePanel{
		PanelLayout:  o.PanelLayout,
		PortIndexing: o.PortIndexing,
		PortGroups:   portGroups,
	}
}

type rawLogicalDevicePanel struct {
	PanelLayout  LogicalDevicePanelLayout    `json:"panel_layout"`
	PortIndexing LogicalDevicePortIndexing   `json:"port_indexing"`
	PortGroups   []rawLogicalDevicePortGroup `json:"port_groups"`
}

func (o rawLogicalDevicePanel) parse() (*LogicalDevicePanel, error) {
	var portGroups []LogicalDevicePortGroup
	for _, pg := range o.PortGroups {
		p, err := pg.parse()
		if err != nil {
			return nil, err
		}
		portGroups = append(portGroups, *p)
	}

	return &LogicalDevicePanel{
		PanelLayout:  LogicalDevicePanelLayout{},
		PortIndexing: LogicalDevicePortIndexing{},
		PortGroups:   portGroups,
	}, nil
}

type LogicalDevice struct {
	DisplayName    string
	Id             ObjectId
	Panels         []LogicalDevicePanel
	CreatedAt      time.Time
	LastModifiedAt time.Time
}

func (o LogicalDevice) raw() *rawLogicalDevice {
	var panels []rawLogicalDevicePanel
	for _, panel := range o.Panels {
		panels = append(panels, *panel.raw())
	}

	return &rawLogicalDevice{
		DisplayName:    o.DisplayName,
		Id:             o.Id,
		Panels:         panels,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
	}
}

type rawLogicalDevice struct {
	DisplayName    string                  `json:"display_name"`
	Id             ObjectId                `json:"id,omitempty"`
	Panels         []rawLogicalDevicePanel `json:"panels"`
	CreatedAt      time.Time               `json:"created_at"`
	LastModifiedAt time.Time               `json:"last_modified_at"`
}

func (o rawLogicalDevice) parse() (*LogicalDevice, error) {
	var panels []LogicalDevicePanel
	for _, panel := range o.Panels {
		parsed, err := panel.parse()
		if err != nil {
			return nil, err
		}
		panels = append(panels, *parsed)
	}

	return &LogicalDevice{
		DisplayName:    o.DisplayName,
		Id:             o.Id,
		Panels:         panels,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
	}, nil
}

func (o *Client) listLogicalDeviceIds(ctx context.Context) ([]ObjectId, error) {
	response := &optionsLogicalDevicesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodOptions,
		urlStr:      apiUrlDesignLogicalDevices,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.Items, nil
}

func (o *Client) getLogicalDevice(ctx context.Context, id ObjectId) (*LogicalDevice, error) {
	response := &rawLogicalDevice{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response.parse()
}

func (o *Client) createLogicalDevice(ctx context.Context, in *LogicalDevice) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignLogicalDevices,
		apiInput:    in.raw(),
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateLogicalDevice(ctx context.Context, id ObjectId, in *LogicalDevice) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
		apiInput: in.raw(),
	})
}

func (o *Client) deleteLogicalDevice(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
	})
}
