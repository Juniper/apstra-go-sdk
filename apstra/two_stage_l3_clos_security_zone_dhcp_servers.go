package apstra

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// GetSecurityZoneDhcpServers returns []net.IP representing the DHCP relay
// targets for the security zone specified by id.
func (o TwoStageL3ClosClient) GetSecurityZoneDhcpServers(ctx context.Context, id string) ([]net.IP, error) {
	response := &struct {
		Items []string `json:"items"`
	}{}
	err := o.client.talkToApstra(ctx, &talkToApstraIn{
		method:      http.MethodGet,
		urlStr:      fmt.Sprintf(apiUrlBlueprintSecurityZoneByIdDhcpServers, o.blueprintId, id),
		apiResponse: response,
	})
	if err != nil {
		return nil, fmt.Errorf("getting dhcp servers for blueprint %q zone %q: %w", o.Id(), id, err)
	}

	result := make([]net.IP, len(response.Items))
	for i, s := range response.Items {
		result[i] = net.ParseIP(s)
		if result[i] == nil {
			err = errors.Join(err, fmt.Errorf("failed to parse blueprint %s security zone %s dhcp server"+
				" at index %d; expected an IP address, got %q", o.blueprintId, id, i, s))
		}
	}

	return result, err
}

// SetSecurityZoneDhcpServers assigns the []net.IP as DHCP relay targets for
// the specified security zone, overwriting whatever is there. On the Apstra
// side, the servers seem to be maintained as an ordered list with duplicates
// permitted (though the web UI sorts the data prior to display)
func (o TwoStageL3ClosClient) SetSecurityZoneDhcpServers(ctx context.Context, id string, ips []net.IP) error {
	items := make([]string, len(ips))
	for i, ip := range ips {
		items[i] = ip.String()
	}

	return convertTtaeToAceWherePossible(o.client.talkToApstra(ctx, &talkToApstraIn{
		method: http.MethodPut,
		urlStr: fmt.Sprintf(apiUrlBlueprintSecurityZoneByIdDhcpServers, o.blueprintId, id),
		apiInput: &struct {
			Items []string `json:"items"`
		}{
			Items: items,
		},
	}))
}
