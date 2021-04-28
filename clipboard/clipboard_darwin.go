// +build darwin

package clipboard

import (
	"errors"

	"github.com/cxfksword/fnsync-desktop/clipboard/clip"
)

func set(text string) error {
	ok := clip.Set(text)
	if !ok {
		return errors.New("nothing to set string")
	}
	return nil
}

func get() (string, error) {
	return clip.Get(), nil
}

func onChange(f func(string)) {
	go func() {
		for {
			text := <-clip.GetNotifyChangeCh()
			f(text)
		}
	}()
}
