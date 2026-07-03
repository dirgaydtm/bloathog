package components

import "github.com/guptarohit/asciigraph"

var steppedBounds = []float64{100, 250, 500, 1024, 2048, 4096, 8192, 16384, 32768, 65536}

func getSteppedUpperBound(maxVal float64) float64 {
	for _, step := range steppedBounds {
		if maxVal <= step {
			return step
		}
	}
	return maxVal * 1.2
}

// RenderGraph renders the RAM or CPU usage graph using asciigraph.
// activeGraph: 0 for RAM, 1 for CPU
func RenderGraph(data []float64, width, height int, activeGraph int) string {
	if len(data) == 0 {
		data = []float64{0}
	}
	if width < 10 {
		width = 10
	}
	if height < 3 {
		height = 3
	}

	// asciigraph adds 1 line for the x-axis (since caption was removed).
	plotHeight := height - 1
	if plotHeight < 1 {
		plotHeight = 1
	}

	maxVal := data[0]
	for _, v := range data {
		if v > maxVal {
			maxVal = v
		}
	}

	// asciigraph width is the number of data columns, not total string width.
	// The Y-axis labels take ~8 chars. Reserve them.
	plotWidth := width - 10
	if plotWidth < 4 {
		plotWidth = 4
	}

	color := asciigraph.LightBlue
	if activeGraph == 1 {
		color = asciigraph.LightYellow
	}

	return asciigraph.Plot(data,
		asciigraph.Height(plotHeight),
		asciigraph.Width(plotWidth),
		asciigraph.SeriesColors(color),
		asciigraph.LowerBound(0),
		asciigraph.UpperBound(getSteppedUpperBound(maxVal)),
	)
}
