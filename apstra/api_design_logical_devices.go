package apstra

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
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

// class PortGroupSchema (scotch/schemas/logical_device.py) specifies:
// validate=validate.OneOf(
//
//	['spine', 'superspine', 'leaf', 'access',
//	 'l3_server', 'peer', 'unused', 'generic']
const (
	LogicalDevicePortRoleSpine = LogicalDevicePortRoleFlags(1 << iota)
	LogicalDevicePortRoleSuperspine
	LogicalDevicePortRoleLeaf
	LogicalDevicePortRoleAccess
	LogicalDevicePortRoleL3Server
	LogicalDevicePortRolePeer
	LogicalDevicePortRoleUnused
	LogicalDevicePortRoleGeneric
	LogicalDevicePortRoleUnknown = "unknown logical device port role '%s'"

	logicalDevicePortRoleSpine      = logicalDevicePortRole("spine")
	logicalDevicePortRoleSuperspine = logicalDevicePortRole("superspine")
	logicalDevicePortRoleLeaf       = logicalDevicePortRole("leaf")
	logicalDevicePortRoleAccess     = logicalDevicePortRole("access")
	logicalDevicePortRoleL3Server   = logicalDevicePortRole("l3_server")
	logicalDevicePortRolePeer       = logicalDevicePortRole("peer")
	logicalDevicePortRoleUnused     = logicalDevicePortRole("unused")
	logicalDevicePortRoleGeneric    = logicalDevicePortRole("generic")
)

func (o *LogicalDevicePortRoleFlags) raw() []logicalDevicePortRole {
	// instantiate as zero-length rather than nil so we send "[]"
	// rather than "null" when no roles are specified.
	result := make([]logicalDevicePortRole, 0)

	if *o&LogicalDevicePortRoleSpine != 0 {
		result = append(result, logicalDevicePortRoleSpine)
	}
	if *o&LogicalDevicePortRoleSuperspine != 0 {
		result = append(result, logicalDevicePortRoleSuperspine)
	}
	if *o&LogicalDevicePortRoleLeaf != 0 {
		result = append(result, logicalDevicePortRoleLeaf)
	}
	if *o&LogicalDevicePortRoleAccess != 0 {
		result = append(result, logicalDevicePortRoleAccess)
	}
	if *o&LogicalDevicePortRoleL3Server != 0 {
		result = append(result, logicalDevicePortRoleL3Server)
	}
	if *o&LogicalDevicePortRolePeer != 0 {
		result = append(result, logicalDevicePortRolePeer)
	}
	if *o&LogicalDevicePortRoleUnused != 0 {
		result = append(result, logicalDevicePortRoleUnused)
	}
	if *o&LogicalDevicePortRoleGeneric != 0 {
		result = append(result, logicalDevicePortRoleGeneric)
	}
	return result
}

func (o *LogicalDevicePortRoleFlags) Strings() []string {
	var result []string
	for _, role := range o.raw() {
		result = append(result, string(role))
	}
	return result
}

func (o *LogicalDevicePortRoleFlags) FromStrings(in []string) error {
	*o = 0
	for _, s := range in {
		f, err := logicalDevicePortRole(s).parse()
		if err != nil {
			return err
		}
		*o = *o | f
	}
	return nil
}

func (o *LogicalDevicePortRoleFlags) SetAll() {
	*o = LogicalDevicePortRoleFlags(math.MaxUint16)
}

func (o logicalDevicePortRole) parse() (LogicalDevicePortRoleFlags, error) {
	switch o {
	case logicalDevicePortRoleSpine:
		return LogicalDevicePortRoleSpine, nil
	case logicalDevicePortRoleSuperspine:
		return LogicalDevicePortRoleSuperspine, nil
	case logicalDevicePortRoleLeaf:
		return LogicalDevicePortRoleLeaf, nil
	case logicalDevicePortRoleAccess:
		return LogicalDevicePortRoleAccess, nil
	case logicalDevicePortRoleL3Server:
		return LogicalDevicePortRoleL3Server, nil
	case logicalDevicePortRolePeer:
		return LogicalDevicePortRolePeer, nil
	case logicalDevicePortRoleUnused:
		return LogicalDevicePortRoleUnused, nil
	case logicalDevicePortRoleGeneric:
		return LogicalDevicePortRoleGeneric, nil
	default:
		return 0, fmt.Errorf(LogicalDevicePortRoleUnknown, o)
	}
}

