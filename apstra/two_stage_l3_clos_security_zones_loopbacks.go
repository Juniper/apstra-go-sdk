package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Juniper/apstra-go-sdk/apstra/compatibility"
	"net/http"
	"net/netip"
)

const apiUrlBlueprintSecurityZoneLoopbacksById = apiUrlBlueprintSecurityZoneById + apiUrlPathDelim + "loopbacks"

var _ json.Marshaler = (*SecurityZoneLoopback)(nil)

// SecurityZoneLoopback is intended to be used with the SetSecurityZoneLoopbacks() method
// and the apiUrlBlueprintSecurityZoneLoopbacksById API endpoint. It is possible to produce
// three different outcomes in the rendered JSON for both IPv4Addr and IPv6Addr elements:
//
//  1. When the element and its IP and Mask elements are non-nil, a string will be produced
//     when the struct is marshaled as JSON.
//  2. When the element is non-nil but contains a nil IP or Mask, a `null` will be produced
//     when the struct is marshaled as JSON.
//  3. When the element is nil, no output will be produced for that element when the struct
//     is marshaled as JSON.
//
// Example:
//
//	aVal := netip.MustParsePrefix("192.0.2.0/32")
//	a := apstra.SecurityZoneLoopback{IPv4Addr: &aVal}
//	b := apstra.SecurityZoneLoopback{IPv4Addr: &netip.Prefix{}}
//	c := apstra.SecurityZoneLoopback{IPv4Addr: nil}
//
//	aJson, _ := json.Marshal(a)
//	bJson, _ := json.Marshal(b)
//	cJson, _ := json.Marshal(c)
//
//	fmt.Print(string(aJson) + "\n" + string(bJson) + "\n" + string(cJson) + "\n")
//
// Output:
//
//	{"ipv4_addr":"192.0.2.0/32"}
//	{"ipv4_addr":null}
//	{}
type SecurityZoneLoopback struct {
	IPv4Addr *netip.Prefix
	IPv6Addr *netip.Prefix
}

func (o SecurityZoneLoopback) MarshalJSON() ([]byte, error) {
	ipInfo := make(map[string]*string)

	if o.IPv4Addr != nil {
		if o.IPv4Addr.IsValid() {
			//if o.IPv4Addr.IP != nil || o.IPv4Addr.Mask != nil {
			ipInfo["ipv4_addr"] = toPtr(o.IPv4Addr.String())
		} else {
			ipInfo["ipv4_addr"] = nil
		}
	}

	if o.IPv6Addr != nil {
		if o.IPv6Addr.IsValid() {
			//if o.IPv6Addr.IP != nil || o.IPv6Addr.Mask != nil {
			ipInfo["ipv6_addr"] = toPtr(o.IPv6Addr.String())
		} else {
			ipInfo["ipv6_addr"] = nil
		}
	}

	return json.Marshal(ipInfo)
}

// SetSecurityZoneLoopbacks takes a map of SecurityZoneLoopback keyed by the loopback interface graph node ID.
func (o TwoStageL3ClosClient) SetSecurityZoneLoopbacks(ctx context.Context, szId ObjectId, loopbacks map[ObjectId]SecurityZoneLoopback) error {
	if !compatibility.SecurityZoneLoopbackApiSupported.Check(o.client.apiVersion) {
		return fmt.Errorf("SetSecurityZoneLoopbacks requires Apstra version %s, have version %s",
			compatibility.SecurityZoneLoopbackApiSupported, o.client.apiVersion,
		)
	}

	var apiInput struct {
		Loopbacks map[ObjectId]json.RawMessage `json:"loopbacks"`
	}
	apiInput.Loopbacks = make(map[ObjectId]json.RawMessage)

	for k, v := range loopbacks {
		rawJson, err := json.Marshal(v)
		if err != nil {
			return err
		}
		apiInput.Loopbacks[k] = rawJson
	}

	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:   http.MethodPatch,
		urlStr:   fmt.Sprintf(apiUrlBlueprintSecurityZoneLoopbacksById, o.blueprintId, szId),
		apiInput: apiInput,
	})
	if err != nil {
		return convertTtaeToAceWherePossible(err)
	}

	return nil
}
