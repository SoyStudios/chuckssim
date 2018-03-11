// This file is required by the index.html file and will
// be executed in the renderer process for that window.
// All of the Node.js APIs are available in this process.


const CHUCKSSIM_CONFIG = {
    simulatorUrl: '127.0.0.1:8080',
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
    myGameArea.start();
    tickerHandler = renderTicker(50);
};

// Log errors
connection.onerror = function (error) {
    isConnected = false;
    console.log('WebSocket Error ' + error);
};

// Log messages from the server
connection.onmessage = function (e) {
    console.log('Server: ' + e.data);
    addStateToQueue(JSON.parse(e.data).bots);
    // if (e.data.type === 'state') {
    // }
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
    },
    clear : function() {
        this.context.clearRect(0, 0, this.canvas.width, this.canvas.height);
    }
};

renderBot = (bot) => {
    const radius = 10;
    // this.width = 15;
    // this.height = 15;
    this.x = bot.x;
    this.y = bot.y;
    let color = bot.dna ? bot.dna * 10 % 493 + 300 : 'rgba(0, 200, 0, 0.6)';
    ctx = myGameArea.context;
    ctx.beginPath();
    ctx.arc(this.x, this.y, radius, 0, 2*Math.PI);
    ctx.fillStyle = bot.dna ? '#' + color.toString(16) : color;
    ctx.fill();

    // Draw eye
    const rad = bot.a * Math.PI / 180;
    let h2 = Math.sin(rad);
    let w2 = Math.sin(rad + Math.PI/2) * -1;
    ctx.beginPath();
    ctx.arc(this.x + (h2 * radius / 2), this.y + (w2 * radius / 2), radius * 0.1, 0, 2*Math.PI);
    ctx.fillStyle = 'black';
    ctx.fill();
};


renderState = (state) => {
    myGameArea.clear();
    if (state && state.length) {
        for (const bot of state) {
            // console.log(bot);
            renderBot(bot);
        }
    }
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
            x: 100.0,
            y: 100.6,
            ang: 0,
            dna: 793
        },
        {
            id: 3,
            x: 200.0,
            y: 100.6,
            ang: 90,
            dna: 2006
        },
        {
            id: 3,
            x: 300.0,
            y: 100.6,
            ang: 180,
            dna: 2007
        },
        {
            id: 4,
            x: 400.0,
            y: 100.6,
            ang: 270,
            dna: 2009
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
    let i = 0;
    while (i < 500) {
        newState = changeStateMock(newState);
        stateQueue.push(newState);
        i++;
    }
    setTimeout(
        () => {
            stopRenderTicker();
        }, 100000
    )
};

// tickerHandler = renderTicker(20);
// startMockSim();