type logicalDevicePortRoles []logicalDevicePortRole

func (o logicalDevicePortRoles) parse() (LogicalDevicePortRoleFlags, error) {
	var result LogicalDevicePortRoleFlags
	for _, r := range o {
		roleFlag, err := r.parse()
		if err != nil {
			return result, err
		}
		result = result | roleFlag
	}
	return result, nil
}

type optionsLogicalDevicesResponse struct {
	Items   []ObjectId `json:"items"`
	Methods []string   `json:"methods"`
}

type getLogicalDevicesResponse struct {
	Items []rawLogicalDevice `json:"items"`
}

type LogicalDeviceData struct {
	DisplayName string
	Panels      []LogicalDevicePanel
}

func (o *LogicalDeviceData) raw() *rawLogicalDeviceData {
	panels := make([]rawLogicalDevicePanel, len(o.Panels))
	for i := range o.Panels {
		panels[i] = *o.Panels[i].raw()
	}

	return &rawLogicalDeviceData{
		DisplayName: o.DisplayName,
		Panels:      panels,
	}
}

type rawLogicalDeviceData struct {
	DisplayName string                  `json:"display_name"`
	Panels      []rawLogicalDevicePanel `json:"panels"`
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
		Speed: *o.Speed.raw(),
		Roles: o.Roles.raw(),
	}
}

type rawLogicalDevicePortGroup struct {
	Count int                       `json:"count"`
	Speed rawLogicalDevicePortSpeed `json:"speed"`
	Roles logicalDevicePortRoles    `json:"roles"`
}

func (o *rawLogicalDevicePortGroup) parse() (*LogicalDevicePortGroup, error) {
	roles, err := o.Roles.parse()
	if err != nil {
		return nil, err
	}
	return &LogicalDevicePortGroup{
		Count: o.Count,
		Speed: o.Speed.parse(),
		Roles: roles,
	}, nil
}

type LogicalDevicePortSpeed string

func (o LogicalDevicePortSpeed) raw() *rawLogicalDevicePortSpeed {
	if o == "" {
		return nil
	}
	defaultSpeed := rawLogicalDevicePortSpeed{
		Unit:  "G",
		Value: 1,
	}
	lower := strings.ToLower(string(o))
	lower = strings.TrimSpace(lower)
	lower = strings.TrimSuffix(lower, "bps")
	lower = strings.TrimSuffix(lower, "b/s")
	var factor int64
	var trimmed string
	switch {
	case strings.HasSuffix(lower, "m"):
		trimmed = strings.TrimSuffix(lower, "m")
		factor = 1000 * 1000
	case strings.HasSuffix(lower, "g"):
		trimmed = strings.TrimSuffix(lower, "g")
		factor = 1000 * 1000 * 1000
	case strings.HasSuffix(lower, "t"):
		trimmed = strings.TrimSuffix(lower, "t")
		factor = 1000 * 1000 * 1000 * 1000
	default:
		trimmed = lower
		factor = 1
	}
	trimmedInt, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return &defaultSpeed
	}
	bps := trimmedInt * factor
	switch {
	case bps >= 1000*1000*1000: // at least 1Gbps
		return &rawLogicalDevicePortSpeed{
			Unit:  "G",
			Value: int(bps / 1000 / 1000 / 1000),
		}
	case bps >= 10*1000*1000: // at least 10Mbps
		return &rawLogicalDevicePortSpeed{
			Unit:  "M",
			Value: int(bps / 1000 / 1000),
		}
	default:
		return &defaultSpeed
	}
}

func (o LogicalDevicePortSpeed) BitsPerSecond() int64 {
	return o.raw().BitsPerSecond()
}

