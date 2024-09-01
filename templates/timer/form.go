package timer

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/trevorgrabham/webserver/webserver/lib/templates"
)

type FormData struct {
	Plus	templates.ElementInfo
	Success	templates.ElementInfo
	Cancel	templates.ElementInfo
}

// var formTemplate = `
//   <form
//     id="timer-form"
//     hx-post="/submitActivity"
//     hx-target="#timer-container"
//     action=""
// 		autocomplete="off"
//   >
//     <input id="hidden-timer" name="timer" type="hidden" />
// 		<div id="activity-input-container" class="timer-form-input-row">

// 			<div id="activity-input-wrapper">
// 				<input type="text" 
// 					id="activity-input"
// 					class="timer-form-input"
// 					hx-get="activitySuggestions"
// 					hx-trigger="keyup changed delay:500ms"
// 					hx-target="#activity-suggestions"
// 					hx-swap="outerHTML"
// 					name="activity" 
// 					placeholder="What were you doing?" 
// 					required
// 					minlength="2"
// 					maxlength="255"
// 				/>
// 				<div id="activity-suggestions"></div>
// 			</div>
	
// 			<div id="tags-wrapper">
// 				<div id="tags-container"></div>
// 			</div>

// 		</div>


// 		<div 
// 			id="tag-input-container" 
// 			class="timer-form-input-row"
// 		>

// 			<input type="text"
// 				id="tag-input" 
// 				class="timer-form-input"
// 				name="temporary-tag" 
// 				placeholder="tags" 
// 				maxlength="255"
// 			/>
			
// 				{{template "plusbutton" .Plus}}
// 			</button>

// 		</div>

// 		<div id="timer-form-submit-container" class="timer-form-input-row">
	
// 				{{template "successbutton" .Success}}
		
// 				{{template "cancelbutton" .Cancel}}

// 		</div>
// 	</form>`

func Form(w io.Writer) {
	// var main = template.Must(template.New("form-template").ParseFiles(""))
	// main = template.Must(main.AddParseTree("plusbutton", PlusButton.Tree))
	// main = template.Must(main.AddParseTree("successbutton", SuccessButton.Tree))
	// main = template.Must(main.AddParseTree("cancelbutton", CancelButton.Tree))
	fmt.Println(os.Getwd())
	main := template.Must(template.ParseFiles("static/html/timer/form.gohtml", "static/html/timer/plusbutton.gohtml", "static/html/timer/successbutton.gohtml", "static/html/timer/cancelbutton.gohtml", "static/html/svg/cancelsvg.gohtml", "static/html/svg/successsvg.gohtml", "static/html/svg/plussvg.gohtml"))
	err := main.ExecuteTemplate(w, "form", map[string]templates.ElementInfo{
		"Plus": templates.NewElementInfo("add-tag-svg", []string{"timer-button", "button-sub-form"}, nil, ""),
		"Success": templates.NewElementInfo("", []string{"timer-button", "button-form"}, nil, ""),
		"Cancel": templates.NewElementInfo("reset-button", []string{"timer-button", "button-form"}, nil, ""),
	})
	if err != nil { panic(err) }
}