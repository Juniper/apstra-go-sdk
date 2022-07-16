package goapstra

import (
	"context"
	"fmt"
	"strings"
)

const (
	apiUrlBlueprintQueryEngine = apiUrlBlueprintById + apiUrlPathDelim + "qe"
	qEElementAttributeSep      = ","

	queryEngineQueryTypeUrlParam = "type"
)

type QEEType int

const (
	qEETypeNode = "node"
	qEETypeIn   = "in_"
	qEETypeOut  = "out"
)

type QEQueryType int

const (
	QEQueryTypeNone = QEQueryType(iota)
	QEQueryTypeConfig
	QEQueryTypeDeployed
	QEQueryTypeOperation
	QEQueryTypeStaging

	qEQueryTypeNone      = ""
	qEQueryTypeConfig    = "config"
	qEQueryTypeDeployed  = "deployed"
	qEQueryTypeOperation = "operation"
	qEQueryTypeStaging   = "staging"
	qEQueryTypeUnknown   = "unknown query type %d"
)

func (o QEQueryType) string() string {
	switch o {
	case QEQueryTypeNone:
		return qEQueryTypeNone
	case QEQueryTypeConfig:
		return qEQueryTypeConfig
	case QEQueryTypeDeployed:
		return qEQueryTypeDeployed
	case QEQueryTypeOperation:
		return qEQueryTypeOperation
	case QEQueryTypeStaging:
		return qEQueryTypeStaging
	default:
		return fmt.Sprintf(qEQueryTypeUnknown, o)
	}
}

// per apstra API
type queryEngineQuery struct {
	Query string `json:"query"`
}

// per apstra API
type QueryEngineResponse struct {
	Count int           `json:"count"`
	Items []interface{} `json:"items"`
}

type QEAttrVal interface {
	String() string
}

type QEEAttribute struct {
	key   string
	value QEAttrVal
}

func (o QEEAttribute) String() string {
	return fmt.Sprintf("%s=%s", o.key, o.value.String())
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
		return "true"
	}
	return "false"
}

func (o *Client) newQuery(blueprint ObjectId) *QEQuery {
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
	queryType    QEQueryType
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

func (o *QEQuery) SetType(t QEQueryType) *QEQuery {
	o.queryType = t
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
	resp, err := o.client.runQuery(ctx, o.blueprint, o)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
