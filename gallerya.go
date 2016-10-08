package main

import (
    "github.com/disintegration/imaging"
    "io/ioutil"
    "fmt"
    "os"
    "path/filepath"
    "github.com/flosch/pongo2"
    "flag"
)
    
/** Configuration Struct >>**/
type GalleryaConfiguration struct {
    thumb_directory string
    original_directory  string
    
    title string


    workers int
    skip_image_processing bool

    files []os.FileInfo
}
    
func (config *GalleryaConfiguration) preCheck() {
    r,_ := fexists(config.thumb_directory)
    if(!r) {
        os.MkdirAll(config.thumb_directory, os.ModePerm)
    }
    config.get_original_files()
}

func (config *GalleryaConfiguration) thumb_file(filename string) string {
    return filepath.Join(config.thumb_directory,filename)   
}

func (config *GalleryaConfiguration) original_file(filename string) string {
    return filepath.Join(config.original_directory,filename)   
}

func (config *GalleryaConfiguration) get_original_files() {
    config.files,_ = ioutil.ReadDir(config.original_directory)
}
/**<< Configuration Struct **/

var templateIndex = pongo2.Must(pongo2.FromFile("./index.template"))

func worker(config *GalleryaConfiguration, jobs <-chan string, results chan<- string) {
    for j := range jobs {
        fmt.Println("Processing image", j)
        do_thumb(j,config)
        results <- j
    }
}
func image_processing(config *GalleryaConfiguration) {
    jobs := make(chan string, 1000)
    results := make(chan string, 1000)
   
    // Init workers
    for w := 1; w <= config.workers; w++ {
        go worker(config, jobs, results)
    }
    // Give job to workers
    for _, f := range config.files {
            jobs <- f.Name()

    }
    close(jobs)

    //Harverst wokers results
    for a := 1; a <=  len(config.files); a++ {
        r := <-results
        fmt.Println("Finished "+r)
    }

}

func generate(config *GalleryaConfiguration)  {
    if(!config.skip_image_processing) {
        image_processing(config)
    }
    generate_html(config)
}

func generate_html(config *GalleryaConfiguration) {
    f, _ := os.Create("./index.html")
    defer f.Close()
    err := templateIndex.ExecuteWriter(pongo2.Context{"config":config,"files": config.files}, f)
    if err != nil {
        panic(err)
    }
}

func do_thumb(file string, config *GalleryaConfiguration) {
    img, err := imaging.Open(config.original_file(file))
    dst := imaging.Fill(img, 360, 247, imaging.Center, imaging.Lanczos)
    err2 := imaging.Save(dst, config.thumb_file(file))
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
    config := GalleryaConfiguration{}

    flag.StringVar(&config.thumb_directory,"thumbnail", "./thumb", "Thumbnail directory")
    flag.StringVar(&config.original_directory ,"original", "./original", "Original Photo directory")

    flag.StringVar(&config.title ,"title", "Gallerya", "Title of the gallery")

    flag.IntVar(&config.workers,"workers",4,"Number of workers")
    flag.BoolVar(&config.skip_image_processing,"skip-image-process",false,"Skip the image transformation")
    
    flag.Parse()
    config.preCheck()
    generate(&config)
}