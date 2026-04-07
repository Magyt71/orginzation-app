package main

import "image/color"

// 1. ألوان الثيم الأساسية
var (
	// لون الخلفية الداكنة (Dark Background) - #14191c
	BackgroundColor = color.NRGBA{R: 20, G: 25, B: 28, A: 255}

	// اللون الأساسي السماوي (Primary Cyan) - #00bcd4
	PrimaryColor = color.NRGBA{R: 0, G: 188, B: 212, A: 255}

	// لون النصوص الأساسية (White Text) - #ffffff
	TextColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	// لون النصوص الثانوية (Gray Text) - #969696
	SecondaryTextColor = color.NRGBA{R: 150, G: 150, B: 150, A: 255}
)

// 2. ألوان الحالات والعناصر (الموجودة أسفل الصورة)
var (
	// لون الأزرار العادية (Teal) - #009688
	ButtonColor = color.NRGBA{R: 0, G: 150, B: 136, A: 255}

	// لون خلفية حقل الإدخال (Input Background) - #1a2428
	InputBgColor = color.NRGBA{R: 26, G: 36, B: 40, A: 255}

	// لون النجاح (Success Green) - #00e676
	SuccessColor = color.NRGBA{R: 0, G: 230, B: 118, A: 255}

	// لون الخطأ (Error Red) - #ff5252
	ErrorColor = color.NRGBA{R: 255, G: 82, B: 82, A: 255}

	// لون التحذير (Warning Yellow) - #ffc107
	WarningColor = color.NRGBA{R: 255, G: 193, B: 7, A: 255}
)
