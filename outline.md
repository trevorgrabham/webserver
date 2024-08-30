## Features

- [x] Dashboard

- [ ] Autocomplete suggestions for the activity and tags

- [ ] Nav bar

- [ ] Accounts

- [x] Tag aggregate section

- [ ] Graphical tabs for all activities, or filter for specific tags

---

## Fixes / Bugs

- [x] Stop the `#timer-display` from overflowing on smaller screens

- [x] Fix the nested forms so that they work properly again. (Look into htmx post methods within forms)

- [x] Style the `.timer-form-input` elements

- [x] Fix the form, not currently sending the timer data

- [ ] Fix the background color when auto-selecting the `#activity-input-container` input

- [x] Make the `#dashboard-section` reset the scroll to the top when a new element is added

- [x] Add the activity bar to each `.card-container` dashboard item

- [x] Add variable colors to the `.card-container` activity bar

- [x] Add the tag scroll container to each `.card-container` dashboard item

- [x] Add the activity tags to each activity on `:hover`

- [x] Add different colors to the `.card-total-hours` text

- [x] Fix the width of the `.duration-bar` for when the total hours = 0

- [x] Get each duration in the `.duration-bar` to expand out from its center, not just left to right

- [ ] Make the colors for the `#dashboard-section` a bit prettier

- [ ] Make the timer auto-reset if timer < 1 min

- [ ] Make the `.tag-summary-tag-title` text expand if on hover if they are hidden

- [ ] Animate the scrolling, so that when adding a new activity, it slowly scrolls to the top of the `#dashboard-section` and when loading more tags, it slowly scrolls to the bottom of the `#tag-summary-data-container`

- [ ] Add some text shadow to the `#timer-display` to make it seem neon

- [ ] Add some keypress listeners to the `#activity-input` so that we can change the focus on our autocomplete suggestions by pressing tab, enter, etc.

- [ ] Deal with the empty space that is created when a use does not have any data yet to display. Careful to make sure that things are styled properly when new data is created for the user

- [ ] Add a tag to the `.tag-summary-section` when we submit a new timer to the database
