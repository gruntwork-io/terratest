$(document).ready(function () {
  $('.youtube-player').on('click', function() {
    if ($(this).find('.youtube-player__thumb').length > 0) {
      $(this).addClass('played')
      console.log('clic', $(this).width(), $(this).height())
      const video_url = $(this).data('video-url')
      $(this).append('<iframe ' +
        'width="' + $(this).width() + 'px"' +
        'height="' + $(this).height() + 'px"' +
        'allowfullscreen ' +
        'src="'+ video_url +'"></iframe>')
      $(this).find('.youtube-player__thumb').remove()
    }
  })
})
