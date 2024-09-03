package templateutil

var ParseFiles = map[string][]string{
	"card":          {"static/html/dashboard/card.html"},
	"cards":         {"static/html/dashboard/cards.html", "static/html/dashboard/card.html"},
	"nav":           {"static/html/nav/nav.html"},
	"profile":       {"static/html/profile/profile.html", "static/html/profile/showpic.html", "static/html/profile/showname.html", "static/html/profile/showemail.html"},
	"editpic":       {"static/html/profile/editpic.html"},
	"savepic":       {"static/html/profile/savepic.html"},
	"editname":      {"static/html/profile/editname.html"},
	"savename":      {"static/html/profile/savename.html"},
	"editemail":     {"static/html/profile/editemail.html"},
	"saveemail":     {"static/html/profile/saveemail.html"},
	"tagscontainer": {"static/html/tags/tagscontainer.html", "static/html/tags/tags.html", "static/html/tags/tag.html", "static/html/tags/disabledmoretagsbutton.html"},
	"tags":          {"static/html/tags/tags.html", "static/html/tags/tag.html", "static/html/tags/disabledmoretagsbutton.html"},
	"pausebutton":   {"static/html/timer/pausebutton.html", "static/html/svg/pausesvg.html"},
	"stopbutton":    {"static/html/timer/stopbutton.html", "static/html/svg/stopsvg.html"},
	"startbutton":   {"static/html/timer/startbutton.html", "static/html/svg/startsvg.html"},
	"form":          {"static/html/timer/form.html", "static/html/timer/successbutton.html", "static/html/timer/plusbutton.html", "static/html/timer/cancelbutton.html", "static/html/svg/successsvg.html", "static/html/svg/plussvg.html", "static/html/svg/cancelsvg.html"},
	"autocomplete":  {"static/html/timer/autocomplete.html"},
	"newtag":        {"static/html/timer/newtag.html", "static/html/svg/removesvg.html"},
	"defaulttimer":  {"static/html/timer/defaulttimer.html", "static/html/timer/startbutton.html", "static/html/svg/startsvg.html"},
	"chart":         {"static/html/chart/chart.html", "static/html/chart/chartbar.html"},
}