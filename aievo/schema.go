package aievo

import (
	"github.com/antgroup/aievo/environment"
)

type AIEvo struct {
	Handler Handler
	*environment.Environment
}
