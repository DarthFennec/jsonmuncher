.PHONY: build bench

build: benchmark/fixture_huge.json
	docker build -t jsonmuncher -f benchmark/Dockerfile .

bench:
	docker run --rm -v "$(pwd):/go/src/github.com/darthfennec/jsonmuncher" -it jsonmuncher go test -test.benchmem -bench . ./benchmark/ -benchtime 5s -v

benchmark/fixture_huge.json:
	cd benchmark && go run gen/fixture_huge_gen.go
