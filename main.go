package main

import (
	"flag"
	"github.com/hajimehoshi/ebiten"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

type Grid struct {
	grid []float64
	x, y int

	nt      int
	eqSteps int
	mcSteps int

	T []float64

	counter float64
}

func (g *Grid) Draw(pixels []byte) {
	for i, f := range g.grid {
		if f == 1 {
			pixels[4*i+0] = 0xff
			pixels[4*i+1] = 0xff
			pixels[4*i+2] = 0xff
		} else if f == -1 {
			pixels[4*i+0] = 0x00
			pixels[4*i+1] = 0x00
			pixels[4*i+2] = 0x00
		} else {
			pixels[4*i+0] = 0xff
			pixels[4*i+1] = 0x00
			pixels[4*i+2] = 0xff
		}
	}
}

func (g *Grid) Update(beta float64) {
	for i := 0; i < g.x; i++ {
		for j := 0; j < g.y; j++ {
			a := rand.Int() % g.x
			b := rand.Int() % g.y
			s := g.grid[b*g.x+a]
			nb := g.grid[b*g.x+FW(a+1, g.x)] + g.grid[b*g.x+FW(a-1, g.x)] + g.grid[FW(b+1, g.y)*g.x+a] + g.grid[FW(b-1, g.y)*g.x+a]
			cost := 2 * s * nb

			if cost < 0 {
				s *= -1
			} else if rand.Float64() < math.Exp(-cost*beta) {
				s *= -1
			}

			g.grid[b*g.x+a] = s
		}
	}
}

type Game struct {
	pixels []byte
	grid   Grid

	x, y     int
	winScale int
}

func (g *Game) Update(_ *ebiten.Image) error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.grid.counter += 0.01
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.grid.counter > 0.01 {
		g.grid.counter -= 0.01
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) || (ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyC)) {
		os.Exit(1)
	}

	iT := 1.0 / g.grid.counter
	g.grid.Update(iT)

	return nil
}

func (g *Game) Layout(_, _ int) (int, int) {
	return g.x, g.y
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.pixels == nil {
		g.pixels = make([]byte, g.x*g.y*4)

		for i := 0; i < g.x*g.y; i++ {
			g.pixels[4*i+3] = 0xff
		}
	}

	g.grid.Draw(g.pixels)

	err := screen.ReplacePixels(g.pixels)

	if err != nil {
		return
	}
}

func FW(a, b int) int {
	if a < 0 {
		return a + b
	} else if a >= b {
		return a - b
	}
	return a
}

func main() {
	var (
		xs int
		ys int
		ss int
		fs bool
	)

	flag.IntVar(&xs, "x", 1920/8, "X resolution")
	flag.IntVar(&ys, "y", 1920/8, "Y resolution")
	flag.IntVar(&ss, "s", 4, "Scale (only works if not fullscreen)")
	flag.BoolVar(&fs, "f", false, "Fullscreen")

	rand.Seed(time.Now().UnixNano())

	g := &Game{
		x:        xs,
		y:        ys,
		winScale: ss,
	}

	g.grid = Grid{
		x: g.x,
		y: g.y,

		grid: make([]float64, g.x*g.y),

		nt:      2048,
		eqSteps: 1 << 8,
		mcSteps: 1 << 9,

		counter: rand.Float64() * 3,
	}

	g.grid.T = make([]float64, g.grid.nt)

	for i := 0; i < g.grid.nt; i++ {
		g.grid.T[i] = (float64(g.grid.nt-i-1)/float64(g.grid.nt))*1.75 + 1.53
	}

	for i := range g.grid.grid {
		if rand.Int()%2 == 1 {
			g.grid.grid[i] = 1
		} else {
			g.grid.grid[i] = -1
		}
	}

	ebiten.SetWindowSize(g.x*g.winScale, g.y*g.winScale)
	ebiten.SetWindowTitle("Ising Model")
	ebiten.SetFullscreen(fs)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
