const electron = require('electron');
const app = electron.app;
const BrowserWindow = electron.BrowserWindow;

const path = require('path');
const isDev = require('electron-is-dev');

require('electron-reload')(__dirname);

let mainWindow;

function createWindow() {

    mainWindow = new BrowserWindow({
        minWidth: 1000,
        minHeight: 600,
        center: true,
        title: 'Bench-Routes - Mark your routes',
        webPreferences: {
          nodeIntegration: true
        },
        hasShadow: true,
        autoHideMenuBar: true,
        transparent: false
    });

    mainWindow.loadURL(isDev ? 'http://localhost:3000' : `file://${path.join(__dirname, '../build/index.html')}`);

    mainWindow.webContents.openDevTools();

    mainWindow.on('closed', () => mainWindow = null);

}

app.on('ready', createWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (mainWindow === null) {
    createWindow();
  }
});
