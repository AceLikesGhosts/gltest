package ui

import (
	"image/color"
	"log"
	"strings"
	"sync"
	"syscall"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/acelikesghosts/gltest/utils"
)

// var titles, err = utils.GetWindowTitles()

var (
	titles     []utils.WindowInfo
	titlesErr  error
	titlesOnce sync.Once
)

var editor widget.Editor
var list widget.List
var clickables []widget.Clickable

func onTitleClicked(hwnd uintptr) {
	log.Printf("on title clicked called %x", hwnd)
	utils.FocusWindow(syscall.Handle(hwnd))
}

func loadTitles() {
	titlesOnce.Do(func() {
		titles, titlesErr = utils.GetWindowTitles()
	})
}

func ListWindows(gtx layout.Context, theme *material.Theme) {
	// log.Printf("getting window titles")
	loadTitles()
	if titlesErr != nil {
		panic(titlesErr)
	}

	// log.Printf("iterating titles")
	// log.Printf("titles %v", titles)

	query := strings.ToLower(editor.Text())

	if len(clickables) < len(titles) {
		clickables = make([]widget.Clickable, len(titles))
	}

	// Filter titles
	var filtered []struct {
		index int
		title string
		hwnd  uintptr
	}
	for i, title := range titles {
		if strings.Contains(strings.ToLower(title.Title), query) {
			filtered = append(filtered, struct {
				index int
				title string
				hwnd  uintptr
			}{
				index: i,
				title: title.Title,
				hwnd:  uintptr(title.Hwnd),
			})
		}
	}

	// layout.Flex{Axis: layout.Vertical}.Layout(gtx,
	// 	layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	// 		return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	// 			return material.Editor(theme, &editor, "Search...").Layout(gtx)
	// 		})
	// 	}),
	// 	layout.Rigid(func(gtx layout.Context) layout.Dimensions {
	// 		return material.List(theme, &list).Layout(gtx, len(filtered), func(gtx layout.Context, index int) layout.Dimensions {
	// 			return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
	// 				return material.Body1(theme, filtered[index]).Layout(gtx)
	// 			})
	// 		})
	// 	}),
	// )

	layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			// search input
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Editor(theme, &editor, "Search...").Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Spacer{Height: unit.Dp(16)}.Layout(gtx)
			}),
			// list
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				list := layout.List{Axis: layout.Vertical}
				return list.Layout(gtx, len(filtered), func(gtx layout.Context, i int) layout.Dimensions {
					item := filtered[i]
					click := &clickables[item.index]

					if click.Clicked(gtx) {
						onTitleClicked(item.hwnd)
					}

					return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							lbl := material.Body1(theme, item.title)
							lbl.Color = color.NRGBA{A: 255}
							return lbl.Layout(gtx)
						})
					})
				})
			}),
		)
	})
}
