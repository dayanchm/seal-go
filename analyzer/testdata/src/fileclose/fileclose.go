package fileclose

import "os"

func missingClose() error {
	f, err := os.Open("example.txt") // want `file "f" is not closed`
	if err != nil {
		return err
	}
	_ = f
	return nil
}

func hasClose() error {
	f, err := os.Open("example.txt")
	if err != nil {
		return err
	}

	defer f.Close()
	return nil
}
