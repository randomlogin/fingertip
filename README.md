# Fingertip

**Note:** This project is experimental use at your own risk.

Fingertip is a menubar app that runs a [lightweight decentralized resolver](https://github.com/handshake-org/hnsd) to resolve names from the [Handshake](https://handshake.org) root zone. It can also resolve names from external namespaces such as the Ethereum Name System. Fingertip integrates with [sane](https://github.com/randomlogin/sane) to provide TLS support without relying on a centralized certificate authority. 

For handshake domains fingertip can be thought as a user-friendly wrapper of SANE, it uses hardcoded community-hosted external proof services, and DNS over HTTPS for name resolution. An advanced user is welcome to use [sane](https://github.com/randomlogin/sane) directly.


<img width="600" src="https://user-images.githubusercontent.com/41967894/127166063-fedf072c-fa5e-45e3-acac-bfb46f256831.png" />

## Install

You can use a pre-built binary from releases or build your own from source.

To run pre-build AppImage on Linux you might need `libfuse`:

```
apt install libfuse2`
```

## Configuration
You can set these as environment variables prefixed with `FINGERTIP_` or store it in the app config directory as `fingertip.env`

```
# sane proxy address
PROXY_ADDRESS=127.0.0.1:9590
# hnsd root server address
ROOT_ADDRESS=127.0.0.1:9591
# hnsd recursive resolver address
RECURSIVE_ADDRESS=127.0.0.1:9592
# Connect your own Ethereum full node/or blockchain provider such as Infura
#ETHEREUM_ENDPOINT=/home/user/.ethereum/geth.ipc or
#ETHEREUM_ENDPOINT=https://mainnet.infura.io/v3/YOUR-PROJECT-ID
```

## Build from source

Go 1.16+ is required.

```
$ git clone https://github.com/imperviousinc/fingertip
```

### MacOS

```
$ brew install dylibbundler git automake autoconf libtool unbound
$ git clone https://github.com/randomlogin/fingertip
$ cd fingertip && ./builds/macos/build.sh
```

For development, you can run fingertip from the following path:
```
$ ./builds/macos/Fingertip.app/Contents/MacOS/fingertip
```
        
Configure your IDE to output to this directory or continue to use `build.sh` when making changes (it will only build hnsd once).

### Linux

Follow [hnsd](https://github.com/handshake-org/hnsd) build instructions for Linux. Copy hnsd binary into the `fingertip/builds/linux/appdir/usr/bin` directory.

```
$ go build -trimpath -o ./builds/linux/appdir/usr/bin/
```

To create an AppImage run 

```
bash builds/linux/create_appimage.sh 
```

### Windows

Due to the [difference](https://github.com/handshake-org/hnsd/issues/128) in hnsd behaviour on Windows and other platforms (and overall complexity of building for windows),
Windows is not supported. This may change in future. 


## Credits
Fingertip uses [hnsd](https://github.com/handshake-org/hnsd) a lightweight Handshake resolver, [sane](https://github.com/randomlogin/sane) and [getdns](https://getdnsapi.net/) for TLS support and [go-ethereum](https://github.com/ethereum/go-ethereum) for .eth and Ethereum [HIP-5](https://github.com/handshake-org/HIPs/blob/master/HIP-0005.md) lookups.

The name "fingertip" was stolen from [@pinheadmz](https://github.com/pinheadmz)