func (o LogicalDevicePortSpeed) IsEqual(in LogicalDevicePortSpeed) bool {
	return o.BitsPerSecond() == in.BitsPerSecond()
}

type rawLogicalDevicePortSpeed struct {
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

func (o rawLogicalDevicePortSpeed) parse() LogicalDevicePortSpeed {
	return LogicalDevicePortSpeed(fmt.Sprintf("%d%s", o.Value, o.Unit))
}

func (o *rawLogicalDevicePortSpeed) BitsPerSecond() int64 {
	switch o.Unit {
	case "M":
		return int64(o.Value * 1000 * 1000)
	case "G":
		return int64(o.Value * 1000 * 1000 * 1000)
	default:
		return int64(0)
	}
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
		PanelLayout:  o.PanelLayout,
		PortIndexing: o.PortIndexing,
		PortGroups:   portGroups,
	}, nil
}

type LogicalDevice struct {
	Id             ObjectId
	CreatedAt      *time.Time
	LastModifiedAt *time.Time
	Data           *LogicalDeviceData
}

func (o LogicalDevice) raw() *rawLogicalDevice {
	var panels []rawLogicalDevicePanel
	for _, panel := range o.Data.Panels {
		panels = append(panels, *panel.raw())
	}

	return &rawLogicalDevice{
		DisplayName:    o.Data.DisplayName,
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
	CreatedAt      *time.Time              `json:"created_at,omitempty"`
	LastModifiedAt *time.Time              `json:"last_modified_at,omitempty"`
}

func (o rawLogicalDevice) polish() (*LogicalDevice, error) {
	panels := make([]LogicalDevicePanel, len(o.Panels))
	for i, panel := range o.Panels {
		parsed, err := panel.parse()
		if err != nil {
			return nil, err
		}
		panels[i] = *parsed
	}

	return &LogicalDevice{
		Id:             o.Id,
		CreatedAt:      o.CreatedAt,
		LastModifiedAt: o.LastModifiedAt,
		Data: &LogicalDeviceData{
			DisplayName: o.DisplayName,
			Panels:      panels,
		},
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

func (o *Client) getAllLogicalDevices(ctx context.Context) ([]LogicalDevice, error) {
	response := &getLogicalDevicesResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      apiUrlDesignLogicalDevices,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	var result []LogicalDevice
	for _, raw := range response.Items {
		ld, err := raw.polish()
		if err != nil {
			return nil, err
		}
		result = append(result, *ld)
	}
	return result, nil
}

func (o *Client) getLogicalDevice(ctx context.Context, id ObjectId) (*rawLogicalDevice, error) {
	response := &rawLogicalDevice{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *Client) getLogicalDeviceByName(ctx context.Context, name string) (*LogicalDevice, error) {
	logicalDevices, err := o.getAllLogicalDevices(ctx)
	if err != nil {
		return nil, err
	}

	var result LogicalDevice
	var found bool

	for _, ld := range logicalDevices {
		foo := &ld
		_ = foo
		if ld.Data.DisplayName == name {
			if found {
				return nil, ApstraClientErr{
					errType: ErrMultipleMatch,
					err:     fmt.Errorf("found multiple logical devices named '%s' found", name),
				}
			}
			result = ld
			found = true
		}
	}
	if found {
		return &result, nil
	}
	return nil, ApstraClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("no logical device named '%s' found", name),
	}
}

func (o *Client) createLogicalDevice(ctx context.Context, in *rawLogicalDeviceData) (ObjectId, error) {
	response := &objectIdResponse{}
	err := o.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      apiUrlDesignLogicalDevices,
		apiInput:    in,
		apiResponse: response,
	})
	if err != nil {
		return "", convertTtaeToAceWherePossible(err)
	}
	return response.Id, nil
}

func (o *Client) updateLogicalDevice(ctx context.Context, id ObjectId, in *rawLogicalDeviceData) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
		apiInput: in,
	})
}

func (o *Client) deleteLogicalDevice(ctx context.Context, id ObjectId) error {
	return o.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlDesignLogicalDeviceById, id),
	})
}
