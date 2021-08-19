
// GO Lang
// go build test-webview.go (on openbsd may need to: CGO_LDFLAGS_ALLOW='-Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib' go build wifi-manager.go)
// (c) 2020 unix-world.org
// version: 20200517

package main

/*
#cgo openbsd LDFLAGS: -Wl,-z,wxneeded|-Wl,-rpath-link,/usr/X11R6/lib
*/


import (
	"os"
	"log"

	"strings"
	"encoding/hex"
	"crypto/sha1"
	"golang.org/x/crypto/pbkdf2"

	smart "github.com/unix-world/smartgo"
	webview "github.com/unix-world/smartgo/webview2"
)


const (

	THE_VERSION = "r.20200517.2223"

	WIFI_INTERFACE = "iwx0"

	TPL_DOC = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>[###TITLE|html###]</title>
<script>
function getFormInputValue(id) {
	var val = '';
	try {
		var el = document.getElementById(id);
		val = el.value;
	} catch(err) {
		val = '';
	}
	return String(val);
}
function getFormSelectValue(id) {
	var val = '';
	try {
		var el = document.getElementById(id);
		val = el.options[el.selectedIndex].value;
	} catch(err) {
		val = '';
	}
	return String(val);
}
</script>
<script>
function call__quit(){
	goQuit(); // go RPC
}
function call__home(){
	goHome(); // go RPC
}
function call__scan(){
	goScan(); // go RPC
}
function call__off(){
	goOff(); // go RPC
}
function call__off(){
	goOn(); // go RPC
}
function call__settform() {
	goSettForm(); // go RPC
}
function call__settapply() {
	var ssid = getFormInputValue('ssid');
	var mode = getFormSelectValue('mode');
	var proto = getFormSelectValue('proto');
	var cipher = getFormSelectValue('cipher');
	var auth = getFormSelectValue('auth');
	var pass = getFormInputValue('pass');
	goSettApply(ssid, mode, proto, cipher, auth, pass); // go RPC
}
</script>
<style>
* {
	font-smooth: always;
	font-family: 'IBM Plex Sans',arial,sans-serif;
}

html, body{
	width: 96% auto;
	margin: 0px;
	padding: 10px;
}

body {
	background: #FFFFFF;
	color: #111111;
}

body {
	font-size: 16px;
	line-height: 1.5;
}

div, p, ol, ul, form, fieldset, blockquote, address, dl {
	font-size: 0.9375rem;
	line-height: 1.5;
}

li {
	list-style-position: inside;
}

table, pre, code, code * {
	font-size: 0.875rem;
	line-height: 1.25;
}

h1 {
	font-size: 2rem;
	margin-top: 0.625rem;
	margin-bottom: 0.625rem;
}
h2 {
	font-size: 1.75rem;
	margin-top: 0.5rem;
	margin-bottom: 0.5rem;
}
h3 {
	font-size: 1.5rem;
	margin-top: 0.375rem;
	margin-bottom: 0.375rem;
}
h4 {
	font-size: 1.25rem;
	margin-top: 0.3125rem;
	margin-bottom: 0.3125rem;
}
h5 {
	font-size: 1.125rem;
	margin-top: 0.25rem;
	margin-bottom: 0.25rem;
}
h6 {
	font-size: 1rem;
	margin-top: 0.1875rem;
	margin-bottom: 0.1875rem;
}

img {
	border: 0;
}

hr {
	border: none 0;
	border-top: 1px solid #CCCCCC;
	height: 1px;
}

table {
	font-weight: normal;
	text-align:left;
}
th {
	font-weight: bold;
}
td {
	font-weight: normal;
}

a, a:link, a:visited, a:hover {
	color: #000000;
}

input[type="text"]:disabled {
	background-color: #FFFFFF;
}

div#growl {
	width: 250px;
	position: fixed;
	top: 5px;
	right: 5px;
	padding: 5px;
	padding-left: 10px;
	padding-right: 10px;
	background: #003366;
	border-radius: 5px;
	opacity: 0.95;
	line-height: 36px;
	font-size: 28px;
}
div#growl * {
	color: #FFFFFF;
}
div#growl span.gtxt {
	display: inline-block;
	margin-top:  7px;
	margin-left: 5px;
	float: right;
}
div#growl img {
	display: inline-block;
	margin-top: 10px;
	width: 32px;
	height: 32px;
}
div.text-ok, div.text-err {
	background: #ECECEC;
	border: 1px solid #CCCCCC;
}
div.text-ok * {
	font-size: 1rem;
	font-weight: bold;
	color: #006699;
}
div.text-err * {
	font-size: 1rem;
	font-weight: bold;
	color: #FF5500;
}
div.text-ok pre, div.text-err pre {
	overflow-x: auto;
	padding-left:  10px;
	padding-right: 10px;
}

