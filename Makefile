export GOPATH=$(HOME)/gocode

all: registry-viewer registry-gc

registry-viewer: registry-viewer.go registry.go registryclient.go sh.go manifest.go imagestream.go ego.go
	go build $^

registry-gc: registry-gc.go registry.go registryclient.go sh.go manifest.go imagestream.go ego.go
	go build $^

$(GOPATH)/bin/ego:
	go get github.com/benbjohnson/ego/cmd/ego

ego.go: $(GOPATH)/bin/ego registry.ego
	$(GOPATH)/bin/ego --package=main .
