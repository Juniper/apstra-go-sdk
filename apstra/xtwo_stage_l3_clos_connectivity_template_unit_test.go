package apstra

import "testing"

func TestThing(t *testing.T) {
	vnNodeId := ObjectId("abc")

	xa := ConnectivityTemplatePrimitiveAttributesAttachSingleVlan{
		Tagged:   true,
		VnNodeId: &vnNodeId,
	}

	x := xConnectivityTemplatePrimitive{
		id:          nil,
		userData:    nil,
		attributes:  &xa,
		subpolicies: nil,
		batchId:     nil,
		pipelineId:  nil,
	}

	p, err := x.rawPipeline()
	if err != nil {
		t.Fatal(err)
	}
	_ = p
}