form.settings-form label {
	font-weight: bold;
}
form.settings-form input,
form.settings-form select {
	width: 500px;
}

.ux-button {
	display: inline-block;
	padding: 7px 15px 7px 15px !important;
	line-height: normal !important;
	white-space: nowrap !important;
	vertical-align: middle !important;
	text-align: center !important;
	font-family: inherit !important;
	font-size: 0.9375rem !important;
	font-weight: bold !important;
	text-decoration: none !important;
	border-radius: 3px !important;
	border: 1px solid #DDDDDD !important;
	background-color: #FCFCFC !important;
	color: #317AB9 !important;
	cursor: pointer !important;
	overflow: hidden;
	text-overflow: ellipsis;
	margin-bottom: 5px;
	margin-right: 5px;
	box-sizing: border-box !important;
	-moz-box-sizing: border-box !important;
	-webkit-box-sizing: border-box !important;
	appearance: none !important;
	-moz-appearance: none !important;
	-webkit-appearance: none !important;
	user-select: none !important;
	zoom: 1 !important;
}
.ux-button:hover,
.ux-button:focus {
	border: 1px solid #CCCCCC !important;
	background-color: #F3F3F3 !important;
	color: #2069A9 !important;
}
.ux-button:focus {
	outline: 0 !important;
}
.ux-button[disabled],
.ux-button-disabled,
.ux-button-disabled:hover,
.ux-button-disabled:focus,
.ux-button-disabled:active,
.ux-button-noaction {
	background-image: none !important;
	opacity: 0.75 !important;
	cursor: not-allowed !important;
	box-shadow: none !important;
}
.ux-button[disabled].ux-button-noaction,
.ux-button-noaction {
	cursor: default !important;
}
.ux-button-dark {
	border: 1px solid #333333 !important;
	background-color: #444444 !important;
	color: #FEFEFE !important;
}
.ux-button-dark:hover,
.ux-button-dark:focus {
	border: 1px solid #222222 !important;
	background-color: #333333 !important;
	color: #FFFFFF !important;
}
.ux-button-primary {
	border: 1px solid #20639E !important;
	background-color: #337AB7 !important;
	color: #FFFFFF !important;
}
.ux-button-primary:hover,
.ux-button-primary:focus {
	border: 1px solid #10528D !important;
	background-color: #2269A6 !important;
	color: #FFFFFF !important;
}
.ux-button-secondary {
	border: 1px solid #667788 !important;
	background-color: #778899 !important;
	color: #FFFFFF !important;
}
.ux-button-secondary:hover,
.ux-button-secondary:focus {
	border: 1px solid #5B646C !important;
	background-color: #6C757D !important;
	color: #FFFFFF !important;
}
.ux-button-regular {
	border: 1px solid #507100 !important;
	background-color: #719300 !important;
	color: #FFFFFF !important;
}
.ux-button-regular:hover,
.ux-button-regular:focus {
	border: 1px solid #406000 !important;
	background-color: #608200 !important;
	color: #FFFFFF !important;
}
.ux-button-highlight {
	border: 1px solid #DB860E !important;
	background-color: #FF9900 !important;
	color: #FFFFFF !important;
}
.ux-button-highlight:hover,
.ux-button-highlight:focus {
	border: 1px solid #CA750D !important;
	background-color: #EE8800 !important;
	color: #FFFFFF !important;
}
.ux-button-special {
	border: 1px solid #D21D1D !important;
	background-color: #F43F3F !important;
	color: #FFFFFF !important;
}
.ux-button-special:hover,
.ux-button-special:focus {
	border: 1px solid #C10C0C !important;
	background-color: #E32E2E !important;
	color: #FFFFFF !important;
}
.ux-button-info {
	border: 1px solid #178DB0 !important;
	background-color: #31B0D5 !important;
	color: #FFFFFF !important;
}
.ux-button-info:hover,
.ux-button-info:focus {
	border: 1px solid #067CA0 !important;
	background-color: #20A0C4 !important;
	color: #FFFFFF !important;
}

