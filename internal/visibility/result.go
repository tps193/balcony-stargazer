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

func (output *ConsoleOutput) Get(astroObject *AstroObject, visibilityWindows *[]VisibilityWindow) string {
	res := make([]byte, 0)
	res = fmt.Appendf(res, "Visibility of %s:\n", astroObject.Name)
	for i, window := range *visibilityWindows {
		res = fmt.Appendf(res, "%d: %s\n", i, window.EndTime.Sub(window.StartTime))
		res = fmt.Appendf(res, "\tStart: %s (%f°)\n", window.StartTime, window.StartAlt)
		res = fmt.Appendf(res, "\tEnd: %s (%f°)\n", window.EndTime, window.EndAlt)
	}
	return string(res)
}

type JsonOutput struct {
	ObjectName        string             `json:"name"`
	VisibilityWindows []VisibilityWindow `json:"windows"`
}

func NewJsonOutput() *JsonOutput {
	return &JsonOutput{}
}

func (output *JsonOutput) Get(astroObject *AstroObject, visibilityWindows []VisibilityWindow) string {
	output.ObjectName = astroObject.Name
	output.VisibilityWindows = visibilityWindows
	res, err := json.Marshal(output)
	if err != nil {
		res = fmt.Append(nil, "Error creating json", err.Error())
	}
	return string(res)
}
