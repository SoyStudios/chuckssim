// This file is required by the index.html file and will
// be executed in the renderer process for that window.
// All of the Node.js APIs are available in this process.


const CHUCKSSIM_CONFIG = {
    simulatorUrl: 'localhost:5000',
    startServerCommand: 'do something'
};
const backendUrl = CHUCKSSIM_CONFIG.simulatorUrl;


let connection;

let isConnected = false;

const stateQueue = [];

let tickerHandler;

let currentBot = {
    id: '',
    x: 0,
    y: 0,
    ang: '',
    desc: '',
    code: ''
};


// let electronify = require('electronify-server');
// const startServerCommand = CHUCKSSIM_CONFIG.startServerCommand;

// electronify({
//     command: startServerCommand,
//     // url: 'http://127.0.0.1:8000',
//     debug: true,
//     window: {height: 768, width: 1024},
//     ready: function(app){
//         // application event listeners could be added here
//     },
//     preLoad: function(app, window){
//         // window event listeners could be added here
//     },
//     postLoad: function(app, window){
//         // url finished loading
//     },
//     showDevTools: false
// }).on('child-started', function(child) {
//     // child process has started
//     console.log('PID: ' + child.pid);
//
//     // setup logging on child process
//     child.stdout.on('data', console.log);
//     child.stderr.on('data', console.log);
//     connection = new WebSocket('ws://' + backendUrl, []);
//
// }).on('child-closed', function(app, stderr, stdout) {
//     // the child process has finished
//
// }).on('child-error', function(err, app) {
//     // close electron if the child crashes
//     // app.quit();
// });

connection = new WebSocket('ws://' + backendUrl, []);
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
    if (e.data.type === 'bot_details') {
        currentBot = e.data.bot;
    }
};

const myGameArea = {
    canvas : document.getElementById("canvas"),
    start : function() {
        this.canvas.width = window.innerWidth;
        this.canvas.height = window.innerHeight;
        this.context = this.canvas.getContext("2d");
        // document.body.insertBefore(this.canvas, document.body.childNodes[0]);
        // this.interval = setInterval(updateGameArea, 20);
    },
    clear : function() {
        this.context.clearRect(0, 0, this.canvas.width, this.canvas.height);
    }
};

renderBot = (bot) => {
    console.log('here');
    this.width = 15;
    this.height = 15;
    this.x = bot.x;
    this.y = bot.y;
    // this.update = function(){
        ctx = myGameArea.context;
        ctx.fillStyle = 'red';
        ctx.fillRect(this.x, this.y, this.width, this.height);
    // };
};


renderState = (state) => {
    myGameArea.clear();
    for (const bot of state.data) {
        // console.log(bot);
        renderBot(bot);
    }
  // do stuff to render on canvas
};

togglePause = () => {
    if (isConnected) {
        connection.send('togglePause')
    }
};


renderTicker = (interval) => {
  console.log('here');
  return setInterval(
      () => {
          if (stateQueue.length > 0){
              renderState(stateQueue.shift());
          }
      }, interval
  )
};

stopRenderTicker = () => {
  if (tickerHandler) {
      window.clearInterval(tickerHandler);
  }
};

getBotDetails = (botId) => {
  if (isConnected) {
      connection.send('bot/' + botId);
  }
};

// only for development
let mockInitialState = {
    type: 'state',
    data: [
        {
            id: 1,
            x: 1.0,
            y: 0.6,
            ang: 240,
            dna: 2006
        },
        {
            id: 3,
            x: 4.0,
            y: 1.6,
            ang: 120,
            dna: 2006
        },

        {
            id: 3,
            x: 40.0,
            y: 12.6,
            ang: 120,
            dna: 2012
        },
    ]
};


function precisionRound(number, precision) {
    var factor = Math.pow(10, precision);
    return Math.round(number * factor) / factor;
}

changeStateMock = (prevState) => {
    const newState = {
        type: 'state',
        data: [
        ]
    };
    modifyCoordinates = (coord) => {
        return coord + precisionRound(-2 * Math.random() + 2 * Math.random(), 2);
    };
    modifyAngle = (angle) => {
        return (angle + Math.floor(-720 * Math.random() + 720 * Math.random())) % 360;
    };
    for (let i = 0; i < prevState.data.length; i++) {
        const bot = prevState.data[i];
        const newData = {
            id: bot.id,
            x: modifyCoordinates(bot.x),
            y: modifyCoordinates(bot.y),
            ang: modifyAngle(bot.ang),
            dna: bot.dna
        };
        newState.data.push(newData);
    }
    return newState;
};

startMockSim = () => {
    let newState = mockInitialState;
    stateQueue.push(newState);
    newState = changeStateMock(newState);
    stateQueue.push(newState);
    newState = changeStateMock(newState);
    stateQueue.push(newState);
    newState = changeStateMock(newState);
    stateQueue.push(newState);
    newState = changeStateMock(newState);
    stateQueue.push(newState);
    let i = 0;
    while (i < 500) {
        newState = changeStateMock(newState);
        stateQueue.push(newState);
        i++;
    }
    myGameArea.start();
    tickerHandler = renderTicker(20);
    setTimeout(
        () => {
            stopRenderTicker();
        }, 10000
    )
};

startMockSim();



