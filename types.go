package goapstra

import "fmt"

const (
	vlanMin = 1
	vlanMax = 4094
)

type Vlan uint16

//lint:ignore U1000 keep for future use
func (o Vlan) validate() error {
	if o < vlanMin || o > vlanMax {
		return fmt.Errorf("Vlan ID %d out of range", o)
	}
	return nil
}
