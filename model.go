package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
)

type Object struct {
	width, height float64
	name          string
	id            string
}

func getImage(o Object) image.Image {
	fn := o.id + ".png"
	if _, e := os.Stat(fn); e == nil {
		f, _ := os.Open(fn)
		im, _ := png.Decode(f)
		f.Close()
		return im
	} else if errors.Is(e, os.ErrNotExist) {
		panic(fmt.Sprintf("File %s does not exist", fn))
	} else {
		panic(fmt.Sprintf("Cannot access file %s", fn))
	}
}

var objects = []Object{
	{3.8, 10, "M2 8mm Socket Head Screw", "91292A832"},
	{7, 18, "M4 14mm Socket Head Screw", "91292A038"},
	{6.86, 20.04, "8-32 5/8in Socket Head Screw", "92196A196"},
	{8, 8, "M4 Washer", "98689A113"},
	{9, 9, "M5 Washer", "98689A114"},
	{11, 11, "M6 Washer", "98689A115"},
	{7, 8.08, "M4 Nut", "91828A231"},
	{5.5, 6.35, "M3 Nut", "91828A211"},
	{5.5, 6.35, "M3 Nyloc Nut", "93625A100"},
	{8.5, 21, "M5 16mm Socket Head Screw", "91292A126"},
	{8.5, 25, "M5 20mm Socket Head Screw", "91292A128"},
	{8.5, 17, "M5 12mm Socket Head Screw", "91292A125"},
	{9.5, 14.75, "M5 12mm Button Head Screw", "92095A210"},
	{10, 16, "M5 16mm Flat Head Screw", "92125A212"},
}
