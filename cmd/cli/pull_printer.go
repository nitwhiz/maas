package main

import (
	"encoding/json"
	"fmt"
	"github.com/nitwhiz/maas/internal/cursor"
	"io"
	"strings"
)

type pullEvent struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Error          string `json:"error,omitempty"`
	Progress       string `json:"progress,omitempty"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func printPullProgress(reader *io.ReadCloser) error {
	terminalCursor := cursor.Cursor{}
	layers := make([]string, 0)
	oldIndex := len(layers)

	var event *pullEvent

	decoder := json.NewDecoder(*reader)

	terminalCursor.Hide()

	for {
		if err := decoder.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		imageID := event.ID

		if strings.HasPrefix(event.Status, "Digest:") || strings.HasPrefix(event.Status, "Status:") {
			fmt.Printf("%s\n", event.Status)
			continue
		}

		index := 0
		for i, v := range layers {
			if v == imageID {
				index = i + 1
				break
			}
		}

		if index > 0 {
			diff := index - oldIndex

			if diff > 1 {
				down := diff - 1
				terminalCursor.MoveDown(down)
			} else if diff < 1 {
				up := diff*(-1) + 1
				terminalCursor.MoveUp(up)
			}

			oldIndex = index
		} else {
			layers = append(layers, event.ID)
			diff := len(layers) - oldIndex

			if diff > 1 {
				terminalCursor.MoveDown(diff)
			}

			oldIndex = len(layers)
		}

		terminalCursor.ClearLine()

		if event.Status == "Pull complete" {
			fmt.Printf("%s: %s\n", event.ID, event.Status)
		} else {
			fmt.Printf("%s: %s %s\n", event.ID, event.Status, event.Progress)
		}
	}

	terminalCursor.Show()

	return nil
}
