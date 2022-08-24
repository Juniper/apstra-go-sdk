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

func getTestClientCfgs() ([]testClientCfg, error) {
	var testClientCfgs []testClientCfg

	// add cloudlabs clients to testClients slice
	clTestClientCfgs, err := getCloudlabsTestClientCfgs()
	if err != nil {
		return nil, err
	}
	testClientCfgs = append(testClientCfgs, clTestClientCfgs...)

	// add future type clients (slicer?) to testClients slice here

	return testClientCfgs, nil
}
