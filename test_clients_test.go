package goapstra

const (
	clientTypeCloudlabs = "cloudlabs"
	clientTypeSlicer    = "slicer"
)

var testClients []testClient

func getTestClients() ([]testClient, error) {
	if testClients != nil {
		return testClients, nil
	}

	// add cloudlabs clients to testClients slice
	clTestClients, err := getCloudlabsTestClients()
	if err != nil {
		return nil, err
	}
	testClients = append(testClients, clTestClients...)

	// add future type clients (slicer?) to testClients slice here

	return testClients, nil
}
