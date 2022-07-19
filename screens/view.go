package screens

func NewView(caption, text string) View {
	return View{caption: caption, text: text}
}

func (view View) Text() string {
	if len(view.caption) > 0 {
		//return "<b><i>" + view.caption + "</i></b>\n═══════════════════\n" + view.text
		return "<b>" + view.caption + "</b>\n═════════════════\n" + view.text
	}
	return view.text
}

type View struct {
	caption, text string
}
