GIT_COMMIT := $(shell git describe --tags | tr -d v)
GIT_DIRTY := $(if $(shell git status --porcelain | grep -qE .),~$(shell git rev-parse --short HEAD))
VERSION := $(GIT_COMMIT)$(GIT_DIRTY)
ARCH ?= amd64
PACKAGE := drx

.PHONY: all
all: clean package

test:
	go vet ./...
	go test -race ./...

$(PACKAGE):
	CGO_ENABLED=1 go build cmd/$(PACKAGE)/$(PACKAGE).go

web/dist:
	cd web && npm install && npm run build

.PHONY: clean
clean:
	rm -rf web/dist $(PACKAGE) build

.PHONY: package
package: $(PACKAGE) web/dist
	mkdir -p build build/opt/pi6atv-$(PACKAGE)/web build/etc/nginx/sites-enabled \
			 build/etc/nginx/snippets build/var/lib/grafana/dashboards build/etc/prometheus/targets
	cp $(PACKAGE) build/opt/pi6atv-$(PACKAGE)
	cp nginx-site.conf build/etc/nginx/sites-enabled/$(PACKAGE).conf
	cp nginx-proxy.conf build/etc/nginx/snippets/$(PACKAGE)-proxy.conf
	cp systemd.service build/$(PACKAGE).service
	cp config/$(PACKAGE).yaml build/opt/pi6atv-$(PACKAGE)/$(PACKAGE).yaml
	cp -a web/dist  build/opt/pi6atv-$(PACKAGE)/web/$(PACKAGE)
	cp grafana-dashboard.json build/var/lib/grafana/dashboards/$(PACKAGE)-dashboard.json
	cp prometheus.yaml build/etc/prometheus/targets/$(PACKAGE).yaml
	cd build && \
		fpm -s dir -t deb -n pi6atv-$(PACKAGE) -v "$(VERSION)" \
			-d nginx \
			--deb-systemd $(PACKAGE).service \
			--deb-systemd-enable --deb-systemd-auto-start --deb-systemd-restart-after-upgrade \
			-a $(ARCH) -m "Wim Fournier <debian@fournier.nl>" \
			.
