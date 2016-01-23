var login = require("facebook-chat-api"),
    prompt = require('prompt'),
    argv = require('yargs').argv;

var commands = [
	"\"olaf help\": get help",
	"\"hi/hey olaf\": say hello to olaf",
	"\"olaf insult <name>\": get olaf to insult <name>"
];

prompt.override = argv;
prompt.start();
/* Meet Olaf */
prompt.get(['email', {
    name: 'password',
    hidden: true,
    conform: function (value) {
        return true; 
    }
}], function (err, result) {
	login({email: result.email, password: result.password}, function callback (err, api) {
		if(err) return console.error(err);

		api.setOptions({listenEvents: true});

		var stopListening = api.listen(function(err, event) {
			if(err) return console.error(err);

			switch(event.type) {
				case "message":
					var cmd = event.body.trim().toLowerCase();
					var res = "";

					if(cmd.indexOf("olaf help") > -1) {
						res = "list of commands:\n\n";
						for (var i in commands) {
  						res += commands[i] + "\n";
						}
					} else if(cmd.indexOf("hey olaf") > -1 || cmd.indexOf("hi olaf") > -1) {
						res = "sup?";
					} else if(cmd.indexOf("olaf insult") > -1) {
						cmd = cmd.split(" ");
						if(cmd.length >= 3) {
							var name = "";
							for (var i = 2; i < cmd.length; i++) {
								name += " " + cmd[i];
							}
							res = "Fuck you," + name + "\n - " + event.senderName.split(" ")[0];
						} else {
							res = "wtf dude?"
						}
					}

					console.log("command:" + cmd);
					console.log("sending:" + res);
					api.sendMessage(res, event.threadID);

					api.markAsRead(event.threadID, function(err) {
						if(err) console.log(err);
					});
					break;
				case "event":
					console.log(event);
					break;
			}
		});
	});
});
