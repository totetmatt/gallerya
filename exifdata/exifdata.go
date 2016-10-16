package exifdata

import (
    "io"
    "github.com/rwcarlsen/goexif/exif"
    "github.com/rwcarlsen/goexif/mknote"

)
type ExifData struct {
    focal_length string
    date string
    fnumber string
    exposure_time string
    iso string

}
func (d *ExifData) Grab_data(r io.Reader) {
    exif.RegisterParsers(mknote.All...)
    x, err := exif.Decode(r)
    if err != nil {
        d.focal_length =""
        d.date =""
        d.fnumber =""
        d.exposure_time =""
        d.iso =""
    } else {
        d.focal_data(x)
        d.date_data(x)
        d.fnumber_data(x)
        d.exposureTime_data(x)
        d.iso_data(x)
    }
}
func (d *ExifData) focal_data(e *exif.Exif) {
    focal, err := e.Get(exif.FocalLength)
    if(err != nil) {
        d.focal_length= ""
    } else {
        numer , _ := focal.Rat(0) 
        d.focal_length = numer.FloatString(0)
    }
}
func (d *ExifData) date_data(e *exif.Exif) {
    tm, err := e.DateTime()
    if(err != nil) {
        d.date = ""
    } else {
        d.date = tm.String()
    }
    
}
func (d *ExifData) fnumber_data(e *exif.Exif)  {
    fnumber ,err := e.Get(exif.FNumber)
    if(err != nil) {
        d.fnumber = ""
    } else {
        fn,_ := fnumber.Rat(0)
        d.fnumber = fn.FloatString(1)
    }
}
func (d *ExifData)  exposureTime_data(e *exif.Exif)  {
    exposureTime ,err := e.Get(exif.ExposureTime )
    if(err != nil) {
            d.exposure_time = ""
        } else {
            etime,_ := exposureTime.Rat(0)
            d.exposure_time = etime.RatString()   
        }

}
func (d *ExifData) iso_data(e *exif.Exif) {
    iso ,err := e.Get(exif.ISOSpeedRatings )
    if(err != nil) {
        d.iso = ""
    } else {
        d.iso = iso.String()
    }
    
}
func (d *ExifData) String() string {
    return "F/"+d.fnumber+" - ISO/"+d.iso+" - "+d.focal_length+ "mm - "+d.exposure_time +"s"
}