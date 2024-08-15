include Makefile.env

## run: Compile and run server
run: go-compile start

## start: Start in development mode. Auto-starts when code changes.
start: start-etcd-proxy

## stop: Stop development mode. GO_ETCD_PROXY
stop: stop-etcd-proxy

start-etcd-proxy: stop-etcd-proxy
	@echo "  >  $(PROJECTNAME) is available at $(HTTP_ADDRESS) and gRPC at $(GRPC_ADDRESS)"
	@-cd ./$(DIR_ETCD_PROXY) && (./etcd-proxy -h $(HTTP_ADDRESS) -g $(GRPC_ADDRESS) && echo $$! > $(PID_GO_ETCD_PROXY))
	@cat $(PID_GO_ETCD_PROXY) | sed "/^/s/^/  \>  PID: /"

stop-etcd-proxy:
	@echo "  >  stop by $(PID_GO_ETCD_PROXY)"
	@-touch $(PID_GO_ETCD_PROXY)
	@-kill `cat $(PID_GO_ETCD_PROXY)` 2> /dev/null || true
	@-rm $(PID_GO_ETCD_PROXY)

restart-etcd-proxy: stop-etcd-proxy start-etcd-proxy

## build: Build and the binary compile server
build: go-build-etcd-proxy

## clean: Clean build files. Runs `go clean` internally.
clean:
	@(MAKEFILE) go-clean

go-compile: go-build-etcd-proxy

go-build-etcd-proxy:
	@echo "  >  Building GO_ETCD_PROXY binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) cd ./$(DIR_ETCD_PROXY) && go build -o ./etcd-proxy $(GOFILES)

go-generate:
	@echo "  >  Generating dependency files..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get $(get)

.PHONY: go-update-deps
go-update-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif

go-install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-swag:
	swag init -g $(MAIN_GO) --output $(DOCS_DIR)
	sed -i 's/"localhost:8080",/env.GetConfig().Address(),/' $(DOCS_GO)
	goimports -w $(DOCS_GO)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

cert:
	@cd cert; openssl req -x509 -newkey rsa:1024 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Tech School/OU=Education/CN=localhost/emailAddress=none@gmail.com"
	@echo "CA's self-signed certificate"
	@cd cert; openssl x509 -in ca-cert.pem -noout -text
	@cd cert; openssl req -newkey rsa:1024 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=Tech School/OU=Education/CN=localhost/emailAddress=none@gmail.com"
	@cd cert; openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf
	@echo "Server's signed certificate"
	@cd cert; openssl x509 -in server-cert.pem -noout -text

##################
# Implicit targets
##################

# This rulle is used to generate the message source files based
# on the *.proto files.
%.pb.go: %.proto
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./$<

####################################
# Major source code-generate targets
####################################
generate: $(PROTO_PB_GO)
	@echo "  >  Done generating source files based on *.proto and Mock files."

test:
	@echo "  > Test Iteration ..."
	go vet -vettool=$(which statictest) ./...
	cd cmd/etcd-proxy

.PHONY: cert help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
