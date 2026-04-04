package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// DarkCyanTheme - ثيم داكن بلون Cyan
type DarkCyanTheme struct{}

var _ fyne.Theme = (*DarkCyanTheme)(nil)

func (t *DarkCyanTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {

	// ألوان الخلفية
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x0A, G: 0x12, B: 0x14, A: 0xFF} // خلفية داكنة جداً
	case theme.ColorNameMenuBackground:
		return color.NRGBA{R: 0x0D, G: 0x1A, B: 0x1E, A: 0xFF}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 0x0A, G: 0x12, B: 0x14, A: 0xEE}
	case theme.ColorNameButton:
		return color.NRGBA{R: 0x00, G: 0x7A, B: 0x8A, A: 0xFF} // Cyan داكن
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x1A, G: 0x2E, B: 0x32, A: 0xFF}

	// ألوان النص
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xB2, G: 0xEF, B: 0xF4, A: 0xFF} // Cyan فاتح للنص
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x3A, G: 0x5A, B: 0x60, A: 0xFF}
	case theme.ColorNamePlaceHolder:
		return color.NRGBA{R: 0x4A, G: 0x7A, B: 0x82, A: 0xFF}

	// اللون الأساسي (Primary / Accent)
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x00, G: 0xBF, B: 0xD8, A: 0xFF} // Cyan مشرق

	// ألوان التركيز والتحديد
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0x00, G: 0xBF, B: 0xD8, A: 0xAA}
	case theme.ColorNameSelection:
		return color.NRGBA{R: 0x00, G: 0x9A, B: 0xB0, A: 0x55}
	case theme.ColorNameHover:
		return color.NRGBA{R: 0x00, G: 0xBF, B: 0xD8, A: 0x22}
	case theme.ColorNamePressed:
		return color.NRGBA{R: 0x00, G: 0x8A, B: 0xA0, A: 0x44}

	// ألوان الحالة
	case theme.ColorNameSuccess:
		return color.NRGBA{R: 0x00, G: 0xD4, B: 0x9A, A: 0xFF}
	case theme.ColorNameWarning:
		return color.NRGBA{R: 0xFF, G: 0xC1, B: 0x07, A: 0xFF}
	case theme.ColorNameError:
		return color.NRGBA{R: 0xFF, G: 0x45, B: 0x5A, A: 0xFF}

	// ألوان الإدخال والحدود
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 0x0D, G: 0x1E, B: 0x22, A: 0xFF}
	case theme.ColorNameInputBorder:
		return color.NRGBA{R: 0x00, G: 0x6A, B: 0x7A, A: 0xFF}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 0x00, G: 0x4A, B: 0x58, A: 0xFF}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x00, G: 0x7A, B: 0x8A, A: 0x88}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x88}
	case theme.ColorNameHeaderBackground:
		return color.NRGBA{R: 0x05, G: 0x14, B: 0x18, A: 0xFF}
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (t *DarkCyanTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *DarkCyanTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *DarkCyanTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameInnerPadding:
		return 8
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 24
	case theme.SizeNameSubHeadingText:
		return 18
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInputRadius:
		return 6
	case theme.SizeNameSelectionRadius:
		return 4
	case theme.SizeNameScrollBar:
		return 4
	case theme.SizeNameScrollBarRadius:
		return 2
	case theme.SizeNameScrollBarSmall:
		return 2
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameInlineIcon:
		return 20
	}
	return theme.DefaultTheme().Size(name)
}
