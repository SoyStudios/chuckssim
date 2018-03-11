// This file is required by the index.html file and will
// be executed in the renderer process for that window.
// All of the Node.js APIs are available in this process.

import {CHUCKSSIM_CONFIG} from "./config";

const backendUrl = CHUCKSSIM_CONFIG.simulatorUrl;

const connection = new WebSocket('ws://' + backendUrl, []);

let isConnected = false;

const stateQueue = [];

let tickerHandler;

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



