package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/mitchellh/cli"
)

var install = `#!/bin/bash

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

function install_global ()
{
    export GOPATH="$(echo ~)/.go-env"
    export GOBIN="$(echo ~)/.go-env/bin"

    PKGFOLDS="$(find $(echo ~)/.go-env/pkg/* -maxdepth 0 -type d)"
    BINARY="${REPOURL##*/}"

    printf "${BLUE}[install${NC}${YELLOW}@${SCOPE}${NC}${BLUE}]${NC} $REPOURL -> $(echo ~)/.go-env\n"

    echo " - installing dependency $REPOURL"
    go get -v "$REPOURL"

    echo " - creating symlink $(echo ~)/.go-env/src/$REPOURL -> $(pwd)/src/$REPOURL"
    mkdir -p "$(pwd)/src/$REPOURL"
    cp -an "$(echo ~)/.go-env/src/$REPOURL" "$(pwd)/src/$REPOURL"

    pkgdir="$(find $(echo ~)/.go-env/pkg/* -maxdepth 0 -type d)"
    for arch in $PKGFOLDS; do
        if [ -d "${arch}" ]; then

            PKG="${arch##*/}"
            PKGHOST=${REPOARR[0]}
            PKGUSER=${REPOARR[1]}

            echo " - creating symlink $(echo ~)/.go-env/pkg/${PKG}/${PKGHOST}/${PKGUSER} -> $(pwd)/pkg/${PKG}/${PKGHOST}/${PKGUSER}"
            mkdir -p "$(pwd)/pkg/${PKG}/${PKGHOST}/${PKGUSER}"
            cp -an "$(echo ~)/.go-env/pkg/${PKG}/${PKGHOST}/${PKGUSER}" "$(pwd)/pkg/${PKG}/${PKGHOST}"
        fi
    done

    if [ -f "$(echo ~)/.go-env/bin/$BINARY" ] ; then
        echo " - creating symlink $(echo ~)/.go-env/bin/$BINARY -> $(pwd)/bin/$BINARY"
        mkdir -p "$(pwd)/bin"
        cp -an "$(echo ~)/.go-env/bin/$BINARY" "$(pwd)/bin/$BINARY"
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
    export GOPATH="$(pwd)"
    export GOBIN="$(pwd)/bin"

    printf "${BLUE}[install${NC}${YELLOW}@${SCOPE}${NC}${BLUE}]${NC} ${REPOURL} -> $(pwd)/\n"

    echo " - installing dependency ${REPOURL}"
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


`

var build = `#!/bin/bash

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


export GOPATH="$(pwd)"
export GOBIN="$(pwd)/bin"

`

var src = `
package main

import "fmt"

func main() {
	fmt.Printf("Project set up correctly\n")
}
`

// CreateFolderStructure creates folders returning errors with the
// path where the error occured; returning "", nil is the desired
// case
func CreateFolderStructure(folders []string) (string, error) {
	for _, path := range folders {
		fmt.Printf("%s\n", path)

		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return path, err
		}
	}

	return "", nil
}

// CreateFileStructure creates files returning errors with the
// file where the error occured; returning nil, nil is the desired
// case
func CreateFileStructure(files []*TemplateFile) (*TemplateFile, error) {
	for _, file := range files {
		fmt.Printf("%s\n", file.String())

		err := file.WriteToDisk()
		if err != nil {
			return file, err
		}
	}

	return nil, nil
}

// TemplateFile represents a file to be written at the path
// with the content to the disk
type TemplateFile struct {
	path    string
	content string
}

// WriteToDisk writes the TemplateFile's contents into a file
// at the TemplateFile's path
func (file *TemplateFile) WriteToDisk() error {
	return ioutil.WriteFile(file.path, []byte(file.content), 0644)
}

// Returns a string representation of the TemplateFile
// which is just the path
func (file *TemplateFile) String() string {
	return file.path
}

// InitCommand represents the data structure for the
// init command
type InitCommand struct {
}

// Help returns the Help Text for the InitCommand
func (*InitCommand) Help() string {
	return "initialize new godpkg project structure"
}

// Run creates the Folder Structure and packages/scripts
func (*InitCommand) Run(args []string) int {

	pathAppend := func(p string) string {
		return "./" + args[0] + p
	}

	path, perr := CreateFolderStructure([]string{pathAppend("/bin"), pathAppend("/src"), pathAppend("/scripts"), pathAppend("/pkg")})
	if perr != nil {
		log.Fatalf("Error creating Directory %s: %s\n", path, perr)
		return 1
	}

	packages := &TemplateFile{
		path:    pathAppend("/packages"),
		content: "\n",
	}

	install := &TemplateFile{
		path:    pathAppend("/scripts/install.sh"),
		content: install,
	}

	build := &TemplateFile{
		path:    pathAppend("/scripts/build.sh"),
		content: build + "\ngo install \"$(pwd)/src/" + args[0] + ".go\"",
	}

	init := &TemplateFile{
		path:    pathAppend("/src/" + args[0] + ".go"),
		content: src,
	}

	file, ferr := CreateFileStructure([]*TemplateFile{packages, install, build, init})
	if ferr != nil {
		log.Fatalf("Error creating File %s: %s\n", file.String(), ferr)
		return 1
	}

	return 0
}

// Synopsis returns the Help Text of the InitCommand
func (command *InitCommand) Synopsis() string {
	return command.Help()
}

// InstallCommand represents the data structure for
// the install command
type InstallCommand struct {
}

// Help returns the Help Text for the InstallCommand
func (*InstallCommand) Help() string {
	return "installs a dependency; download the dependency into the project (local) or use ~./go-env as cache and symlink (global)"
}

// Run installs the go package
func (command *InstallCommand) Run(args []string) int {
	arguments := append([]string{"./scripts/install.sh"}, args...)

	cmd := exec.Command("/bin/bash", arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error during Install: %s\n", err)
		return 1
	}

	return 0
}

// Synopsis returns the Help Text of the InstallCommand
func (command *InstallCommand) Synopsis() string {
	return command.Help()
}

// BuildCommand represents the data structure for
// the build command
type BuildCommand struct {
}

// Help returns the Help text for the BuildCommand
func (*BuildCommand) Help() string {
	return "builds the project"
}

// Run builds the project
func (*BuildCommand) Run(args []string) int {

	cmd := exec.Command("/bin/sh", "./scripts/build.sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error during Build: %s\n", err)
		return 1
	}

	return 0
}

// Synopsis returns the Help Text of the BuildCommand
func (command *BuildCommand) Synopsis() string {
	return command.Help()
}

func main() {
	c := cli.NewCLI("godpkg", "1.0.0")
	c.Args = os.Args[1:]

	install := func() (cli.Command, error) {
		return &InstallCommand{}, nil
	}

	init := func() (cli.Command, error) {
		return &InitCommand{}, nil
	}

	build := func() (cli.Command, error) {
		return &BuildCommand{}, nil
	}

	c.Commands = map[string]cli.CommandFactory{
		"install": install,
		"init":    init,
		"build":   build,
	}

	status, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(status)
}
