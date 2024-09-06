package apstra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Juniper/apstra-go-sdk/apstra/enum"
)

const (
	apiUrlBlueprintSecurityZones               = apiUrlBlueprintById + apiUrlPathDelim + "security-zones"
	apiUrlBlueprintSecurityZonesPrefix         = apiUrlBlueprintSecurityZones + apiUrlPathDelim
	apiUrlBlueprintSecurityZoneById            = apiUrlBlueprintSecurityZonesPrefix + "%s"
	apiUrlBlueprintSecurityZoneByIdDhcpServers = apiUrlBlueprintSecurityZoneById + apiUrlPathDelim + "dhcp-servers"
)

type (
	SecurityZoneType int
	securityZoneType string
)

const (
	SecurityZoneTypeNone = SecurityZoneType(iota)
	SecurityZoneTypeL3Fabric
	SecurityZoneTypeL3FabricVirtual
	SecurityZoneTypeEVPN
	SecurityZoneTypeUnknown = "unknown security zone type %s"

	securityZoneTypeNone            = securityZoneType("")
	securityZoneTypeL3Fabric        = securityZoneType("l3_fabric")
	securityZoneTypeL3FabricVirtual = securityZoneType("virtual_l3_fabric")
	securityZoneTypeEVPN            = securityZoneType("evpn")
	securityZoneTypeUnknown         = "unknown security zone type %d"
)

func (o SecurityZoneType) Int() int {
	return int(o)
}

func (o SecurityZoneType) String() string {
	switch o {
	case SecurityZoneTypeNone:
		return string(securityZoneTypeNone)
	case SecurityZoneTypeL3Fabric:
		return string(securityZoneTypeL3Fabric)
	case SecurityZoneTypeL3FabricVirtual:
		return string(securityZoneTypeL3FabricVirtual)
	case SecurityZoneTypeEVPN:
		return string(securityZoneTypeEVPN)
	default:
		return fmt.Sprintf(securityZoneTypeUnknown, o)
	}
}

func (o *SecurityZoneType) FromString(in string) error {
	i, err := securityZoneType(in).parse()
	if err != nil {
		return err
	}
	*o = SecurityZoneType(i)
	return nil
}

func (o SecurityZoneType) raw() securityZoneType {
	return securityZoneType(o.String())
}

func (o securityZoneType) string() string {
	return string(o)
}

func (o securityZoneType) parse() (int, error) {
	switch o {
	case securityZoneTypeNone:
		return int(SecurityZoneTypeNone), nil
	case securityZoneTypeL3Fabric:
		return int(SecurityZoneTypeL3Fabric), nil
	case securityZoneTypeL3FabricVirtual:
		return int(SecurityZoneTypeL3FabricVirtual), nil
	case securityZoneTypeEVPN:
		return int(SecurityZoneTypeEVPN), nil
	default:
		return 0, fmt.Errorf(SecurityZoneTypeUnknown, o)
	}
}

type SecurityZone struct {
	Id   ObjectId
	Data *SecurityZoneData
}

type SecurityZoneData struct {
	Label   string
	SzType  SecurityZoneType
	VrfName string

	RoutingPolicyId  ObjectId               // automatically assigned
	RouteTarget      *string                // can be null
	RtPolicy         *RtPolicy              // can be null
	VlanId           *Vlan                  // can be null
	VniId            *int                   // can be null
	JunosEvpnIrbMode *enum.JunosEvpnIrbMode // Apstra 4.2+ only
}

func (o SecurityZoneData) raw() *rawSecurityZone {
	var routeTarget string
	if o.RouteTarget != nil {
		routeTarget = *o.RouteTarget
	}

	var junosEvpnIrbMode string
	if o.JunosEvpnIrbMode != nil {
		junosEvpnIrbMode = o.JunosEvpnIrbMode.Value
	}

	return &rawSecurityZone{
		Label:            o.Label,
		SzType:           o.SzType.raw(),
		VrfName:          o.VrfName,
		RoutingPolicyId:  o.RoutingPolicyId,
		RouteTarget:      routeTarget,
		RtPolicy:         o.RtPolicy,
		VlanId:           o.VlanId,
		VniId:            o.VniId,
		JunosEvpnIrbMode: junosEvpnIrbMode,
	}
}

