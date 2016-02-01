var logger = exports;
var levels = ['ERROR', 'WARN', 'INFO'];
logger.debugLevel = 'WARN';
logger.log = function(level, message) {
  if (levels.indexOf(level) >= levels.indexOf(logger.debugLevel)) {
    if (typeof message !== 'string') {
      message = JSON.stringify(message);
    };
    console.log(level + ': ' + message);
  }
}
