// terminal.js — streaming terminal panel for towr dashboard
(function() {
  "use strict";
  var MAX_LINES = 500;
  var RETRY_MS = 3000;
  var ANSI_RE = /\x1B\[[0-9;]*[a-zA-Z]/g;
  var OSC_RE = /\x1B\].*?\x07/g;

  var evtSource = null;
  var retryTimer = null;
  var lines = [];

  var panel = document.getElementById("termPanel");
  var body = document.getElementById("termBody");
  var titleEl = document.getElementById("termTitle");
  var closeBtn = document.getElementById("termClose");

  function stripAnsi(s) { return s.replace(ANSI_RE, "").replace(OSC_RE, ""); }

  function esc(s) {
    var d = document.createElement("span");
    d.textContent = s;
    return d.innerHTML;
  }

  function classifyLine(text) {
    if (/Error|FAIL/i.test(text)) return "term-error";
    if (/✓|PASS/.test(text)) return "term-pass";
    if (/^-{3,}$/.test(text.trim())) return "term-dim";
    return "";
  }

  function renderLines() {
    var html = "";
    for (var i = 0; i < lines.length; i++) {
      var cls = classifyLine(lines[i]);
      html += cls
        ? '<div class="' + cls + '">' + esc(lines[i]) + "</div>"
        : "<div>" + esc(lines[i]) + "</div>";
    }
    body.innerHTML = html;
    body.scrollTop = body.scrollHeight;
  }

  function replaceSnapshot(raw) {
    var clean = stripAnsi(raw);
    lines = clean.split("\n");
    if (lines.length > MAX_LINES) lines = lines.slice(lines.length - MAX_LINES);
    renderLines();
  }

  function setStatus(msg) {
    body.innerHTML = '<div class="term-dim">' + esc(msg) + "</div>";
  }

  function disconnect() {
    if (retryTimer) { clearTimeout(retryTimer); retryTimer = null; }
    if (evtSource) { evtSource.close(); evtSource = null; }
  }

  function connect(id) {
    disconnect();
    lines = [];
    setStatus("Connecting\u2026");

    evtSource = new EventSource("/stream/" + encodeURIComponent(id));
    evtSource.onopen = function() {
      lines = [];
      body.innerHTML = "";
    };
    evtSource.onmessage = function(e) { replaceSnapshot(e.data); };
    evtSource.onerror = function() {
      if (evtSource.readyState === EventSource.CLOSED) {
        setStatus("Stream ended");
        evtSource = null;
      } else {
        setStatus("Disconnected — retrying\u2026");
        evtSource.close();
        evtSource = null;
        retryTimer = setTimeout(function() { connect(id); }, RETRY_MS);
      }
    };
  }

  // --- public API exposed on window ---
  window.openTerminal = function(id, name) {
    panel.classList.add("open");
    titleEl.innerHTML =
      '<span class="dot"></span> ' + esc(name || id);
    connect(id);
    // highlight active card
    document.querySelectorAll(".card").forEach(function(el) {
      el.classList.toggle("active", el.getAttribute("data-id") === id);
    });
  };

  // close button
  closeBtn.addEventListener("click", closeTerminal);

  function closeTerminal() {
    panel.classList.remove("open");
    disconnect();
    lines = [];
    body.innerHTML = "";
    document.querySelectorAll(".card.active").forEach(function(el) {
      el.classList.remove("active");
    });
    window.activeTerminalId = null;
  }

  // Escape key closes terminal
  document.addEventListener("keydown", function(e) {
    if (e.key === "Escape" && panel.classList.contains("open")) closeTerminal();
  });

  // Clear button
  var clearBtn = document.getElementById("termClear");
  if (clearBtn) {
    clearBtn.addEventListener("click", function() {
      lines = [];
      body.innerHTML = "";
    });
  }
})();
