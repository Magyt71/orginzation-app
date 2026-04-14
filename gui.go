package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"os"
	"runtime"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/sqweek/dialog"
)

var (
	deleteButtons []widget.Clickable
	StartStopBtn  widget.Clickable
	AddFolderBtn  widget.Clickable
	SettingsBtn   widget.Clickable
	moveList      widget.List
	pathList      widget.List
	appImageOp    paint.ImageOp
	imageLoaded   bool
	showSettings  bool
)

func init() {
	loadAppImage()
}

func loadAppImage() {
	f, err := os.Open("image/lobster.png")
	if err != nil {
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return
	}

	appImageOp = paint.NewImageOp(img)
	imageLoaded = true
}

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
	theme.Palette.Fg = TextColor
	theme.Palette.ContrastFg = TextColor
	theme.Palette.ContrastBg = PrimaryColor

	var ops op.Ops

	Org.logCallback = func(msg string) {
		Window.Invalidate()
	}

	for {
		switch e := Window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Full-window background
			fillShape(gtx.Ops, BackgroundColor, gtx.Constraints.Max)

			// Button logic
			handleEvents(gtx, Window)

			// Root layout
			drawRoot(gtx, theme)

			e.Frame(gtx.Ops)
		}
	}
}

// ── handleEvents ─────────────────────────────────────────────────────────────

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
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			d := dialog.Directory().Title("Select Folder to Watch")
			Dir, err := d.Browse()
			if err == nil && Dir != "" {
				Org.AddPath(Dir)
				w.Invalidate()
			}
		}()
	}
	if SettingsBtn.Clicked(gtx) {
		showSettings = !showSettings
	}
}

// ── Root Layout ──────────────────────────────────────────────────────────────
//
//	┌────────────────────┬──────────────────────┐
//	│                    │ Organization         │
//	│     Image          │ Moderator            │
//	│     Area           │  [Start]             │
//	│                    │  [Add Folder]        │
//	│                    │  [Settings]          │
//	│                    ├──────────────────────┤
//	│                    │   Last move files    │
//	├────────────────────┤                      │
//	│ Status bar         │                      │
//	└────────────────────┴──────────────────────┘

func drawRoot(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// Outer bordered container
		return drawBorderedPanel(gtx, BorderColor, unit.Dp(2), unit.Dp(10), BackgroundColor,
			func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,

						// ── Left column: Image + Status bar (45%)
						layout.Flexed(0.45, func(gtx layout.Context) layout.Dimensions {
							return drawLeftColumn(gtx, theme)
						}),

						// Spacer between columns
						layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),

						// ── Right column: Title + Buttons + Last moves (55%)
						layout.Flexed(0.55, func(gtx layout.Context) layout.Dimensions {
							return drawRightColumn(gtx, theme)
						}),
					)
				})
			},
		)
	})
}

// ── Left Column ──────────────────────────────────────────────────────────────

func drawLeftColumn(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

		// Image area (takes remaining space)
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return drawBorderedPanel(gtx, BorderColor, unit.Dp(2), unit.Dp(8), SurfaceColor,
				func(gtx layout.Context) layout.Dimensions {
					return drawImageArea(gtx, theme)
				},
			)
		}),

		layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),

		// Status bar at bottom
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return drawBorderedPanel(gtx, BorderColor, unit.Dp(2), unit.Dp(6), InputBgColor,
				func(gtx layout.Context) layout.Dimensions {
					return drawStatusBar(gtx, theme)
				},
			)
		}),
	)
}

// ── Right Column ─────────────────────────────────────────────────────────────

