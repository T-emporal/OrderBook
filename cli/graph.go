package cli

import (
	"fmt"
	"image/color"
	"log"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func PlotGraph() {

	err := plotData("out.png")
	if err != nil {
		log.Fatalf("could not plot data: %v", err)
	}

}

func plotData(path string) error {

	p := plot.New()

	p.X.Label.Text = "Duration"
	p.Y.Label.Text = "Price"
	p.Add(plotter.NewGrid())

	s, err := plotter.NewScatter(_xys)
	if err != nil {
		return fmt.Errorf("could not create scatter: %v", err)
	}
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}

	p.Add(s)

	if err := p.Save(6*vg.Inch, 6*vg.Inch, path); err != nil {
		panic(err)
	}
	return nil

}
