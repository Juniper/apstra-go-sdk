package goapstra

import (
	"context"
	"fmt"
	"strings"
)

const (
	apiUrlBlueprintQueryEngine = apiUrlBlueprintById + apiUrlPathDelim + "qe"
	qEElementAttributeSep      = ","
)

type QEEType int

const (
	QEETypeNode = QEEType(iota)
	QEETypeIn
	QEETypeout

	qEETypeNode    = "node"
	qEETypeIn      = "in_"
	qEETypeOut     = "out"
	qEETypeUnknown = "unknown QueryEngine element type '%d'"
)

func (o QEEType) String() string {
	switch o {
	case QEETypeNode:
		return qEETypeNode
	case QEETypeIn:
		return qEETypeIn
	case QEETypeout:
		return qEETypeOut
	default:
		return fmt.Sprintf(qEETypeUnknown, o)
	}
}

type QEAttrVal interface {
	String() string
}

type QEEAttributes struct {
	key   string
	value QEAttrVal
}

func (o QEEAttributes) String() string {
	return fmt.Sprintf("%s=%s", o.key, o.value.String())
}

type QEElement struct {
	qeeType    string
	attributes []QEEAttributes
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
		return "true"
	}
	return "false"
}

func (o *Client) NewQuery(blueprint ObjectId) *QEQuery {
	return &QEQuery{
		client:    o,
		blueprint: blueprint,
	}
}

type QEQuery struct {
	firstElement *QEElement
	client       *Client
	context      context.Context
	blueprint    ObjectId
}

func (o *QEQuery) addElement(elementType string, attributes []QEEAttributes) *QEQuery {
	new := QEElement{
		qeeType:    elementType,
		attributes: attributes,
	}
	if o.firstElement == nil {
		o.firstElement = &new
		return o
	}
	o.firstElement.getLast().next = &new
	return o
}

func (o *QEQuery) Node(attributes []QEEAttributes) *QEQuery {
	return o.addElement(qEETypeNode, attributes)
}
func (o *QEQuery) Out(attributes []QEEAttributes) *QEQuery {
	return o.addElement(qEETypeOut, attributes)
}
func (o *QEQuery) In(attributes []QEEAttributes) *QEQuery {
	return o.addElement(qEETypeIn, attributes)
}

func (o *QEQuery) Context(ctx context.Context) *QEQuery {
	o.context = ctx
	return o
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

func (o *QEQuery) Do() (interface{}, error) {
	ctx := o.context
	if o.context == nil {
		ctx = context.TODO()
	}
	resp, err := o.client.runQuery(ctx, o.blueprint, &QueryEngineQuery{Query: o.string()})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
