# Installing LXD

The easiest way to install LXD is to install one of the available packages, but you can also install LXD from the sources.

## Installing LXD from packages

% Include some content from [../README.md](../README.md)
```{include} ../README.md
    :start-after: Installing LXD from packages
    :end-before: To install LXD from source, see
```

## Installing LXD from source
We recommend having the latest versions of liblxc (>= 4.0.0 required)
available for LXD development. Additionally, LXD requires Golang 1.16 or
later to work. On ubuntu, you can get those with:

```bash
sudo apt update
sudo apt install acl attr autoconf dnsmasq-base git golang libacl1-dev libcap-dev liblxc1 liblxc-dev libsqlite3-dev libtool libudev-dev liblz4-dev libuv1-dev make pkg-config rsync squashfs-tools tar tcl xz-utils ebtables
```

There are a few storage backends for LXD besides the default "directory" backend.
Installing these tools adds a bit to initramfs and may slow down your
host boot, but are needed if you'd like to use a particular backend:

```bash
sudo apt install lvm2 thin-provisioning-tools
sudo apt install btrfs-progs
```

To run the testsuite, you'll also need:

```bash
sudo apt install curl gettext jq sqlite3 uuid-runtime socat bind9-dnsutils
```

### From Source: Building the latest version

These instructions for building from source are suitable for individual developers who want to build the latest version
of LXD, or build a specific release of LXD which may not be offered by their Linux distribution. Source builds for
integration into Linux distributions are not covered here and may be covered in detail in a separate document in the
future.

```bash
git clone https://github.com/lxc/lxd
cd lxd
```

This will download the current development tree of LXD and place you in the source tree.
Then proceed to the instructions below to actually build and install LXD.

### From Source: Building a Release

The LXD release tarballs bundle a complete dependency tree as well as a
local copy of libraft and libdqlite for LXD's database setup.

```bash
tar zxvf lxd-4.18.tar.gz
cd lxd-4.18
```

This will unpack the release tarball and place you inside of the source tree.
Then proceed to the instructions below to actually build and install LXD.

### Starting the Build

The actual building is done by two separate invocations of the Makefile: `make deps` -- which builds libraries required
by LXD -- and `make`, which builds LXD itself. At the end of `make deps`, a message will be displayed which will specify environment variables that should be set prior to invoking `make`. As new versions of LXD are released, these environment
variable settings may change, so be sure to use the ones displayed at the end of the `make deps` process, as the ones
below (shown for example purposes) may not exactly match what your version of LXD requires:

We recommend having at least 2GB of RAM to allow the build to complete.

```bash
make deps
# Follow the instructions from `make deps` to export the required environment variables.
# For example:
#  export CGO_CFLAGS="${CGO_CFLAGS} -I$(go env GOPATH)/deps/dqlite/include/ -I$(go env GOPATH)/deps/raft/include/"
#  export CGO_LDFLAGS="${CGO_LDFLAGS} -L$(go env GOPATH)/deps/dqlite/.libs/ -L$(go env GOPATH)/deps/raft/.libs/"
#  export LD_LIBRARY_PATH="$(go env GOPATH)/deps/dqlite/.libs/:$(go env GOPATH)/deps/raft/.libs/:${LD_LIBRARY_PATH}"
#  export CGO_LDFLAGS_ALLOW="(-Wl,-wrap,pthread_create)|(-Wl,-z,now)"
make
```

### From Source: Installing

Once the build completes, you simply keep the source tree, add the directory referenced by `$(go env GOPATH)/bin` to
your shell path, and set the `LD_LIBRARY_PATH` variable printed by `make deps` to your environment. This might look
something like this for a `~/.bashrc` file:

```bash
export PATH="${PATH}:$(go env GOPATH)/bin"
export LD_LIBRARY_PATH="$(go env GOPATH)/deps/dqlite/.libs/:$(go env GOPATH)/deps/raft/.libs/:${LD_LIBRARY_PATH}"
```

Now, the `lxd` and `lxc` binaries will be available to you and can be used to set up LXD. The binaries will automatically find and use the dependencies built in `$(go env GOPATH)/deps` thanks to the `LD_LIBRARY_PATH` environment variable.

### Machine Setup
You'll need sub{u,g}ids for root, so that LXD can create the unprivileged containers:

```bash
echo "root:1000000:1000000000" | sudo tee -a /etc/subuid /etc/subgid
```

Now you can run the daemon (the `--group sudo` bit allows everyone in the `sudo`
group to talk to LXD; you can create your own group if you want):

```bash
sudo -E PATH=${PATH} LD_LIBRARY_PATH=${LD_LIBRARY_PATH} $(go env GOPATH)/bin/lxd --group sudo
```
