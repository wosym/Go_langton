package main

import (
	"fmt"
    "time"
    "math"
    "math/rand"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/yakshaveinc/go-keycodes"
)

const winWidth, winHeight int = 1500, 1500
const gridDim int = 200     //size of grid
const stepTime = 0        //time between ant steps in ms

//Pattern: L: true, R: false
var pattern = []bool{true,false,false,false,false,false,true,true,false}   //TODO: selectable?
//var pattern = []bool{true,false,false,false,false,false,true,true,false,true,false,false,false,false,false,true,true,false,true,false,false,false,false,false,true,true,false,true,false,false,false,false,false,true,true,false}
var colorList = []sdl.Color{}
var shuffleColors = true

type position struct {
    x,y int
}

const (
    NORTH = iota
    EAST
    SOUTH
    WEST
    OF
)

//TODO: move these color-util functions to seperate file
type HSV struct {
	H, S, V float64
}

type RGB struct {
	R, G, B float64
}

func (c HSV) RGB() sdl.Color {
	var r, g, b float64
	if c.S == 0 { //HSV from 0 to 1
		r = c.V * 255
		g = c.V * 255
		b = c.V * 255
	} else {
		h := c.H/360 * 6
		if h == 6 {
			h = 0
		} //H must be < 1
		i := math.Floor(h) //Or ... var_i = floor( var_h )
		v1 := c.V * (1 - c.S)
		v2 := c.V * (1 - c.S*(h-i))
		v3 := c.V * (1 - c.S*(1-(h-i)))

		if i == 0 {
			r = c.V
			g = v3
			b = v1
		} else if i == 1 {
			r = v2
			g = c.V
			b = v1
		} else if i == 2 {
			r = v1
			g = c.V
			b = v3
		} else if i == 3 {
			r = v1
			g = v2
			b = c.V
		} else if i == 4 {
			r = v3
			g = v1
			b = c.V
		} else {
			r = c.V
			g = v1
			b = v2
		}

		r = r * 255 //RGB results from 0 to 255
		g = g * 255
		b = b * 255
	}
	return sdl.Color{uint8(r), uint8(g), uint8(b), 255}

}

func printGrid(grid [][]int) {
    for y := 0; y < gridDim; y++ {
        fmt.Println(grid[y])
    }
}
func drawAnt(renderer *sdl.Renderer, nx, ny, visRad int32, antpos position) {
    bw := int32(winWidth)/(visRad*2)
    bh := int32(winHeight)/(visRad*2)
    startX := int32(nx/2) - visRad
    startY := int32(ny/2) - visRad
    cx := bw * (int32(antpos.x)-startX) + bw/2
    cy := bh * (int32(antpos.y)-startY) + bh/2
    //TODO: make ant ellipse-form based on direction?

    ret := gfx.FilledEllipseColor(renderer, cx, cy, bw/2, bh/2, sdl.Color{100, 50, 0, 255})
    if !ret {
        fmt.Println("Error while drawing box")
    }
}
func drawGrid(renderer *sdl.Renderer, grid [][]int, nx, ny, visRad int32) {
    bw := int32(math.Round(float64(winWidth)/float64(visRad*2)))  //TODO: still not 100% correct. math.Round didn't solve everything. There still is a black band with some gridDims. math.Round made it a bit better though.
    bh := int32(math.Round(float64(winHeight)/float64(visRad*2)))
    var ret bool

    startX := int32(nx/2) - visRad
    endX := int32(nx/2) + visRad
    startY := int32(ny/2) - visRad
    endY := int32(nx/2) + visRad

    for y  := startY; y < endY; y++ {
        for x := startX; x < endX; x++ {
            ret = gfx.BoxColor(renderer, (x-startX)*bw, (y-startY)*bh, ((x-startX)+1)*bw, ((y-startY)+1)*bh, colorList[grid[y][x]])
            if !ret {
                fmt.Println("Error while drawing box")
            }




            ret = gfx.RectangleColor(renderer, (x-startX)*bw, (y-startY)*bh, ((x-startX)+1)*bw, ((y-startY)+1)*bh, sdl.Color{0, 0, 0, 255})
            if !ret {
                fmt.Println("Error while drawing rect")
            }
        }
    }

    

}

