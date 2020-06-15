window.history.pushState('page2', 'Title', '/page2.php');

function showHidePasswordParent(passwordInputId) {
  return function() {
    let x = document.getElementById(urlPassword);
    if (x.type === "password") {
      x.type = "text";
    } else {
      x.type = "password";
    }
  }
}

document.addEventListener("load", "read", function() {
  var showHidePassword = showHidePasswordParent("urlPassword");

  document.getElementById("hidePassword").addEventListener("toggle", function() {
    showHidePassword();
  })
});