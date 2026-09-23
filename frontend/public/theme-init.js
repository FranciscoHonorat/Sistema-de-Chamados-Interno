// Applies the saved (or the system) theme before the first paint, so a dark
// theme never flashes white. Kept out of index.html: the CSP forbids inline scripts.
;(function () {
  var mode = 'system'
  try {
    mode = localStorage.getItem('sys-called:theme') || 'system'
  } catch {
    /* storage unavailable: follow the system */
  }
  var dark = mode === 'dark' || (mode === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)
  document.documentElement.classList.toggle('dark', dark)
})()
