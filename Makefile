PLUGIN_ID=com.pointplus.reactification
PLUGIN_VERSION=1.0.0
BUNDLE_NAME=$(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz

.PHONY: all build dist clean

all: dist

build:
	cd server && GOOS=linux GOARCH=amd64 go build -o dist/plugin-linux-amd64 .

dist: build
	rm -rf dist/$(PLUGIN_ID)
	mkdir -p dist/$(PLUGIN_ID)/server/dist
	cp plugin.json dist/$(PLUGIN_ID)/
	cp server/dist/plugin-linux-amd64 dist/$(PLUGIN_ID)/server/dist/
	cd dist && tar -czf $(BUNDLE_NAME) $(PLUGIN_ID)
	rm -rf dist/$(PLUGIN_ID)

clean:
	rm -rf server/dist dist $(PLUGIN_ID).tar.gz
