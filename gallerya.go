package main

import (
    "github.com/disintegration/imaging"
    "io/ioutil"
    "fmt"
    "sync"
    "os"
    "path/filepath"

)
var dir_thumb = "./thumb"
var dir_original = "./original"

var wg sync.WaitGroup

func thumb_from_dir(path string)  {
    files, _ := ioutil.ReadDir(path)
    wg.Add(len(files))
    for _, f := range files {
            go do_thumb(path,f.Name())
            
    }
}
func do_thumb(path string,file string) {
    defer wg.Done()
    fmt.Println(filepath.Join(path, file))
    img, err := imaging.Open(filepath.Join(path, file))
    dst := imaging.Fill(img, 100, 100, imaging.Center, imaging.Lanczos)
    err2 := imaging.Save(dst, filepath.Join(dir_thumb,file))
    if err2 != nil {
        panic(err)
    }
}

func fexists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func main() {
    r,_ := fexists(dir_thumb)
    if(r) {
        fmt.Println(filepath.Join(dir_thumb, "file"))
    } else {
        os.MkdirAll(dir_thumb, os.ModePerm)
    }
    thumb_from_dir(dir_original)
    wg.Wait()
}