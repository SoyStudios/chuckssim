// This file is required by the index.html file and will
// be executed in the renderer process for that window.
// All of the Node.js APIs are available in this process.

let electronify = require('electronify-server');

import {CHUCKSSIM_CONFIG} from "./config";

const backendUrl = CHUCKSSIM_CONFIG.simulatorUrl;

const startServerCommand = CHUCKSSIM_CONFIG.startServerCommand;

let connection;

let isConnected = false;

const stateQueue = [];

let tickerHandler;


electronify({
    command: startServerCommand,
    // url: 'http://127.0.0.1:8000',
    debug: true,
    window: {height: 768, width: 1024},
    ready: function(app){
        // application event listeners could be added here
    },
    preLoad: function(app, window){
        // window event listeners could be added here
    },
    postLoad: function(app, window){
        // url finished loading
    },
    showDevTools: false
}).on('child-started', function(child) {
    // child process has started
    console.log('PID: ' + child.pid);

    // setup logging on child process
    child.stdout.on('data', console.log);
    child.stderr.on('data', console.log);
    connection = new WebSocket('ws://' + backendUrl, []);

}).on('child-closed', function(app, stderr, stdout) {
    // the child process has finished

}).on('child-error', function(err, app) {
    // close electron if the child crashes
    // app.quit();
});

addStateToQueue = (state) => {
    stateQueue.length = 5;
    stateQueue.unshift(state);
};

// When the connection is open, send some data to the server
connection.onopen = function () {
    // connection.send('Ping'); // Send the message 'Ping' to the server
    isConnected = true;
    tickerHandler = renderTicker();
};

// Log errors
connection.onerror = function (error) {
    isConnected = false;
    console.log('WebSocket Error ' + error);
};

// Log messages from the server
connection.onmessage = function (e) {
    console.log('Server: ' + e.data);
    if (e.data.type === 'state') {
        addStateToQueue(e.data);
    }
};


renderState = (state) => {
  // do stuff to render on canvas
};

togglePause = () => {
    if (isConnected) {
        connection.send('togglePause')
    }
};


renderTicker = () => {
  return setInterval(
      () => {
          if (stateQueue.length > 0){
              renderState(stateQueue.shift());
          }
      }, 10000
  )
};

stopRenderTicker = () => {
  if (tickerHandler) {
      window.clearInterval(tickerHandler);
  }
};



