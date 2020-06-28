<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <script src="/pub/js/home.js"></script>
  <title>Web Chat</title>
</head>
<body>
<table>
  <tr><td valign="top" width="50%">
      <p>Send a message!
      <p>
      <form onsubmit="return false;">
        <label for="input">Message</label><input id="input" type="text" value="Hello world!"><button type="button" id="sendmsg">Send Message </button>
      </form>
    </td><td valign="top" width="50%">
      <div id="output"></div>
    </td></tr></table>
</body>
</html>