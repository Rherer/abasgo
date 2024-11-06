# abasgo - a simple Notification Wrapper

This is a simple wrapper for the abasgui.exe, that notifies the user if changes in a file were made.
This is done for the purpose of notifying the programmer, when the automatic builds in his clients are finished/started.

## Installation

Place the compiled executable as well as the abasgo.yml into the network share, containing your GUI-Client.
Simply adjust the settings in abasgo.yml for your use-case.
If you do not provide a abasgo.yml settings file, a default one will be generated.

## Self building

#### Build for current os

```
git clone https://github.com/Rherer/abasgo.git
go build 
```

#### Build for windows on linux

```
GOOS=windows GOARCH=386 \
go build
```

## Usage

Just launch this executable instead of the main abasgui.exe.
The format of the build status file can be seen in the example file.

### Features

- Adjustable paths
- Adjustable launch parameters
- OS-Agnostic (Tested on Windows + Linux)

### Contributing

You can contribute to this project if you want to.
Just create a Pull-Request with your changes and a description and wait for approval.