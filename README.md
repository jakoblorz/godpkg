# godpkg
Ever wondered why setting up Go Environments is so complicated? Then this Package Manager is for you! godpkg (pronounced *god*-pkg or *go*-dpkg? :stuck_out_tongue_winking_eye:) keeps track of the dependencies you install using a plain file (*./packages*) which allows you the reinstall dependecies later. Furthermore, the scripts to **build** and **install** set their own GOPATH thus keep all dependencies in your folder structure (just like in *node_modules* when using *npm*).

With this package manager I want to lower the entry barrier for newcomers who might be overwhelmed by the complicated and centralized setup process for Go Environments, but please (Quoting @nathany) [Recognize that **Go is different**](https://nathany.com/go-packages/). I advise you to switch to more common managers embraced by the Go Community later.
People should be free to choose whatever package manager they might like which is why I will list popular go package managers:
  - [dep](https://github.com/golang/dep)
  - [godep](https://github.com/tools/godep)
  
 If you want to add more, feel free to PR and if you like this project, :star: it!

## Install
```bash
curl -sL http://github.com/jakoblorz/godpkg/raw/master/bin/godpkg > /usr/local/bin/godpkg && chmod +x /usr/local/bin/godpkg
```

## Usage
### Initialize Project Folder
```bash
godpkg init <name>
```
You can then `cd` into the created folder and execute all the other command with the shell's `cwd` pointing to the root of the project.

### (Re-) Install Dependencies
To install specific dependencies, you can use the install command. Further arguments (e.g. `github.com/user/...`) will be piped over to `go get` which is used under the hood.
```bash
godpkg install github.com/user/...
```
If you omit any further arguments, the `packages` file will be read and each command listed there will be piped over to `go get`.

### Build your Project
If you call `build`, the file in `src/<foldername>.go` will be built.
```bash
godpkg build
```
