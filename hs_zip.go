/*
The MIT License (MIT)

Copyright (c) 2015 David Schmidt ()

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
    "archive/zip"
    "fmt"
    "github.com/nfnt/resize"
    "github.com/oliamb/cutter"
    "image/jpeg"
    "image/png"
    "io"
    "io/ioutil"
    "os"
    "time"
)

const (
	TargetWidth = 1024
	TargetHeight = 768
)
// 90% jpeg
var JpegQuality *jpeg.Options = &jpeg.Options{ 90 }



// read, resize, crop and zip a 'screenshot' named 'name' into 'zipArchive'
func convert(zipFile *zip.Writer, infile io.Reader, filename string) error {
    img, err := png.Decode(infile)
    if err != nil {
        return err
    }

    resized := resize.Resize(0, TargetHeight, img, resize.Bicubic)
    croppedImg, _ := cutter.Crop(resized, cutter.Config{
        Width: TargetWidth,
        Height: TargetHeight,
        Mode: cutter.Centered,
    })

    zippedFile, err := zipFile.Create(filename)
    if err != nil {
        return err
    }

    return jpeg.Encode(zippedFile, croppedImg, JpegQuality)
}


// find hearthstone screenshots in current directory
func findScreenshots() []string {
    screenshots := make([]string, 0, 8)
    // get all files from current directory
    files, err := ioutil.ReadDir("./")
    if err != nil {
        fmt.Printf("failed to read current directory - error: %v", err)
        return screenshots
    }

    // filter Hearthstone*.png files
    for _, file := range files {
        if ! file.IsDir() &&
                len(file.Name()) > 15 &&
                file.Name()[:11] == "Hearthstone" &&
                file.Name()[len(file.Name()) - 4:] == ".png" {
            screenshots = append(screenshots, file.Name())
        }
    }
    
    return screenshots
}


// main entry
func main() {
    // find screenshot file names
    screenshots := findScreenshots()
    if len(screenshots) == 0 {
        fmt.Printf("\nno screenshots found")
        return
    }
    
    // create zip archive file if there are screenshots
    outfile, err := os.Create("Hearthstone Screenshots.zip")
    if err != nil {
        fmt.Printf("could not create zip file to write to - error: %v", err)
        return
    }
    defer outfile.Close()

    zipwriter := zip.NewWriter(outfile)
    defer zipwriter.Close()

    fmt.Printf("\nstart zippping: %s\n", time.Now().Format("15:04:05.999"))
    // process all screenshots
    for _, screenshot := range screenshots {
        // try to open file
        fmt.Printf("processing: %v\n", screenshot)
        file, err := os.Open(screenshot)
        if err != nil {
            fmt.Printf("could not open file %v error: %v", screenshot, err)
            return
        }
        defer file.Close()

        jpegName := screenshot[:len(screenshot) - 4] + ".jpg"
        // zip it as smaller jpeg
        err = convert(zipwriter, file, jpegName)
        if err != nil {
            fmt.Printf("could not resize file %v error: %v", screenshot, err)
            return
        }
    }
    fmt.Printf("finished: %s", time.Now().Format("15:04:05.999"))
}