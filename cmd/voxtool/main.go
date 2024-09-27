package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/qbradq/cubit/internal/util"
)

var inputFileName = flag.String("in", "in.vox",
	"File name of the .vox model to inspect.")

func main() {
	flag.Parse()
	if len(*inputFileName) < 1 {
		flag.PrintDefaults()
		os.Exit(1)
	}
	f, err := os.Open(*inputFileName)
	if err != nil {
		panic(err)
	}
	vox, err := util.NewVoxFromReader(f)
	f.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s:\nW=%d H=%d D=%d\n", *inputFileName, vox.Width, vox.Height,
		vox.Depth)
	min := [3]int{math.MaxInt, math.MaxInt, math.MaxInt}
	max := [3]int{math.MinInt, math.MinInt, math.MinInt}
	for iy := 0; iy < vox.Height; iy++ {
		for iz := 0; iz < vox.Depth; iz++ {
			for ix := 0; ix < vox.Width; ix++ {
				i := iy*vox.Depth*vox.Width + iz*vox.Width + ix
				if vox.Voxels[i][3] == 0 {

				}
			}
		}
	}
}
