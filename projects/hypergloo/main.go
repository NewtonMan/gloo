package main

import (
	"context"
	"flag"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	fdssetup "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/setup"
	uds "github.com/solo-io/solo-projects/projects/discovery/pkg/uds/setup"
	gatewaysetup "github.com/solo-io/solo-projects/projects/gateway/pkg/setup"
	gloosetup "github.com/solo-io/solo-projects/projects/gloo/pkg/setup"
	sqoopsetup "github.com/solo-io/solo-projects/projects/sqoop/pkg/setup"
)

func main() {
	stats.StartStatsServer()
	if err := run(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}

func run() error {
	contextutils.LoggerFrom(context.TODO()).Infof("hypergloo!")
	flag.Parse()
	errs := make(chan error)
	go func() {
		errs <- gloosetup.Main()
	}()
	go func() {
		errs <- gatewaysetup.Main()
	}()
	go func() {
		errs <- sqoopsetup.Main()
	}()
	go func() {
		errs <- uds.Main()
	}()
	go func() {
		errs <- fdssetup.Main()
	}()
	return <-errs
}
