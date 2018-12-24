package main

import (
    "../util"
    "./conversion"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run cli.go <raw db path> <converted db directory path>")
        return
    }

    rawDbPath := os.Args[1]
    convertedDirPath := os.Args[2]

    t1 := util.CurTime()
    conversion.ConvertRawDatabase(rawDbPath, convertedDirPath)
    t2 := util.CurTime()

    fmt.Printf("Database wasconverted in %.3f sec", float32((t2 - t1) / 1000000) / 1000)
}
