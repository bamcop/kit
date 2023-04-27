package tui

import (
	"github.com/manifoldco/promptui"
)

func Confirm(title string, f func() error) error {
	prompt := promptui.Prompt{
		Label:     title,
		IsConfirm: true,
	}
	if result, err := prompt.Run(); err != nil {
		return nil
	} else {
		if result != "y" {
			return nil
		}

		if err := f(); err != nil {
			return err
		}
		return nil
	}
}
