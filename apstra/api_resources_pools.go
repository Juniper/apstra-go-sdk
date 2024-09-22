package apstra

import "fmt"

type (
	PoolStatus int
	poolStatus string
)

const (
	PoolStatusUnused = PoolStatus(iota)
	PoolStatusInUse
	PoolStatusDeleting
	PoolStatusCreating
	PoolStatusUnknown = "unknown redundancy protocol '%s'"

	poolStatusUnused   = poolStatus("not_in_use")
	poolStatusInUse    = poolStatus("in_use")
	poolStatusDeleting = poolStatus("marked_for_deletion")
	poolStatusCreating = poolStatus("creation_in_progress")
	poolStatusUnknown  = "unknown redundancy protocol '%d'"
)

func (o PoolStatus) Int() int {
	return int(o)
}

func (o PoolStatus) String() string {
	switch o {
	case PoolStatusUnused:
		return string(poolStatusUnused)
	case PoolStatusInUse:
		return string(poolStatusInUse)
	case PoolStatusDeleting:
		return string(poolStatusDeleting)
	case PoolStatusCreating:
		return string(poolStatusCreating)
	default:
		return fmt.Sprintf(poolStatusUnknown, o)
	}
}

func (o poolStatus) string() string {
	return string(o)
}

func (o poolStatus) parse() (int, error) {
	switch o {
	case poolStatusUnused:
		return int(PoolStatusUnused), nil
	case poolStatusInUse:
		return int(PoolStatusInUse), nil
	case poolStatusDeleting:
		return int(PoolStatusDeleting), nil
	case poolStatusCreating:
		return int(PoolStatusCreating), nil
	default:
		return 0, fmt.Errorf(PoolStatusUnknown, o)
	}
}
