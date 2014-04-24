package qstr

import "testing"

func TestHSL(t *testing.T) {
    rgbcolors := []RGBColor {
        {0, 0, 0},
        {255, 0, 0},
        {0, 255, 0},
        {0, 0, 255},
        {0, 255, 255},
        {255, 255, 0},
        {255, 0, 255},
        {255, 255, 255},
    }

    hslcolors := []HSLColor {
        {0.0, 0.0, 0.0},
        {0, 1, 0.5},
        {0.3333333333333333, 1, 0.5},
        {0.6666666666666666, 1, 0.5},
        {0.5, 1, 0.5},
        {0.16666666666666666, 1, 0.5},
        {0.8333333333333334, 1, 0.5},
        {0, 0, 1},
    }

    for i, v := range rgbcolors {
        expected := hslcolors[i]
        received := v.HSL()

        if received.H != expected.H || 
           received.S != expected.S || 
           received.L != expected.L {
           t.Errorf("Incorrect HSL translation for RGB color %v. Expected: %v, Got: %v.", v, expected, received)
        }
    }
}
