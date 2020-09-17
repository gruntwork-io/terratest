---
---
$(document).ready(function () {

  const CODE_LINE_HEIGHT = 22
  const CODE_BLOCK_PADDING = 10

  window.examples = {
    tags: {},
    nav: {}
  }

  initExamplesNav()

  $(window).resize($.debounce(250, function() {
     buildExamplesNav()
  }))

  // Activate first example
  $('.examples__container').each(function (i, ec) {
    // Find first element:
    const firstElementId = $(ec).find('.examples__nav-item').data('id')
    // Open first element:
    openExample($(ec).attr('id'), firstElementId)

    // Open example when user clicks on tab
    $('.navs').on('click', '.examples__nav-item:not(.static-link)', function() {
      openExample($(ec).attr('id'), $(this).data('id'))
      $('.navs__dropdown-menu').removeClass('active')
    })
  })

  // Open example and scroll to examples section when user clicks on
  // tech in the header
  $('.link-to-test-with-terratest').on('click', function() {
    // Find any containting the keyword from data-target
    const found = $('.navs .examples__nav-item[data-id*="'+$(this).data('target')+'"]')
    if (found && found.length > 0) {
      openExample('index_page', $(found[0]).data('id'))
    } else {
      // RESCUE: If none found, open any (first available):
      openExample('index_page', $($('.navs .examples__nav-item')[0]).data('id'))
    }
    scrollToTests()
  })

  // Switch between code snippets (files)
  $('.examples__tabs .tab').on('click', function() {
    $(this).parents('.examples__tabs').find('.tab').removeClass('active')
    $(this).addClass('active')

    $(this).parents('.examples__block').find('.examples__code').removeClass('active')
    $($(this).data('target')).addClass('active')

    loadCodeSnippet()
  })

  // Open dropdown of technologies to select
  $('.navs__dropdown-arrow').on('click', function() {
    $('.navs__dropdown-menu').toggleClass('active')
  })

  // Open popup when user click on circle with the number
  $('.examples__container').on('click', '.code-popup-handler', function() {
    const isActive = $(this).hasClass('active')
    $('.code-popup-handler').removeClass('active')
    if (!isActive) {
      $(this).addClass('active')
    }
  })

  function scrollToTests() {
    $([document.documentElement, document.body]).animate({
        scrollTop: $('#index-page__test-with-terratest').offset().top
    }, 500)
  }

  function openExample(exampleContainerId, target) {
    // Change active nav in window state and rebuild navigation first
    const $ecId = $('#'+exampleContainerId)
    window.examples.nav[exampleContainerId].current = target
    buildExamplesNav()

    // Change active tab in navigation
    $ecId.find('.examples__nav-item').removeClass('active')
    const jTarget = $('.navs .examples__nav-item[data-id="'+target+'"]')
    jTarget.addClass('active')

    // Change the block below navigation (with code snippets)
    $ecId.find('.examples__block').removeClass('active')
    $ecId.find('#example__block-' + target).addClass('active')

    // Set current tab
    $ecId.find('.examples__nav .navs').removeClass('active')

    loadCodeSnippet()
  }

  function loadCodeSnippet() {
    $('.examples__block.active .examples__code.active').each(async function (i, activeCodeSnippet) {
      const $activeCodeSnippet = $(activeCodeSnippet)
      const exampleTarget = $(this).data('example')
      const fileId = $(this).data('target')
      if (!$activeCodeSnippet.data('loaded')) {
        try {
          const response = await fetch($activeCodeSnippet.data('url'))
          let content = await response.text()
          $activeCodeSnippet.attr('data-loaded', true)
          if ($activeCodeSnippet.data('skip-tags')) {
            // Remove the website::tag::xxx:: prefix from the code snippet
            content = content.replace(/website::tag::.*?:: ?/mg, '')
          } else {
            findTags(content, exampleTarget, fileId)
            // Remove the website::tag::xxx:: comment entirely from the code snippet
            content = content.replace(/^.*website::tag.*\n?/mg, '')
          }
          $activeCodeSnippet.find('code').text(content)
          Prism.highlightAll()
        } catch(err) {
          $activeCodeSnippet.find('code').text('Resource could not be loaded.')
          console.error(err)
        }
      }
      updatePopups()
      openPopup(exampleTarget, 1)
    })
  }

  function findTags(content, exampleTarget, fileId) {
    let tags = []
    let regexpTags = /website::tag::(\d)::\s*(.*)/mg
    let match = regexpTags.exec(content)
    do {
      if (match && match.length > 0) {
        tags.push({
          text: match[2],
          tag: match[0],
          step: match[1],
          line: findLineNumber(content, match[0])
        })
      }
    } while((match = regexpTags.exec(content)) !== null)
    window.examples.tags[exampleTarget] = Object.assign({
        [fileId]: tags
      },
      window.examples.tags[exampleTarget]
    )
  }

  function findLineNumber(content, text) {
    let tagIndex = content.indexOf(text)
    let tempString = content.substring(0, tagIndex)
    let lineNumber = tempString.split('\n').length
    return lineNumber
  }

  function updatePopups() {
    $('.code-popup-handler').remove()
    const activeCode = $('.examples__block.active .examples__code.active')
    const exampleTarget = activeCode.data('example')
    const fileId = activeCode.data('target')
    const exampleTargetTags = window.examples.tags[exampleTarget] || {};
    const fileTags = exampleTargetTags[fileId];

    if (fileTags) {
      const tagsLen = fileTags.length

      fileTags.map( function(v,k) {
        const top = (CODE_LINE_HEIGHT * (v.line - k)) + CODE_BLOCK_PADDING;

        // If two pop-ups are close to each other, add CSS class that will scale them down
        let scaleClass = ''
        if (
            (k > 0 && Math.abs(v.line - fileTags[k-1].line) < 3 )
            || (k < tagsLen - 1 && Math.abs(v.line - fileTags[k+1].line) < 3 )
        ) {
          scaleClass = 'sm-scale'
        }

        const elToAppend =
            '<div class="code-popup-handler '+scaleClass+'" style="top: '+top+'px" data-step="'+v.step+'">' +
            '<span class="number">' + v.step + '</span>' +
            '<div class="shadow-bg-1"></div><div class="shadow-bg-2"></div>' +
            '<div class="popup">' +
            '<div class="left-border"></div>' +
            '<div class="content">' +
            '<p class="text">' + v.text + '</p>' +
            '</div>' +
            '</div>'
        const code = $("#example__code-"+exampleTarget+"-"+fileId)
        code.append(elToAppend)
      })
    }

    openPopup(exampleTarget, 0)
  }

  function openPopup(techName, step) {
    $('.code-popup-handler').removeClass('active')
    $('#example__block-'+techName).find('.code-popup-handler[data-step="'+step+'"]').addClass('active')
  }

  function loadExampleDescription(name) {
    return $('#index-page__examples').find('#example__block-'+name+' .description').html()
  }

  function initExamplesNav() {
    window.examples.nav = {}
    $('.examples__container').each(function(eci, ec) {
      $(ec).find('.examples__nav .hidden-navs').each(function(rni, refNavs) {
        let navsArr = []
        let currentNav
        $(refNavs).find('.examples__nav-item').each( function(ni, nav) {
          if ($(nav).hasClass('active')) {
            currentNav = $(nav).data('id')
          }
          navsArr.push($(nav))
        })
        window.examples.nav = Object.assign({
          [$(ec).attr('id')]: {
            current: currentNav,
            items: navsArr
          }
        }, window.examples.nav)
      })
    })
  }

  function buildExamplesNav() {
    $('.examples__container').each(function(eci, ec) {
      const ecId = $(ec).attr('id')
      const containerWidth = $(ec).width()
      const NAV_WIDTH = 150
      const ARROW_SLOT_WIDTH = 100

      const noOfVisible = Math.floor((containerWidth - NAV_WIDTH - ARROW_SLOT_WIDTH) / 150)

      const $visibleBar = $($(ec).find('.navs__visible-bar'))
      const $dropdownInput = $($(ec).find('.navs__dropdown-input'))
      const $dropdownMenu = $($(ec).find('.navs__dropdown-menu'))

      $visibleBar.html('')
      $dropdownInput.html('')
      $dropdownMenu.html('')

      let settingCurrent = false

      // Build initial a navigation bar
      if (window.examples.nav
        && ecId in window.examples.nav
        && window.examples.nav[ecId].items) {

        // Visible elements
        let breakSlice = noOfVisible > window.examples.nav[ecId].items.length ? window.examples.nav[ecId].items.length : noOfVisible
        let visibleEls = window.examples.nav[ecId].items.slice(0, breakSlice)
        let hiddenEls = window.examples.nav[ecId].items.slice(breakSlice, window.examples.nav[ecId].items.length)

        let visibleNavIsActive = false
        let hiddenNavIsActive = -1

        if (window.examples.nav[ecId].current) {
          visibleEls.map( function(x,i) {
            if(x.data('id') === window.examples.nav[ecId].current) {
              visibleNavIsActive = true
              x.addClass('active')
            }
          })
          hiddenEls.map( function(x,i) {
            if(x.data('id') === window.examples.nav[ecId].current) {
              hiddenNavIsActive = i
              x.addClass('active')
            }
          })
        }

        visibleEls.map(function(nav,i) {
          $visibleBar.append($(nav).clone())
        })

        if (hiddenNavIsActive > -1) {
          const sliced = hiddenEls.splice(hiddenNavIsActive, 1)
          $dropdownInput.append($(sliced[0]).clone())
        } else {
          $dropdownInput.append($(hiddenEls.shift()).clone())
        }

        hiddenEls.map(function(nav,i) {
          $dropdownMenu.append($(nav).clone())
        })

        // Add static links
        $dropdownMenu.append($(ec).find('.hidden-navs__static-links').html())
      }
    })
  }

})
