package apstra

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiUrlBlueprintExperienceWeb = apiUrlBlueprintById + apiUrlPathDelim + "experience/web/"
	apiUrlBlueprintCablingMap    = apiUrlBlueprintExperienceWeb + "cabling-map"

	includeLagParam    = "aggregate_links"
	linksBySystemParam = "system_node_id"
)

type InterfaceType int
type interfaceType string

const (
	InterfaceTypeEthernet = InterfaceType(iota)
	InterfaceTypeIp
	InterfaceTypeLoopback
	InterfaceTypePortChannel
	InterfaceTypeSvi
	InterfaceTypeLogicalVtep
	InterfaceTypeAnycastVtep
	InterfaceTypeUnicastVtep
	InterfaceTypeGlobalAnycastVtep
	InterfaceTypeSubinterface
	InterfaceTypeUnknown = "unknown interface type '%s'"

	interfaceTypeEthernet          = interfaceType("ethernet")
	interfaceTypeIp                = interfaceType("ip")
	interfaceTypeLoopback          = interfaceType("loopback")
	interfaceTypePortChannel       = interfaceType("port_channel")
	interfaceTypeSvi               = interfaceType("svi")
	interfaceTypeLogicalVtep       = interfaceType("logical_vtep")
	interfaceTypeAnycastVtep       = interfaceType("anycast_vtep")
	interfaceTypeUnicastVtep       = interfaceType("unicast_vtep")
	interfaceTypeGlobalAnycastVtep = interfaceType("global_anycast_vtep")
	interfaceTypeSubinterface      = interfaceType("subinterface")
	interfaceTypeUnknown           = "unknown interface type '%d'"
)

func (o InterfaceType) Int() int {
	return int(o)
}

func (o InterfaceType) String() string {
	switch o {
	case InterfaceTypeEthernet:
		return string(interfaceTypeEthernet)
	case InterfaceTypeIp:
		return string(interfaceTypeIp)
	case InterfaceTypeLoopback:
		return string(interfaceTypeLoopback)
	case InterfaceTypePortChannel:
		return string(interfaceTypePortChannel)
	case InterfaceTypeSvi:
		return string(interfaceTypeSvi)
	case InterfaceTypeLogicalVtep:
		return string(interfaceTypeLogicalVtep)
	case InterfaceTypeAnycastVtep:
		return string(interfaceTypeAnycastVtep)
	case InterfaceTypeUnicastVtep:
		return string(interfaceTypeUnicastVtep)
	case InterfaceTypeGlobalAnycastVtep:
		return string(interfaceTypeGlobalAnycastVtep)
	case InterfaceTypeSubinterface:
		return string(interfaceTypeSubinterface)
	default:
		return fmt.Sprintf(interfaceTypeUnknown, o)
	}
}

func (o InterfaceType) raw() interfaceType {
	return interfaceType(o.String())
}

func (o interfaceType) string() string {
	return string(o)
}

func (o interfaceType) parse() (int, error) {
	switch o {
	case interfaceTypeEthernet:
		return int(InterfaceTypeEthernet), nil
	case interfaceTypeIp:
		return int(InterfaceTypeIp), nil
	case interfaceTypeLoopback:
		return int(InterfaceTypeLoopback), nil
	case interfaceTypePortChannel:
		return int(InterfaceTypePortChannel), nil
	case interfaceTypeSvi:
		return int(InterfaceTypeSvi), nil
	case interfaceTypeLogicalVtep:
		return int(InterfaceTypeLogicalVtep), nil
	case interfaceTypeAnycastVtep:
		return int(InterfaceTypeAnycastVtep), nil
	case interfaceTypeUnicastVtep:
		return int(InterfaceTypeUnicastVtep), nil
	case interfaceTypeGlobalAnycastVtep:
		return int(InterfaceTypeGlobalAnycastVtep), nil
	case interfaceTypeSubinterface:
		return int(InterfaceTypeSubinterface), nil
	default:
		return 0, fmt.Errorf(InterfaceTypeUnknown, o)
	}
}

type InterfaceOperationState int
type interfaceOperationState string

const (
	InterfaceOperationStateAdminDown = InterfaceOperationState(iota)
	InterfaceOperationStateDown
	InterfaceOperationStateUp
	InterfaceOperationStateUnknown = "unknown interface operation state '%s'"

	interfaceOperationStateAdminDown = interfaceOperationState("admin_down")
	interfaceOperationStateDown      = interfaceOperationState("deduced_down")
	interfaceOperationStateUp        = interfaceOperationState("up")
	interfaceOperationStateUnknown   = "unknown interface operation state '%d'"
)

