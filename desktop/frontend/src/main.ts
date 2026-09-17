import './style.css';
import { SelectFolder, SelectGoFile, AnalyzeProject } from '../wailsjs/go/main/App';
document.documentElement.lang = 'en';
document.title = 'Cleanupcheck';

// Restore the appearance before rendering the interface.
type Theme = 'dark' | 'light';
let theme: Theme = 'dark';
try {
  const saved = localStorage.getItem('cleanupcheck.theme');
  if (saved === 'light' || saved === 'dark') theme = saved;
} catch {
  // The theme still works when persistent storage is unavailable.
}
document.documentElement.dataset.theme = theme;




// Render the interface before attaching button handlers.
document.querySelector<HTMLDivElement>('#app')!.innerHTML = `
  <div class="workspace">
    <aside class="sidebar" aria-label="Workspace">
      <div class="brand">cleanupcheck<span class="edition">desktop</span></div>
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
      <footer><span>cleanupcheck</span><span>HTTP response body checks</span></footer>
    </div>
  </div>
`;

const themeButtons = document.querySelectorAll<HTMLButtonElement>('[data-theme-choice]');
function applyTheme(next: Theme) {
  theme = next;
  document.documentElement.dataset.theme = next;
  themeButtons.forEach(button => {
    button.setAttribute('aria-pressed', String(button.dataset.themeChoice === next));
  });
}
applyTheme(theme);
themeButtons.forEach(button => {
  button.addEventListener('click', () => {
    const next = button.dataset.themeChoice;
    if (next !== 'dark' && next !== 'light') return;
    applyTheme(next);
    try {
      localStorage.setItem('cleanupcheck.theme', next);
    } catch {
      // Keep the selected appearance for this session.
    }
  });
});

const folderButton = document.querySelector<HTMLButtonElement>('#select-folder')!;
const fileButton = document.querySelector<HTMLButtonElement>('#select-file')!;
const projectPath = document.querySelector<HTMLInputElement>('#project-path')!;
const connectionNote = document.querySelector<HTMLParagraphElement>('#connection-note')!;

// Set up handlers only after the interface exists.
const analyzeButton =
  document.querySelector<HTMLButtonElement>('#analyze')!;

const results = document.querySelector<HTMLDivElement>('#results')!;
const status = document.querySelector<HTMLElement>('#analysis-status')!;
const count = document.querySelector<HTMLElement>('#finding-count')!;
const detail = document.querySelector<HTMLElement>('#finding-detail')!;

analyzeButton.addEventListener('click', async () => {
  if (!projectPath.value) return;

  analyzeButton.disabled = true;
  folderButton.disabled = true;
  fileButton.disabled = true;
  analyzeButton.textContent = 'Analyzing…';
  results.textContent = 'Checking your project…';
  results.setAttribute('aria-busy', 'true');
  status.textContent = 'Running';
  count.textContent = '—';
  detail.textContent = 'Waiting for results.';

  try {
    const findings = await AnalyzeProject(projectPath.value);

    results.replaceChildren();
    status.textContent = 'Complete';
    count.textContent = String(findings.length);
    detail.textContent = findings.length ? 'Select a finding to inspect its location and message.' : 'No findings to inspect.';

    if (findings.length === 0) {
      results.textContent = 'No issues found.';
    }

    for (const finding of findings) {
      const row = document.createElement('button');
      row.type = 'button';
      row.className = 'finding-row';
      row.setAttribute('aria-pressed', 'false');
      row.addEventListener('click', () => {
        results.querySelectorAll('button').forEach(button => button.setAttribute('aria-pressed', 'false'));
        row.setAttribute('aria-pressed', 'true');
        detail.textContent = `${finding.file}\nLine ${finding.line}\n\n${finding.message}`;
      });

      row.textContent =
        `${finding.file}:${finding.line} — ${finding.message}`;

      results.appendChild(row);
    }
  } catch (error) {
    results.textContent = `Analysis failed: ${String(error)}`;
    status.textContent = 'Failed';
    detail.textContent = 'Analysis could not be completed.';
  } finally {
    results.setAttribute('aria-busy', 'false');
    analyzeButton.disabled = false;
    folderButton.disabled = false;
    fileButton.disabled = false;
    analyzeButton.textContent = 'Analyze';
  }
});

function setSelectedPath(path: string, kind: 'folder' | 'file') {
  projectPath.value = path;
  projectPath.title = path;
  document.querySelector<HTMLElement>('#project-name')!.textContent = path.split(/[\\\\/]/).filter(Boolean).pop() || path;
  results.textContent = `${kind === 'file' ? 'Go file' : 'Project'} selected. Run an analysis to view findings.`;
  status.textContent = 'Not run';
  count.textContent = '—';
  detail.textContent = 'Select a finding to inspect its location and message.';
  analyzeButton.disabled = false;
  connectionNote.textContent = `${kind === 'file' ? 'Go file' : 'Folder'} selected. Ready to analyze.`;
}

folderButton.addEventListener('click', async () => {
  // Prevent opening multiple dialogs while waiting for the user's selection.
  folderButton.disabled = true;
  fileButton.disabled = true;
  folderButton.textContent = 'Choosing…';

  try {
    const folder = await SelectFolder();

    // Cancelling the dialog preserves any previously selected folder.
    if (folder) {
      setSelectedPath(folder, 'folder');
    }
    if (!projectPath.value) connectionNote.textContent = 'No folder selected. Choose a Go project folder or Go file.';
  } catch (error) {
    connectionNote.textContent = `Could not open folder selection: ${String(error)}`;
  } finally {
    folderButton.disabled = false;
    fileButton.disabled = false;
    folderButton.textContent = 'Choose folder';
  }
});

fileButton.addEventListener('click', async () => {
  folderButton.disabled = true;
  fileButton.disabled = true;
  fileButton.textContent = 'Choosing…';

  try {
    const file = await SelectGoFile();
    if (file) {
      setSelectedPath(file, 'file');
    }
    if (!projectPath.value) connectionNote.textContent = 'No file selected. Choose a Go project folder or Go file.';
  } catch (error) {
    connectionNote.textContent = `Could not open file selection: ${String(error)}`;
  } finally {
    folderButton.disabled = false;
    fileButton.disabled = false;
    fileButton.textContent = 'Choose .go file';
  }
});
