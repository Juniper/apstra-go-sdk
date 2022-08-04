package goapstra

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	portAny       = "any"
	portRangeSep  = "-"
	portRangesSep = ","
)

// RULE_SCHEMA = {
//    'id': s.Optional(s.NodeId(description='ID of the rule node')),
//    'label': s.GenericName(description='Unique user-friendly name of the rule'),
//    'protocol': s.SecurityRuleProtocol(),
//    'src_port': s.PortSetOrAny(),
//    'dst_port': s.PortSetOrAny(),
//    'description': s.Optional(s.Description(), load_default=''),
//    'action': s.SecurityRuleAction() //             ['deny', 'deny_log', 'permit', 'permit_log'],
//}

type PolicyRuleAction int
type policyRuleAction string

const (
	PolicyRuleActionDeny = iota
	PolicyRuleActionDenyLog
	PolicyRuleActionPermit
	PolicyRuleActionPermitLog
	PolicyRuleActionUnknown = "unknown policy action '%s'"

	policyRuleActionDeny      = policyRuleAction("deny")
	policyRuleActionDenyLog   = policyRuleAction("deny_log")
	policyRuleActionPermit    = policyRuleAction("permit")
	policyRuleActionPermitLog = policyRuleAction("permit_log")
	policyRuleActionUnknown   = "unknown policy action %d"
)

func (o PolicyRuleAction) Int() int {
	return int(o)
}

func (o PolicyRuleAction) String() string {
	return string(o.raw())
}

func (o PolicyRuleAction) raw() policyRuleAction {
	switch o {
	case PolicyRuleActionDeny:
		return policyRuleActionDeny
	case PolicyRuleActionDenyLog:
		return policyRuleActionDenyLog
	case PolicyRuleActionPermit:
		return policyRuleActionPermit
	case PolicyRuleActionPermitLog:
		return policyRuleActionPermitLog
	default:
		return policyRuleAction(fmt.Sprintf(policyRuleActionUnknown, o))
	}
}

func (o policyRuleAction) string() string {
	return string(o)
}

func (o policyRuleAction) parse() (int, error) {
	switch o {
	case policyRuleActionDeny:
		return PolicyRuleActionDeny, nil
	case policyRuleActionDenyLog:
		return PolicyRuleActionDenyLog, nil
	case policyRuleActionPermit:
		return PolicyRuleActionPermit, nil
	case policyRuleActionPermitLog:
		return PolicyRuleActionPermitLog, nil
	default:
		return 0, fmt.Errorf(PolicyRuleActionUnknown, o)
	}
}

type PortRange struct {
	first uint16
	last  uint16
}

func (o PortRange) string() string {
	if o.first < o.last {
		return strconv.Itoa(int(o.first)) + portRangeSep + strconv.Itoa(int(o.last))
	} else {
		return strconv.Itoa(int(o.last)) + portRangeSep + strconv.Itoa(int(o.first))
	}
}

type rawPortRanges string

func (o rawPortRanges) parse() (PortRanges, error) {
	if o == portAny {
		return []PortRange{}, nil
	}

	rawRangeSlice := strings.Split(portRangesSep, string(o))
	result := make([]PortRange, len(rawRangeSlice))
	for i, raw := range rawRangeSlice {
		var first, last uint64
		var err error
		portStrs := strings.Split(portRangeSep, raw)
		switch len(portStrs) {
		case 1:
			first, err = strconv.ParseUint(raw, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing port range '%s' - %w", raw, err)
			}
			last = first
		case 2:
			first, err = strconv.ParseUint(portStrs[0], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing first element of port range '%s' - %w", raw, err)
			}
			last, err = strconv.ParseUint(portStrs[0], 10, 16)
			if err != nil {
				return nil, fmt.Errorf("error parsing last element of port range '%s' - %w", raw, err)
			}
		default:
			return nil, fmt.Errorf("cannot parse port range '%s'", raw)
		}
		if first > math.MaxUint16 || last > math.MaxUint16 {
			return nil, fmt.Errorf("port spec '%s' falls outside of range %d-%d", raw, 0, math.MaxUint16)
		}
		result[i] = PortRange{
			first: uint16(first),
			last:  uint16(last),
		}
	}
	return result, nil
}

type PortRanges []PortRange

func (o PortRanges) string() string {
	if len(o) == 0 {
		return portAny
	}
	sb := strings.Builder{}
	sb.WriteString(o[0].string())
	for _, pr := range o[1:] {
		sb.WriteString(portRangesSep + pr.string())
	}
	return sb.String()
}

type PolicyRule struct {
	Id          ObjectId         `json:"id"`
	Label       string           `json:"label"`
	Description string           `json:"description"`
	Protocol    string           `json:"protocol"`
	Action      PolicyRuleAction `json:"action"`
	SrcPort     PortRanges       `json:"src_port"`
	DstPort     PortRanges       `json:"dst_port"`
}

func (o PolicyRule) raw() *rawPolicyRule {
	return &rawPolicyRule{
		Id:          o.Id,
		Label:       o.Label,
		Description: o.Description,
		Protocol:    o.Protocol,
		Action:      o.Action.raw(),
		SrcPort:     rawPortRanges(o.SrcPort.string()),
		DstPort:     rawPortRanges(o.DstPort.string()),
	}
}

type rawPolicyRule struct {
	Id          ObjectId         `json:"id,omitempty"`
	Label       string           `json:"label"`
	Description string           `json:"description"`
	Protocol    string           `json:"protocol"`
	Action      policyRuleAction `json:"action"`
	SrcPort     rawPortRanges    `json:"src_port"`
	DstPort     rawPortRanges    `json:"dst_port"`
}

func (o rawPolicyRule) polish() (*PolicyRule, error) {
	action, err := o.Action.parse()
	if err != nil {
		return nil, err
	}
	srcPort, err := o.SrcPort.parse()
	if err != nil {
		return nil, err
	}
	dstPort, err := o.SrcPort.parse()
	if err != nil {
		return nil, err
	}
	return &PolicyRule{
		Id:          o.Id,
		Label:       o.Label,
		Description: o.Description,
		Protocol:    o.Protocol,
		Action:      PolicyRuleAction(action),
		SrcPort:     srcPort,
		DstPort:     dstPort,
	}, nil
}

type Policy struct {
	Enabled             bool         `json:"enabled"`
	Label               string       `json:"label"`
	Description         string       `json:"description"`
	SrcApplicationPoint string       `json:"src_application_point"`
	DstApplicationPoint string       `json:"dst_application_point"`
	Rules               []PolicyRule `json:"rules"`
	Tags                []TagLabel   `json:"tags"`
}
