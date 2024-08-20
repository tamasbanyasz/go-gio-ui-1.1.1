package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		window := new(app.Window)
		window.Option(
			app.Title("Bányász Tamás"),
			app.Size(unit.Dp(640), unit.Dp(480)),
		)
		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

type userInput struct {
	input     string
	addedTime time.Time
}

func run(window *app.Window) error {

	theme := material.NewTheme()

	fixedWidth := unit.Dp(640)
	fixedHeight := unit.Dp(480)

	pageButtonText := "Turn page"

	var ops op.Ops

	var button widget.Clickable
	var turnPage widget.Clickable
	var input widget.Editor

	slider := &widget.Float{}
	var list widget.List
	list.Axis = layout.Vertical

	var text string
	var page int
	buttonClicked := false
	var userInputs []userInput

	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)

		case app.FrameEvent:

			window.Option(app.Size(fixedWidth, fixedHeight))

			gtx := app.NewContext(&ops, e)

			if turnPage.Clicked(gtx) {
				if page == 0 {
					page = 1
					pageButtonText = "Turn back"
				} else {
					pageButtonText = "Turn page"
					page = 0
				}
				fmt.Println("Page turned!")
			}
			layout.Stack{}.Layout(gtx,
				layout.Expanded(func(gtx layout.Context) layout.Dimensions {
					paint.Fill(gtx.Ops, color.NRGBA{R: 0, G: 237, B: 245, A: 255})
					return layout.Dimensions{Size: gtx.Constraints.Max}
				}),
				layout.Stacked(func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:   unit.Dp(430),
						Left:  unit.Dp(520),
						Right: unit.Dp(10),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						btn := material.Button(theme, &turnPage, pageButtonText)
						btn.Background = color.NRGBA{R: 33, G: 80, B: 255, A: 255}

						return btn.Layout(gtx)
					})
				}),
			)

			switch page {
			case 0:

				if button.Clicked(gtx) {
					if input.Text() != "" {
						text = input.Text()
						userInputs = append(userInputs, userInput{input: text, addedTime: time.Now()})
					} else {
						text = ""
					}

					fmt.Println(userInputs)

					slider.Value = 0.2
					fmt.Println(text)
					fmt.Println("Button pressed!")
					input.SetText("")
					buttonClicked = true
				}

				if buttonClicked && text != "" {
					layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{
							Axis:    layout.Vertical,
							Spacing: layout.SpaceEvenly,
						}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								return layout.Inset{
									Top: unit.Dp(400),
								}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									sliderWidth := gtx.Dp(unit.Dp(300))
									sliderHeight := gtx.Dp(unit.Dp(16))

									gtx.Constraints.Min.X = gtx.Dp(unit.Dp(sliderWidth))
									gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(float32(sliderHeight)))
									s := material.Slider(theme, slider)

									return s.Layout(gtx)
								})
							}),
						)
					})
				}

				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:   unit.Dp(40),
						Left:  unit.Dp(50),
						Right: unit.Dp(50),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						editor := material.Editor(theme, &input, "Enter text")
						editor.Color = color.NRGBA{R: 0, G: 100, B: 0, A: 255}
						return editor.Layout(gtx)
					})
				})

				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{
							Top:   unit.Dp(300),
							Left:  unit.Dp(160),
							Right: unit.Dp(160),
						}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return material.Button(theme, &button, "Click Me").Layout(gtx)
						})
					}),
				)

				layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Inset{
						Top:   unit.Dp(-200),
						Left:  unit.Dp(50),
						Right: unit.Dp(50),
					}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.Label(theme, unit.Sp(slider.Value*100), text).Layout(gtx)
					})
				})

			case 1:

				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: layout.SpaceEvenly,
				}.Layout(gtx,

					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return list.Layout(gtx, len(userInputs), func(gtx layout.Context, index int) layout.Dimensions {
							texts := userInputs[index]
							return material.Body1(theme, fmt.Sprintf("Input was: %s\nDate: %s", texts.input,
								texts.addedTime)).Layout(gtx)
						})
					}),
				)
			}

			e.Frame(gtx.Ops)
		}
	}
}
