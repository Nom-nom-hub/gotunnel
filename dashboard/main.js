const { app, BrowserWindow, ipcMain, dialog, shell } = require('electron');
const path = require('path');
const { spawn } = require('child_process');
const fs = require('fs');

let mainWindow;
let serverProcess = null;
let clientProcesses = new Map();

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false
    },
    icon: path.join(__dirname, 'assets/icon.png'),
    title: 'GoTunnel Dashboard - The open-source ngrok killer',
    show: false
  });

  mainWindow.loadFile('index.html');

  mainWindow.once('ready-to-show', () => {
    mainWindow.show();
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
  });
}

// IPC Handlers
ipcMain.handle('start-server', async (event, config) => {
  try {
    const serverPath = path.join(__dirname, '..', 'gotunnel-server.exe');
    
    if (!fs.existsSync(serverPath)) {
      throw new Error('Server binary not found. Please build the project first.');
    }

    const args = [
      '--port', config.port || '8080',
      '--allowed-tokens', config.token || 'default-token',
      '--tls', config.tls ? 'true' : 'false',
      '--log-level', config.logLevel || 'info'
    ];

    serverProcess = spawn(serverPath, args);

    serverProcess.stdout.on('data', (data) => {
      mainWindow.webContents.send('server-log', data.toString());
    });

    serverProcess.stderr.on('data', (data) => {
      mainWindow.webContents.send('server-error', data.toString());
    });

    serverProcess.on('close', (code) => {
      mainWindow.webContents.send('server-closed', code);
      serverProcess = null;
    });

    return { success: true, message: 'Server started successfully' };
  } catch (error) {
    return { success: false, message: error.message };
  }
});

ipcMain.handle('stop-server', async () => {
  if (serverProcess) {
    serverProcess.kill();
    serverProcess = null;
    return { success: true, message: 'Server stopped' };
  }
  return { success: false, message: 'No server running' };
});

ipcMain.handle('start-client', async (event, config) => {
  try {
    const clientPath = path.join(__dirname, '..', 'gotunnel-client.exe');
    
    if (!fs.existsSync(clientPath)) {
      throw new Error('Client binary not found. Please build the project first.');
    }

    const args = [
      '--server', config.server || 'localhost:8080',
      '--subdomain', config.subdomain,
      '--local-port', config.localPort.toString(),
      '--local-host', config.localHost || 'localhost',
      '--token', config.token || 'default-token',
      '--tls', config.tls ? 'true' : 'false',
      '--log-level', config.logLevel || 'info'
    ];

    const clientProcess = spawn(clientPath, args);
    const clientId = `${config.subdomain}-${Date.now()}`;
    
    clientProcesses.set(clientId, clientProcess);

    clientProcess.stdout.on('data', (data) => {
      mainWindow.webContents.send('client-log', { id: clientId, data: data.toString() });
    });

    clientProcess.stderr.on('data', (data) => {
      mainWindow.webContents.send('client-error', { id: clientId, data: data.toString() });
    });

    clientProcess.on('close', (code) => {
      mainWindow.webContents.send('client-closed', { id: clientId, code });
      clientProcesses.delete(clientId);
    });

    return { success: true, message: 'Client started successfully', id: clientId };
  } catch (error) {
    return { success: false, message: error.message };
  }
});

ipcMain.handle('stop-client', async (event, clientId) => {
  const clientProcess = clientProcesses.get(clientId);
  if (clientProcess) {
    clientProcess.kill();
    clientProcesses.delete(clientId);
    return { success: true, message: 'Client stopped' };
  }
  return { success: false, message: 'Client not found' };
});

ipcMain.handle('get-tunnel-url', async (event, subdomain) => {
  return `http://${subdomain}.localhost:8080`;
});

ipcMain.handle('open-url', async (event, url) => {
  await shell.openExternal(url);
});

ipcMain.handle('select-file', async () => {
  const result = await dialog.showOpenDialog(mainWindow, {
    properties: ['openFile'],
    filters: [
      { name: 'Executables', extensions: ['exe'] },
      { name: 'All Files', extensions: ['*'] }
    ]
  });
  return result.filePaths[0];
});

ipcMain.handle('get-config-path', async () => {
  const configDir = path.join(process.env.USERPROFILE || process.env.HOME, '.gotunnel');
  return path.join(configDir, 'config.yaml');
});

ipcMain.handle('load-config', async (event, configPath) => {
  try {
    if (fs.existsSync(configPath)) {
      const yaml = require('js-yaml');
      const configContent = fs.readFileSync(configPath, 'utf8');
      return yaml.load(configContent);
    }
    return null;
  } catch (error) {
    console.error('Failed to load config:', error);
    return null;
  }
});

ipcMain.handle('save-config', async (event, config) => {
  try {
    const yaml = require('js-yaml');
    const configDir = path.join(process.env.USERPROFILE || process.env.HOME, '.gotunnel');
    const configPath = path.join(configDir, 'config.yaml');
    
    if (!fs.existsSync(configDir)) {
      fs.mkdirSync(configDir, { recursive: true });
    }
    
    const configContent = yaml.dump(config);
    fs.writeFileSync(configPath, configContent, 'utf8');
    return { success: true };
  } catch (error) {
    console.error('Failed to save config:', error);
    return { success: false, error: error.message };
  }
});

ipcMain.handle('load-profiles', async () => {
  try {
    const profilesDir = path.join(process.env.USERPROFILE || process.env.HOME, '.gotunnel', 'profiles');
    const profilesPath = path.join(profilesDir, 'profiles.json');
    
    if (fs.existsSync(profilesPath)) {
      const profilesContent = fs.readFileSync(profilesPath, 'utf8');
      return JSON.parse(profilesContent);
    }
    return [];
  } catch (error) {
    console.error('Failed to load profiles:', error);
    return [];
  }
});

ipcMain.handle('save-profile', async (event, profile) => {
  try {
    const profilesDir = path.join(process.env.USERPROFILE || process.env.HOME, '.gotunnel', 'profiles');
    const profilesPath = path.join(profilesDir, 'profiles.json');
    
    if (!fs.existsSync(profilesDir)) {
      fs.mkdirSync(profilesDir, { recursive: true });
    }
    
    let profiles = [];
    if (fs.existsSync(profilesPath)) {
      const profilesContent = fs.readFileSync(profilesPath, 'utf8');
      profiles = JSON.parse(profilesContent);
    }
    
    profiles.push(profile);
    fs.writeFileSync(profilesPath, JSON.stringify(profiles, null, 2), 'utf8');
    return { success: true };
  } catch (error) {
    console.error('Failed to save profile:', error);
    return { success: false, error: error.message };
  }
});

ipcMain.handle('delete-profile', async (event, index) => {
  try {
    const profilesDir = path.join(process.env.USERPROFILE || process.env.HOME, '.gotunnel', 'profiles');
    const profilesPath = path.join(profilesDir, 'profiles.json');
    
    if (fs.existsSync(profilesPath)) {
      const profilesContent = fs.readFileSync(profilesPath, 'utf8');
      let profiles = JSON.parse(profilesContent);
      
      if (index >= 0 && index < profiles.length) {
        profiles.splice(index, 1);
        fs.writeFileSync(profilesPath, JSON.stringify(profiles, null, 2), 'utf8');
      }
    }
    return { success: true };
  } catch (error) {
    console.error('Failed to delete profile:', error);
    return { success: false, error: error.message };
  }
});

app.whenReady().then(createWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

app.on('before-quit', () => {
  // Clean up processes
  if (serverProcess) {
    serverProcess.kill();
  }
  clientProcesses.forEach(process => process.kill());
}); 