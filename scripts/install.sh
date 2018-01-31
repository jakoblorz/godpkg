#!/bin/bash

# MIT License

# Copyright (c) 2018 Jakob Lorz

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

function set_REPOARR () 
{
    OIFS="$IFS"
    IFS="/"
    read -a REPOARR <<< "${REPOURL}"
    IFS="$OIFS"
}

function set_SCOPE () 
{
    SCOPE=$1
}

function set_REPOURL () 
{
    REPOURL=$2
}

function set_HEAD ()
{
    local HEADSEARCH=($(git ls-remote https://$REPOURL | grep "HEAD$"))
    HEAD=${HEADSEARCH[0]}
}

function set_GOPATH ()
{
    export GOPATH=$1
    export GOBIN=$2
}

function set_ARCHITECTURES ()
{
    ARCHITECTURES="$(find ${1}/pkg/* -maxdepth 0 -type d)"
}

function set_BINARY ()
{
    BINARY="${REPOURL##*/}"
}

function install_global ()
{

    set_GOPATH "$(echo ~)/.go-env" "$(echo ~)/.go-env/bin"
    set_ARCHITECTURES "$(echo ~)/.go-env"
    set_BINARY

    printf "${BLUE}[install${NC}${YELLOW}@${SCOPE}${NC}${BLUE}]${NC} $REPOURL -> $(echo ~)/.go-env\n"

    echo " - installing dependency $REPOURL[$HEAD]"
    go get -v "$REPOURL"

    echo " - creating symlink $(echo ~)/.go-env/src/$REPOURL -> $(pwd)/src/$REPOURL"
    mkdir -p "$(pwd)/src/$REPOURL"
    cp -ans "$(echo ~)/.go-env/src/$REPOURL" "$(pwd)/src/$REPOURL"

    pkgdir="$(find $(echo ~)/.go-env/pkg/* -maxdepth 0 -type d)"
    for ARCH in $ARCHITECTURES; do
        if [ -d "${ARCH}" ]; then

            PKG="${ARCH##*/}"
            PKGHOST=${REPOARR[0]}
            PKGUSER=${REPOARR[1]}

            echo " - creating symlink $(echo ~)/.go-env/pkg/${PKG}/${PKGHOST}/${PKGUSER} -> $(pwd)/pkg/${PKG}/${PKGHOST}/${PKGUSER}"
            mkdir -p "$(pwd)/pkg/${PKG}/${PKGHOST}/${PKGUSER}"
            cp -ans "$(echo ~)/.go-env/pkg/${PKG}/${PKGHOST}/${PKGUSER}" "$(pwd)/pkg/${PKG}/${PKGHOST}"
        fi
    done

    if [ -f "$(echo ~)/.go-env/bin/$BINARY" ] ; then
        echo " - creating symlink $(echo ~)/.go-env/bin/$BINARY -> $(pwd)/bin/$BINARY"
        mkdir -p "$(pwd)/bin"
        cp -ans "$(echo ~)/.go-env/bin/$BINARY" "$(pwd)/bin/$BINARY"
    fi

    if [ $3 ] ; then

        echo " - adding dependency to $(pwd)/packages"
        printf "\n${SCOPE} ${REPOURL}" >> "packages"
        cat "packages" >> "packages.temp"
        cat "packages.temp" | sed '/^$/d' > "packages"
        rm "packages.temp"
    fi
}

function install_local () 
{
    set_GOPATH "$(pwd)" "$(pwd)/bin"

    printf "${BLUE}[install${NC}${YELLOW}@${SCOPE}${NC}${BLUE}]${NC} ${REPOURL} -> $(pwd)/\n"

    echo " - installing dependency $REPOURL[$HEAD]"
    go get -v "${REPOURL}"

    if [ $3 ] ; then

        echo " - adding dependency to $(pwd)/packages"
        printf "\n${SCOPE} ${REPOURL}" >> "packages"
        cat "packages" >> "packages.temp"
        cat "packages.temp" | sed '/^$/d' > "packages"
        rm "packages.temp"
    fi

}

if ! [ -f "packages" ] ; then
    touch "packages"
fi

if [ $# -eq 0 ] ; then

    export GOPATH="$(pwd)"
    export GOBIN="$(pwd)/bin"

    cat "packages" | while read in; do
        if [ -n "$in" ] ; then

            line=($in)

            set_SCOPE ${line[0]} ${line[1]}
            set_REPOURL ${line[0]} ${line[1]}
            set_REPOARR ${line[0]} ${line[1]}

            if [ $SCOPE == "global" ] ; then

                install_global $1 $2 false
            fi

            if [ $SCOPE == "local" ] ; then

                install_local $1 $2 false
            fi
        fi
    done

    
    cat "packages" >> "packages.temp"
    cat "packages.temp" | sed '/^$/d' > "packages"
    rm "packages.temp"
    
    exit 0
fi

if [ $# -eq 2 ] ; then

    set_SCOPE $1 $2
    set_REPOURL $1 $2
    set_REPOARR $1 $2
    set_HEAD $1 $2

    if [ $SCOPE == "global" ] ; then

        install_global $1 $2 true
        exit 0
    fi

    if [ $1 == "local" ] ; then

        install_local $1 $2 true
        exit 0
    fi
fi

printf "${RED}[error]${NC} missing argument: global or local?\n"
exit 1
