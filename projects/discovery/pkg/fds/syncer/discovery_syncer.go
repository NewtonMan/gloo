package syncer

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1"
)

type syncer struct {
	fd *fds.FunctionDiscovery
}

func NewDiscoverySyncer(fd *fds.FunctionDiscovery) v1.DiscoverySyncer {
	s := &syncer{
		fd: fd,
	}
	return s
}

func (s *syncer) Sync(ctx context.Context, snap *v1.DiscoverySnapshot) error {
	ctx = contextutils.WithLogger(ctx, "syncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Infof("begin sync %v (%v upstreams)", snap.Hash(), len(snap.Upstreams.List()))
	defer logger.Infof("end sync %v", snap.Hash())

	logger.Debugf("%v", snap)

	return s.fd.Update(snap.Upstreams.List(), snap.Secrets.List())
}
