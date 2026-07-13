import type { Slide } from '../types'
import { CDN_REVEAL_CSS, CDN_REVEAL_THEME_CSS, CDN_REVEAL_JS, POST_MSG_TYPE } from '../cfg'

export function buildPresentHtml(slides: Slide[], initialSlide: number): string {
  const slidesHtml = slides.map(s => `<section>${s.content}</section>`).join('\n')
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <link rel="stylesheet" href="${CDN_REVEAL_CSS}">
  <link rel="stylesheet" href="${CDN_REVEAL_THEME_CSS}">
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    html, body { width: 100%; height: 100%; overflow: hidden; }
  </style>
</head>
<body>
  <div class="reveal" id="reveal">
    <div class="slides">${slidesHtml}</div>
  </div>
  <script src="${CDN_REVEAL_JS}"></script>
  <script>
    Reveal.initialize({ transition: 'slide', progress: false, controls: false, touch: false, scrollActivationWidth: null, }).then(function() {
      Reveal.slide(${initialSlide});
    });
    window.addEventListener('message', function(e) {
      if (e.data.type === '${POST_MSG_TYPE.Navigate}') Reveal.slide(e.data.index);
    });
  </script>
</body>
</html>`
}

export function buildPreviewHtml(allSlides: Slide[], targetIndex: number): string {
  const slidesHtml = allSlides.map((s) => `<section>${s.content}</section>`).join('\n')
  return `<!DOCTYPE html>
<html>
<head>
  <link rel="stylesheet" href="${CDN_REVEAL_CSS}">
  <link rel="stylesheet" href="${CDN_REVEAL_THEME_CSS}">
</head>
<body>
  <div class="reveal" id="reveal" style="opacity: 0;">
    <div class="slides">
${slidesHtml}
    </div>
  </div>
  <script src="${CDN_REVEAL_JS}"></script>
  <script>
    Reveal.initialize({transition: 'none', progress: false}).then(() => {
      Reveal.slide(${targetIndex});
      Reveal.configure({ transition: 'slide'})
      document.getElementById('reveal').style.opacity = '1';
    });
  </script>
</body>
</html>`
}
