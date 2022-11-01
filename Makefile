# run make clean before run make release
release:
	./build/release.sh
clean:
	rm -rf tmp/
	rm -rf src/service/dist
.PHONY: release