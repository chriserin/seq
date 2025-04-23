package themes

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func Darken(lgColor lipgloss.Color, percent float64) lipgloss.Color {
	return Lighten(lgColor, -1.0*percent)
}
func Lighten(lgColor lipgloss.Color, percent float64) lipgloss.Color {
	hexColor := string(lgColor)

	// Parse the hex color into RGB components
	r, g, b, err := parseHexColor(hexColor)
	if err != nil {
		fmt.Println(err)
		panic("can't parse color")
	}

	// Calculate the factor to increase the brightness
	factor := 1 + percent/100.0

	// Increase each component by the factor and ensure it does not exceed 255
	r = int(float64(r) * factor)
	g = int(float64(g) * factor)
	b = int(float64(b) * factor)

	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}

	// Convert the adjusted RGB values back to a hex string
	return lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", r, g, b))
}

func parseHexColor(hexColor string) (int, int, int, error) {
	r, err := strconv.ParseInt(hexColor[1:3], 16, 16)
	if err != nil {
		return 0, 0, 0, err
	}
	g, err := strconv.ParseInt(hexColor[3:5], 16, 16)
	if err != nil {
		return 0, 0, 0, err
	}
	b, err := strconv.ParseInt(hexColor[5:7], 16, 16)
	if err != nil {
		return 0, 0, 0, err
	}

	return int(r), int(g), int(b), nil
}