.ux-inline {
	display: inline;
}

.ux-inline-block {
	display: inline-block;
}

.ux-form input[type="text"],
.ux-form input[type="password"],
.ux-form input[type="email"],
.ux-form input[type="url"],
.ux-form input[type="date"],
.ux-form input[type="month"],
.ux-form input[type="time"],
.ux-form input[type="datetime"],
.ux-form input[type="datetime-local"],
.ux-form input[type="week"],
.ux-form input[type="number"],
.ux-form input[type="search"],
.ux-form input[type="tel"],
.ux-form input[type="color"],
.ux-form textarea,
.ux-form select,
.ux-field {
	background-color: #FFFFFF;
	color: #222222;
	font-weight: normal;
	font-size: 1rem;
	padding: 3px;
	display: inline-block;
	border: 1px solid #ccc;
	margin-bottom: 5px;
	margin-right: 5px;
	border-radius: 3px;
	vertical-align: middle;
	box-sizing: border-box;
	-moz-box-sizing: border-box;
	-webkit-box-sizing: border-box;
}

.ux-form select option,
.ux-select option {
	font-weight: normal !important;
	font-size: 1rem;
	padding: 1px !important;
}

.ux-form input:not([type]) {
	padding: 0.5rem 0.6rem;
	display: inline-block;
	border: 1px solid #ccc;
	box-shadow: inset 0 1px 3px #ddd;
	border-radius: 3px;
	box-sizing: border-box;
	-moz-box-sizing: border-box;
	-webkit-box-sizing: border-box;
}

.ux-form input[type="color"] {
	padding: 0.2rem 0.5rem;
}

.ux-form input[type="text"]:focus,
.ux-form input[type="password"]:focus,
.ux-form input[type="email"]:focus,
.ux-form input[type="url"]:focus,
.ux-form input[type="date"]:focus,
.ux-form input[type="month"]:focus,
.ux-form input[type="time"]:focus,
.ux-form input[type="datetime"]:focus,
.ux-form input[type="datetime-local"]:focus,
.ux-form input[type="week"]:focus,
.ux-form input[type="number"]:focus,
.ux-form input[type="search"]:focus,
.ux-form input[type="tel"]:focus,
.ux-form input[type="color"]:focus,
.ux-form textarea:focus,
.ux-form select:focus,
.ux-select:focus,
.ux-field:focus {
	outline: 0;
	border-color: #129FEA;
}

.ux-form input:not([type]):focus {
	outline: 0;
	border-color: #129FEA;
}

.ux-form input[type="file"]:focus,
.ux-form input[type="radio"]:focus,
.ux-form input[type="checkbox"]:focus {
	outline: thin solid #129FEA;
	/* @alternate */ outline: 1px auto #129FEA;
}
.ux-form .ux-checkbox,
.ux-form .ux-radio {
	margin: 0.5rem 0;
	display: block;
}

.ux-form input[type="text"][disabled],
.ux-form input[type="password"][disabled],
.ux-form input[type="email"][disabled],
.ux-form input[type="url"][disabled],
.ux-form input[type="date"][disabled],
.ux-form input[type="month"][disabled],
.ux-form input[type="time"][disabled],
.ux-form input[type="datetime"][disabled],
.ux-form input[type="datetime-local"][disabled],
.ux-form input[type="week"][disabled],
.ux-form input[type="number"][disabled],
.ux-form input[type="search"][disabled],
.ux-form input[type="tel"][disabled],
.ux-form input[type="color"][disabled],
.ux-form textarea[disabled],
.ux-form select[disabled],
.ux-select[disabled],
.ux-field[disabled] {
	cursor: not-allowed;
	background-color: #ECECEC;
	color: #999999;
}

.ux-form input:not([type])[disabled] {
	cursor: not-allowed;
	background-color: #ECECEC;
	color: #999999;
}
.ux-form input[readonly],
.ux-form select[readonly],
.ux-form textarea[readonly],
.ux-select[readonly],
.ux-field[readonly] {
	cursor: pointer;
	background-color: #FAFAFA;
	color: #777;
	border-color: #ccc;
}

