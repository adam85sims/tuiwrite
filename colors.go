package main

import "fmt"

// Catppuccin Mocha color palette
// https://github.com/catppuccin/catppuccin

// Color represents an RGB color value
type Color struct {
	R, G, B uint8
}

// Catppuccin Mocha colors
var (
	// Accent colors
	Rosewater = Color{245, 224, 220} // #f5e0dc
	Flamingo  = Color{242, 205, 205} // #f2cdcd
	Pink      = Color{245, 194, 231} // #f5c2e7
	Mauve     = Color{203, 166, 247} // #cba6f7
	Red       = Color{243, 139, 168} // #f38ba8
	Maroon    = Color{235, 160, 172} // #eba0ac
	Peach     = Color{250, 179, 135} // #fab387
	Yellow    = Color{249, 226, 175} // #f9e2af
	Green     = Color{166, 227, 161} // #a6e3a1
	Teal      = Color{148, 226, 213} // #94e2d5
	Sky       = Color{137, 220, 235} // #89dceb
	Sapphire  = Color{116, 199, 236} // #74c7ec
	Blue      = Color{137, 180, 250} // #89b4fa
	Lavender  = Color{180, 190, 254} // #b4befe

	// Monochromatic colors
	Text     = Color{205, 214, 244} // #cdd6f4
	Subtext1 = Color{186, 194, 222} // #bac2de
	Subtext0 = Color{166, 173, 200} // #a6adc8
	Overlay2 = Color{147, 153, 178} // #9399b2
	Overlay1 = Color{127, 132, 156} // #7f849c
	Overlay0 = Color{108, 112, 134} // #6c7086
	Surface2 = Color{88, 91, 112}   // #585b70
	Surface1 = Color{69, 71, 90}    // #45475a
	Surface0 = Color{49, 50, 68}    // #313244
	Base     = Color{30, 30, 46}    // #1e1e2e
	Mantle   = Color{24, 24, 37}    // #181825
	Crust    = Color{17, 17, 27}    // #11111b
)

// ColorToANSI converts an RGB color to an ANSI escape sequence for 24-bit color
func ColorToANSI(c Color, foreground bool) string {
	if foreground {
		return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", c.R, c.G, c.B)
	}
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm", c.R, c.G, c.B)
}
