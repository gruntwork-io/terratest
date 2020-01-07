/* TEMPORARY LOGIN BOX */
$(document).ready(function () {
  if (!localStorage.getItem('temp_login_box')) {
    $('body').append(
      "<div id='temp-login-box' style='position:absolute;width:100vw;height:100vh;top:0;left:0;background:#f2f2f2;z-index:1000000;display: flex;align-items:center;justify-content:center;flex-direction:column'>" +
      "<div>" +
        "<label>Password:</label></br>" +
        "<input id='temp-login-pass' style='width: 300px' type='password' />" +
      "</div><br/>" +
      "<button id='temp-login-btn' class='btn btn-primary'>Login</button>" +
      "</div>"
    )
    $('body').css("overflow", "hidden")
    $('body').css("max-width", "100vw")
    $('body').css("max-height", "100vh")

    $('body').on('click', '#temp-login-btn', function() {
      if ($('#temp-login-pass').val() === 'gruntwork') {
        localStorage.setItem('temp_login_box', true)
        $('#temp-login-box').remove()
        $('body').css("overflow", "")
        $('body').css("max-width", "")
        $('body').css("max-height", "")
      }
    })
  }
})
/* END: TEMPORARY LOGIN BOX */
