cmake_minimum_required(VERSION 3.13)
project(custom_ffmpeg)
set(CMAKE_CXX_STANDARD 17)

## options
set(BMF_ROOT "/opt/tiger/bmf/" CACHE STRING "Bmf output path")
set(FFMPEG_ROOT "/opt/")


## BMF setup
set(BMF_INCLUDE_PATH ${BMF_ROOT}/bmf/include)
set(BMF_LIB_PATH ${BMF_ROOT}/bmf/lib/)
find_library(BMF_LIB bmf_module_sdk PATHS ${BMF_LIB_PATH})
message("Find BMF SDK Library: ${BMF_LIB}")

## FFMPEG setup
find_path(FFMPEG_INCLUDE_DIR libavcodec/avcodec.h REQUIRED
            PATHS ${FFMPEG_ROOT} /opt/conda/include /usr/include /usr/local/include)
message("Find FFMPEG include DIR: ${FFMPEG_INCLUDE_DIR}")

## 
add_library(custom_ffmpeg SHARED custom_ffmpeg.cpp)

target_include_directories(custom_ffmpeg 
	PRIVATE 
	${BMF_INCLUDE_PATH}
	${FFMPEG_INCLUDE_DIR}
	)
target_link_libraries(custom_ffmpeg PRIVATE ${BMF_LIB})