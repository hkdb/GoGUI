### PROJECT: GoGUI
###
### MAINTAINED BY: hkdb <hkdb@3df.io>
###
### SPONSORED BY: 3DF OSI - https://osi.3df.io
###
### DISCLAIMER:
### This application is maintained by volunteers and in no way
### do the maintainers make any guarantees. Use at your own risk.

ASSETS_DIR = "assets"
build:
	@export GOPATH=$${GOPATH-~/go} && \
	go get github.com/jteeuwen/go-bindata... && \
	$$GOPATH/bin/go-bindata -o bindata.go -tags builtinassets ${ASSETS_DIR}/... && \
	go build -tags builtinassets -ldflags "-X main.builtinAssets=${ASSETS_DIR}"