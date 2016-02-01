var login = require("facebook-chat-api"),
  prompt = require('prompt'),
  argv = require('yargs').argv,
  http = require('http'),
  cheerio = require('cheerio'),
  logger = require('./logger');

logger.debugLevel = 'INFO';

var commands = [
  "\"Olaf help\": get help",
  "\"Olaf hi/hey\": say hello to olaf",
  "\"Olaf insult <name>\": olaf to insult <name>",
  "\"Olaf insult <keyword> <name>\": olaf to special insult <name>"
];

var get_url = function(params) {
  var p, url, _i, _len;
  url = "";
  for (_i = 0, _len = params.length; _i < _len; _i++) {
    p = params[_i];
    url += "/" + (encodeURIComponent(p));
  }
  return url;
};

var fuckoff_request = function(params, callback) {
  var service_host = "foaas.com";
  var path;
  if (params.name == null) {
    path = get_url([params.keyword, params.from]);
  } else {
    path = get_url([params.keyword, params.name, params.from]);
  }

  logger.log("INFO", "Requesting insult from: " + service_host + path);
  var options = {
    host: service_host,
    port: 80,
    path: path,
    method: 'GET'
  };
  return http.request(options, function(res) {
    res.setEncoding('utf8');
    return res.on('data', function(data) {
      data = cheerio.load(data)("h1").text();
      logger.log("INFO", "Insult response: " + data);
      return callback(data);
    });
  }).on('error', function(e) {
    return console.error(e.message);
  }).end();
};

prompt.override = argv;
prompt.start();
prompt.get(['email', {
  name: 'password',
  hidden: true,
  conform: function(value) {
    return true;
  }
}], function(err, result) {
  login({
    email: result.email,
    password: result.password
  }, function callback(err, api) {
    if (err) return console.error(err);

    api.setOptions({
      listenEvents: true
    });

    var stopListening = api.listen(function(err, event) {
      if (err) return console.error(err);

      switch (event.type) {
        case "message":
          var cmdStr = event.body.trim().toLowerCase();
          logger.log("INFO", "Command: " + cmdStr);
          var cmd = cmdStr.split(" ");
          var res = "";

          if (cmdStr.indexOf("olaf") > -1) {
            if (cmd[0] === "olaf") {
              // initiate command recognition
              switch (cmd[1]) {
                case "help":
                  res = "list of commands:\n\n";
                  for (var i in commands) {
                    res += commands[i] + "\n";
                  }
                  api.sendMessage(res, event.threadID);
                  break;
                case "hi":
                  //fall through
                case "hey":
                  res = "sup " + event.senderName.split(" ")[0] + "?";
                  api.sendMessage(res, event.threadID);
                  break;
                case "insult":
                  if (cmd.length == 4) {
                    var params = {
                      keyword: cmd[2],
                      name: cmd[3],
                      from: event.senderName.split(" ")[0]
                    };
                    fuckoff_request(params, function(result) {
                      res = result;
                      api.sendMessage(res, event.threadID);
                    });
                  } else if (cmd.length >= 3) {
                    for (var i = 3; i < cmd.length; i++) {
                      cmd[2] += " " + cmd[i];
                    }
                    var params = {
                      keyword: cmd[2],
                      name: null,
                      from: event.senderName.split(" ")[0]
                    };
                    fuckoff_request(params, function(result) {
                      res = result;
                      api.sendMessage(res, event.threadID);
                    });
                  } else {
                    res = "wtf dude?";
                    api.sendMessage(res, event.threadID);
                  }
                  break;
                default:
                  res = "I'm sorry, I can't hear you over the sound of how awesome I am";
                  api.sendMessage(res, event.threadID);
              }
            } else {
              //unrecognized command
              res = "wtf do you want " + event.senderName.split(" ")[0] + "?"
              api.sendMessage(res, event.threadID);
            }
          }

          api.markAsRead(event.threadID, function(err) {
            if (err) console.log(err);
          });
          break;
        case "event":
          logger.log("INFO", event);
          break;
      }
    });
  });
});
