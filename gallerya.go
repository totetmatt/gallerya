package main

import (
    "github.com/disintegration/imaging"
    "io/ioutil"
    "fmt"
    "os"
    "path/filepath"

)
var dir_thumb = "./thumb"
var dir_original = "./original"

func worker(id int, jobs <-chan string, results chan<- string) {
    for j := range jobs {
        fmt.Println(">> worker", id, "processing job", j)
        do_thumb(j)
        results <- j
    }
}

func thumb_from_dir(path string)  {
    
    jobs := make(chan string, 1000)
    results := make(chan string, 1000)

    for w := 1; w <= 8; w++ {
        go worker(w, jobs, results)
    }
  
    
    files, _ := ioutil.ReadDir(path)
    for _, f := range files {
            jobs <- filepath.Join(path,f.Name())    
    }
    close(jobs)
    

    for a := 1; a <=  len(files); a++ {
        r := <-results
        fmt.Println("Finished "+r)
    }
}

func do_thumb(path string) {
    img, err := imaging.Open(path)
    dst := imaging.Fill(img, 100, 100, imaging.Center, imaging.Lanczos)
    _,file := filepath.Split(path)
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
    if(!r) {
        os.MkdirAll(dir_thumb, os.ModePerm)
    }
    thumb_from_dir(dir_original)
}