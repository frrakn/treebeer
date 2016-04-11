'use strict';

const WebSocketServer = require('websocket').server;
const http = require('http');

const AcceptPacket = require('./accept.json');
const MessagePacket = require('./message.json');
const MessageInterval = 1000;

function formatMessage(message) {
  return JSON.stringify(message);
}

const server = http.createServer(function(request, response) {
  console.log((new Date()) + ' Received request for ' + request.url);
  response.writeHead(404);
  response.end();
});

server.listen(8080, function() {
  console.log((new Date()) + ' Server is listening on port 8080');
});

const wsServer = new WebSocketServer({
  httpServer: server,
  autoAcceptConnections: false
});

wsServer.on('request', function(request) {
  const connection = request.accept('emitter-protocol', request.origin);
  console.log((new Date()) + ' Connection accepted.');

  connection.send(formatMessage(AcceptPacket), function() {
    let UpdatePacket = MessagePacket;

    setInterval(function () {
      // Simulate time passing
      UpdatePacket["1001600049"]["teamStats"]["100"]["gameLength"] += 1000;
      UpdatePacket["1001600049"]["teamStats"]["200"]["gameLength"] += 1000;

      connection.send(formatMessage(UpdatePacket))
    }, MessageInterval);
  });
});
