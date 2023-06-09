package apstra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	apiUrlBlueprintQueryEngine = apiUrlBlueprintById + apiUrlPathDelim + "qe"
	qEElementAttributeSep      = ","
)

type QEQuery interface {
	Do(context.Context, interface{}) error
	String() string
	getBlueprintType() BlueprintType
}

var _ QEQuery = &PathQuery{}
var _ QEQuery = &MatchQuery{}

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

func (o *QEElement) String() string {
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
	if len(o) == 0 { // handle <nil> gracefully
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

type QEIntVal int

func (o QEIntVal) String() string {
	return strconv.Itoa(int(o))
}

func (o *Client) newQuery(blueprintId ObjectId) *PathQuery {
	return &PathQuery{
		client:      o,
		blueprintId: blueprintId,
	}
}

type PathQuery struct {
	firstElement  *QEElement
	client        *Client
	context       context.Context
	blueprintId   ObjectId
	blueprintType BlueprintType
	where         []string
}

func (o *PathQuery) getBlueprintType() BlueprintType {
	return o.blueprintType
}

func (o *PathQuery) Do(ctx context.Context, response interface{}) error {
	return o.client.runQuery(ctx, o.blueprintId, o, response)
}

func (o *PathQuery) SetBlueprintId(id ObjectId) *PathQuery {
	o.blueprintId = id
	return o
}

func (o *PathQuery) SetBlueprintType(t BlueprintType) *PathQuery {
	o.blueprintType = t
	return o
}

func (o *PathQuery) SetClient(client *Client) *PathQuery {
	o.client = client
	return o
}

func (o *PathQuery) String() string {
	sb := strings.Builder{}

	var next *QEElement
	if o.firstElement != nil {
		sb.WriteString(o.firstElement.String())
		next = o.firstElement.getNext()
	}
	for next != nil {
		sb.WriteString(".")
		sb.WriteString(next.String())
		next = next.next
	}
	for _, where := range o.where {
		sb.WriteString(".where(" + where + ")")
	}
	return sb.String()
}

func (o *PathQuery) Where(where string) *PathQuery {
	o.where = append(o.where, where)
	return o
}

func (o *PathQuery) addElement(elementType string, attributes []QEEAttribute) *PathQuery {
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

func (o *PathQuery) Node(attributes []QEEAttribute) *PathQuery {
	return o.addElement(qEETypeNode, attributes)
}
func (o *PathQuery) Out(attributes []QEEAttribute) *PathQuery {
	return o.addElement(qEETypeOut, attributes)
}
func (o *PathQuery) In(attributes []QEEAttribute) *PathQuery {
	return o.addElement(qEETypeIn, attributes)
}

type MatchQueryElement struct {
	mqeType string
	value   QEAttrVal
	next    *MatchQueryElement
}

func (o *MatchQueryElement) String() string {
	return fmt.Sprintf("%s(%s)", o.mqeType, o.value.String())
}

func (o *MatchQueryElement) getNext() *MatchQueryElement {
	return o.next
}

func (o *MatchQueryElement) getLast() *MatchQueryElement {
	last := o
	next := last.getNext()
	for next != nil {
		last = next
		next = last.getNext()
	}
	return last
}

type MatchQueryDistinct []string

func (o MatchQueryDistinct) String() string {
	if len(o) == 0 { // handle <nil> gracefully
		return "[]"
	}
	return "['" + strings.Join(o, "','") + "']"
}

type MatchQuery struct {
	client        *Client
	context       context.Context
	blueprintId   ObjectId
	blueprintType BlueprintType
	match         []PathQuery
	firstElement  *MatchQueryElement
	where         []string
}

//func (o *MatchQuery) Having(v QEAttrVal) *MatchQuery          {} // todo
//func (o *MatchQuery) Where(v QEAttrVal) *MatchQuery           {} // todo
//func (o *MatchQuery) EnsureDifferent(v QEAttrVal) *MatchQuery {} // todo

func (o *MatchQuery) Distinct(distinct MatchQueryDistinct) *MatchQuery {
	o.addElement("distinct", distinct)
	return o
}

func (o *MatchQuery) addElement(t string, v QEAttrVal) *MatchQuery {
	newElement := MatchQueryElement{
		mqeType: t,
		value:   v,
	}
	if o.firstElement == nil {
		o.firstElement = &newElement
		return o
	}
	o.firstElement.getLast().next = &newElement
	return o

}

func (o *MatchQuery) getBlueprintType() BlueprintType {
	return o.blueprintType
}

func (o *MatchQuery) Do(ctx context.Context, response interface{}) error {
	if o.client == nil {
		return errors.New("attempt to execute query without setting client")
	}
	return o.client.runQuery(ctx, o.blueprintId, o, response)
}

func (o *MatchQuery) SetBlueprintId(id ObjectId) *MatchQuery {
	o.blueprintId = id
	return o
}

func (o *MatchQuery) SetBlueprintType(t BlueprintType) *MatchQuery {
	o.blueprintType = t
	return o
}

func (o *MatchQuery) SetClient(client *Client) *MatchQuery {
	o.client = client
	return o
}

func (o *MatchQuery) String() string {
	var sb strings.Builder
	sb.WriteString("match(")
	for i := range o.match {
		if i != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(o.match[i].String())
	}
	sb.WriteString(")")

	var next *MatchQueryElement
	if o.firstElement != nil {
		sb.WriteString(".")
		sb.WriteString(o.firstElement.String())
		next = o.firstElement.getNext()
	}
	for next != nil {
		sb.WriteString(".")
		sb.WriteString(next.String())
		next = next.next
	}

	for _, where := range o.where {
		sb.WriteString(".where(" + where + ")")
	}

	return sb.String()
}

func (o *MatchQuery) Where(where string) *MatchQuery {
	o.where = append(o.where, where)
	return o
}

func (o *MatchQuery) Match(q *PathQuery) *MatchQuery {
	o.match = append(o.match, *q)
	return o
}
