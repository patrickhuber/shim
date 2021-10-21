# Shim

Shim is a standalone binary used for shimming executables instead of relying on symlinks. 

Shim is based on the chocolatey shim https://docs.chocolatey.org/en-us/features/shim

# Running

Shim consists of two binaries, shim and patch. 

Shim will try to read the exe path information in one of two ways:

## Config File

Read a &lt;name&gt;.yml or &lt;name&gt;.json file in the same directory with the same name as the shim. 

```yaml
name: go
```

```json
{ "name": "go" }
```

## Patch

Read a patch appended to it with the same schema as above. 

In order to apply a patch, you need to write the patch to a yml file (json not supported for patching) and use the patch binary to generate a patched copy of shim. 

> go.yml

```yaml
name: go
```

```bash
./cmd/patch/patch -s ./cmd/shim/shim -t ./cmd/shim/test -d go.yml
./cmd/shim/test
```

> output

```    
Go is a tool for managing Go source code.

Usage:

        go <command> [arguments]

The commands are:

        bug         start a bug report
        build       compile packages and dependencies
        clean       remove object files and cached files
        doc         show documentation for package or symbol
        env         print Go environment information
        fix         update packages to use new APIs
        fmt         gofmt (reformat) package sources
        generate    generate Go files by processing source
        get         add dependencies to current module and install them
        install     compile and install packages and dependencies
        list        list packages or modules
        mod         module maintenance
        run         compile and run Go program
        test        test packages
        tool        run specified go tool
        version     print Go version
        vet         report likely mistakes in packages

Use "go help <command>" for more information about a command.

Additional help topics:

        buildconstraint build constraints
        buildmode       build modes
        c               calling between Go and C
        cache           build and test caching
        environment     environment variables
        filetype        file types
        go.mod          the go.mod file
        gopath          GOPATH environment variable
        gopath-get      legacy GOPATH go get
        goproxy         module proxy protocol
        importpath      import path syntax
        modules         modules, module versions, and more
        module-get      module-aware go get
        module-auth     module authentication using go.sum
        packages        package lists and patterns
        private         configuration for downloading non-public code
        testflag        testing flags
        testfunc        testing functions
        vcs             controlling version control with GOVCS

Use "go help <topic>" for more information about that topic.
```

# Build

> powershell

```powershell
.\build.ps1
```

> bash

```bash
./build.sh
```