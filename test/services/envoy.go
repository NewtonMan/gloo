package services

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
	"time"

	"bytes"
	"io"

	"github.com/onsi/ginkgo"
	"github.com/pkg/errors"

	"github.com/solo-io/go-utils/log"
)

const (
	containerName = "e2e_envoy"
)

var adminPort = 19000

func (ei *EnvoyInstance) buildBootstrap() string {
	var b bytes.Buffer
	parsedTemplate.Execute(&b, ei)
	return b.String()
}

const envoyConfigTemplate = `
node:
 cluster: ingress
 id: {{.ID}}
 metadata:
  role: {{.Role}}

static_resources:
  clusters:
  - name: xds_cluster
    connect_timeout: 5.000s
    load_assignment:
      cluster_name: xds_cluster
      endpoints:
        - lb_endpoints:
            - endpoint:
                address:
                  socket_address:
                    address: {{.GlooAddr}}
                    port_value: {{.Port}}
    http2_protocol_options: {}
    type: STATIC

layered_runtime:
  layers:
  - name: static_layer
    static_layer:
      upstream:
        healthy_panic_threshold:
          value: 0
  - name: admin_layer
    admin_layer: {}

dynamic_resources:
  ads_config:
    transport_api_version: {{ .ApiVersion }}
    api_type: GRPC
    grpc_services:
    - envoy_grpc: {cluster_name: xds_cluster}
  cds_config:
    resource_api_version: {{ .ApiVersion }}
    ads: {}
  lds_config:
    resource_api_version: {{ .ApiVersion }}
    ads: {}
  
admin:
  access_log_path: /dev/null
  address:
    socket_address:
      address: 0.0.0.0
      port_value: {{.AdminPort}}

`

var parsedTemplate = template.Must(template.New("bootstrap").Parse(envoyConfigTemplate))

type EnvoyFactory struct {
	envoypath string
	tmpdir    string
	useDocker bool
	instances []*EnvoyInstance
}

func NewEnvoyFactory() (*EnvoyFactory, error) {
	// if an envoy binary is explicitly specified
	// use it
	envoypath := os.Getenv("ENVOY_BINARY")
	if envoypath != "" {
		log.Printf("Using envoy from environment variable: %s", envoypath)
		return &EnvoyFactory{
			envoypath: envoypath,
		}, nil
	}

	// if ENVOY_IMAGE_TAG is specified, always use docker image
	if imageTag := os.Getenv("ENVOY_IMAGE_TAG"); imageTag != "" {
		log.Printf("Using docker to run envoy")
		return &EnvoyFactory{useDocker: true}, nil
	}

	// maybe it is in the path?!
	envoypath, err := exec.LookPath("envoy")
	if err == nil {
		log.Printf("Using envoy from PATH: %s", envoypath)
		return &EnvoyFactory{
			envoypath: envoypath,
		}, nil
	}

	switch runtime.GOOS {
	case "darwin":
		log.Printf("Using docker to run envoy")

		return &EnvoyFactory{useDocker: true}, nil
	case "linux":
		// try to grab one form docker...
		tmpdir, err := ioutil.TempDir(os.Getenv("HELPER_TMP"), "envoy")
		if err != nil {
			return nil, err
		}

		envoyImageTag := os.Getenv("ENVOY_IMAGE_TAG")
		if envoyImageTag == "" {
			panic("The ENVOY_IMAGE_TAG env var is not set. Find valid tag names here https://quay.io/repository/solo-io/gloo-ee-envoy-wrapper?tab=tags")
		}
		log.Printf("Using envoy docker image tag: %s", envoyImageTag)

		bash := fmt.Sprintf(`
set -ex
CID=$(docker run -d  soloio/envoy:%s /bin/bash -c exit)

# just print the image sha for repoducibility
echo "Using Envoy Image:"
docker inspect soloio/envoy:%s -f "{{.RepoDigests}}"

docker cp $CID:/usr/local/bin/envoy .
docker rm $CID
    `, envoyImageTag, envoyImageTag)
		scriptfile := filepath.Join(tmpdir, "getenvoy.sh")

		ioutil.WriteFile(scriptfile, []byte(bash), 0755)

		cmd := exec.Command("bash", scriptfile)
		cmd.Dir = tmpdir
		cmd.Stdout = ginkgo.GinkgoWriter
		cmd.Stderr = ginkgo.GinkgoWriter
		if err := cmd.Run(); err != nil {
			return nil, err
		}

		return &EnvoyFactory{
			envoypath: filepath.Join(tmpdir, "envoy"),
			tmpdir:    tmpdir,
		}, nil

	default:
		return nil, errors.New("Unsupported OS: " + runtime.GOOS)
	}
}

