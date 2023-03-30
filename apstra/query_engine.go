package apstra

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	apiUrlBlueprintQueryEngine = apiUrlBlueprintById + apiUrlPathDelim + "qe"
	qEElementAttributeSep      = ","
)

type QEEType int

const (
	qEETypeNode = "node"
	qEETypeIn   = "in_"
	qEETypeOut  = "out"
)

type QueryEngineResponse []json.RawMessage

type QEAttrVal interface {
	String() string
}

type QEEAttribute struct {
	Key   string
	Value QEAttrVal
}

func (o QEEAttribute) String() string {
	return fmt.Sprintf("%s=%s", o.Key, o.Value.String())
}

type QEElement struct {
	qeeType    string
	attributes []QEEAttribute
	next       *QEElement
}

func (o *QEElement) getNext() *QEElement {
	return o.next
}

func (o *QEElement) getLast() *QEElement {
	last := o
	next := last.getNext()
	for next != nil {
		last = next
		next = last.getNext()
	}
	return last
}

func (o QEElement) String() string {
	attrsSB := strings.Builder{}

	// add first attribute to string builder without leading separator
	if len(o.attributes) > 0 {
		attrsSB.WriteString(o.attributes[0].String())
	}

	// remaining attributes added with leading separator
	for _, a := range o.attributes[1:] {
		attrsSB.WriteString(qEElementAttributeSep)
		attrsSB.WriteString(a.String())
	}
	return fmt.Sprintf("%s(%s)", o.qeeType, attrsSB.String())
}

type QEStringValIsIn []string

func (o QEStringValIsIn) String() string {
	if len(o) == 0 {
		return "is_in([])"
	}
	return "is_in(['" + strings.Join(o, "','") + "'])"
}

type QEStringValNotIn []string

func (o QEStringValNotIn) String() string {
	if len(o) == 0 {
		return "not_in([])"
	}
	return "not_in(['" + strings.Join(o, "','") + "'])"
}

type QEStringVal string

func (o QEStringVal) String() string {
	return fmt.Sprintf("'%s'", string(o))
}

type QEBoolVal bool

func (o QEBoolVal) String() string {
	if o {
		return "True"
	}
	return "False"
}

func (o *Client) newQuery(blueprint ObjectId) *QEQuery {
	return &QEQuery{
		client:    o,
		blueprint: blueprint,
	}
}

type QEQuery struct {
	firstElement  *QEElement
	client        *Client
	context       context.Context
	blueprint     ObjectId
	blueprintType BlueprintType
}

func (o *QEQuery) addElement(elementType string, attributes []QEEAttribute) *QEQuery {
	newElement := QEElement{
		qeeType:    elementType,
		attributes: attributes,
	}
	if o.firstElement == nil {
		o.firstElement = &newElement
		return o
	}
	o.firstElement.getLast().next = &newElement
	return o
}

func (o *QEQuery) Node(attributes []QEEAttribute) *QEQuery {
	return o.addElement(qEETypeNode, attributes)
}
func (o *QEQuery) Out(attributes []QEEAttribute) *QEQuery {
	return o.addElement(qEETypeOut, attributes)
}
func (o *QEQuery) In(attributes []QEEAttribute) *QEQuery {
	return o.addElement(qEETypeIn, attributes)
}

func (o *QEQuery) SetContext(ctx context.Context) *QEQuery {
	o.context = ctx
	return o
}

func (o *QEQuery) SetType(t BlueprintType) *QEQuery {
	o.blueprintType = t
	return o
}

func (o *QEQuery) String() string {
	return o.string()
}

func (o *QEQuery) string() string {
	sb := strings.Builder{}
	if o.firstElement != nil {
		sb.WriteString(o.firstElement.String())
	}
	next := o.firstElement.getNext()
	for next != nil {
		sb.WriteString(".")
		sb.WriteString(next.String())
		next = next.next
	}
	return sb.String()
}

func (o *QEQuery) Do(response interface{}) error {
	ctx := o.context
	if o.context == nil {
		ctx = context.TODO()
	}
	return o.client.runQuery(ctx, o.blueprint, o, response)
}