func (o InterfaceOperationState) Int() int {
	return int(o)
}

func (o InterfaceOperationState) String() string {
	switch o {
	case InterfaceOperationStateAdminDown:
		return string(interfaceOperationStateAdminDown)
	case InterfaceOperationStateDown:
		return string(interfaceOperationStateDown)
	case InterfaceOperationStateUp:
		return string(interfaceOperationStateUp)
	default:
		return fmt.Sprintf(interfaceOperationStateUnknown, o)
	}
}

func (o InterfaceOperationState) raw() interfaceOperationState {
	return interfaceOperationState(o.String())
}

func (o interfaceOperationState) string() string {
	return string(o)
}

func (o interfaceOperationState) parse() (int, error) {
	switch o {
	case interfaceOperationStateAdminDown:
		return int(InterfaceOperationStateAdminDown), nil
	case interfaceOperationStateDown:
		return int(InterfaceOperationStateDown), nil
	case interfaceOperationStateUp:
		return int(InterfaceOperationStateUp), nil
	default:
		return 0, fmt.Errorf(InterfaceOperationStateUnknown, o)
	}
}

type LinkRole int
type linkRole string

const (
	LinkRoleAccessL3PeerLink = LinkRole(iota)
	LinkRoleAccessServer
	LinkRoleLeafAccess
	LinkRoleLeafL2Server
	LinkRoleLeafL3PeerLink
	LinkRoleLeafL3Server
	LinkRoleLeafLeaf
	LinkRoleLeafPairAccess
	LinkRoleLeafPairAccessPair
	LinkRoleLeafPairL2Server
	LinkRoleLeafPeerLink
	LinkRoleSpineLeaf
	LinkRoleSpineSuperspine
	LinkRoleToExternalRouter
	LinkRoleToGeneric
	LinkRoleUnknown = "unknown link role '%s'"

	linkRoleAccessL3PeerLink   = linkRole("access_l3_peer_link")
	linkRoleAccessServer       = linkRole("access_server")
	linkRoleLeafAccess         = linkRole("leaf_access")
	linkRoleLeafL2Server       = linkRole("leaf_l2_server")
	linkRoleLeafL3PeerLink     = linkRole("leaf_l3_peer_link")
	linkRoleLeafL3Server       = linkRole("leaf_l3_server")
	linkRoleLeafLeaf           = linkRole("leaf_leaf")
	linkRoleLeafPairAccess     = linkRole("leaf_pair_access")
	linkRoleLeafPairAccessPair = linkRole("leaf_pair_access_pair")
	linkRoleLeafPairL2Server   = linkRole("leaf_pair_l2_server")
	linkRoleLeafPeerLink       = linkRole("leaf_peer_link")
	linkRoleSpineLeaf          = linkRole("spine_leaf")
	linkRoleSpineSuperspine    = linkRole("spine_superspine")
	linkRoleToExternalRouter   = linkRole("to_external_router")
	linkRoleToGeneric          = linkRole("to_generic")
	linkRoleUnknown            = "unknown link role '%d'"
)

func (o LinkRole) Int() int {
	return int(o)
}

func (o LinkRole) String() string {
	switch o {
	case LinkRoleAccessL3PeerLink:
		return string(linkRoleAccessL3PeerLink)
	case LinkRoleAccessServer:
		return string(linkRoleAccessServer)
	case LinkRoleLeafAccess:
		return string(linkRoleLeafAccess)
	case LinkRoleLeafL2Server:
		return string(linkRoleLeafL2Server)
	case LinkRoleLeafL3PeerLink:
		return string(linkRoleLeafL3PeerLink)
	case LinkRoleLeafL3Server:
		return string(linkRoleLeafL3Server)
	case LinkRoleLeafLeaf:
		return string(linkRoleLeafLeaf)
	case LinkRoleLeafPairAccess:
		return string(linkRoleLeafPairAccess)
	case LinkRoleLeafPairAccessPair:
		return string(linkRoleLeafPairAccessPair)
	case LinkRoleLeafPairL2Server:
		return string(linkRoleLeafPairL2Server)
	case LinkRoleLeafPeerLink:
		return string(linkRoleLeafPeerLink)
	case LinkRoleSpineLeaf:
		return string(linkRoleSpineLeaf)
	case LinkRoleSpineSuperspine:
		return string(linkRoleSpineSuperspine)
	case LinkRoleToExternalRouter:
		return string(linkRoleToExternalRouter)
	case LinkRoleToGeneric:
		return string(linkRoleToGeneric)
	default:
		return fmt.Sprintf(linkRoleUnknown, o)
	}
}

