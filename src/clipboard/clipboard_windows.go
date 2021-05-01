// +build windows

package clipboard

import "github.com/lxn/walk"

func get() (string, error) {
	c := walk.Clipboard()
	return c.Text()
}

func set(text string) error {
	c := walk.Clipboard()
	return c.SetText(text)
}

func onChange(f func(string)) {
	c := walk.Clipboard()
	// c.ContentsChanged().Attach()
}
