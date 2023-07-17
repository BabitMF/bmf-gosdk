# Go SDK for BMF

This repo implements some basic data structure used by BMF processing.

In addition, it also contains a pre-written CGO exporting tool to compile the user-written module into BMF-recognizable library.

The following code shows a module copying input packets to downstreams.

Ref: [pass_through](example/pass_through.go)


## Run tests
1. setup envs
> export CGO_FLAGS=xxx
> export CGO_LDFLAGS=xxx
> export LD_LIBRARY_PATH=xxx

2. build example module
> cd example && go build -buildmode c-shared -o go_pass_through.so pass_through.go && cd -

3. run unit tests
> cd bmf && go test -v .
