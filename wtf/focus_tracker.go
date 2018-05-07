package wtf

import (
	"github.com/rivo/tview"
)

// FocusTracker is used by the app to track which onscreen widget currently has focus,
// and to move focus between widgets.
type FocusTracker struct {
	App     *tview.Application
	Idx     int
	Widgets []TextViewer
}

/* -------------------- Exported Functions -------------------- */

// Next sets the focus on the next widget in the widget list. If the current widget is
// the last widget, sets focus on the first widget.
func (tracker *FocusTracker) Next() {
	if tracker.widgetHasFocus() {
		tracker.blur(tracker.Idx)
		tracker.increment()
		tracker.focus(tracker.Idx)
	}
}

// None removes focus from the currently-focused widget.
func (tracker *FocusTracker) None() {
	tracker.blur(tracker.Idx)
}

// Prev sets the focus on the previous widget in the widget list. If the current widget is
// the last widget, sets focus on the last widget.
func (tracker *FocusTracker) Prev() {
	if tracker.widgetHasFocus() {
		tracker.blur(tracker.Idx)
		tracker.decrement()
		tracker.focus(tracker.Idx)
	}
}

func (tracker *FocusTracker) Refocus() {
	tracker.focus(tracker.Idx)
}

/* -------------------- Unexported Functions -------------------- */

func (tracker *FocusTracker) blur(idx int) {
	widget := tracker.focusableAt(idx)
	if widget == nil {
		return
	}

	view := widget.TextView()
	view.Blur()

	view.SetBorderColor(ColorFor(widget.BorderColor()))
}

func (tracker *FocusTracker) decrement() {
	tracker.Idx = tracker.Idx - 1

	if tracker.Idx < 0 {
		tracker.Idx = len(tracker.focusables()) - 1
	}
}

func (tracker *FocusTracker) focus(idx int) {
	widget := tracker.focusableAt(idx)
	if widget == nil {
		return
	}

	view := widget.TextView()

	tracker.App.SetFocus(view)
	view.SetBorderColor(ColorFor(Config.UString("wtf.colors.border.focused", "gray")))
}

func (tracker *FocusTracker) focusables() []TextViewer {
	focusable := []TextViewer{}

	for _, widget := range tracker.Widgets {
		if widget.Focusable() {
			focusable = append(focusable, widget)
		}
	}

	return focusable
}

func (tracker *FocusTracker) focusableAt(idx int) TextViewer {
	if idx < 0 || idx >= len(tracker.focusables()) {
		return nil
	}

	return tracker.focusables()[idx]
}

func (tracker *FocusTracker) increment() {
	tracker.Idx = tracker.Idx + 1

	if tracker.Idx == len(tracker.focusables()) {
		tracker.Idx = 0
	}
}

// widgetHasFocus returns true if one of the widgets currently has the app's focus,
// false if none of them do (ie: perhaps a modal dialog currently has it instead)
func (tracker *FocusTracker) widgetHasFocus() bool {
	if tracker.Idx < 0 {
		return true
	}

	for _, widget := range tracker.Widgets {
		if widget.TextView() == tracker.App.GetFocus() {
			return true
		}
	}

	return false
}