type rawSecurityZone struct {
	Id               ObjectId         `json:"id,omitempty"`
	Label            string           `json:"label"`
	SzType           securityZoneType `json:"sz_type"`
	VrfName          string           `json:"vrf_name"`
	RoutingPolicyId  ObjectId         `json:"routing_policy_id,omitempty"`
	RouteTarget      string           `json:"route_target,omitempty"`
	RtPolicy         *RtPolicy        `json:"rt_policy,omitempty"`
	VlanId           *Vlan            `json:"vlan_id,omitempty"`
	VniId            *int             `json:"vni_id,omitempty"`
	JunosEvpnIrbMode string           `json:"junos_evpn_irb_mode,omitempty"`
}

func (o rawSecurityZone) polish() (*SecurityZone, error) {
	szType, err := o.SzType.parse()
	if err != nil {
		return nil, fmt.Errorf("error parsing security zone type %q - %w", o.SzType, err)
	}

	var routeTarget *string
	if o.RouteTarget != "" {
		rt := o.RouteTarget
		routeTarget = &rt
	}

	return &SecurityZone{
		Id: o.Id,
		Data: &SecurityZoneData{
			Label:            o.Label,
			SzType:           SecurityZoneType(szType),
			VrfName:          o.VrfName,
			RoutingPolicyId:  o.RoutingPolicyId,
			RouteTarget:      routeTarget,
			RtPolicy:         o.RtPolicy,
			VlanId:           o.VlanId,
			VniId:            o.VniId,
			JunosEvpnIrbMode: enum.JunosEvpnIrbModes.Parse(o.JunosEvpnIrbMode),
		},
	}, nil
}

func (o *TwoStageL3ClosClient) createSecurityZone(ctx context.Context, cfg *rawSecurityZone) (*objectIdResponse, error) {
	response := &objectIdResponse{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodPost,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZones, o.blueprintId),
		apiInput:    cfg,
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *TwoStageL3ClosClient) getSecurityZone(ctx context.Context, zoneId ObjectId) (*rawSecurityZone, error) {
	response := &rawSecurityZone{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, zoneId),
		apiResponse: response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}
	return response, nil
}

func (o *TwoStageL3ClosClient) getSecurityZoneByVrfName(ctx context.Context, vrfName string) (*rawSecurityZone, error) {
	zones, err := o.getAllSecurityZones(ctx)
	if err != nil {
		return nil, err
	}

	for _, zone := range zones {
		if zone.VrfName == vrfName {
			return &zone, nil
		}
	}

	return nil, ClientErr{
		errType: ErrNotfound,
		err:     fmt.Errorf("security zone with vrf name %q not found in blueprint %q", vrfName, o.blueprintId),
	}
}

func (o *TwoStageL3ClosClient) getAllSecurityZones(ctx context.Context) (map[string]rawSecurityZone, error) {
	response := &struct {
		Items map[string]rawSecurityZone `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZones, o.blueprintId),
		apiResponse: response,
	})
	if err != nil {
		return nil, err
	}

	return response.Items, nil
}

func (o *TwoStageL3ClosClient) updateSecurityZone(ctx context.Context, zoneId ObjectId, cfg *rawSecurityZone) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPut,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, zoneId),
		apiInput: cfg,
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) deleteSecurityZone(ctx context.Context, zoneId ObjectId) error {
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodDelete,
		urlStr: fmt.Sprintf(apiUrlBlueprintSecurityZoneById, o.blueprintId, zoneId),
	})
	return convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) getSecurityZoneDhcpServers(ctx context.Context, zoneId ObjectId) ([]string, error) {
	response := &struct {
		Items []string `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZoneByIdDhcpServers, o.blueprintId, zoneId),
		apiResponse: response,
	})
	return response.Items, convertTtaeToAceWherePossible(err)
}

func (o *TwoStageL3ClosClient) setSecurityZoneDhcpServers(ctx context.Context, zoneId ObjectId, ips []string) error {
	return convertTtaeToAceWherePossible(o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPut,
		urlStr: fmt.Sprintf(apiUrlBlueprintSecurityZoneByIdDhcpServers, o.blueprintId, zoneId),
		apiInput: &struct {
			Items []string `json:"items"`
		}{
			Items: ips,
		},
	}))
}
