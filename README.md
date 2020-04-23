# Usage

## Install a new glooe distribution from scratch on Kubernetes

```bash
# Setup your repo
make init
make update-deps
make allprojects

# for a new UI: update the version in solo-projects/install/helm/gloo-ee/generate.go

# at this point you should have gloo built to you ./_output/ directory
# make the manifest
VERSION="1.10.0" make manifest # note that there is no "v" in the version, version pertains to the solo-projects version. Use "dev" or something if you want to use local images
eval $(minikube docker-env) # so minikube can use local images
make docker -B # creates all your images locally and tags them as "dev" by default

# install
# prep: create a secret with you docker credentials
./_output/glooctl install kube -f ./install/manifest/glooe-distribution.yaml
# NOTE: glooe-distribution.yaml is the same as glooe-release.yaml except that "distribution" uses an IfNotPresent pull policy
```

## Updated instructions for the grpcserver

### prep

- get the right version of protoc (3.6.1)
  - the make target below will warn you if you need to update

### build

```bash
make update-ui-deps
make generated-ui
make run-apiserver
```

## Building `extauth` components locally
We build the `extauth` binaries inside a [docker container](projects/extauth/cmd/Dockerfile) for reproducibility. 
Since it needs to access private git repositories, the container relies on a GitHub token to be provided via the 
`GITHUB_TOKEN` environment variable. If you need to build the `extauth` binaries locally, you have to generate a token. 
You can do that by opening the settings page for your GitHub account and navigating to 
`Developer Settings > Personal access tokens > Generate new token`. Once you have a token, you can export it to your 
environment and run the desired `make` target, e.g.:

``` 
export GITHUB_TOKEN=<your token> 
make extauth
```

## Noteworthy make targets

- `build-test-assets`: pushes all images and creates the zipped helm chart
  - requires `BUILD_ID` and `GCLOUD_PROJECT_ID` set
  - zipped helm chart saved in the `_test` dir
  - when running locally, should set `LOCAL_BUILD=1` in order to build the ui resources

## Additional Notes

- Shared projects across Solo.io.
- This repo contains the git history for Gloo and Solo-Kit.

## Helm Repositories
- [GlooE](https://console.cloud.google.com/storage/browser/gloo-ee-helm)
- [Gloo with read-only UI](https://console.cloud.google.com/storage/browser/gloo-os-ui-helm)
- [Dev portal](https://console.cloud.google.com/storage/browser/dev-portal-helm)
