package barnard

import (
	"fmt"
	"strings"
	"time"

	"github.com/bontibon/gumble/barnard/uiterm"
	"github.com/bontibon/gumble/gumble"
	"github.com/kennygrant/sanitize"
)

const (
	uiViewLogo        = "logo"
	uiViewTop         = "top"
	uiViewStatus      = "status"
	uiViewInput       = "input"
	uiViewInputStatus = "inputstatus"
	uiViewOutput      = "output"
	uiViewTree        = "tree"
)

func esc(str string) string {
	return sanitize.HTML(str)
}

func (b *Barnard) UpdateInputStatus(status string) {
	b.UiInputStatus.Text = status
	b.UiTree.Rebuild()
	b.Ui.Refresh()
}

func (b *Barnard) AddOutputLine(line string) {
	now := time.Now()
	b.UiOutput.AddLine(fmt.Sprintf("[%02d:%02d:%02d] %s", now.Hour(), now.Minute(), now.Second(), line))
	b.Ui.Refresh()
}

func (b *Barnard) AddOutputMessage(sender *gumble.User, message string) {
	if sender == nil {
		b.AddOutputLine(message)
	} else {
		b.AddOutputLine(fmt.Sprintf("%s: %s", sender.Name(), strings.TrimSpace(esc(message))))
	}
}

func (b *Barnard) OnVoiceToggle(ui *uiterm.Ui, key uiterm.Key) {
	if b.UiStatus.Text == "  Tx  " {
		b.UiStatus.Text = " Idle "
		b.UiStatus.Fg = uiterm.ColorBlack
		b.UiStatus.Bg = uiterm.ColorWhite
		b.Stream.StopSource()
	} else {
		b.UiStatus.Fg = uiterm.ColorWhite | uiterm.AttrBold
		b.UiStatus.Bg = uiterm.ColorRed
		b.UiStatus.Text = "  Tx  "
		b.Stream.StartSource()
	}
	ui.Refresh()
}

func (b *Barnard) OnQuitPress(ui *uiterm.Ui, key uiterm.Key) {
	b.Client.Close()
	b.Ui.Close()
}

func (b *Barnard) OnClearPress(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.Clear()
	b.Ui.Refresh()
}

func (b *Barnard) OnScrollOutputUp(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.ScrollUp()
	b.Ui.Refresh()
}

func (b *Barnard) OnScrollOutputDown(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.ScrollDown()
	b.Ui.Refresh()
}

func (b *Barnard) OnScrollOutputTop(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.ScrollTop()
	b.Ui.Refresh()
}

func (b *Barnard) OnScrollOutputBottom(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.ScrollBottom()
	b.Ui.Refresh()
}

func (b *Barnard) OnFocusPress(ui *uiterm.Ui, key uiterm.Key) {
	active := b.Ui.Active()
	if active == &b.UiInput {
		b.Ui.SetActive(uiViewTree)
	} else if active == &b.UiTree {
		b.Ui.SetActive(uiViewInput)
	}
}

func (b *Barnard) OnTextInput(ui *uiterm.Ui, textbox *uiterm.Textbox, text string) {
	if text == "" {
		return
	}
	if b.Client != nil && b.Client.Self() != nil {
		b.Client.Self().Channel().Send(text, false)
		b.AddOutputMessage(b.Client.Self(), text)
	}
}

func (b *Barnard) OnUiInitialize(ui *uiterm.Ui) {
	ui.SetView(uiViewLogo, 0, 0, 0, 0, &uiterm.Label{
		Text: " barnard ",
		Fg:   uiterm.ColorWhite | uiterm.AttrBold,
		Bg:   uiterm.ColorMagenta,
	})

	ui.SetView(uiViewTop, 0, 0, 0, 0, &uiterm.Label{
		Fg: uiterm.ColorWhite,
		Bg: uiterm.ColorBlue,
	})

	b.UiStatus = uiterm.Label{
		Text: " Idle ",
		Fg:   uiterm.ColorBlack,
		Bg:   uiterm.ColorWhite,
	}
	ui.SetView(uiViewStatus, 0, 0, 0, 0, &b.UiStatus)

	b.UiInput = uiterm.Textbox{
		Fg:    uiterm.ColorWhite,
		Bg:    uiterm.ColorBlack,
		Input: b.OnTextInput,
	}
	ui.SetView(uiViewInput, 0, 0, 0, 0, &b.UiInput)

	b.UiInputStatus = uiterm.Label{
		Fg: uiterm.ColorBlack,
		Bg: uiterm.ColorWhite,
	}
	ui.SetView(uiViewInputStatus, 0, 0, 0, 0, &b.UiInputStatus)

	b.UiOutput = uiterm.Textview{
		Fg: uiterm.ColorWhite,
		Bg: uiterm.ColorBlack,
	}
	ui.SetView(uiViewOutput, 0, 0, 0, 0, &b.UiOutput)

	b.UiTree = uiterm.Tree{
		Generator: b.TreeItem,
		Listener:  b.TreeItemSelect,
		Fg:        uiterm.ColorWhite,
		Bg:        uiterm.ColorBlack,
	}
	ui.SetView(uiViewTree, 0, 0, 0, 0, &b.UiTree)

	b.Ui.AddKeyListener(b.OnFocusPress, uiterm.KeyTab)
	b.Ui.AddKeyListener(b.OnVoiceToggle, uiterm.KeyF1)
	b.Ui.AddKeyListener(b.OnQuitPress, uiterm.KeyF10)
	b.Ui.AddKeyListener(b.OnClearPress, uiterm.KeyCtrlL)
	b.Ui.AddKeyListener(b.OnScrollOutputUp, uiterm.KeyPgup)
	b.Ui.AddKeyListener(b.OnScrollOutputDown, uiterm.KeyPgdn)
	b.Ui.AddKeyListener(b.OnScrollOutputTop, uiterm.KeyHome)
	b.Ui.AddKeyListener(b.OnScrollOutputBottom, uiterm.KeyEnd)
}

func (b *Barnard) OnUiResize(ui *uiterm.Ui, width, height int) {
	ui.SetView(uiViewLogo, 0, 0, 9, 1, nil)
	ui.SetView(uiViewTop, 9, 0, width-6, 1, nil)
	ui.SetView(uiViewStatus, width-6, 0, width, 1, nil)
	ui.SetView(uiViewInput, 0, height-1, width, height, nil)
	ui.SetView(uiViewInputStatus, 0, height-2, width, height-1, nil)
	ui.SetView(uiViewOutput, 0, 1, width-20, height-2, nil)
	ui.SetView(uiViewTree, width-20, 1, width, height-2, nil)
}
