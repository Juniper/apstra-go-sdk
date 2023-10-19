package apstra

import "fmt"

const (
	vlanMin = 1
	vlanMax = 4094

	vniMin = 4096
	vniMax = 16777214
)

type Vlan uint16

//lint:ignore U1000 keep for future use
func (o Vlan) validate() error {
	if o < vlanMin || o > vlanMax {
		return fmt.Errorf("VLAN %d out of range", o)
	}
	return nil
}

type VNI uint32

//lint:ignore U1000 keep for future use
func (o VNI) validate() error {
	if o < vniMin || o > vniMax {
		return fmt.Errorf("VNI %d out of range", o)
	}
	return nil
}

type RtPolicy struct {
	ImportRTs []string `json:"import_RTs"`
	ExportRTs []string `json:"export_RTs"`
}
