## Features

- [x] Dashboard

- [x] Autocomplete suggestions for the activity and tags

- [x] Nav bar

- [x] Accounts

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

- [x] Fix the cacheing of the profile pic so that when a new profile pic with the same file extension is uploaded, it doesn't reuse the older cahced version

- [ ] Add the profile page to the browser history, and make it it's own HTML page

- [ ] Figure out why the timer isn't working if i go for a second time

- [ ] Add a confirm to cancel pop up, when the timer is running, and when they want to cancel at the `timer-form`

- [ ] Set up the browser history when moving to the profile page, or the stats page

- [ ] Set up the link account form, and make sure that when users to try change their email address, that it doesn't already exist in the database. We will need to be able to send back error messages in this case, and we can prompt the user to link it to their account

- [ ] Update the database so that user.email is a unique field

- [ ] Add a 'last_seen' field in the users DB table, to track when the last time the user visited was, and if it is longer than 400 days ago, we can assume that the cookie has expired, and if they don't have an email associated with the user, we can delete the entry since once the cookie is lost, there is no way to retrieve the data ever again
