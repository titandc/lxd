package zone

import (
	"github.com/lxc/lxd/lxd/cluster/request"
	"github.com/lxc/lxd/lxd/state"
	"github.com/lxc/lxd/shared/api"

	"strings"
)

// NetworkZone represents a Network zone.
type NetworkZone interface {
	// Initialise.
	init(state *state.State, id int64, projectName string, zoneInfo *api.NetworkZone)

	// Info.
	ID() int64
	Project() string
	Info() *api.NetworkZone
	Etag() []interface{}
	UsedBy() ([]string, error)
	Content() (*strings.Builder, error)

	// Internal validation.
	validateName(name string) error
	validateConfig(config *api.NetworkZonePut) error

	// Modifications.
	Update(config *api.NetworkZonePut, clientType request.ClientType) error
	Delete() error
}
