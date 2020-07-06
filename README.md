# rifx
Binary parsing module for RIFX files (Big-Endian variant of the Resource Interchange File Format).

The Resource Interchange File Format (RIFF) is a generic file container format for storing data in tagged chunks. It is primarily used to store multimedia such as sound and video, though it may also be used to store any arbitrary data ([Read more on Wikipedia](https://en.wikipedia.org/wiki/Resource_Interchange_File_Format))

### Quick Start

```bash
go get -u github.com/rioam2/rifx
```

```go
package main

import (
    "os"
    "fmt"
    "github.com/rioam2/rifx"
)

func main() {
    file, err := os.Open("my-rifx-file.wav")
    if err != nil {
        panic(err)
    }
    rootList, err := rifx.FromReader(file)
    if err != nil {
        panic(err)
    }
    rootList.ForEach(func(block *rifx.Block) {
        fmt.Println(block)
    })
}
```