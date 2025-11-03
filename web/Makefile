
all: build

build:
	npm run build
	make gzip

gzip:
	for f in static/js/gen/*.js; do \
		echo Compressing $$f ; \
		gzip -9 -k -f $$f ; \
	done

test:
	npm run test:e2e
	# npm run test
