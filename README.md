# mp3vis

POC / starting point for a mp3 visualizer akin to [projectM](https://github.com/projectM-visualizer/projectm), written in Go.

## Prerequisites

* Your OS must use [PulseAudio](https://www.freedesktop.org/wiki/Software/PulseAudio/) for sound management.  
* Your graphics card must support OpenGL 4.1.
* For compilation, the [Go v1.20](https://go.dev/doc/install) compiler is needed.

This project was developed/tested on Ubuntu 22.04.  OpenGL libs installed via:
```shell
# apt install libgl1-mesa-dev xorg-dev
```

## Compiling

Ensure Go is installed and in your PATH, then run:
```shell
chmod +x build.sh
./build.sh
```

## Running

To run with the bundled image/mp3 and default window size:
```shell
chmod +x run.sh
./run.sh
```

To run with a custom image/mp3 and window size :
```shell
./mp3vis <path_to_mp3> <path_to_png> <window_width> <window_height>
```
