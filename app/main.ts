import { BrowserWindow, app, dialog, ipcMain } from 'electron';
import * as electronReload from 'electron-reload';
import * as path from 'path';
import * as url from 'url';

let window: BrowserWindow;
electronReload(__dirname);

function createWindow () {

  window = new BrowserWindow({
    minWidth: 800,
    minHeight: 600,
    center: true,
    title: 'Bench-Routes - Mark your routes',
    webPreferences: {
      nodeIntegration: true
    },
    hasShadow: true,
    autoHideMenuBar: true
  })

  window.loadURL(url.format({
    pathname: path.join(__dirname, './renderer/templates/index.html'),
    protocol: 'file:',
    slashes: true
  }));

  window.webContents.openDevTools();

  window.on('closed', () => {
    window = null
  });
}

app.on('ready', createWindow);


app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (window === null) {
    createWindow();
  }
});