func (o *LinkRole) FromString(s string) error {
	i, err := linkRole(s).parse()
	if err != nil {
		return err
	}
	*o = LinkRole(i)
	return nil
}

func (o LinkRole) raw() linkRole {
	return linkRole(o.String())
}

func (o linkRole) string() string {
	return string(o)
}

func (o linkRole) parse() (int, error) {
	switch o {
	case linkRoleAccessL3PeerLink:
		return int(LinkRoleAccessL3PeerLink), nil
	case linkRoleAccessServer:
		return int(LinkRoleAccessServer), nil
	case linkRoleLeafAccess:
		return int(LinkRoleLeafAccess), nil
	case linkRoleLeafL2Server:
		return int(LinkRoleLeafL2Server), nil
	case linkRoleLeafL3PeerLink:
		return int(LinkRoleLeafL3PeerLink), nil
	case linkRoleLeafL3Server:
		return int(LinkRoleLeafL3Server), nil
	case linkRoleLeafLeaf:
		return int(LinkRoleLeafLeaf), nil
	case linkRoleLeafPairAccess:
		return int(LinkRoleLeafPairAccess), nil
	case linkRoleLeafPairAccessPair:
		return int(LinkRoleLeafPairAccessPair), nil
	case linkRoleLeafPairL2Server:
		return int(LinkRoleLeafPairL2Server), nil
	case linkRoleLeafPeerLink:
		return int(LinkRoleLeafPeerLink), nil
	case linkRoleSpineLeaf:
		return int(LinkRoleSpineLeaf), nil
	case linkRoleSpineSuperspine:
		return int(LinkRoleSpineSuperspine), nil
	case linkRoleToExternalRouter:
		return int(LinkRoleToExternalRouter), nil
	case linkRoleToGeneric:
		return int(LinkRoleToGeneric), nil
	default:
		return 0, fmt.Errorf(LinkRoleUnknown, o)
	}
}

type SystemNodeRole int
type systemNodeRole string

const (
	SystemNodeRoleNone = SystemNodeRole(iota)
	SystemNodeRoleAccess
	SystemNodeRoleGeneric
	SystemNodeRoleL3Server
	SystemNodeRoleLeaf
	SystemNodeRoleRemoteGateway
	SystemNodeRoleSpine
	SystemNodeRoleSuperspine
	SystemNodeRoleUnknown = "unknown system node role '%s'"

	systemNodeRoleNone          = systemNodeRole("")
	systemNodeRoleAccess        = systemNodeRole("access")
	systemNodeRoleGeneric       = systemNodeRole("generic")
	systemNodeRoleL3Server      = systemNodeRole("l3_server")
	systemNodeRoleLeaf          = systemNodeRole("leaf")
	systemNodeRoleRemoteGateway = systemNodeRole("remote_gateway")
	systemNodeRoleSpine         = systemNodeRole("spine")
	systemNodeRoleSuperspine    = systemNodeRole("superspine")
	systemNodeRoleUnknown       = "unknown system node role '%d'"
)

func (o SystemNodeRole) Int() int {
	return int(o)
}

func (o SystemNodeRole) String() string {
	switch o {
	case SystemNodeRoleNone:
		return string(systemNodeRoleNone)
	case SystemNodeRoleAccess:
		return string(systemNodeRoleAccess)
	case SystemNodeRoleGeneric:
		return string(systemNodeRoleGeneric)
	case SystemNodeRoleL3Server:
		return string(systemNodeRoleL3Server)
	case SystemNodeRoleLeaf:
		return string(systemNodeRoleLeaf)
	case SystemNodeRoleRemoteGateway:
		return string(systemNodeRoleRemoteGateway)
	case SystemNodeRoleSpine:
		return string(systemNodeRoleSpine)
	case SystemNodeRoleSuperspine:
		return string(systemNodeRoleSuperspine)
	default:
		return fmt.Sprintf(systemNodeRoleUnknown, o)
	}
}

func (o SystemNodeRole) raw() systemNodeRole {
	return systemNodeRole(o.String())
}

func (o SystemNodeRole) QEAttribute() QEEAttribute {
	return QEEAttribute{"type", QEStringVal(o.String())}
}

func (o systemNodeRole) string() string {
	return string(o)
}