.ux-form input:focus:invalid,
.ux-form textarea:focus:invalid,
.ux-form select:focus:invalid {
	color: #b94a48;
	border-color: #e9322d;
}
.ux-form input[type="file"]:focus:invalid:focus,
.ux-form input[type="radio"]:focus:invalid:focus,
.ux-form input[type="checkbox"]:focus:invalid:focus {
	outline-color: #e9322d;
}
.ux-form select {
	height: 2rem;
	border: 1px solid #ccc;
	background-color: white;
}
.ux-form select[multiple] {
	height: auto;
}
.ux-form label {
	margin: 0.5rem 0 0.2rem;
}
.ux-form fieldset {
	margin: 0;
	padding: 0.35rem 0 0.75rem;
	border: 0;
}
.ux-form legend {
	display: block;
	width: 100%;
	padding: 0.3rem 0;
	margin-bottom: 0.3rem;
	color: #333;
	border-bottom: 1px solid #e5e5e5;
}

.ux-form-stacked input[type="text"],
.ux-form-stacked input[type="password"],
.ux-form-stacked input[type="email"],
.ux-form-stacked input[type="url"],
.ux-form-stacked input[type="date"],
.ux-form-stacked input[type="month"],
.ux-form-stacked input[type="time"],
.ux-form-stacked input[type="datetime"],
.ux-form-stacked input[type="datetime-local"],
.ux-form-stacked input[type="week"],
.ux-form-stacked input[type="number"],
.ux-form-stacked input[type="search"],
.ux-form-stacked input[type="tel"],
.ux-form-stacked input[type="color"],
.ux-form-stacked input[type="file"],
.ux-form-stacked select,
.ux-form-stacked label,
.ux-form-stacked textarea {
	display: block;
	margin: 0.25rem 0;
}

.ux-form-stacked input:not([type]) {
	display: block;
	margin: 0.25rem 0;
}
.ux-form-aligned input,
.ux-form-aligned textarea,
.ux-form-aligned select,
.ux-form-aligned .ux-help-inline,
.ux-form-message-inline {
	display: inline-block;
	*display: inline;
	*zoom: 1;
	vertical-align: middle;
}
.ux-form-aligned textarea,
.ux-form-aligned label {
	vertical-align: top !important;
}

.ux-form-aligned .ux-control-group {
	margin-bottom: 0.5rem;
}
.ux-form-aligned .ux-control-group label {
	text-align: right;
	display: inline-block;
	vertical-align: middle;
	width: 10rem;
	margin: 0 1rem 0 0;
}
.ux-form-aligned .ux-controls {
	margin: 1.5rem 0 0 11rem;
}

.ux-form input.ux-input-rounded,
.ux-form .ux-input-rounded {
	border-radius: 9px;
	padding: 0.5rem 1rem;
}

.ux-form .ux-group fieldset {
	margin-bottom: 10px;
}
.ux-form .ux-group input,
.ux-form .ux-group textarea {
	display: block;
	padding: 7px;
	margin: 0 0 -1px;
	border-radius: 0;
	position: relative;
	top: -1px;
}
.ux-form .ux-group input:focus,
.ux-form .ux-group textarea:focus {
	z-index: 3;
}
.ux-form .ux-group input:first-child,
.ux-form .ux-group textarea:first-child {
	top: 1px;
	border-radius: 3px 3px 0 0;
	margin: 0;
}
.ux-form .ux-group input:first-child:last-child,
.ux-form .ux-group textarea:first-child:last-child {
	top: 1px;
	border-radius: 3px;
	margin: 0;
}
.ux-form .ux-group input:last-child,
.ux-form .ux-group textarea:last-child {
	top: -2px;
	border-radius: 0 0 3px 3px;
	margin: 0;
}
.ux-form .ux-group button {
	margin: 0.35rem 0;
}

.ux-form .ux-input-1-1 {
	width: 600px;
}
.ux-form .ux-input-2-3 {
	width: 400px;
}
.ux-form .ux-input-1-2 {
	width: 300px;
}
.ux-form .ux-input-1-3 {
	width: 200px;
}
.ux-form .ux-input-1-4 {
	width: 150px;
}

.ux-form .ux-input-def {
	padding: 1px 3px 1px 3px !important;
	height: 2rem !important;
	font-family: inherit !important;
	font-size: 1rem !important;
	margin-bottom: 5px;
}

