
# Image URL to use all building/pushing image targets
IMG_MANAGER ?= rykren/openelb:v0.5.0
#IMG_AGENT ?= rykren/openelb-agent:v0.5.0
IMG_PROXY ?= rykren/openelb-proxy:v0.5.0
IMG_FORWARD ?= rykren/openelb-forward:v0.5.0
BRANCH ?= release

CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...   ./api/... ./pkg/controllers/...

# Run go vet against code
vet:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go vet ./pkg/... ./cmd/...  ./pkg/controllers/...

# Run tests
test: fmt vet
	KUBEBUILDER_ASSETS=$(shell $(GOBIN)/setup-envtest use -p path 1.19.x) go test -v  ./api/... ./pkg/controllers/... ./pkg/...  -coverprofile cover.out

# Build manager binary
manager: fmt vet
	#CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o bin/manager github.com/openelb/openelb/cmd/manager
	CGO_ENABLED=0 go build  -o bin/manager github.com/openelb/openelb/cmd/manager


deploy: generate
ifeq ($(uname), Darwin)
	sed -i '' -e 's@image: .*@image: '"${IMG_AGENT}"'@' ./config/${BRANCH}/agent_image_patch.yaml
	sed -i '' -e 's@image: .*@image: '"${IMG_MANAGER}"'@' ./config/${BRANCH}/manager_image_patch.yaml
	sed -i '' -e 's@NodeProxyDefaultForwardImage      string = \".*\"@NodeProxyDefaultForwardImage      string = \"'"${IMG_FORWARD}"'\"@' ./pkg/constant/constants.go
	sed -i '' -e 's@NodeProxyDefaultProxyImage        string = \".*\"@NodeProxyDefaultProxyImage        string = \"'"${IMG_PROXY}"'\"@' ./pkg/constant/constants.go
else
	sed -i -e 's@image: .*@image: '"${IMG_AGENT}"'@' ./config/${BRANCH}/agent_image_patch.yaml
	sed -i -e 's@image: .*@image: '"${IMG_MANAGER}"'@' ./config/${BRANCH}/manager_image_patch.yaml
	sed -i -e 's@NodeProxyDefaultForwardImage      string = \".*\"@NodeProxyDefaultForwardImage      string = \"'"${IMG_FORWARD}"'\"@' ./pkg/constant/constants.go
	sed -i -e 's@NodeProxyDefaultProxyImage        string = \".*\"@NodeProxyDefaultProxyImage        string = \"'"${IMG_PROXY}"'\"@' ./pkg/constant/constants.go
endif
	kustomize build config/${BRANCH} -o deploy/openelb.yaml
	@echo "Done, the yaml is in deploy folder named 'openelb.yaml'"

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths=./api/...
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=openelb-manager-role webhook paths="./api/..." paths="./pkg/controllers/..." output:crd:artifacts:config=config/crd/bases

controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.4.0
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

clean-up:
	./hack/cleanup.sh

release: deploy
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/manager-linux-amd64 github.com/openelb/openelb/cmd/manager
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/agent-linux-amd64 github.com/openelb/openelb/cmd/agent
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o bin/gobgp-linux-amd64 github.com/osrg/gobgp/cmd/gobgp
#	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/manager-linux-arm64 github.com/openelb/openelb/cmd/manager
#	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/agent-linux-arm64 github.com/openelb/openelb/cmd/agent
#	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build  -o bin/gobgp-linux-arm64 github.com/osrg/gobgp/cmd/gobgp
#	DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/amd64,linux/arm64 -t ${IMG_AGENT} -f ./cmd/agent/Dockerfile .  --push
#	DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/amd64,linux/arm64 -t ${IMG_MANAGER} -f ./cmd/manager/Dockerfile .  --push
#	DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/amd64,linux/arm64 -t ${IMG_PROXY} -f ./images/proxy/Dockerfile . --push
#	DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/amd64,linux/arm64 -t ${IMG_FORWARD} -f ./images/forward/Dockerfile . --push
	DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/amd64 -t ${IMG_MANAGER} -f ./cmd/manager/Dockerfile .  --push
	DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/amd64 -t ${IMG_PROXY} -f ./images/proxy/Dockerfile . --push
	DOCKER_CLI_EXPERIMENTAL=enabled docker buildx build --platform linux/amd64 -t ${IMG_FORWARD} -f ./images/forward/Dockerfile . --push

install-travis:
	echo "install kubebuilder/kustomize etc."
	chmod +x ./hack/*.sh
	./hack/install_tools.sh
