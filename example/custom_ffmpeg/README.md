
> This example demostrate how to integrate custom ffmpeg library to bmf-gosdk

## How to build?

1. Setup BMF envrionment

2. Build custom_ffmpeg library(cpp)
    > cd cpp && mkdir build && cd build 
    > cmake -DBMF_ROOT=/opt/tiger/bmf -DFFMPEG_ROOT=/usr/local ..
    > make -j12

3. Run GO tests
    > export CGO_FLAGS=-I/path/to/bmf/include -I/path/to/ffmpeg/include
    > export CGO_LDFLAGS=-L/path/to/bmf/lib
    > export LD_LIBRARY_PATH=/path/to/bmf/lib:/path/to/ffmpeg_lib:/path/to/ffmpeg_custom/lib
    > cd ../../go/
    > go test -v