func (o systemNodeRole) parse() (int, error) {
	switch o {
	case systemNodeRoleNone:
		return int(SystemNodeRoleNone), nil
	case systemNodeRoleAccess:
		return int(SystemNodeRoleAccess), nil
	case systemNodeRoleGeneric:
		return int(SystemNodeRoleGeneric), nil
	case systemNodeRoleL3Server:
		return int(SystemNodeRoleL3Server), nil
	case systemNodeRoleLeaf:
		return int(SystemNodeRoleLeaf), nil
	case systemNodeRoleRemoteGateway:
		return int(SystemNodeRoleRemoteGateway), nil
	case systemNodeRoleSpine:
		return int(SystemNodeRoleSpine), nil
	case systemNodeRoleSuperspine:
		return int(SystemNodeRoleSuperspine), nil
	default:
		return 0, fmt.Errorf(SystemNodeRoleUnknown, o)
	}
}

type LinkType int
type linkType string

const (
	LinkTypeAggregateLink = LinkType(iota)
	LinkTypeEthernet
	LinkTypeLogicalLink
	LinkTypeUnknown = "unknown link type '%s'"

	linkTypeAggregateLink = linkType("aggregate_link")
	linkTypeEthernet      = linkType("ethernet")
	linkTypeLogicalLink   = linkType("logical_link")
	linkTypeUnknown       = "unknown link type '%d'"
)

func (o LinkType) Int() int {
	return int(o)
}

func (o LinkType) String() string {
	switch o {
	case LinkTypeAggregateLink:
		return string(linkTypeAggregateLink)
	case LinkTypeEthernet:
		return string(linkTypeEthernet)
	case LinkTypeLogicalLink:
		return string(linkTypeLogicalLink)
	default:
		return fmt.Sprintf(linkTypeUnknown, o)
	}
}

func (o LinkType) QEAttributee() QEEAttribute {
	return QEEAttribute{"type", QEStringVal(o.String())}
}

func (o LinkType) raw() linkType {
	return linkType(o.String())
}

func (o linkType) string() string {
	return string(o)
}

func (o linkType) parse() (int, error) {
	switch o {
	case linkTypeAggregateLink:
		return int(LinkTypeAggregateLink), nil
	case linkTypeEthernet:
		return int(LinkTypeEthernet), nil
	case linkTypeLogicalLink:
		return int(LinkTypeLogicalLink), nil
	default:
		return 0, fmt.Errorf(LinkTypeUnknown, o)
	}
}

type CablingMapLinkEndpoint struct {
	Interface *CablingMapLinkEndpointInterface
	System    *CablingMapLinkEndpointSystem
}

type rawCablingMapLinkEndpoint struct {
	Interface *rawCablingMapLinkEndpointInterface `json:"interface"`
	System    *rawCablingMapLinkEndpointSystem    `json:"system"`
}

func (o *rawCablingMapLinkEndpoint) polish() (*CablingMapLinkEndpoint, error) {
	var err error
	var result CablingMapLinkEndpoint

	if o.Interface != nil {
		if result.Interface, err = o.Interface.polish(); err != nil {
			return nil, err
		}
	}

	if o.System != nil {
		if result.System, err = o.System.polish(); err != nil {
			return nil, err
		}
	}

	return &result, nil
}

type CablingMapLinkEndpointInterface struct {
	OperationState InterfaceOperationState
	IfName         *string
	PortChannelId  *int
	IfType         InterfaceType
	Id             ObjectId
	LagMode        RackLinkLagMode
}

type rawCablingMapLinkEndpointInterface struct {
	OperationState interfaceOperationState `json:"operation_state"`
	IfName         *string                 `json:"if_name"`
	PortChannelId  *int                    `json:"port_channel_id"`
	IfType         interfaceType           `json:"if_type"`
	Id             ObjectId                `json:"id"`
	LagMode        rackLinkLagMode         `json:"lag_mode"`
}

func (o *rawCablingMapLinkEndpointInterface) polish() (*CablingMapLinkEndpointInterface, error) {
	operationState, err := o.OperationState.parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse operation state %q - %w", o.OperationState, err)
	}

	ifType, err := o.IfType.parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse interface type %q - %w", o.IfType, err)
	}

	lagMode, err := o.LagMode.parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse LAG mode %q - %w", o.LagMode, err)
	}

	return &CablingMapLinkEndpointInterface{
		OperationState: InterfaceOperationState(operationState),
		IfName:         o.IfName,
		PortChannelId:  o.PortChannelId,
		IfType:         InterfaceType(ifType),
		Id:             o.Id,
		LagMode:        RackLinkLagMode(lagMode),
	}, nil
}

