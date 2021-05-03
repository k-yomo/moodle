package moodle

import (
	"fmt"
	"strings"
)

type Warning struct {
	Item        string `json:"item"`
	ItemID      int    `json:"itemid"`
	WarningCode string `json:"warningcode"`
	Message     string `json:"message"`
}

type Warnings []*Warning

func (l Warnings) Error() string {
	warnings := make([]string, 0, len(l))
	for _, w := range l {
		warnings = append(warnings, fmt.Sprintf("item: %s, itemID: %d, warningCode: %s, message: %s", w.Item, w.ItemID, w.WarningCode, w.Message))
	}
	return strings.Join(warnings, "\n")
}
