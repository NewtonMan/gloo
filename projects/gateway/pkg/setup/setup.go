package setup

import (
	"github.com/solo-io/solo-projects/pkg/utils/setuputils"
	"github.com/solo-io/solo-projects/projects/gateway/pkg/syncer"
)

func Main() error {
	return setuputils.Main("gateway", syncer.Setup)
}