func drawRightColumn(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

		// Top section: Title + Buttons
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,

				// Title
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Top: unit.Dp(16), Bottom: unit.Dp(20)}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							lbl := material.H5(theme, "Organization Moderator")
							lbl.Color = PrimaryColor
							lbl.Alignment = text.Middle
							return lbl.Layout(gtx)
						})
				}),

				// Start / Stop button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					btnText := "▶  Start"
					bgColor := SuccessColor
					if Org.Config.IsRunning {
						btnText = "■  Stop"
						bgColor = ErrorColor
					}
					return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return styledButton(gtx, theme, &StartStopBtn, btnText, bgColor)
						})
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

				// Add Folder button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20)}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return styledButton(gtx, theme, &AddFolderBtn, "＋  Add Folder", ButtonColor)
						})
				}),

				layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

				// Settings button
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					settingsBg := SettingsBtnColor
					if showSettings {
						settingsBg = PrimaryColor
					}
					return layout.Inset{Left: unit.Dp(20), Right: unit.Dp(20), Bottom: unit.Dp(16)}.Layout(gtx,
						func(gtx layout.Context) layout.Dimensions {
							return styledButton(gtx, theme, &SettingsBtn, "⚙  Settings", settingsBg)
						})
				}),
			)
		})
	}),
		// ── Bottom section: Last move files or Settings panel
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return drawBorderedPanel(gtx, BorderColor, unit.Dp(2), unit.Dp(8), SurfaceColor,
				func(gtx layout.Context) layout.Dimensions {
					if showSettings {
						return drawSettingsPanel(gtx, theme)
					}
					return drawLastMoveFiles(gtx, theme)
				},
			)
		}),
	)
}

// ── Image Area ───────────────────────────────────────────────────────────────

func drawImageArea(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	size := gtx.Constraints.Max

	if imageLoaded {
		imgSz := appImageOp.Size()

		// Calculate scale to fit while maintaining aspect ratio
		scaleX := float32(size.X) / float32(imgSz.X)
		scaleY := float32(size.Y) / float32(imgSz.Y)
		scale := scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		// Calculate centered position
		scaledW := int(float32(imgSz.X) * scale)
		scaledH := int(float32(imgSz.Y) * scale)
		offsetX := (size.X - scaledW) / 2
		offsetY := (size.Y - scaledH) / 2

		// Clip to container bounds
		clipStack := clip.Rect{Max: size}.Push(gtx.Ops)

		// Offset to center, then scale
		offStack := op.Offset(image.Pt(offsetX, offsetY)).Push(gtx.Ops)
		aff := f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(scale, scale))
		affStack := op.Affine(aff).Push(gtx.Ops)

		appImageOp.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		affStack.Pop()
		offStack.Pop()
		clipStack.Pop()
	} else {
		// Placeholder
		layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			lbl := material.H3(theme, "Image")
			lbl.Color = SecondaryTextColor
			lbl.Alignment = text.Middle
			return lbl.Layout(gtx)
		})
	}

	return layout.Dimensions{Size: size}
}

// ── Last Move Files Panel ────────────────────────────────────────────────────

func drawLastMoveFiles(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	moveList.Axis = layout.Vertical
	moves := Org.RecentMoves

	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

			// Section header
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							dot := PrimaryColor
							if len(moves) == 0 {
								dot = SecondaryTextColor
							}
							return drawDot(gtx, dot, 8)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							title := material.H6(theme, "Last move files")
							title.Color = PrimaryColor
							return title.Layout(gtx)
						}),
					)
				})
			}),

			// Separator
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return drawHLine(gtx, BorderColor)
				})
			}),

			// List or empty state
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if len(moves) == 0 {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body1(theme, "◌")
								lbl.Color = color.NRGBA{R: 40, G: 60, B: 70, A: 255}
								lbl.TextSize = unit.Sp(36)
								lbl.Alignment = text.Middle
								return lbl.Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								lbl := material.Body2(theme, "Waiting for file events…")
								lbl.Color = SecondaryTextColor
								lbl.Alignment = text.Middle
								return lbl.Layout(gtx)
							}),
						)
					})
				}

				return material.List(theme, &moveList).Layout(gtx, len(moves), func(gtx layout.Context, i int) layout.Dimensions {
					return layout.Inset{Bottom: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return drawActivityCard(gtx, theme, moves[i])
					})
				})
			}),
		)
	})
}

// ── Settings Panel (Watched Directories) ─────────────────────────────────────

