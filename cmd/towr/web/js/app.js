(function() {
  "use strict";
  var STATUS_COLORS = {
    RUNNING: "#58a6ff", SPAWNED: "#58a6ff",
    READY: "#3fb950", MERGED: "#3fb950", LANDED: "#3fb950",
    BLOCKED: "#f85149", FAILED: "#f85149", ERROR: "#f85149",
    STALE: "#8b949e", ORPHANED: "#8b949e", IDLE: "#8b949e"
  };
  var DEFAULT_COLOR = "#8b949e";
  var PAGE_LOAD = Date.now();
  var lastJSON = "";
  var safetyCache = {};

  function statusColor(s) { return STATUS_COLORS[(s||"").toUpperCase()] || DEFAULT_COLOR; }

  function zone(status) {
    var s = (status||"").toUpperCase();
    if (s === "RUNNING" || s === "SPAWNED") return "working";
    if (s === "BLOCKED" || s === "FAILED" || s === "ERROR" || s === "STALE" || s === "ORPHANED") return "attention";
    if (s === "READY" || s === "MERGED" || s === "LANDED") return "completed";
    return "working";
  }

  function badgeBg(color) { return color + "22"; }

  function esc(s) {
    var d = document.createElement("span");
    d.textContent = s;
    return d.innerHTML;
  }

  function sortPriority(status) {
    var s = (status || "").toUpperCase();
    if (s === "BLOCKED" || s === "FAILED" || s === "ERROR") return 0;
    if (s === "RUNNING" || s === "SPAWNED") return 1;
    if (s === "STALE" || s === "ORPHANED" || s === "IDLE") return 2;
    return 3;
  }

  function isCompleted(status) {
    var s = (status || "").toUpperCase();
    return s === "READY" || s === "MERGED" || s === "LANDED";
  }

  function parseAgeMinutes(age) {
    if (!age) return 0;
    var m = 0;
    var hMatch = age.match(/(\d+)h/);
    var mMatch = age.match(/(\d+)m/);
    if (hMatch) m += parseInt(hMatch[1], 10) * 60;
    if (mMatch) m += parseInt(mMatch[1], 10);
    return m;
  }

  function completedOpacity(age) {
    var mins = parseAgeMinutes(age);
    if (mins < 30) return 0.7;
    if (mins <= 120) return 0.5;
    return 0.3;
  }

  function uptimeStr() {
    var sec = Math.floor((Date.now() - PAGE_LOAD) / 1000);
    var m = Math.floor(sec / 60);
    var s = sec % 60;
    return m > 0 ? m + "m " + s + "s" : s + "s";
  }

  // Safety shields
  function fetchSafety(id) {
    fetch("/api/workspace/" + encodeURIComponent(id) + "/safety")
      .then(function(r) { return r.json(); })
      .then(function(s) {
        safetyCache[id] = s;
        var el = document.querySelector('[data-shield="' + id + '"]');
        if (el) renderShield(el, s);
      })
      .catch(function() {});
  }

  function renderShield(el, s) {
    if (s.bypasses > 0) {
      el.textContent = "\uD83D\uDEE1\uFE0F bypass";
      el.style.color = "#f85149";
    } else if (s.approvals > 0) {
      el.textContent = "\uD83D\uDEE1\uFE0F " + s.approvals + " approved";
      el.style.color = "#d29922";
    } else {
      el.textContent = "\uD83D\uDEE1\uFE0F sandboxed";
      el.style.color = "#3fb950";
    }
  }

  // Score cards
  function renderCounters(data) {
    var total = data.length;
    var working = 0, blocked = 0, completed = 0;
    var totalApprovals = 0, totalBypasses = 0;
    data.forEach(function(ws) {
      var s = (ws.status||"").toUpperCase();
      if (s === "RUNNING" || s === "SPAWNED") working++;
      else if (s === "BLOCKED" || s === "FAILED" || s === "ERROR") blocked++;
      else if (isCompleted(ws.status)) completed++;
    });
    Object.keys(safetyCache).forEach(function(k) {
      totalApprovals += safetyCache[k].approvals || 0;
      totalBypasses += safetyCache[k].bypasses || 0;
    });
    var bypassColor = totalBypasses > 0 ? "#f85149" : "#3fb950";
    var bypassClass = totalBypasses > 0 ? " stat-pulse" : "";
    var bar = document.getElementById("statsBar");
    bar.innerHTML =
      '<span class="stat-pill" style="color:#c9d1d9;background:#30363d">' + total + " total</span>" +
      '<span class="stat-pill" style="color:#58a6ff;background:#58a6ff22">' + working + " working</span>" +
      '<span class="stat-pill" style="color:#f85149;background:#f8514922">' + blocked + " blocked</span>" +
      '<span class="stat-pill" style="color:#3fb950;background:#3fb95022">' + completed + " completed</span>" +
      '<span class="stat-pill' + bypassClass + '" style="color:' + bypassColor + ";background:" + bypassColor + '22">' + totalBypasses + " bypasses</span>" +
      '<span class="stat-pill" style="color:#3fb950;background:#3fb95022">' + totalApprovals + " approvals</span>" +
      '<span class="stat-meta">uptime ' + esc(uptimeStr()) + "</span>";
  }

  // Workspace list with zones
  function renderWorkspaces(data) {
    var sorted = data.slice().sort(function(a, b) {
      return sortPriority(a.status) - sortPriority(b.status);
    });

    var groups = { working: [], attention: [], completed: [] };
    sorted.forEach(function(ws) { groups[zone(ws.status)].push(ws); });

    var zoneConfig = [
      { key: "attention", title: "Needs Attention", color: "#f85149" },
      { key: "working", title: "Working", color: "#58a6ff" },
      { key: "completed", title: "Completed", color: "#3fb950" }
    ];

    var html = "";
    var hasAny = false;
    zoneConfig.forEach(function(z) {
      var items = groups[z.key];
      if (items.length === 0) return;
      hasAny = true;
      html += '<div class="zone">';
      html += '<div class="zone-title" style="color:' + esc(z.color) + '">' + esc(z.title) +
              ' <span class="count">(' + items.length + ')</span></div>';
      html += '<div class="cards">';
      items.forEach(function(ws) {
        var c = statusColor(ws.status);
        var opacity = isCompleted(ws.status) ? completedOpacity(ws.age) : 1;
        html += '<div class="card" data-id="' + esc(ws.id) + '" style="border-left:3px solid ' + c +
          ";opacity:" + opacity + '">';
        html += '<div class="card-top">';
        html += '<span class="card-id">' + esc(ws.id) + '</span>';
        html += '<span class="shield" data-shield="' + esc(ws.id) + '"></span>';
        html += '<span class="badge" style="color:' + esc(c) + ';background:' + badgeBg(c) + '">' + esc(ws.status) + '</span>';
        html += '</div>';
        html += '<div class="card-details">';
        html += '<span><span class="label">task</span> ' + esc(ws.task) + '</span>';
        html += '<span><span class="label">diff</span> ' + esc(ws.diff) + '</span>';
        html += '<span><span class="label">age</span> ' + esc(ws.age) + '</span>';
        html += '</div>';
        html += '</div>';
      });
      html += '</div></div>';
    });
    if (!hasAny) {
      html = '<div class="empty-state">No workspaces found.</div>';
    }
    document.getElementById("sidebar").innerHTML = html;

    // Bind card clicks
    document.querySelectorAll(".card").forEach(function(el) {
      el.addEventListener("click", function() {
        var id = el.getAttribute("data-id");
        if (typeof window.expandWorkspace === "function") {
          window.expandWorkspace(id);
        }
      });
    });
  }

  // Main render with DOM caching
  function render(data) {
    var json = JSON.stringify(data);
    if (json === lastJSON) {
      var meta = document.querySelector(".stat-meta");
      if (meta) meta.textContent = "uptime " + uptimeStr();
      return;
    }
    lastJSON = json;
    renderCounters(data);
    renderWorkspaces(data);
    (data || []).forEach(function(ws) { fetchSafety(ws.id); });
  }

  // Activity log
  var EVENT_COLORS = {
    "task.completed": "#3fb950", "task.failed": "#f85149",
    "task.dispatched": "#58a6ff", "task.blocked": "#d29922"
  };

  document.getElementById("actToggle").addEventListener("click", function() {
    this.classList.toggle("open");
    document.getElementById("actFeed").classList.toggle("open");
  });

  function renderEvents(events) {
    var feed = document.getElementById("actFeed");
    var countEl = document.getElementById("actCount");
    countEl.textContent = "(" + (events||[]).length + " events)";
    var html = "";
    (events||[]).forEach(function(ev) {
      var ts = new Date(ev.ts).toLocaleTimeString();
      var c = EVENT_COLORS[ev.kind] || "#8b949e";
      var summary = "";
      if (ev.data && ev.data.summary) summary = ev.data.summary;
      else if (ev.data && ev.data.message) summary = ev.data.message;
      html += '<div class="evt-row">';
      html += '<span class="evt-ts">' + esc(ts) + '</span>';
      html += '<span class="evt-ws">' + esc(ev.workspace_id||"-") + '</span>';
      html += '<span class="evt-kind" style="color:' + c + '">' + esc(ev.kind) + '</span>';
      html += '<span>' + esc(summary) + '</span>';
      html += '</div>';
    });
    feed.innerHTML = html;
  }

  // Export audit
  var exportBtn = document.getElementById("exportAudit");
  if (exportBtn) {
    exportBtn.addEventListener("click", function() {
      window.location.href = "/api/audit/export?format=csv&since=7d";
    });
  }

  // Header meta
  var headerMeta = document.querySelector(".header-meta");
  if (headerMeta) {
    headerMeta.innerHTML = '<span class="dot"></span>live &middot; refreshing every 4s';
  }

  // Poll loop
  function poll() {
    fetch("/api/workspaces").then(function(r) { return r.json(); }).then(render).catch(function() {});
    fetch("/api/events").then(function(r) { return r.json(); }).then(renderEvents).catch(function() {});
    setTimeout(poll, 4000);
  }
  poll();
})();
