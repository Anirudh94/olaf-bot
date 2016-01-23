var login = require("facebook-chat-api"),
    prompt = require('prompt'),
    argv = require('yargs').argv,
    logger = require('./logger');

logger.debugLevel = 'info';

var commands = [
    "\"Olaf help\": get help",
    "\"Olaf hi/hey\": say hello to olaf",
    "\"Olaf insult <name>\": get olaf to insult <name>"
];

// TODO: this is a temporary store of insults
// will be moved to JSON/DB or use online API  
var insults = [
    "Fuck you",
    "Suck ma snowy dick"
];

prompt.override = argv;
prompt.start();
/* Meet Olaf */
prompt.get(['email', {
    name: 'password',
    hidden: true,
    conform: function(value) {
        return true;
    }
}], function(err, result) {
    /* Meet Olaf */
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
                                    break;
                                case "hi":
                                    //fall through
                                case "hey":
                                    res = "sup " + event.senderName.split(" ")[0] + "?";
                                    break;
                                case "insult":
                                    if (cmd.length >= 3) {
                                        var name = "";
                                        for (var i = 2; i < cmd.length; i++) {
                                            name += " " + cmd[i];
                                        }
                                        var rndIndex = Math.floor((Math.random() * 1000) % insults.length);
                                        var insult = insults[rndIndex];
                                        res = insult + name;
                                    } else {
                                        res = "wtf dude?";
                                    }
                                    break;
                                default:
                                    res = "I'm sorry, I can't hear you over the sound of how awesome I am";
                            }
                        } else {
                            //unrecognized command
                            res = "wtf do you want " + event.senderName.split(" ")[0] + "?"
                        }

                        console.log("command:" + cmd);
                        console.log("sending:" + res);
                        api.sendMessage(res, event.threadID);
                    }

                    api.markAsRead(event.threadID, function(err) {
                        if (err) console.log(err);
                    });
                    break;
                case "event":
                    console.log(event);
                    break;
            }
        });
    });
});