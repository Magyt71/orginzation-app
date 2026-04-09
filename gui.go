package main

import (
	"fmt"
	"image"
	"log"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sqweek/dialog"
)

var (
	deleteButtons []widget.Clickable
	StartStopBtn  widget.Clickable
	AddFolderBtn  widget.Clickable
	moveList      widget.List
	pathList      widget.List
)

func GuiLoop() {
	for range showGuiCh {
		OpenWindow()
	}
}

func OpenWindow() {
	W := new(app.Window)
	err := Run(W)
	if err != nil {
		log.Fatal(err)
	}
}

func Run(Window *app.Window) error {
	theme := material.NewTheme()

	// تخصيص ألوان الـ Material Theme بناءً على طلبك
	theme.Palette.Fg = TextColor
	theme.Palette.ContrastFg = TextColor
	theme.Palette.ContrastBg = PrimaryColor

	var ops op.Ops

	// تحديث الواجهة عند حدوث تغيير في المحرك
	Org.logCallback = func(msg string) {
		Window.Invalidate()
	}

	for {
		switch e := Window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// تلوين الخلفية الكلية
			paint.FillShape(gtx.Ops, BackgroundColor, clip.Rect{Max: gtx.Constraints.Max}.Op())

			// منطق الأزرار
			handleEvents(gtx, Window)

			layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				// القسم الجانبي
				layout.Flexed(0.35, func(gtx layout.Context) layout.Dimensions {
					return SideBar(gtx, theme)
				}),

				// فاصل عمودي بلون الـ InputBgColor
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					size := image.Pt(gtx.Dp(unit.Dp(2)), gtx.Constraints.Max.Y)
					paint.FillShape(gtx.Ops, InputBgColor, clip.Rect{Max: size}.Op())
					return layout.Dimensions{Size: size}
				}),

				// قسم النشاطات
				layout.Flexed(0.65, func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return ActivityLogs(gtx, theme)
					})
				}),
			)

			e.Frame(gtx.Ops)
		}
	}
}

func handleEvents(gtx layout.Context, w *app.Window) {
	if StartStopBtn.Clicked(gtx) {
		if !Org.Config.IsRunning {
			Org.Start()
		} else {
			Org.Stop()
		}
	}
	if AddFolderBtn.Clicked(gtx) {
		go func() {
			dir, err := dialog.Directory().Title("Select Folder").Browse()
			if err == nil && dir != "" {
				Org.AddPath(dir)
				w.Invalidate()
			}
		}()
	}
}

func SideBar(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.UniformInset(unit.Dp(20)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			// زر التشغيل/الإيقاف مع تغيير اللون حسب الحالة
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btnText := "START ENGINE"
				bgColor := SuccessColor
				if Org.Config.IsRunning {
					btnText = "STOP ENGINE"
					bgColor = ErrorColor
				}
				btn := material.Button(theme, &StartStopBtn, btnText)
				btn.Background = bgColor
				btn.TextSize = unit.Sp(14)
				return btn.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),

			// زر إضافة مجلد باللون Teal (ButtonColor)
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(theme, &AddFolderBtn, "ADD FOLDER")
				btn.Background = ButtonColor
				return btn.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(30)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(theme, "WATCHING DIRECTORIES")
				lbl.Color = PrimaryColor // استعمال السماوي للعناوين
				return lbl.Layout(gtx)
			}),

			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return PathListUI(gtx, theme)
			}),
		)
	})
}

func PathListUI(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	pathList.Axis = layout.Vertical
	paths := Org.Config.WatchPaths

	return material.List(theme, &pathList).Layout(gtx, len(paths), func(gtx layout.Context, i int) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			// محاكاة حقل إدخال لعرض المسار
			return layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					paint.FillShape(gtx.Ops, InputBgColor, clip.UniformRRect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Dp(unit.Dp(30))), 4).Op(gtx.Ops))
					return layout.Dimensions{}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						txt := material.Caption(theme, paths[i])
						txt.Color = SecondaryTextColor
						return txt.Layout(gtx)
					})
				}),
			)
		})
	})
}

func ActivityLogs(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	moveList.Axis = layout.Vertical
	moves := Org.RecentMoves

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H6(theme, "LIVE ACTIVITY")
			title.Color = PrimaryColor
			return title.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			if len(moves) == 0 {
				txt := material.Body2(theme, "Waiting for file events...")
				txt.Color = SecondaryTextColor
				return txt.Layout(gtx)
			}
			return material.List(theme, &moveList).Layout(gtx, len(moves), func(gtx layout.Context, i int) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Stack{Alignment: layout.Center}.Layout(gtx,
						// خلفية الكارت (InputBgColor)
						layout.Expanded(func(gtx layout.Context) layout.Dimensions {
							paint.FillShape(gtx.Ops, InputBgColor, clip.UniformRRect(image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Dp(unit.Dp(45))), 8).Op(gtx.Ops))
							return layout.Dimensions{}
						}),
						layout.Stacked(func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
									layout.Rigid(func(gtx layout.Context) layout.Dimensions {
										t := material.Caption(theme, moves[i].Time)
										t.Color = PrimaryColor
										return t.Layout(gtx)
									}),
									layout.Rigid(layout.Spacer{Width: unit.Dp(15)}.Layout),
									layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
										txt := fmt.Sprintf("%s  →  %s", moves[i].FileName, moves[i].Dest)
										l := material.Body1(theme, txt)
										l.TextSize = unit.Sp(13)
										return l.Layout(gtx)
									}),
								)
							})
						}),
					)
				})
			})
		}),
	)
}
