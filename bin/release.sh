#! bash

OS_TARGETS=( darwin linux )
ARCH_TARGETS=( amd64 arm64 )

for OS_TARGET in "${OS_TARGETS[@]}"; do
  for ARCH_TARGET in "${ARCH_TARGETS[@]}"; do
    GOOS=$OS_TARGET GOARCH=$ARCH_TARGET go build -o dist/nv-$OS_TARGET-$ARCH_TARGET main.go
  done
done
