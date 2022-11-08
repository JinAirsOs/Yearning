# run make clean before run make release
release:
	./build/release.sh
clean:
	rm -rf tmp/
	rm -rf src/service/dist
test:
	go test ./...
.PHONY: release