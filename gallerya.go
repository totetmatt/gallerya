package main

import (
    "github.com/disintegration/imaging"
    "io/ioutil"
    "fmt"
    "os"
    "path/filepath"
    "github.com/flosch/pongo2"
    "flag"
    "totetmatt/gallerya/exifdata"
    "log"
)
type Media struct {
    name string
    metadata string
}    
/** Configuration Struct >>**/
type GalleryaConfiguration struct {
    thumb_directory string
    original_directory  string
    static_path string

    title string


    workers int
    skip_image_thumb_processing bool

    medias []Media
    metadata map[string]string
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
    files,_ := ioutil.ReadDir(config.original_directory)
    config.medias =  make([]Media,len(files))

    for i := 0; i < len(files); i++ {
        config.medias[i] = Media{}
        config.medias[i].name = files[i].Name()
    }
}
/**<< Configuration Struct **/

var templateIndex = pongo2.Must(pongo2.FromFile("./index.template"))

func worker(config *GalleryaConfiguration, jobs <-chan string, results chan<- string) {
    for j := range jobs {
        fmt.Println("Processing image", j)
        if(!config.skip_image_thumb_processing) {
            process_image(j,config)
        }
        

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
    for _, f := range config.medias{
            jobs <- f.name

    }
    close(jobs)

    //Harverst wokers results
    for a := 1; a <=  len(config.medias); a++ {
        r := <-results
        fmt.Println("Finished "+r)
    }

    // Metadata 
    for i := 0; i < len(config.medias); i++ {
        config.medias[i].metadata = extract_data(config.medias[i].name,config)
    }
}

func generate(config *GalleryaConfiguration)  {
   
    image_processing(config)
    
    generate_html(config)
}

func generate_html(config *GalleryaConfiguration) {
    f, _ := os.Create("./index.html")
    defer f.Close()
    err := templateIndex.ExecuteWriter(pongo2.Context{"config":config}, f)
    if err != nil {
        panic(err)
    }
}

func process_image(file string, config *GalleryaConfiguration) {
    img, err := imaging.Open(config.original_file(file))
    dst := imaging.Fill(img, 360, 247, imaging.Center, imaging.Lanczos)
    err2 := imaging.Save(dst, config.thumb_file(file))
    if err2 != nil {
        panic(err)
    }
}
func extract_data(file string,config *GalleryaConfiguration) string {
    f, err := os.Open(config.original_file(file))
    if err != nil {
        log.Fatal(err)
    }
    data :=exifdata.ExifData{}
    data.Grab_data(f)
    return data.String()
}

func fexists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func main() {
    config := GalleryaConfiguration{}
    config.metadata = make(map[string]string)
    flag.StringVar(&config.thumb_directory,"thumbnail", "./thumb", "Thumbnail directory")
    flag.StringVar(&config.original_directory ,"original", "./original", "Original Photo directory")
    flag.StringVar(&config.static_path ,"static", "./static", "Path to Static File")
    flag.StringVar(&config.title ,"title", "Gallerya", "Title of the gallery")

    flag.IntVar(&config.workers,"workers",4,"Number of workers")
    flag.BoolVar(&config.skip_image_thumb_processing,"skip-image-thumb-processing",false,"Skip the image transformation")
    
    flag.Parse()
    config.preCheck()
    generate(&config)
}