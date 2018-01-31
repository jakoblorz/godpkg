# godpkg
Go Environments use a opinionated project structure, which is why I scripted my Environment to manually switch between GOPATHs and symlink already existing versions of a dependency. `godpkg` is a cli-app which allows you to use the same scripts without having to manually copy them. Decide if you want to install a dependency in *local* scope (downloads / `go get`) right into your project folder or install it globally (*global* scope, downloads it into `~/.go-env`, symlinks into your project). Get started with the `godpkg init` command.

## Status
At this point this project is just a compilation of scripts I personally used into a small cli-app. There are more features planned but if you want to get a proper package manager people use in production environments today, please switch to more popular / established package managers:
  - [dep](https://github.com/golang/dep)
  - [gb](https://github.com/constabulary/gb) - Thank you, [MoneyWorthington](https://www.reddit.com/r/programming/comments/7u4eyz/gopath_independent_go_package_manager_wip/dthnbz9/)
  - [glide](https://github.com/Masterminds/glide)
  
### Planned Features
 - [ ] Snapshoting Versioned Dependencies when installing in Global Scope
 - [ ] Windows Support
 - [ ] Release Command to copy dependencies into `/vendor`

## Install
```bash
curl -sL http://github.com/jakoblorz/godpkg/raw/master/bin/godpkg > /usr/local/bin/godpkg && chmod +x /usr/local/bin/godpkg
```

## Usage
### Initialize Project Folder
```bash
godpkg init <name>
```
You can then `cd` into the folder and execute all the other commands.

### (Re-) Install Dependencies
To install specific dependencies, you can use the install command. `go get ...` is used under the hood. If you want to install the dependency directly (into the current project structure), choose **local**. Otherwise (**global**) the dependencies will be installed into `~/.go-env` and symlinked into the project folder (linking *bin, pkg and src*). Similar to a "normal" Go Environment with fixed GOPATH, dependencies can be shared between different projects which reduces the footprint significantly.
```bash
godpkg install <local|global> github.com/user/...
```
When omitting all arguments (`godpkg install`), godpkg will read the `./packages` file and install the listed dependencies (even with the original local/global scope).

### Build your Project
If you call `build`, the file in `src/<foldername>.go` will be built.
```bash
godpkg build
```
