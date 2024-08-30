package timer

import "fmt"

func Form() string {
	return fmt.Sprintf(`
  <form
    id="timer-form"
    hx-post="/submitActivity"
    hx-target="#timer-container"
    action=""
		autocomplete="off"
  >
    <input id="hidden-timer" name="timer" type="hidden" />
		<div id="activity-input-container" class="timer-form-input-row">

			<div id="activity-input-wrapper">
				<input type="text" 
					id="activity-input"
					class="timer-form-input"
					hx-get="activitySuggestions"
					hx-trigger="keyup changed delay:500ms"
					hx-target="#activity-suggestions"
					hx-swap="outerHTML"
					name="activity" 
					placeholder="What were you doing?" 
					required
					minlength="2"
					maxlength="255"
				/>
				<div id="activity-suggestions"></div>
			</div>
	
			<div id="tags-wrapper">
				<div id="tags-container"></div>
			</div>

		</div>


		<div 
			id="tag-input-container" 
			class="timer-form-input-row"
		>

			<input type="text"
				id="tag-input" 
				class="timer-form-input"
				name="temporary-tag" 
				placeholder="tags" 
				maxlength="255"
			/>
			
			%s

		</div>

		<div id="timer-form-submit-container" class="timer-form-input-row">
	
			%s
		
			%s

		</div>
	</form>
	`, PlusButton("add-tag-svg", []string{"timer-button", "button-sub-form"}, nil), SuccessButton("", []string{"timer-button", "button-form"}, nil), CancelButton("reset-button", []string{"timer-button", "button-form"}, nil))
}