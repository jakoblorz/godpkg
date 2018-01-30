package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/cli"
)

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
	fmt.Printf("%s\n", file.path)
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

	path, perr := CreateFolderStructure([]string{pathAppend("/bin"), pathAppend("/src"), pathAppend("/scripts")})
	if perr != nil {
		log.Fatalf("Error creating Directory %s: %s\n", path, perr)
		return 1
	}

	packages := &TemplateFile{
		path:    pathAppend("/packages"),
		content: "-v github.com/jakoblorz/godpkg",
	}

	file, ferr := CreateFileStructure([]*TemplateFile{packages})
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
	return "install dependency"
}

// Run installs the go package
func (*InstallCommand) Run(args []string) int {
	fmt.Printf("install, %v\n", args)
	return 0
}

// Synopsis returns the Help Text of the InstallCommand
func (command *InstallCommand) Synopsis() string {
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

	c.Commands = map[string]cli.CommandFactory{
		"install": install,
		"init":    init,
	}

	status, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(status)
}