func drawSettingsPanel(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	pathList.Axis = layout.Vertical
	paths := Org.Config.WatchPaths

	if len(deleteButtons) != len(paths) {
		deleteButtons = make([]widget.Clickable, len(paths))
	}

	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,

			// Section header
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return drawDot(gtx, PrimaryColor, 8)
						}),
						layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							title := material.H6(theme, "Watching Directories")
							title.Color = PrimaryColor
							return title.Layout(gtx)
						}),
					)
				})
			}),

			// Separator
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Bottom: unit.Dp(10)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return drawHLine(gtx, BorderColor)
				})
			}),

			// Path list
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if len(paths) == 0 {
					return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Caption(theme, "No folders added yet")
						lbl.Color = SecondaryTextColor
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					})
				}

				return material.List(theme, &pathList).Layout(gtx, len(paths), func(gtx layout.Context, i int) layout.Dimensions {
					return layout.Inset{Bottom: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {

						// Delete logic
						if deleteButtons[i].Clicked(gtx) {
							Org.Config.mu.Lock()
							Org.Config.WatchPaths = append(Org.Config.WatchPaths[:i], Org.Config.WatchPaths[i+1:]...)
							Org.Config.mu.Unlock()
						}

						return drawPathCard(gtx, theme, paths[i], &deleteButtons[i])
					})
				})
			}),
		)
	})
}

// ── Activity Card ────────────────────────────────────────────────────────────

func drawActivityCard(gtx layout.Context, theme *material.Theme, move MoveRecord) layout.Dimensions {
	return layout.Stack{Alignment: layout.W}.Layout(gtx,
		// Card background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := image.Pt(gtx.Constraints.Min.X, gtx.Constraints.Min.Y)
			rr := gtx.Dp(unit.Dp(6))
			defer clip.RRect{
				Rect: image.Rectangle{Max: size},
				SE:   rr, SW: rr, NE: rr, NW: rr,
			}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: InputBgColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// Top accent strip
			stripH := gtx.Dp(unit.Dp(2))
			stripRect := image.Rectangle{Max: image.Pt(size.X, stripH)}
			defer clip.RRect{Rect: stripRect, NE: rr, NW: rr}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: BorderColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: size}
		}),

		// Content
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(8), Bottom: unit.Dp(8),
				Left: unit.Dp(12), Right: unit.Dp(12),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
							// Timestamp
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								t := material.Caption(theme, move.Time)
								t.Color = PrimaryColor
								return t.Layout(gtx)
							}),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(gtx,
									func(gtx layout.Context) layout.Dimensions {
										return drawDot(gtx, DividerColor, 4)
									})
							}),
							// File name
							layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
								txt := fmt.Sprintf("%s → %s", move.FileName, move.Dest)
								l := material.Caption(theme, txt)
								l.Color = TextColor
								return l.Layout(gtx)
							}),
						)
					}),
				)
			})
		}),
	)
}

// ── Path Card ────────────────────────────────────────────────────────────────

func drawPathCard(gtx layout.Context, theme *material.Theme, path string, delBtn *widget.Clickable) layout.Dimensions {
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		// Card background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := image.Pt(gtx.Constraints.Min.X, gtx.Constraints.Min.Y)
			rr := gtx.Dp(unit.Dp(6))
			defer clip.RRect{
				Rect: image.Rectangle{Max: size},
				SE:   rr, SW: rr, NE: rr, NW: rr,
			}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: InputBgColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			// Left cyan strip
			strip := image.Rectangle{Max: image.Pt(gtx.Dp(unit.Dp(3)), size.Y)}
			defer clip.RRect{Rect: strip, NW: rr, SW: rr}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: PrimaryColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: size}
		}),

		// Content
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Inset{
				Top: unit.Dp(8), Bottom: unit.Dp(8),
				Left: unit.Dp(12), Right: unit.Dp(8),
			}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
					Spacing:   layout.SpaceBetween,
				}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						txt := material.Caption(theme, path)
						txt.Color = SecondaryTextColor
						return txt.Layout(gtx)
					}),
					layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(theme, delBtn, "✕")
						btn.Background = DeleteBtnBg
						btn.Color = ErrorColor
						btn.TextSize = unit.Sp(10)
						btn.Inset = layout.Inset{Top: 2, Bottom: 2, Left: 6, Right: 6}
						return btn.Layout(gtx)
					}),
				)
			})
		}),
	)
}

// ── Styled button with left accent bar ───────────────────────────────────────

