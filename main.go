package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/big"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("Insufficient arguments given")
		printUsage()
	}

	output := args[0]
	var mainColour color.RGBA

	if output == "" {
		fmt.Println("Output argument must be specified")
		printUsage()
	} else if !strings.HasSuffix(output, ".png") {
		fmt.Println("Output file must be .png file")
		printUsage()
	}

	if args[1] == "" {
		fmt.Println("Colour must be specified")
		printUsage()
	} else {
		switch args[1] {
		case "R", "r":
			mainColour = color.RGBA{
				R: math.MaxUint8,
				G: 0,
				B: 0,
				A: math.MaxUint8,
			}
		case "G", "g":
			mainColour = color.RGBA{
				R: 0,
				G: math.MaxUint8,
				B: 0,
				A: math.MaxUint8,
			}
		case "B", "b":
			mainColour = color.RGBA{
				R: 0,
				G: 0,
				B: math.MaxUint8,
				A: math.MaxUint8,
			}
		default:
			fmt.Println("Color must be 'R' 'G' or 'B' (not case sensitive)")
			printUsage()
		}
	}

	bounds := image.Rect(0, 0, 128, 128)
	pic := image.NewRGBA(bounds)

	paintBlock := func(x, y int, color color.Color) {
		for i := 0; i < 16; i++ {
			for j := 0; j < 16; j++ {
				pic.Set(x+i, y+j, color)
			}
		}
	}

	var blocks [4][8]bool
	for x, column := range blocks {
		for y := range column {
			val, err := rand.Int(rand.Reader, big.NewInt(2))
			if err != nil {
				panic(err)
			}

			coloured := val.Int64() == big.NewInt(1).Int64()
			if coloured {
				paintBlock(x*16, y*16, mainColour)
			} else {
				paintBlock(x*16, y*16, color.White)
			}
			blocks[x][y] = coloured
		}
	}

	for x, column := range blocks {
		for y, coloured := range column {
			base := 128 - 16

			if coloured {
				paintBlock(base-(x*16), y*16, mainColour)
			} else {
				paintBlock(base-(x*16), y*16, color.White)
			}
		}
	}

	buffer := bytes.NewBuffer(nil)
	err := png.Encode(buffer, pic)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func printUsage() {
	fmt.Println("Usage: pp-gen <outputFile> <color (R or G or B)>")
	os.Exit(-1)
}
