# !/bin/sh
# This is a simple shell based bot that simplifies
# the release creation through a single command.
# Currently, this is a basic version, and hence a shell script,
# however, we plan to later make this a github robot, that can
# trigger the release creation by a github comment or issue.
#
# bench-release v0.1
#
# Currently supported versions:
# 1. Fedora
# 2. Debian family (ubuntu and its flavours, deepin, zorin, etc)
# 3. MacOS
# 4. Arch (needs confirmation on working)
# 5. BSD  (needs confirmation on working)
#

BOT_VERSION="0.1"

# check installation of necessity
MAKE_="$(command -v make)"
GO_="$(command -v go)"

echo "${MAKE_}"
echo "${GO_}"

if test -z "$MAKE_"
then
    echo "cannot make release due to unavalibility of makefile execution " \
    "cannot execute unit-tests before release" \
    "use --ignore to force build"
fi

if test -z "$GO_"
then
    echo "fatal: make release due to unavalibility of go installation binaries. " \
    "if installed, check the path."
    exit 1
fi

echo "checks done. moving ahead with the build process ...."

# make build directory in the root
mkdir -p $(pwd)/tmp_build
$(pwd)/shell/go-build-all.sh

mkdir -p $(pwd)/build

echo "binding files ...."
BUILD_PATH=$(pwd)/tmp_build

# get all build files
LIST_ALL_BUILD_FILES=$(pwd)/tmp_build/*
# echo $LIST_ALL_BUILD_FILES


for f in $(pwd)/tmp_build/*
do
    echo "packaging ${f}"
    f=$(basename $f)
    mkdir -p $(pwd)/build/$f
    cp $(pwd)/tmp_build/$f $(pwd)/build/$f/$f.sh
    cp $(pwd)/LICENSE $(pwd)/build/$f/LICENSE
    cp $(pwd)/CONTRIBUTING.md $(pwd)/build/$f/CONTRIBUTING.md
    cp $(pwd)/local-config.yml $(pwd)/build/$f/local-config.yml

    # make a package
    echo "compressing ${f}"
    tar -czvf ./build/$f.tar.gz ./build/$f --verbose

    # clean
    echo "cleaning ${f}"
    rm -R $(pwd)/build/$f
done

echo "cleaning remaining files"

rm -R $BUILD_PATH

echo "builds success :)"
