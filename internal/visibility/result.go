package visibility

import (
	"encoding/json"
	"fmt"
)

type Result interface {
	Get(*AstroObject, *[]VisibilityWindow) string
}

type ConsoleOutput struct{}

func NewSimpleOutputResult() *ConsoleOutput {
	return &ConsoleOutput{}
}

func (output *ConsoleOutput) Get(visibilityWindows *[]VisibilityInfo) string {
	res := make([]byte, 0)
	for _, info := range *visibilityWindows {
		res = fmt.Appendf(res, "Visibility of %s (%s):\n", info.Object.Name, info.Object.ObjectType)
		for i, window := range info.VisibilityWindows {
			res = fmt.Appendf(res, "%d: %s\n", i, window.EndTime.Sub(window.StartTime))
			res = fmt.Appendf(res, "\tStart: %s (%f°)\n", window.StartTime, window.StartAlt)
			res = fmt.Appendf(res, "\tEnd: %s (%f°)\n", window.EndTime, window.EndAlt)
		}
	}
	return string(res)
}

type JsonOutput struct {
	VisibilityInfos []VisibilityInfo `json:"windows"`
}

func NewJsonOutput() *JsonOutput {
	return &JsonOutput{}
}

func (output *JsonOutput) Get(visibilityInfos []VisibilityInfo) string {
	output.VisibilityInfos = visibilityInfos
	res, err := json.Marshal(output)
	if err != nil {
		res = fmt.Append(nil, "Error creating json", err.Error())
	}
	return string(res)
}
