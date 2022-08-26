package goapstra

const (
	clientTypeCloudlabs = "cloudlabs"
	clientTypeSlicer    = "slicer"
)

var testClients map[string]testClient

func getTestClients() (map[string]testClient, error) {
	if testClients != nil {
		return testClients, nil
	}
	testClients = make(map[string]testClient)

	// add cloudlabs clients to testClients slice
	clTestClients, err := getCloudlabsTestClients()
	if err != nil {
		return nil, err
	}
	for k, v := range clTestClients {
		testClients[k] = v
	}

	// add future type clients (slicer?) to testClients slice here

	return testClients, nil
}

func getTestClientCfgs() (map[string]testClientCfg, error) {
	var testClientCfgs map[string]testClientCfg

	// add cloudlabs clients to testClients slice
	clTestClientCfgs, err := getCloudlabsTestClientCfgs()
	if err != nil {
		return nil, err
	}
	for k, v := range clTestClientCfgs {
		testClientCfgs[k] = v
	}

	// add future type clients (slicer?) to testClients slice here

	return testClientCfgs, nil
}
