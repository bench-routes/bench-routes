#!/bin/bash

cd src/
PLATFORMS="darwin/amd64" # amd64 only as of go1.5
PLATFORMS="$PLATFORMS linux/amd64 linux/386"
PLATFORMS="$PLATFORMS freebsd/amd64"
PLATFORMS="$PLATFORMS netbsd/amd64" # amd64 only as of go1.6
PLATFORMS="$PLATFORMS openbsd/amd64" # amd64 only as of go1.6

# unsupported platforms for now
_PLATFORMS="$PLATFORMS windows/amd64 windows/386" # arm compilation not available for Windows
_PLATFORMS="$PLATFORMS linux/ppc64 linux/ppc64le"
_PLATFORMS="$PLATFORMS linux/mips64 linux/mips64le" # experimental in go1.6
_PLATFORMS="$PLATFORMS dragonfly/amd64" # amd64 only as of go1.5
_PLATFORMS="$PLATFORMS plan9/amd64 plan9/386" # as of go1.4
_PLATFORMS="$PLATFORMS solaris/amd64" # as of go1.3

# ARMBUILDS lists the platforms that are currently supported.  From this list
# we generate the following architectures:
#
#   ARM64 (aka ARMv8) <- only supported on linux and darwin builds (go1.6)
#   ARMv7
#   ARMv6
#   ARMv5
#
PLATFORMS_ARM="linux freebsd netbsd"

# Shouldn't really need to modify anything below this line.

type setopt >/dev/null 2>&1

SCRIPT_NAME=`basename "$0"`
FAILURES=""
SOURCE_FILE=`echo $@ | sed 's/\.go//'`
CURRENT_DIRECTORY=${PWD##*/}
OUTPUT=${SOURCE_FILE:-$CURRENT_DIRECTORY} # if no src file given, use current dir name

for PLATFORM in $PLATFORMS; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  BIN_FILENAME="${OUTPUT}-${GOOS}-${GOARCH}"
  if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ../tmp_build/${BIN_FILENAME} $@"
  echo "${CMD}"
  eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
done

# ARM builds
if [[ $PLATFORMS_ARM == *"linux"* ]]; then 
  CMD="GOOS=linux GOARCH=arm64 go build -o ../tmp_build/${OUTPUT}-linux-arm64 $@"
  echo "${CMD}"
  eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
fi
for GOOS in $PLATFORMS_ARM; do
  GOARCH="arm"
  # build for each ARM version
  for GOARM in 7 6 5; do
    BIN_FILENAME="${OUTPUT}-${GOOS}-${GOARCH}${GOARM}"
    CMD="GOARM=${GOARM} GOOS=${GOOS} GOARCH=${GOARCH} go build -o ../tmp_build/${BIN_FILENAME} $@"
    echo "${CMD}"
    eval "${CMD}" || FAILURES="${FAILURES} ${GOOS}/${GOARCH}${GOARM}" 
  done
done

# eval errors
if [[ "${FAILURES}" != "" ]]; then
  echo ""
  echo "${SCRIPT_NAME} failed on: ${FAILURES}"
  exit 1
fi
