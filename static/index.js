let timer = document.getElementById("timer-display");
let startTime;
let elapsedHours = 0;
let elapsedMinutes = 0;
let elapsedSeconds = 0;
let startDate = "";

// Resize hook
window.addEventListener("resize", () => {
  // Resize the timer display
  let timer = document.getElementById("timer-display");
  if (timer !== null) {
    resizeFont(timer);
  }
});

// Before htmx event triggers
document.addEventListener("htmx:before-request", function (event) {
  /* 
    Need to have an empty .card-container with the appropriate date id before we submit any timer data.
    The timer form submit uses hx-swap-oob to add the card to the dashboard if it doesn't already exist
  */
  if (event.detail.elt.id === "timer-form") {
    addNewDashboardCard();
  }
});

// After htmx event triggers
document.addEventListener("htmx:after-request", function (event) {
  // If a tag was removed, check to see if the overflow value on the tag container needs to be changed
  if (event.detail.pathInfo.requestPath === "/removeTag") {
    checkOverflowing(document.getElementById("tags-wrapper"));
  }
});

// After the new htmx elements settle
document.addEventListener("htmx:after-settle", function (event) {
  // If we just added a new tag, ensure that the new tag doesn't overflow the container
  if (event.detail.pathInfo.requestPath === "/addTag") {
    checkOverflowing(document.getElementById("tags-wrapper"));
  }
  // Set the value on the hidden timer input if we just stopped the timer
  if (event.detail.pathInfo.requestPath === "/stopTimer") {
    updateTimerInput();
  }
  // Reset the scroll position to the top (new card just added) if we just added new timer data
  if (event.detail.pathInfo.requestPath === "/submitActivity") {
    document.getElementById("dashboard-section").scrollTop = 0;
  }
});

// Works for starting and resuming
function startTimer() {
  if (timer === null) {
    timer = document.getElementById("timer-display");
  }
  startTime =
    Date.now() -
    elapsedSeconds * 1000 -
    elapsedMinutes * 1000 * 60 -
    elapsedHours * 1000 * 60 * 60;
}

function resetTimer() {
  elapsedSeconds = 0;
  elapsedMinutes = 0;
  elapsedHours = 0;
}

// Calculate the sec, min, hr as difference between now and start, not just an increment so we don't have to worry about time creep if the sleep(1s) isn't exact
function formatTime() {
  let elapsedTime = Date.now() - startTime;
  // elapsedTime in ms, so convert to relevant units
  elapsedSeconds = Math.floor(elapsedTime / 1000) % 60;
  elapsedMinutes = Math.floor(elapsedTime / 1000 / 60) % 60;
  elapsedHours = Math.floor(elapsedTime / 1000 / 60 / 60);
  secondsString = formatTimeString(elapsedSeconds);
  minutesString = formatTimeString(elapsedMinutes);
  return elapsedHours < 1
    ? minutesString + ":" + secondsString
    : elapsedHours + ":" + minutesString + ":" + secondsString;
}

function formatTimeString(time) {
  if (time < 10) {
    return "0" + time;
  }
  return time.toString();
}

function updateTimer() {
  if (timer === null) {
    timer = document.getElementById("timer-display");
  }
  timer.innerHTML = formatTime();
  resizeFont(timer);
}

// Starts at 5em font-size and reduces it by 10% until it no longer overflows
function resizeFont(element) {
  element.style.fontSize = "5em";
  while (element.scrollWidth > element.clientWidth) {
    let fontSize = parseInt(
      element.style.fontSize.substring(0, element.style.fontSize.length - 2)
    );
    fontSize *= 0.9;
    element.style.fontSize = fontSize + "em";
  }
}

function checkOverflowing(element) {
  element.scrollHeight > element.clientHeight
    ? element.classList.add("overflowing")
    : element.classList.remove("overflowing");
}

// For after we stop the timer, to populate the hidden timer input
function updateTimerInput() {
  if (timer === null) {
    timer = document.getElementById("timer-display");
  }
  document
    .getElementById("hidden-timer")
    .setAttribute("value", timer.innerHTML);
}

function checkDateExists() {
  if (timer === null) {
    timer = document.getElementById("timer-display");
  }
  let timeSegments = timer.innerHTML.split(":");
  let timeDelta = 0;
  // convert our time into ms
  if (timeSegments.length == 3) {
    // H:MM:SS
    timeDelta += 60 * 60 * 1000 * parseInt(timeSegments[0]);
    timeDelta += 60 * 1000 * parseInt(timeSegments[1]);
    timeDelta += 1000 * parseInt(timeSegments[2]);
  } else {
    // MM:SS
    timeDelta += 60 * 1000 * parseInt(timeSegments[0]);
    timeDelta += 1000 * parseInt(timeSegments[1]);
  }
  startDate = new Date(new Date() - timeDelta).toLocaleDateString("en-CA", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  });
  return document.getElementById(`date-${startDate}`) != null;
}

// So that we can recieve any new dashboard cards if needed
function addNewDashboardCard() {
  if (!checkDateExists()) {
    let parent = document.getElementById("dashboard-section");
    let child = document.createElement("div");
    child.id = `date-${startDate}`;
    startDate = "";
    if (parent.firstChild !== null) {
      parent.insertBefore(child, parent.firstChild);
    } else {
      parent.appendChild(child);
    }
  }
}
