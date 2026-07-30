import { CDN_REVEAL_CSS, CDN_REVEAL_THEME_CSS, CDN_REVEAL_JS, POST_MSG_TYPE } from "../cfg";

export const SLIDE_SEPARATOR = /^---+\s*$/m;

export function splitSlides(content: string): string[] {
  return content.split(SLIDE_SEPARATOR).map((p) => p.replace(/^\n+/, ""));
}

export function buildPresentHtml(slides: string[], initialSlide: number): string {
  const slidesHtml = slides.map((content) => `<section>${content}</section>`).join("\n");
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
</html>`;
}

export function buildPreviewHtml(slides: string[], targetIndex: number): string {
  const slidesHtml = slides.map((content) => `<section>${content}</section>`).join("\n");
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
    let pending = null;
    let ready = false;
    function applySlides(slides, index) {
      const container = document.querySelector('.reveal .slides');
      container.innerHTML = slides.map(function (s) { return '<section>' + s + '</section>'; }).join('\\n');
      Reveal.sync();
      Reveal.slide(Math.min(Math.max(index, 0), Math.max(Reveal.getTotalSlides() - 1, 0)));
    }
    window.addEventListener('message', function (e) {
      if (!e.data || e.data.type !== '${POST_MSG_TYPE.UpdateSlides}') return;
      if (ready) applySlides(e.data.slides, e.data.index);
      else pending = e.data;
    });
    Reveal.initialize({ transition: 'slide', progress: false, controls: false, touch: false, scrollActivationWidth: null, }).then(function () {
      ready = true;
      if (pending) { applySlides(pending.slides, pending.index); pending = null; } else Reveal.slide(${targetIndex});
      Reveal.configure({ transition: 'slide'});
      document.getElementById('reveal').style.opacity = '1';
    });
  </script>
</body>
</html>`;
}
