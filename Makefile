.PHONY: all mod run run-sim clean

EXECUTABLES = bin/binagot bin/binagot-sim
LDFLAGS = -ldflags="-s -w"

all: $(EXECUTABLES)

mod:
	@go mod tidy
	@go mod vendor

bin/%: clean
	go build $(LDFLAGS) -o bin/$* github.com/nmmsantos/binagot/cmd/$*

run:
	@go run github.com/nmmsantos/binagot/cmd/binagot

run-sim:
	@go run github.com/nmmsantos/binagot/cmd/binagot-sim

clean:
	@rm -rf vendor bin
