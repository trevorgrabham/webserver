package html

import (
	"html/template"
	"unicode/utf8"

	"github.com/trevorgrabham/webserver/webserver/lib/dashboard"
	"github.com/trevorgrabham/webserver/webserver/lib/profile"
	"github.com/trevorgrabham/webserver/webserver/lib/tag"
)

type ElementAttributesData []template.HTMLAttr 

type SVGButtonData struct {
	Button ElementAttributesData
	SVG ElementAttributesData
}

// TagValue and TagWidth should only ever be set by calling FormTagData.AddTag(tag)
type FormTagData struct {
	SVG ElementAttributesData
	TagValue string
	TagWidth int
}

type FormData struct {
	TimerValue string 
	ActivitySuggestions []string
	// TODO 
	// TagSuggestions []string
	PlusButton ElementAttributesData 
	SuccessButton ElementAttributesData 
	CancelButton ElementAttributesData
}

type TimerData struct {
	ButtonContainer ElementAttributesData
	StartButton ElementAttributesData
	PauseButton ElementAttributesData
	StopButton ElementAttributesData
}

type ProfileFieldData struct {
	Value interface{}
	Errors []string
}

type ProfileData struct {
	Pic ProfileFieldData
	Name ProfileFieldData
	Email ProfileFieldData
	Link []string
}

type NavData struct {
	TimerData
	ID int64
	Chart ElementAttributesData
}

// 																									|| End types || 

// 																								|| Start methods || 

func (f *FormTagData) AddTag(tag string) {
	f.TagValue = template.HTMLEscapeString(tag)
	f.TagWidth = utf8.RuneCountInString(tag)
}

func (f *FormData) AddSuggestions(activityOrTag string, suggestions []string) {
	switch {
	case activityOrTag == "activity":
		f.ActivitySuggestions = make([]string, len(suggestions))
		for _, suggestion := range suggestions {
			f.ActivitySuggestions = append(f.ActivitySuggestions, template.HTMLEscapeString(suggestion))
		}
	case activityOrTag == "tag":
		// TODO
		// f.TagSuggestions = make([]string, len(suggestions))
		// for _, suggestion := range suggestions {
		// 	f.TagSuggestions = append(f.TagSuggestions, template.HTMLEscapeString(suggestion))
		// }
	}
}

// 																								|| End methods || 

// 																						|| Start builder functions || 

func NewElementAttributes(attrs []string) (attributes ElementAttributesData) {
	attributes = make([]template.HTMLAttr, len(attrs))
	for _, attr := range attrs {
		attributes = append(attributes, template.HTMLAttr(attr))
	}
	return 
}

// 																						|| End builder functions || 

// 																							|| Start Func Maps || 

var TagFuncMap = template.FuncMap{ "normalizeCount": tag.NormalizeCount, "getMax": tag.GetMaxCount }

var ProfileFuncMap = template.FuncMap{ "getPic": profile.GetProfilePic }

var NavFuncMap = template.FuncMap{ "getPic": profile.GetProfilePic }

var DashboardFuncMap = template.FuncMap{ "formatTotalMin": dashboard.FormatTotalMin }