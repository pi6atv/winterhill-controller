#!/bin/bash
docker build -t builder -f Dockerfile.build .
docker run -ti --rm -e ARCH -e GOOS -e GOARCH -e GOARM -e CC -v "$PWD:/build" -w /build builder make $@