func moveAnt(grid [][]int, antpos *position, antdir *int) bool {
    //move
    switch *antdir {
        case NORTH:
            (*antpos).y--   //pixel coordinates origin is in top left corner
        case EAST:
            (*antpos).x++
        case SOUTH:
            (*antpos).y++
        case WEST:
            (*antpos).x--
        default:
            fmt.Println("Error moving ant: illegal direction")
    }
    if (*antpos).x < 0 || (*antpos).y < 0 || (*antpos).x > gridDim-1 || (*antpos).y > gridDim-1 {
        return true     //pause
    }

    //rotate based on the square we land on
    if grid[(*antpos).y][(*antpos).x] >= len(pattern) {
        fmt.Println("cell has invalid value! This should never happen!")
        return true
    }

    if pattern[grid[(*antpos).y][(*antpos).x]] {
        *antdir--;
    } else {
        *antdir++;
    }

    //Check for overflows
    if *antdir >= OF{
        *antdir = NORTH;
    } else if *antdir <= -1 {
        *antdir = WEST;
    }

    //Update cell
    grid[(*antpos).y][(*antpos).x]++
    if grid[(*antpos).y][(*antpos).x] >= len(pattern) {
        grid[(*antpos).y][(*antpos).x] = 0
    }

    return false    //don't pause



}
func IntegerAbs(x int) int {   //TODO: move to util
	if x < 0 {
		return -x
	}
	return x
}

func main() {
    err := sdl.Init(sdl.INIT_EVERYTHING)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer sdl.Quit()

    window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer window.Destroy()

    renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer renderer.Destroy()

    tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
    if err != nil {
        fmt.Println(err)
        return
    }
    defer tex.Destroy()

    //Create grid for the ant
    grid := make([][]int, gridDim)
    for i := 0; i < gridDim; i++ {
        grid[i] = make([]int, gridDim)
    }


    //fill colorList based on patternlength //TODO: seperate function?  //TODO: add shuffle option to make highways more clearly
    spacing := float64(360 / len(pattern))
    var tmpCol = HSV{}
    for i := 0; i<=len(pattern); i++ {
        tmpCol = HSV{float64(i)*spacing, 0.5 ,0.5}
        colorList = append(colorList, tmpCol.RGB())
    }
    if shuffleColors {
        rand.Seed(time.Now().UnixNano())
        rand.Shuffle(len(colorList), func(i, j int) {colorList[i], colorList[j]=colorList[j], colorList[i]})
    }

    antpos := position{gridDim/2, gridDim/2}
    antdir := NORTH
    fmt.Println("Starting position for ant: ", antpos, "in direction: ", antdir)

    var frameStart time.Time
    var elapsedTime float32
    var running bool = true
    var paused bool = false
    var visRad int32 = 10     //inital visible grid
    for running {
        frameStart = time.Now()


        if !paused {
            paused = moveAnt(grid, &antpos, &antdir)
            fmt.Println("Position: ", antpos, "direction: ", antdir)


            drawGrid(renderer, grid, int32(gridDim), int32(gridDim), visRad);
            drawAnt(renderer, int32(gridDim), int32(gridDim), visRad, antpos);
            renderer.Present()

            //calculate ant distance from center
            dx := IntegerAbs(int(gridDim)/2 - antpos.x)
            dy := IntegerAbs(int(gridDim)/2 - antpos.y)
            if dx+2 > int(visRad) || dy+2 > int(visRad) {   //Added +2 so it will zoom out sooner
                visRad += 5;
            }

            //printGrid(grid)
        }

        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch t := event.(type) {
                case *sdl.QuitEvent:
                    return
                case *sdl.KeyboardEvent:
                    if(t.Type == sdl.KEYDOWN) {
                        fmt.Println("key pressed: ", t.Keysym.Scancode)
                        if uint16(t.Keysym.Scancode) == keycodes.KeyEscape{
                            fmt.Println("escape pressed, exiting...")
                            running = false
                        }
                    }
            }
        }


    elapsedTime = float32(time.Since(frameStart).Seconds())
    if elapsedTime < .005 {
        sdl.Delay(5 - uint32(elapsedTime*1000))
        elapsedTime = float32(time.Since(frameStart).Seconds())
    }

    fmt.Println("framerate: ", 1/elapsedTime)

    sdl.Delay(stepTime)

    }



}
