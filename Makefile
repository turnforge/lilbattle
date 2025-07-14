
buildweb:
	cd web ; npm run build

binlocal: 
	go build -ldflags "$(LDFLAGS)" -o /tmp/weewar ./main.go

buf:
	buf generate
