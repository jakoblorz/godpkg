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

if ! [ -f "packages" ] ; then
    touch "packages"
fi

if [ $# -eq 0 ] ; then

    export GOPATH="$(pwd)"
    export GOBIN="$(pwd)/bin"

    cat "packages" | while read in; do
        if [ -n "$in" ] ; then
            printf "${BLUE}[install]${NC} $in -> $(pwd)/\n"

            printf " - installing dependency $in\n"
            go get -v $in | while read line; do
                printf " - ${YELLOW}[message]${NC} $line\n"
            done
        fi
    done

    
    cat "packages" >> "packages.temp"
    cat "packages.temp" | sed '/^$/d' > "packages"
    rm "packages.temp"
    
    exit 0
fi

if [ $# -eq 2 ] ; then

    SCOPE=$1
    REPOURL=${@:2}

    OIFS="$IFS"
    IFS="/"
    read -a REPOARR <<< "${REPOURL}"
    IFS="$OIFS"


    if [ $SCOPE == "global" ] ; then

        export GOPATH="$(echo ~)/.go-env"
        export GOBIN="$(echo ~)/.go-env/bin"

        PKGFOLDS="$(find $(echo ~)/.go-env/pkg/* -maxdepth 0 -type d)"
        BINARY="${REPOURL##*/}"

        printf "${BLUE}[install${NC}${YELLOW}@${SCOPE}${NC}${BLUE}]${NC} $REPOURL -> $(echo ~)/.go-env\n"

        echo " - installing dependency $REPOURL"
        go get -v "$REPOURL"

        echo " - creating symlink $(echo ~)/.go-env/src/$REPOURL -> $(pwd)/src/$REPOURL"
        mkdir -p "$(pwd)/src/$REPOURL"
        ln -sf "$(echo ~)/.go-env/src/$REPOURL" "$(pwd)/src/$REPOURL"

        pkgdir="$(find $(echo ~)/.go-env/pkg/* -maxdepth 0 -type d)"
        for arch in $PKGFOLDS; do
            if [ -d "${arch}" ]; then

                PKG="${arch##*/}"

                echo " - creating symlink $(echo ~)/.go-env/pkg/${PKG}/${REPOARR[0]}/${REPOARR[1]} -> $(pwd)/pkg/${PKG}/${REPOARR[0]}/${REPOARR[1]}"
                mkdir -p "$(pwd)/pkg/${PKG}/${REPOARR[0]}/${REPOARR[1]}"
                ln -sf "$(echo ~)/.go-env/pkg/${PKG}/${REPOARR[0]}/${REPOARR[1]}" "$(pwd)/pkg/${PKG}/${REPOARR[0]}"
            fi
        done

        if [ -f "$(echo ~)/.go-env/bin/$BINARY" ] ; then
            echo " - creating symlink $(echo ~)/.go-env/bin/$BINARY -> $(pwd)/bin/$BINARY"
            ln -sf "$(echo ~)/.go-env/bin/$BINARY" "$(pwd)/bin/$BINARY"
        fi
        exit 0
    fi

    if [ $1 == "local" ] ; then

        printf "${BLUE}[install${NC}${YELLOW}@${SCOPE}${NC}${BLUE}]${NC} ${@:2} -> $(pwd)/\n"

        export GOPATH="$(pwd)"
        export GOBIN="$(pwd)/bin"

        echo " - installing dependency ${REPOURL}"
        go get -v "${REPOURL}"

        echo " - adding dependency to $(pwd)/packages"
        printf "\n${REPOURL}" >> "packages"
        cat "packages" >> "packages.temp"
        cat "packages.temp" | sed '/^$/d' > "packages"
        rm "packages.temp"

        exit 0
    fi
fi

printf "${RED}[error]${NC} missing argument: global or local?\n"
exit 1
