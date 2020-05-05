#!/bin/bash
set -e

if [ -z $PREFIX ]
then
    # PREFIX not set
    ROOT="$PWD"
    PREFIX=$ROOT/../"wallet-core/build/local"
    if  [ ! -d $PREFIX ] || \
        [ ! -d $PREFIX/include ] || \
        [ ! -f $PREFIX/bin/protoc ] || \
        [ ! -f $PREFIX/bin/protoc-gen-c-typedef ]
    then
        echo $PREFIX does not exist or not complete, fallback to /usr/local
        PREFIX=/usr/local
    fi
fi
echo "PREFIX: $PREFIX"
export PATH="$PREFIX/bin":$PATH
# library paths, for proto plugins
export LD_LIBRARY_PATH="$PREFIX/lib":$LD_LIBRARY_PATH
export DYLD_LIBRARY_PATH="$PREFIX/lib":$LD_LIBRARY_PATH
# protoc executable (proto compiler)
PROTOC="$PREFIX/bin/protoc"
which $PROTOC
$PROTOC --version

rm -rf Ethereum
mkdir Ethereum
rm -rf Bitcoin
mkdir Bitcoin


"$PROTOC" -I=../wallet-core/src/proto --plugin=$GOPATH/bin/protoc-gen-go --go_out=./Ethereum ../wallet-core/src/proto/Ethereum.proto
"$PROTOC" -I=../wallet-core/src/proto --plugin=$GOPATH/bin/protoc-gen-go --go_out=./Bitcoin ../wallet-core/src/proto/Bitcoin.proto


