package templateutil

var ParseFiles = map[string][]string{
	"card":          {"static/html/dashboard/card.gohtml"},
	"cards":         {"static/html/dashboard/cards.gohtml", "static/html/dashboard/card.gohtml"},
	"nav":           {"static/html/nav/nav.gohtml"},
	"profile":       {"static/html/profile/profile.gohtml", "static/html/profile/showpic.gohtml", "static/html/profile/showname.gohtml", "static/html/profile/showemail.gohtml"},
	"editpic":       {"static/html/profile/editpic.gohtml"},
	"savepic":       {"static/html/profile/savepic.gohtml"},
	"editname":      {"static/html/profile/editname.gohtml"},
	"savename":      {"static/html/profile/savename.gohtml"},
	"editemail":     {"static/html/profile/editemail.gohtml"},
	"saveemail":     {"static/html/profile/saveemail.gohtml"},
	"tagscontainer": {"static/html/tags/tagscontainer.gohtml", "static/html/tags/tags.gohtml", "static/html/tags/tag.gohtml", "static/html/tags/disabledmoretagsbutton.gohtml"},
	"tags":          {"static/html/tags/tags.gohtml", "static/html/tags/tag.gohtml", "static/html/tags/disabledmoretagsbutton.gohtml"},
	"pausebutton":   {"static/html/timer/pausebutton.gohtml", "static/html/svg/pausesvg.gohtml"},
	"stopbutton":    {"static/html/timer/stopbutton.gohtml", "static/html/svg/stopsvg.gohtml"},
	"startbutton":   {"static/html/timer/startbutton.gohtml", "static/html/svg/startsvg.gohtml"},
	"form":          {"static/html/timer/form.gohtml", "static/html/timer/successbutton.gohtml", "static/html/timer/plusbutton.gohtml", "static/html/timer/cancelbutton.gohtml", "static/html/svg/successsvg.gohtml", "static/html/svg/plussvg.gohtml", "static/html/svg/cancelsvg.gohtml"},
	"autocomplete":  {"static/html/timer/autocomplete.gohtml"},
	"newtag":        {"static/html/timer/newtag.gohtml", "static/html/svg/removesvg.gohtml"},
	"defaulttimer":  {"static/html/timer/defaulttimer.gohtml", "static/html/timer/startbutton.gohtml", "static/html/svg/startsvg.gohtml"},
}