func styledButton(gtx layout.Context, theme *material.Theme, btn *widget.Clickable, label string, bg color.NRGBA) layout.Dimensions {
	return layout.Stack{}.Layout(gtx,
		// Background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rr := gtx.Dp(unit.Dp(7))
			rect := image.Rectangle{Max: gtx.Constraints.Min}
			bgColor := bg
			if btn.Hovered() {
				bgColor = addAlpha(bg, 220)
			}
			defer clip.RRect{Rect: rect, SE: rr, SW: rr, NE: rr, NW: rr}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: bgColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),

		// Left accent bar
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			bw := gtx.Dp(unit.Dp(4))
			rr := gtx.Dp(unit.Dp(7))
			rect := image.Rectangle{Max: image.Pt(bw, gtx.Constraints.Min.Y)}
			defer clip.RRect{Rect: rect, NW: rr, SW: rr}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: color.NRGBA{R: 255, G: 255, B: 255, A: 50}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Min}
		}),

		// Label
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return btn.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top: unit.Dp(11), Bottom: unit.Dp(11),
					Left: unit.Dp(16), Right: unit.Dp(12),
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body1(theme, label)
					lbl.Color = TextColor
					lbl.TextSize = unit.Sp(13)
					lbl.Alignment = text.Middle
					return lbl.Layout(gtx)
				})
			})
		}),
	)
}

// ── Status Bar ───────────────────────────────────────────────────────────────

func drawStatusBar(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	status := "Idle"
	statusColor := SecondaryTextColor
	if Org.Config.IsRunning {
		status = "Running — watching for changes"
		statusColor = SuccessColor
	}
	return layout.Inset{
		Top: unit.Dp(8), Bottom: unit.Dp(8),
		Left: unit.Dp(12), Right: unit.Dp(12),
	}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		dims := layout.Flex{Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return drawDot(gtx, statusColor, 6)
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.Caption(theme, "The app is : "+status)
				lbl.Color = statusColor
				return lbl.Layout(gtx)
			}),
		)
		// 🔴 الجوهر: نُجبر شريط الحالة على ملء المساحة الأفقية بالكامل (نفس تكتيك الصورة)
		dims.Size.X = gtx.Constraints.Max.X
		return dims
	})
}

// ── Bordered Panel Helper ────────────────────────────────────────────────────

func drawBorderedPanel(gtx layout.Context, borderColor color.NRGBA, borderWidth unit.Dp, radius unit.Dp, fillColor color.NRGBA, w layout.Widget) layout.Dimensions {
	bw := gtx.Dp(borderWidth)
	rr := gtx.Dp(radius)

	return layout.Stack{}.Layout(gtx,
		// Border background
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min
			outerRect := image.Rectangle{Max: size}
			defer clip.RRect{Rect: outerRect, SE: rr, SW: rr, NE: rr, NW: rr}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: borderColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: size}
		}),

		// Inner fill
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min
			innerRR := rr - bw
			if innerRR < 0 {
				innerRR = 0
			}
			innerRect := image.Rectangle{
				Min: image.Pt(bw, bw),
				Max: image.Pt(size.X-bw, size.Y-bw),
			}
			defer clip.RRect{Rect: innerRect, SE: innerRR, SW: innerRR, NE: innerRR, NW: innerRR}.Push(gtx.Ops).Pop()
			paint.ColorOp{Color: fillColor}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			return layout.Dimensions{Size: size}
		}),

		// Content (inset by border width)
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(borderWidth).Layout(gtx, w)
		}),
	)
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func fillShape(ops *op.Ops, c color.NRGBA, size image.Point) {
	defer clip.Rect{Max: size}.Push(ops).Pop()
	paint.ColorOp{Color: c}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

func drawHLine(gtx layout.Context, c color.NRGBA) layout.Dimensions {
	h := gtx.Dp(unit.Dp(1))
	rect := image.Rectangle{Max: image.Pt(gtx.Constraints.Max.X, h)}
	defer clip.Rect(rect).Push(gtx.Ops).Pop()
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: rect.Max}
}

func drawDot(gtx layout.Context, c color.NRGBA, dp unit.Dp) layout.Dimensions {
	sz := gtx.Dp(dp)
	dot := image.Rectangle{Max: image.Pt(sz, sz)}
	defer clip.Ellipse{Max: dot.Max}.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return layout.Dimensions{Size: image.Pt(sz+4, sz)}
}

// addAlpha returns the same color with a new alpha (used for hover effect)
func addAlpha(c color.NRGBA, a uint8) color.NRGBA {
	return color.NRGBA{R: c.R, G: c.G, B: c.B, A: a}
}