func (ef *EnvoyFactory) EnvoyPath() string {
	return ef.envoypath
}

func (ef *EnvoyFactory) Clean() error {
	if ef == nil {
		return nil
	}
	if ef.tmpdir != "" {
		os.RemoveAll(ef.tmpdir)
	}
	instances := ef.instances
	ef.instances = nil
	for _, ei := range instances {
		ei.Clean()
	}
	return nil
}

type EnvoyInstance struct {
	RatelimitAddr string
	RatelimitPort uint32
	ID            string
	Role          string
	envoypath     string
	envoycfg      string
	logs          *bytes.Buffer
	cmd           *exec.Cmd
	useDocker     bool
	GlooAddr      string // address for gloo and services
	Port          uint32
	AdminPort     int32
	ApiVersion    string
	// Path to access logs for binary run
	AccessLogs string
}

func (ef *EnvoyFactory) NewEnvoyInstance() (*EnvoyInstance, error) {
	adminPort = adminPort + 1

	gloo := "127.0.0.1"
	var err error

	if ef.useDocker {
		gloo, err = localAddr()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	ei := &EnvoyInstance{
		envoypath:  ef.envoypath,
		useDocker:  ef.useDocker,
		GlooAddr:   gloo,
		AdminPort:  int32(adminPort),
		ApiVersion: "V3",
	}
	ef.instances = append(ef.instances, ei)
	return ei, nil

}

func (ei *EnvoyInstance) RunWithId(id string) error {
	ei.ID = id
	ei.Role = "default~proxy"

	return ei.runWithPort(8081)
}

func (ei *EnvoyInstance) Run(port int) error {
	ei.Role = "default~proxy"

	return ei.runWithPort(uint32(port))
}
func (ei *EnvoyInstance) RunWithRole(role string, port int) error {
	ei.Role = role
	return ei.runWithPort(uint32(port))
}

/*
func (ei *EnvoyInstance) DebugMode() error {

	_, err := http.Get("http://localhost:19000/logging?level=debug")

	return err
}
*/
func (ei *EnvoyInstance) runWithPort(port uint32) error {
	if ei.ID == "" {
		ei.ID = "ingress~for-testing"
	}
	ei.Port = port

	ei.envoycfg = ei.buildBootstrap()
	if ei.useDocker {
		err := ei.runContainer()
		if err != nil {
			return err
		}
		return nil
	}

	args := []string{
		"--config-yaml", ei.envoycfg,
		"--disable-hot-restart",
		"--log-level", "debug",
		"--concurrency", "1",
		"--file-flush-interval-msec", "10",
		"--bootstrap-version", "3",
	}

	// run directly
	cmd := exec.Command(ei.envoypath, args...)

	buf := &bytes.Buffer{}
	ei.logs = buf
	w := io.MultiWriter(ginkgo.GinkgoWriter, buf)
	cmd.Stdout = w
	cmd.Stderr = w

	runner := Runner{Sourcepath: ei.envoypath, ComponentName: "ENVOY"}
	cmd, err := runner.run(cmd)
	if err != nil {
		return err
	}
	ei.cmd = cmd
	return nil
}

func (ei *EnvoyInstance) Binary() string {
	return ei.envoypath
}

func (ei *EnvoyInstance) LocalAddr() string {
	return ei.GlooAddr
}

func (ei *EnvoyInstance) EnablePanicMode() error {
	return ei.setRuntimeConfiguration(fmt.Sprintf("upstream.healthy_panic_threshold=%d", 100))
}

func (ei *EnvoyInstance) DisablePanicMode() error {
	return ei.setRuntimeConfiguration(fmt.Sprintf("upstream.healthy_panic_threshold=%d", 0))
}

func (ei *EnvoyInstance) setRuntimeConfiguration(queryParameters string) error {
	_, err := http.Post(fmt.Sprintf("http://localhost:%d/runtime_modify?%s", ei.AdminPort, queryParameters), "", nil)
	return err
}

func (ei *EnvoyInstance) UseDocker() bool {
	return ei.useDocker
}

func (ei *EnvoyInstance) Clean() error {
	http.Post(fmt.Sprintf("http://localhost:%d/quitquitquit", ei.AdminPort), "", nil)
	if ei.cmd != nil {
		ei.cmd.Process.Kill()
		ei.cmd.Wait()
	}

	if ei.useDocker {
		if err := KillAndRemoveContainer(containerName); err != nil {
			return err
		}
	}

	// Wait till envoy is completely cleaned up
	request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", ei.AdminPort), nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	timeout := 5 // seconds
	timer := timeout
	for timer > 0 {
		_, err = client.Do(request)
		if err != nil {
			break
		}
		if timer == 0 {
			return errors.Errorf("did not shut down envoy succesfully in %d seconds", timeout)
		}
		time.Sleep(1 * time.Second)
		timer -= 1
	}
	return nil
}

func (ei *EnvoyInstance) runContainer() error {
	envoyImageTag := os.Getenv("ENVOY_IMAGE_TAG")
	if envoyImageTag == "" {
		panic("The ENVOY_IMAGE_TAG env var is not set. Find valid tag names here https://quay.io/repository/solo-io/gloo-ee-envoy-wrapper?tab=tags")
	}

	image := "quay.io/solo-io/gloo-ee-envoy-wrapper:" + envoyImageTag
	args := []string{"run", "-d", "--rm", "--name", containerName,
		"-p", "8080:8080",
		"-p", "8083:8083",
		"-p", "8443:8443",
		"-p", fmt.Sprintf("%v:%v", ei.AdminPort, ei.AdminPort),
		"--entrypoint=envoy",
		image,
		"--disable-hot-restart", "--log-level", "debug",
		"--bootstrap-version", "3",
		"--config-yaml", ei.envoycfg,
	}

	_, _ = fmt.Fprintln(ginkgo.GinkgoWriter, args)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = ginkgo.GinkgoWriter
	cmd.Stderr = ginkgo.GinkgoWriter
	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, "Unable to start envoy container")
	}
	return nil
}

func localAddr() (string, error) {
	ip := os.Getenv("GLOO_IP")
	if ip != "" {
		return ip, nil
	}
	// go over network interfaces
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		if (i.Flags&net.FlagUp == 0) ||
			(i.Flags&net.FlagLoopback != 0) ||
			(i.Flags&net.FlagPointToPoint != 0) {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() != nil {
					return v.IP.String(), nil
				}
			case *net.IPAddr:
				if v.IP.To4() != nil {
					return v.IP.String(), nil
				}
			}
		}
	}
	return "", errors.New("unable to find Gloo IP")
}

func (ei *EnvoyInstance) Logs() (string, error) {
	if ei.useDocker {
		logsArgs := []string{"logs", containerName}
		cmd := exec.Command("docker", logsArgs...)
		byt, err := cmd.CombinedOutput()
		if err != nil {
			return "", errors.Wrap(err, "Unable to fetch logs from envoy container")
		}
		return string(byt), nil
	}

	return ei.logs.String(), nil
}

func (ei *EnvoyInstance) GetConfigDump() (*http.Response, error) {
	return http.Get(fmt.Sprintf("http://localhost:%d/config_dump", ei.AdminPort))
}
