$(document).ready(function () {
  $('.video-player').on('click', function() {
    if ($(this).find('.frame').length > 0) {
      $(this).addClass('played')
      const video_url = $(this).data('video-url')
      $(this).append('<iframe ' +
        'width="' + $(this).width() + 'px"' +
        'height="' + $(this).height() + 'px"' +
        'allowfullscreen ' +
        'src="'+ video_url +'"></iframe>')
      $(this).find('.frame').remove()
      $(this).find('.btn-video').remove()
    }
  })
})
