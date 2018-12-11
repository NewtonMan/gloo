package flagutils

import (
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/cmd/options"
	"github.com/spf13/pflag"
)

func AddDockerSecretFlags(set *pflag.FlagSet, install *options.Install) {
	set.StringVar(&install.DockerAuth.Email, "docker-email", "", "Email for docker registry. Use for pulling private images.")
	set.StringVar(&install.DockerAuth.Username, "docker-username", "", "Username for Docker registry authentication. Use for pulling private images.")
	set.StringVar(&install.DockerAuth.Password, "docker-password", "", "Password for docker registry authentication. Use for pulling private images.")
	set.StringVar(&install.DockerAuth.Server, "docker-server", "https://index.docker.io/v1/", "Docker server to use for pulling images")
}