.ux-form .ux-help-inline,
.ux-form-message-inline {
	display: inline-block;
	padding-left: 0.3rem;
	color: #428BCA;
	vertical-align: middle;
	font-size: 0.875rem;
}

.ux-form-message {
	display: block;
	color: #666;
	font-size: 0.875rem;
}

.ux-form ::-webkit-input-placeholder {
	color: #BBBBBB;
}
.ux-form ::-moz-placeholder {
	color: #999999;
}
</style>
</head>
<body>
[###MAIN###]
<br>
<div style="text-align:right; color:#778899;">
<small>OpenBSD Wifi Manafer ` + THE_VERSION + ` &copy;&nbsp;2020&nbsp;unix-world.org</small>
</div>
</body>
</html>
`

	TPL_HOME_OK = `
<button class="ux-button ux-button-dark" onClick="call__quit();">Quit</button>
<!-- <button class="ux-button ux-button-special" onClick="call__off();">Turn OFF Wifi</button> -->
<button class="ux-button ux-button-primary" onClick="call__scan();">Scan Wifi Networks</button>
<hr>
<h1>[###AREA-TTL|html###]</h1>
<div id="ok" class="text-ok"><pre>
[###OUTPUT|html###]
</pre></div>
<script>
setTimeout(function(){ call__home(); }, 5500);
</script>
<div id="growl">
<img src="data:image/svg+xml,[###SVG-IMG|url|html###]"><span class="gtxt">Refreshing ...</span>
</div>
`

	TPL_HOME_ERR = `
<!-- <button class="ux-button ux-button-regular" onClick="call__on();">Turn ON Wifi</button> -->
<hr>
<h1>[###AREA-TTL|html###]</h1>
<div id="ok" class="text-ok"><pre>
[###OUTPUT|html###]
</pre></div>
<br>
<div id="err" class="text-err"><pre>
[###ERRORS|html###]
</pre></div>
`

	TPL_SCAN_OK = `
<button class="ux-button" onClick="call__home();">Go Back</button>
<button class="ux-button ux-button-info" onClick="call__settform();">Connect to a Wifi Network</button>
<hr>
<h1>[###AREA-TTL|html###]</h1>
<div id="ok" class="text-ok"><pre>
[###OUTPUT|html###]
</pre></div>
<script>
setTimeout(function(){ call__scan(); }, 5500);
</script>
<div id="growl">
<img src="data:image/svg+xml,[###SVG-IMG|url|html###]"><span class="gtxt">Scanning ...</span>
</div>
`

	TPL_SCAN_ERR = `
<button class="ux-button" onClick="call__home();">Go Back</button>
<hr>
<h1>[###AREA-TTL|html###]</h1>
<div id="ok" class="text-ok"><pre>
[###OUTPUT|html###]
</pre></div>
<br>
<div id="err" class="text-err"><pre>
[###ERRORS|html###]
</pre></div>
<h3>Scanning ERROR !</h3>
`

	TPL_SETT_FORM = `
<button class="ux-button" onClick="call__home();">Go Back</button>
<hr>
<h1>[###AREA-TTL|html###]</h1>
<form id="settings" class="ux-form ux-form-aligned settings-form">
	<fieldset>
		<div class="ux-control-group">
			<label for="ssid">Network Name</label>
			<input id="ssid" type="text" placeholder="Network Name (SSID)" maxlength="32">
		</div>
		<div class="ux-control-group">
			<label for="mode">Mode</label>
			<select id="mode">
				<option value="11ac">802.11ac (1300 Mbps, 5Ghz only)</option>
				<option value="11ax">802.11ax (1300 Mbps, 2.4Ghz and 5Ghz)</option>
				<option value="11n" selected>802.11n (450 Mbps, 2.4Ghz and 5Ghz)</option>
				<option value="11g">802.11g (54 Mbps, 2.4Ghz only)</option>
				<option value="11a">802.11a (54 Mbps, 5Ghz only)</option>
				<option value="11b">802.11b (11 Mbps, 2.4Ghz only)</option>
			</select>
		</div>
		<div class="ux-control-group">
			<label for="proto">Protocol</label>
			<select id="proto">
				<option value="wpa2" selected>WPA2</option>
				<option value="wpa1">WPA1 (weak)</option>
			</select>
		</div>
		<div class="ux-control-group">
			<label for="cipher">Cipher</label>
			<select id="cipher">
				<option value="ccmp" selected>CCMP (strong, AES)</option>
				<option value="tkip">TKIP (weak)</option>
			</select>
		</div>
		<div class="ux-control-group">
			<label for="auth">Authentication</label>
			<select id="auth">
				<option value="psk" selected>WPA / PSK</option>
			</select>
		</div>
		<div class="ux-control-group">
			<label for="pass">Passphrase</label>
			<input id="pass" type="password" placeholder="Passphrase for WPA/PSK mode" maxlength="63">
		</div>
	</fieldset>
	<button class="ux-button ux-button-highlight" onClick="call__settapply(); return false;">Apply Wifi Connection</button>
</form>
`

	TPL_SETT_APPLY_OK = `
<button class="ux-button" onClick="call__home();">Go Back</button>
<hr>
<h1>[###AREA-TTL|html###]</h1>
<div id="ok" class="text-ok"><pre>
[###OUTPUT|html###]
</pre></div>
`

	TPL_SETT_APPLY_ERR = `
<button class="ux-button" onClick="call__home();">Go Back</button>
<hr>
<h1>[###AREA-TTL|html###]</h1>
<div id="ok" class="text-ok"><pre>
[###OUTPUT|html###]
</pre></div>
<br>
<div id="err" class="text-err"><pre>
[###ERRORS|html###]
</pre></div>
`

	SVG_SPIN = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="32" height="32" fill="grey" id="loading-spin-svg"><path opacity=".25" d="M16 0 A16 16 0 0 0 16 32 A16 16 0 0 0 16 0 M16 4 A12 12 0 0 1 16 28 A12 12 0 0 1 16 4"/><path d="M16 0 A16 16 0 0 1 32 16 L28 16 A12 12 0 0 0 16 4z"><animateTransform attributeName="transform" type="rotate" from="0 16 16" to="360 16 16" dur="0.8s" repeatCount="indefinite" /></path></svg>`

)


func LogToConsoleWithColors() {
	//--
	smart.ClearPrintTerminal()
	//--
	smart.LogToConsole("DEBUG", true)
	//--
} //END FUNCTION


func fatalError(logMessages ...interface{}) {
	//--
	log.Fatal("[ERROR] ", logMessages) // standard logger
	//--
	os.Exit(1)
	//--
} //END FUNCTION


func getErrFromCmdStdErr(stdErr string) string {
	//--
	return smart.StrTrimWhitespaces(smart.StrReplaceAll(stdErr, "[SmartGo:cmdExec:Exit:ERROR]", ""))
	//--
} //END FUNCTION


func renderTplDoc(areaTitleTxt string, areaMainHtml string) string {
	//--
	var arr = map[string]string{
		"TITLE": 	areaTitleTxt,
		"MAIN": 	areaMainHtml,
	}
	//--
	return smart.MarkersTplRender(TPL_DOC, arr, false, false) + "\n"
	//--
} //END FUNCTION


func renderHome() string {
	//--
	isSuccess, outStd, errStd := smart.ExecTimedCmd(30, "capture", "capture", "", "", "doas", "ifconfig", WIFI_INTERFACE)
	//--
	var arr = map[string]string{
		"AREA-TTL": 	"Wifi Status " + smart.DateNowIsoLocal(),
		"OUTPUT": 		outStd,
		"ERRORS": 		getErrFromCmdStdErr(errStd),
		"SVG-IMG": 		SVG_SPIN,
	}
	//--
	var theHtml string = ""
	//--
	if((isSuccess == true) && (errStd == "")) {
		theHtml = smart.MarkersTplRender(TPL_HOME_OK, arr, false, false)
	} else {
		theHtml = smart.MarkersTplRender(TPL_HOME_ERR, arr, false, false)
	} //end if else
	//--
	return renderTplDoc("Wifi Home", theHtml)
	//--
} //END FUNCTION


func renderScan() string {
	//--
	isSuccess, outStd, errStd := smart.ExecTimedCmd(30, "capture", "capture", "", "", "doas", "ifconfig", WIFI_INTERFACE, "scan")
	//--
	var arr = map[string]string{
		"AREA-TTL": 	"Available Wifi Networks " + smart.DateNowIsoLocal(),
		"OUTPUT": 		outStd,
		"ERRORS": 		getErrFromCmdStdErr(errStd),
		"SVG-IMG": 		SVG_SPIN,
	}
	//--
	var theHtml string = ""
	//--
	if((isSuccess == true) && (errStd == "")) {
		theHtml = smart.MarkersTplRender(TPL_SCAN_OK, arr, false, false)
	} else {
		theHtml = smart.MarkersTplRender(TPL_SCAN_ERR, arr, false, false)
	} //end if else
	//--
	return renderTplDoc("Wifi Scan", theHtml)
	//--
} //END FUNCTION


func renderSettingsForm() string {
	//--
	var arr = map[string]string{
		"AREA-TTL": 	"Connect to a Wifi Network (DHCP) - Form",
	}
	//--
	var theHtml string = smart.MarkersTplRender(TPL_SETT_FORM, arr, false, false)
	//--
	return renderTplDoc("Wifi Form", theHtml)
	//--
} //END FUNCTION


func renderSettingsApply(ssid string, mode string, proto string, cipher string, auth string, pass string) string {
	//--
	rawKey := pbkdf2.Key([]byte(pass), []byte(ssid), 4096, 32, sha1.New)
	var hexKey string = "0x" + strings.ToLower(hex.EncodeToString(rawKey))
	//--
	log.Println("[DEBUG] SettingsApply ; SSID=`" + ssid + "` ; HexKey=`" + hexKey + "` ; Mode=`" + mode + "` ; Proto=`" + proto + "` ; Cipher=`" + cipher + "` ; Auth=`" + auth + "`")
	//--
	var arr = map[string]string{
		"AREA-TTL": 	"Connect to a Wifi Network (DHCP) - Apply",
		"OUTPUT": 		"",
		"ERRORS": 		"",
	}
	//--
	var theHtml string = ""
	//--
	if((ssid == "") || (smart.StrTrimWhitespaces(ssid) == "")) {
		arr["ERRORS"] = "ERROR: Empty SSID ..."
		theHtml = smart.MarkersTplRender(TPL_SETT_APPLY_ERR, arr, false, false)
		return renderTplDoc("Wifi Connect", theHtml)
	} //end if
	//--
	isSuccess, outStd, errStd := smart.ExecTimedCmd(30, "capture", "capture", "", "", "doas", "ifconfig", WIFI_INTERFACE, "mode", mode, "wpa", "wpaakms", auth, "wpaprotos", proto, "wpagroupcipher", cipher, "nwid", ssid, "wpakey", hexKey, "autoconf")
	//--
	arr["OUTPUT"] = outStd;
	arr["ERRORS"] = getErrFromCmdStdErr(errStd)
	//--
	if((isSuccess == true) && (errStd == "")) {
		theHtml = smart.MarkersTplRender(TPL_SETT_APPLY_OK, arr, false, false)
	} else {
		theHtml = smart.MarkersTplRender(TPL_SETT_APPLY_ERR, arr, false, false)
	} //end if else
	//--
	return renderTplDoc("Wifi Connect", theHtml)
	//--
} //END FUNCTION


//---


func makeDataUrl(htmlDoc string) string {
	//--
//	return `data:text/html,` + smart.RawUrlEncode(htmlDoc)
	return `data:text/html;base64,` + smart.Base64Encode(htmlDoc)
	//--
} //END FUNCTION


//---


func init() {

	LogToConsoleWithColors()

} //END FUNCTION


func main() {

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("OpenBSD Wifi Manager")
	w.SetSize(980, 720, webview.HintFixed)

	w.Bind("goQuit", func() {
		w.Terminate()
	})
	w.Bind("goHome", func() {
		log.Println("[NOTICE] navigate Home")
		w.Navigate(makeDataUrl(renderHome()))
	})
	w.Bind("goScan", func() {
		log.Println("[NOTICE] navigate Scan")
		w.Navigate(makeDataUrl(renderScan()))
	})
	w.Bind("goSettForm", func() {
		log.Println("[NOTICE] navigate SettingsForm")
		w.Navigate(makeDataUrl(renderSettingsForm()))
	})
	w.Bind("goSettApply", func(ssid string, mode string, proto string, cipher string, auth string, pass string) {
		log.Println("[NOTICE] navigate SettingsApply")
		w.Navigate(makeDataUrl(renderSettingsApply(ssid, mode, proto, cipher, auth, pass)))
	})

	log.Println("[NOTICE] navigate Home")
	w.Navigate(makeDataUrl(renderHome()))

	w.Run()

} //END FUNCTION


// #END
