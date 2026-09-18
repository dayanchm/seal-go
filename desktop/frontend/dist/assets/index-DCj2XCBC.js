(function(){const t=document.createElement("link").relList;if(t&&t.supports&&t.supports("modulepreload"))return;for(const n of document.querySelectorAll('link[rel="modulepreload"]'))p(n);new MutationObserver(n=>{for(const a of n)if(a.type==="childList")for(const m of a.addedNodes)m.tagName==="LINK"&&m.rel==="modulepreload"&&p(m)}).observe(document,{childList:!0,subtree:!0});function o(n){const a={};return n.integrity&&(a.integrity=n.integrity),n.referrerPolicy&&(a.referrerPolicy=n.referrerPolicy),n.crossOrigin==="use-credentials"?a.credentials="include":n.crossOrigin==="anonymous"?a.credentials="omit":a.credentials="same-origin",a}function p(n){if(n.ep)return;n.ep=!0;const a=o(n);fetch(n.href,a)}})();function C(e){return window.go.main.App.AnalyzeProject(e)}function x(){return window.go.main.App.SelectFolder()}function S(){return window.go.main.App.SelectGoFile()}document.documentElement.lang="en";document.title="seal-go";let h="dark";try{const e=localStorage.getItem("seal-go.theme");(e==="light"||e==="dark")&&(h=e)}catch{}document.documentElement.dataset.theme=h;document.querySelector("#app").innerHTML=`
  <div class="workspace">
    <aside class="sidebar" aria-label="Workspace">
      <div class="brand">seal-go<span class="edition">desktop</span></div>
      <div class="sidebar-section">Workspace</div>
      <a class="nav-item active" href="#analysis" aria-current="page">Analysis <span aria-hidden="true">↗</span></a>
      <div class="sidebar-project">
        <span class="sidebar-section">Current project</span>
        <p id="project-name">No project open</p>
      </div>
      <div class="sidebar-bottom">On your machine<span>Go resource analyzer</span></div>
    </aside>
    <div class="main-column">
      <header class="topbar"><span>Workspace <span class="separator">/</span> <strong>Analysis</strong></span><div class="theme-switch" role="group" aria-label="Appearance">
        <button type="button" data-theme-choice="light" aria-pressed="false" aria-label="Light theme" title="Light theme">
          <svg aria-hidden="true" viewBox="0 0 24 24"><circle cx="12" cy="12" r="4"></circle><path d="M12 2v2"></path><path d="M12 20v2"></path><path d="m4.93 4.93 1.41 1.41"></path><path d="m17.66 17.66 1.41 1.41"></path><path d="M2 12h2"></path><path d="M20 12h2"></path><path d="m6.34 17.66-1.41 1.41"></path><path d="m19.07 4.93-1.41 1.41"></path></svg>
        </button>
        <button type="button" data-theme-choice="dark" aria-pressed="false" aria-label="Dark theme" title="Dark theme">
          <svg aria-hidden="true" viewBox="0 0 24 24"><path d="M20.2 14.4A8.4 8.4 0 0 1 9.6 3.8a8.5 8.5 0 1 0 10.6 10.6Z"></path></svg>
        </button>
      </div></header>
      <main id="analysis">
        <div class="page-heading"><h1>Analysis</h1><p>Inspect HTTP response cleanup in your Go project.</p></div>
        <section class="project-panel" aria-labelledby="project-title">
          <label id="project-title" for="project-path">Project path</label>
          <div class="folder-row">
            <input id="project-path" type="text" placeholder="Select a Go project or .go file…" readonly aria-describedby="connection-note" />
            <button id="select-folder" class="button secondary">Choose folder</button>
            <button id="select-file" class="button secondary">Choose .go file</button>
            <button id="analyze" class="button primary" disabled>Analyze</button>
          </div>
          <p id="connection-note" role="status">Choose a folder or Go file to get started.</p>
        </section>
        <section class="results-panel" aria-labelledby="results-title">
          <div class="results-heading"><h2 id="results-title">Findings <span id="finding-count">—</span></h2><span id="analysis-status">Not run</span></div>
          <div class="results-layout">
            <div id="results" aria-live="polite">
              <div class="empty-state"><h3>No analysis yet</h3><p>Select a project to begin.<br>Findings will appear here after an analysis.</p></div>
            </div>
            <aside class="detail-panel" aria-labelledby="detail-title"><h2 id="detail-title">Details</h2><p id="finding-detail">Select a finding to inspect its location and message.</p></aside>
          </div>
        </section>
      </main>
      <footer><span>seal-go</span><span>HTTP response body checks</span></footer>
    </div>
  </div>
`;const g=document.querySelectorAll("[data-theme-choice]");function b(e){h=e,document.documentElement.dataset.theme=e,g.forEach(t=>{t.setAttribute("aria-pressed",String(t.dataset.themeChoice===e))})}b(h);g.forEach(e=>{e.addEventListener("click",()=>{const t=e.dataset.themeChoice;if(!(t!=="dark"&&t!=="light")){b(t);try{localStorage.setItem("seal-go.theme",t)}catch{}}})});const i=document.querySelector("#select-folder"),l=document.querySelector("#select-file"),c=document.querySelector("#project-path"),u=document.querySelector("#connection-note"),r=document.querySelector("#analyze"),s=document.querySelector("#results"),f=document.querySelector("#analysis-status"),y=document.querySelector("#finding-count"),d=document.querySelector("#finding-detail");r.addEventListener("click",async()=>{if(c.value){r.disabled=!0,i.disabled=!0,l.disabled=!0,r.textContent="Analyzing…",s.textContent="Checking your project…",s.setAttribute("aria-busy","true"),f.textContent="Running",y.textContent="—",d.textContent="Waiting for results.";try{const e=await C(c.value);s.replaceChildren(),f.textContent="Complete",y.textContent=String(e.length),d.textContent=e.length?"Select a finding to inspect its location and message.":"No findings to inspect.",e.length===0&&(s.textContent="No issues found.");for(const t of e){const o=document.createElement("button");o.type="button",o.className="finding-row",o.setAttribute("aria-pressed","false"),o.addEventListener("click",()=>{s.querySelectorAll("button").forEach(p=>p.setAttribute("aria-pressed","false")),o.setAttribute("aria-pressed","true"),d.textContent=`${t.file}
Line ${t.line}

${t.message}`}),o.textContent=`${t.file}:${t.line} — ${t.message}`,s.appendChild(o)}}catch(e){s.textContent=`Analysis failed: ${String(e)}`,f.textContent="Failed",d.textContent="Analysis could not be completed."}finally{s.setAttribute("aria-busy","false"),r.disabled=!1,i.disabled=!1,l.disabled=!1,r.textContent="Analyze"}}});function v(e,t){c.value=e,c.title=e,document.querySelector("#project-name").textContent=e.split(/[\\\\/]/).filter(Boolean).pop()||e,s.textContent=`${t==="file"?"Go file":"Project"} selected. Run an analysis to view findings.`,f.textContent="Not run",y.textContent="—",d.textContent="Select a finding to inspect its location and message.",r.disabled=!1,u.textContent=`${t==="file"?"Go file":"Folder"} selected. Ready to analyze.`}i.addEventListener("click",async()=>{i.disabled=!0,l.disabled=!0,i.textContent="Choosing…";try{const e=await x();e&&v(e,"folder"),c.value||(u.textContent="No folder selected. Choose a Go project folder or Go file.")}catch(e){u.textContent=`Could not open folder selection: ${String(e)}`}finally{i.disabled=!1,l.disabled=!1,i.textContent="Choose folder"}});l.addEventListener("click",async()=>{i.disabled=!0,l.disabled=!0,l.textContent="Choosing…";try{const e=await S();e&&v(e,"file"),c.value||(u.textContent="No file selected. Choose a Go project folder or Go file.")}catch(e){u.textContent=`Could not open file selection: ${String(e)}`}finally{i.disabled=!1,l.disabled=!1,l.textContent="Choose .go file"}});