type CablingMapLinkEndpointSystem struct {
	Role  SystemNodeRole
	Id    ObjectId
	Label string
}

type rawCablingMapLinkEndpointSystem struct {
	Role  systemNodeRole `json:"role"`
	Id    ObjectId       `json:"id"`
	Label string         `json:"label"`
}

func (o *rawCablingMapLinkEndpointSystem) polish() (*CablingMapLinkEndpointSystem, error) {
	role, err := o.Role.parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse system role %q - %w", o.Role, err)
	}

	return &CablingMapLinkEndpointSystem{
		Role:  SystemNodeRole(role),
		Id:    o.Id,
		Label: o.Label,
	}, nil
}

type CablingMapLink struct {
	TagLabels       []string
	Speed           LogicalDevicePortSpeed
	AggregateLinkId ObjectId
	GroupLabel      string
	Label           string
	Role            LinkRole
	Endpoints       []CablingMapLinkEndpoint
	Type            LinkType
	Id              ObjectId
}

type rawCablingMapLink struct {
	TagLabels       []string                    `json:"tags"`
	Speed           LogicalDevicePortSpeed      `json:"speed"`
	AggregateLinkId ObjectId                    `json:"aggregate_link_id"`
	GroupLabel      string                      `json:"group_label"`
	Label           string                      `json:"label"`
	Role            linkRole                    `json:"role"`
	Endpoints       []rawCablingMapLinkEndpoint `json:"endpoints"`
	Type            linkType                    `json:"type"`
	Id              ObjectId                    `json:"id"`
}

func (o *rawCablingMapLink) polish() (*CablingMapLink, error) {
	lRole, err := o.Role.parse()
	if err != nil {
		return nil, err
	}

	endpoints := make([]CablingMapLinkEndpoint, len(o.Endpoints))
	for i, endpoint := range o.Endpoints {
		polished, err := endpoint.polish()
		if err != nil {
			return nil, err
		}
		endpoints[i] = *polished
	}

	lType, err := o.Type.parse()
	if err != nil {
		return nil, fmt.Errorf("failed parsing link type %q, - %w", o.Type, err)
	}

	result := &CablingMapLink{
		TagLabels:       o.TagLabels,
		Speed:           o.Speed,
		AggregateLinkId: o.AggregateLinkId,
		GroupLabel:      o.GroupLabel,
		Label:           o.Label,
		Role:            LinkRole(lRole),
		Endpoints:       endpoints,
		Type:            LinkType(lType),
		Id:              o.Id,
	}

	return result, nil
}

// GetCablingMapLinks returns []CablingMapLink representing every link in the blueprint
func (o *TwoStageL3ClosClient) GetCablingMapLinks(ctx context.Context) ([]CablingMapLink, error) {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintCablingMap, o.blueprintId))
	if err != nil {
		return nil, err
	}

	params := apstraUrl.Query()
	params.Set(includeLagParam, "true")
	apstraUrl.RawQuery = params.Encode()

	response := struct {
		Links []rawCablingMapLink `json:"links"`
	}{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]CablingMapLink, len(response.Links))
	for i, rawLink := range response.Links {
		polishedLink, err := rawLink.polish()
		if err != nil {
			return nil, convertTtaeToAceWherePossible(err)
		}
		result[i] = *polishedLink
	}

	return result, nil
}

// GetCablingMapLinksBySystem returns []CablingMapLink representing every link (including LAGs)
func (o *TwoStageL3ClosClient) GetCablingMapLinksBySystem(ctx context.Context, systemNodeId ObjectId) ([]CablingMapLink, error) {
	apstraUrl, err := url.Parse(fmt.Sprintf(apiUrlBlueprintCablingMap, o.blueprintId))
	if err != nil {
		return nil, err
	}

	params := apstraUrl.Query()
	params.Set(includeLagParam, "true")
	params.Set(linksBySystemParam, systemNodeId.String())
	apstraUrl.RawQuery = params.Encode()

	response := struct {
		Links []rawCablingMapLink `json:"links"`
	}{}

	err = o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		url:         apstraUrl,
		apiResponse: &response,
	})
	if err != nil {
		return nil, convertTtaeToAceWherePossible(err)
	}

	result := make([]CablingMapLink, len(response.Links))
	for i, rawLink := range response.Links {
		polishedLink, err := rawLink.polish()
		if err != nil {
			return nil, convertTtaeToAceWherePossible(err)
		}
		result[i] = *polishedLink
	}

	return result, nil
}
