package admin

const configPageHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width,initial-scale=1" />
  <title>Deadliner Config Admin</title>
  <style>
    body { font-family: ui-sans-serif, -apple-system, BlinkMacSystemFont, sans-serif; margin: 0; background: #f4f7f8; color: #10212b; }
    main { max-width: 980px; margin: 0 auto; padding: 24px; }
    .card { background: #fff; border-radius: 16px; padding: 20px; margin-top: 20px; box-shadow: 0 8px 28px rgba(16, 33, 43, 0.08); }
    h1 { margin: 0 0 8px; }
    .grid { display: grid; gap: 16px; grid-template-columns: repeat(auto-fit, minmax(220px, 1fr)); }
    label { display: block; margin-bottom: 6px; font-size: 13px; font-weight: 700; }
    input[type="text"], input[type="number"], input[type="password"] { width: 100%; box-sizing: border-box; padding: 10px 12px; border: 1px solid #c8d4dc; border-radius: 10px; }
    button { border: 0; border-radius: 999px; padding: 12px 18px; background: #0c7a6a; color: #fff; font-weight: 700; cursor: pointer; }
    .status { margin-top: 12px; font-size: 14px; }
    .hint { font-size: 12px; color: #657a87; }
    .mono { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
  </style>
</head>
<body>
<main>
  <h1>Deadliner Config Admin</h1>
  <p>Updates only non-sensitive runtime config from <span class="mono">conf/config.json</span>. Restart is required after saving.</p>

  <section class="card">
    <div class="grid">
      <div>
        <label for="token">Admin Token</label>
        <input id="token" type="password" placeholder="Bearer token for admin API" />
      </div>
    </div>
    <div style="margin-top:16px; display:flex; gap:12px; flex-wrap:wrap;">
      <button id="load">Load Config</button>
      <button id="save">Save Config</button>
    </div>
    <div id="status" class="status"></div>
  </section>

  <section class="card">
    <div class="grid">
      <div><label>Service Name</label><input id="serviceName" type="text" /></div>
      <div><label>Kitex Address</label><input id="serviceAddress" type="text" /></div>
      <div><label>HTTP Address</label><input id="httpAddress" type="text" /></div>
      <div><label>Database Driver</label><input id="databaseDriver" type="text" /></div>
      <div><label>Read Timeout</label><input id="readTimeout" type="number" /></div>
      <div><label>Write Timeout</label><input id="writeTimeout" type="number" /></div>
      <div><label>Idle Timeout</label><input id="idleTimeout" type="number" /></div>
      <div><label>Max Body Bytes</label><input id="maxBodyBytes" type="number" /></div>
      <div><label>Default Pull Limit</label><input id="defaultPullLimit" type="number" /></div>
      <div><label>Max Pull Limit</label><input id="maxPullLimit" type="number" /></div>
      <div><label>HTTP Rate / Minute</label><input id="rateLimitPerMinute" type="number" /></div>
      <div><label>HTTP Burst</label><input id="rateLimitBurst" type="number" /></div>
      <div><label>Auth Rate / Minute</label><input id="authRateLimitPerMinute" type="number" /></div>
      <div><label>Auth Burst</label><input id="authRateLimitBurst" type="number" /></div>
      <div><label>Sync Rate / Minute</label><input id="syncRateLimitPerMinute" type="number" /></div>
      <div><label>Sync Burst</label><input id="syncRateLimitBurst" type="number" /></div>
      <div><label>Admin Base Path</label><input id="adminBasePath" type="text" /></div>
      <div>
        <label>Admin Enabled</label>
        <input id="adminEnabled" type="checkbox" />
      </div>
    </div>
    <p class="hint">Secret values stay in <span class="mono">conf/secret.json</span> and are never returned by this page.</p>
    <pre id="secretStatus" class="mono"></pre>
  </section>
</main>
<script>
const statusEl = document.getElementById("status");
const secretStatusEl = document.getElementById("secretStatus");

function headers() {
  const token = document.getElementById("token").value.trim();
  const result = { "Content-Type": "application/json" };
  if (token) result["Authorization"] = "Bearer " + token;
  return result;
}

function setStatus(message) {
  statusEl.textContent = message;
}

function writeForm(data) {
  document.getElementById("serviceName").value = data.service.name || "";
  document.getElementById("serviceAddress").value = data.service.address || "";
  document.getElementById("httpAddress").value = data.http.address || "";
  document.getElementById("databaseDriver").value = data.database.driver || "";
  document.getElementById("readTimeout").value = data.http.readTimeoutSeconds || 0;
  document.getElementById("writeTimeout").value = data.http.writeTimeoutSeconds || 0;
  document.getElementById("idleTimeout").value = data.http.idleTimeoutSeconds || 0;
  document.getElementById("maxBodyBytes").value = data.http.maxRequestBodyBytes || 0;
  document.getElementById("defaultPullLimit").value = data.sync.defaultPullLimit || 0;
  document.getElementById("maxPullLimit").value = data.sync.maxPullLimit || 0;
  document.getElementById("rateLimitPerMinute").value = data.http.rateLimitPerMinute || 0;
  document.getElementById("rateLimitBurst").value = data.http.rateLimitBurst || 0;
  document.getElementById("authRateLimitPerMinute").value = data.http.authRateLimitPerMinute || 0;
  document.getElementById("authRateLimitBurst").value = data.http.authRateLimitBurst || 0;
  document.getElementById("syncRateLimitPerMinute").value = data.http.syncRateLimitPerMinute || 0;
  document.getElementById("syncRateLimitBurst").value = data.http.syncRateLimitBurst || 0;
  document.getElementById("adminBasePath").value = data.admin.basePath || "/admin";
  document.getElementById("adminEnabled").checked = !!data.admin.enabled;
  secretStatusEl.textContent = JSON.stringify(data.secretStatus, null, 2);
}

function readForm() {
  return {
    service: {
      name: document.getElementById("serviceName").value,
      address: document.getElementById("serviceAddress").value
    },
    http: {
      address: document.getElementById("httpAddress").value,
      readTimeoutSeconds: Number(document.getElementById("readTimeout").value),
      writeTimeoutSeconds: Number(document.getElementById("writeTimeout").value),
      idleTimeoutSeconds: Number(document.getElementById("idleTimeout").value),
      maxRequestBodyBytes: Number(document.getElementById("maxBodyBytes").value),
      rateLimitPerMinute: Number(document.getElementById("rateLimitPerMinute").value),
      rateLimitBurst: Number(document.getElementById("rateLimitBurst").value),
      authRateLimitPerMinute: Number(document.getElementById("authRateLimitPerMinute").value),
      authRateLimitBurst: Number(document.getElementById("authRateLimitBurst").value),
      syncRateLimitPerMinute: Number(document.getElementById("syncRateLimitPerMinute").value),
      syncRateLimitBurst: Number(document.getElementById("syncRateLimitBurst").value)
    },
    database: {
      driver: document.getElementById("databaseDriver").value
    },
    sync: {
      defaultPullLimit: Number(document.getElementById("defaultPullLimit").value),
      maxPullLimit: Number(document.getElementById("maxPullLimit").value)
    },
    admin: {
      enabled: document.getElementById("adminEnabled").checked,
      basePath: document.getElementById("adminBasePath").value
    }
  };
}

async function loadConfig() {
  setStatus("Loading...");
  const response = await fetch("./api/config", { headers: headers() });
  const payload = await response.json();
  if (!response.ok) {
    setStatus(payload.error + (payload.request_id ? " (request_id=" + payload.request_id + ")" : ""));
    return;
  }
  writeForm(payload);
  setStatus("Config loaded.");
}

async function saveConfig() {
  setStatus("Saving...");
  const response = await fetch("./api/config", {
    method: "PUT",
    headers: headers(),
    body: JSON.stringify(readForm())
  });
  const payload = await response.json();
  if (!response.ok) {
    setStatus(payload.error + (payload.request_id ? " (request_id=" + payload.request_id + ")" : ""));
    return;
  }
  writeForm(payload);
  setStatus("Config saved. Restart the server to apply the change.");
}

document.getElementById("load").addEventListener("click", loadConfig);
document.getElementById("save").addEventListener("click", saveConfig);
</script>
</body>
</html>`
