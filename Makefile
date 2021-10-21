plat ?= darwin
plats = linux darwin windows

arch ?= amd64
archs = amd64 arm arm64

all:  server wechat

define build_app
        @echo 'building $(1) ...'
        @GOOS=$(2) GOARCH=$(3) go build -o builder/$(1) ./cmd/$(1)
        @echo 'build $(1) done'
endef


server:
	$(call build_app,server,$(plat),$(arch))
.PHONY: server

wechat:
	$(call build_app,wechat,$(plat),$(arch))
.PHONY: wechat

clean:
	@rm -f builder/*