<?php

// Controller Handler Sample
// (c) 2017-2018 Radu.I
// v.180112.1537

// DO NOT MODIFY, THIS IS JUST A SAMPLE !!!

ini_set('display_errors', '1');	// display runtime errors
error_reporting(E_ALL & ~E_NOTICE & ~E_STRICT & ~E_DEPRECATED); // error reporting

$year = date('Y');
$dtime = date('Y-m-d H:i:s');
$starter = <<<HTML
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>PyWebkitGTK :: Scrapper Service</title>
<script>
	var UxmPageStart = true;
</script>
</head>
<body>
	<h1 id="PyWebkitGTK-INIT-Page">PyWebkitGTK / Scrapper Service :: Sample</h1>
	<h2>{$dtime}</h2>
	<h3>UserAgent:</h3>
	<h4 id="U-A"></h4>
	<br>
	<div align="center">
		<br>
		<img width="256" height="256" src="data:image/svg+xml;base64,PHN2ZyB4bWxuczpkYz0iaHR0cDovL3B1cmwub3JnL2RjL2VsZW1lbnRzLzEuMS8iIHhtbG5zOmNjPSJodHRwOi8vY3JlYXRpdmVjb21tb25zLm9yZy9ucyMiIHhtbG5zOnN2Zz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgaWQ9IndrdCIgdmlld0JveD0iMCAwIDEyOCAxMjgiIHZlcnNpb249IjEuMSIgd2lkdGg9IjEyOCIgaGVpZ2h0PSIxMjgiPgo8ZGVmcyBpZD0iZGVmczQiPgogIDxsaW5lYXJHcmFkaWVudCBncmFkaWVudFVuaXRzPSJ1c2VyU3BhY2VPblVzZSIgeTE9IjAuODI2MzM4MjMiIGdyYWRpZW50VHJhbnNmb3JtPSJzY2FsZSgxLjEyODg3NzQsMC44ODU4MzU4MikiIHkyPSI0NDEuMDg4NSIgeDE9IjIyMC41NzQwMSIgeDI9IjIyMC41NzQwMSIgaWQ9ImEiPgogICAgPHN0b3AgaWQ9InN0b3A3IiBvZmZzZXQ9IjAiIHN0b3AtY29sb3I9IiMzNEFBREMiLz4KICAgIDxzdG9wIGlkPSJzdG9wOSIgb2Zmc2V0PSIxIiBzdG9wLWNvbG9yPSIjMDA3QUZGIi8+CiAgPC9saW5lYXJHcmFkaWVudD4KPC9kZWZzPgo8ZyB0cmFuc2Zvcm09Im1hdHJpeCgwLjIyNSwwLDAsMC4yMjUsNy45NzQ3MTg4LDIuMDcyNDYyNSkiIHN0eWxlPSJmaWxsOiNmZmZmZmY7ZmlsbC1ydWxlOmV2ZW5vZGQiIGlkPSJnMTEiPgogIDxwYXRoIHN0eWxlPSJmaWxsOiNmZjlkMDAiIGlkPSJwYXRoMTMiIGQ9Im0gNDcxLjg2LDMwNS45MyBjIDM0LjE5LDI2LjY2IDM0LjE5LDcwLjUyIDAsOTcuMzkgTCAzMTEuMDIsNTI5Ljc0IGMgLTM0LjE5LDI2LjY2IC04OS44NSwyNi42NiAtMTI0LjA0LDAgTCAyNi4xNCw0MDMuNTQgQyAtOC4wNDYsMzc2Ljg4IC04LjA0NiwzMzMuMDIgMjYuMTM5LDMwNi4xNCBsIDE2MC44NCwtMTI2LjQxIGMgMzQuMTksLTI2LjY2IDg5Ljg1LC0yNi42NiAxMjQuMDQsMCBsIDE2MC44NCwxMjYuMiB6Ii8+CiAgPHBhdGggc3R5bGU9ImZpbGw6I2ZmY2MwMCIgaWQ9InBhdGgxNSIgZD0iTSAxODcuMzcsNDUwLjgzIDI1LjgxLDMyMy44NiBDIDkuNTI0LDMxMS4xIDAuNSwyOTQuMDEgMC41LDI3NS44NCBjIDAsLTE4LjE3IDkuMDI0NCwtMzUuMjYgMjUuMzEyLC00OC4wMiBsIDE2MS41NiwtMTI3LjE5IGMgMTYuMjksLTEyLjc2IDM4LjMsLTE5Ljg5OCA2MS42MywtMTkuODk4IDIzLjMzLDAgNDUuMTIsNy4xMzggNjEuNjMsMTkuODk4IGwgMTYxLjU2LDEyNi45NyBjIDE2LjI5LDEyLjc3IDI1LjMxLDI5Ljg1IDI1LjMxLDQ4LjAyIDAsMTguMTcgLTkuMDIsMzUuMjYgLTI1LjMxLDQ4LjAyIGwgLTE2MS41NiwxMjYuOTggYyAtMTYuNTEsMTIuOTcgLTM4LjMsMjAuMTEgLTYxLjYzLDIwLjExIC0yMy4zMywwIC00NS4xMiwtNy4xNCAtNjEuNjMsLTE5LjkgeiIvPgogIDxwYXRoIHN0eWxlPSJmaWxsOnVybCgjYSkiIGlkPSJwYXRoMTciIGQ9Ik0gMTg3LjM3LDM3MC44MyAyNS44MSwyNDMuODYgQyA5LjUyNCwyMzEuMSAwLjUsMjE0LjAxIDAuNSwxOTUuODQgMC41LDE3Ny42NyA5LjUyNDQsMTYwLjU4IDI1LjgxMiwxNDcuODIgTCAxODcuMzcyLDIwLjYzIGMgMTYuMjksLTEyLjc2IDM4LjMsLTE5Ljg5OCA2MS42MywtMTkuODk4IDIzLjMzLC02ZS01IDQ1LjEyLDcuMTM4IDYxLjYzLDE5LjkgbCAxNjEuNTYsMTI2Ljk3IGMgMTYuMjksMTIuNzcgMjUuMzEsMjkuODUgMjUuMzEsNDguMDIgMCwxOC4xNyAtOS4wMiwzNS4yNiAtMjUuMzEsNDguMDIgbCAtMTYxLjU2LDEyNi45OCBjIC0xNi41MSwxMi45NyAtMzguMywyMC4xMSAtNjEuNjMsMjAuMTEgLTIzLjMzLDAgLTQ1LjEyLC03LjE0IC02MS42MywtMTkuOSB6Ii8+CiAgPHBhdGggaWQ9InBhdGgxOSIgZD0ibSAyNDkuNTYsMzE3LjUyIGMgODIuNTUsMCAxNDkuNDgsLTU0LjcgMTQ5LjQ4LC0xMjIuMTggMCwtNjcuNDkgLTY2LjkzLC0xMjIuMTkgLTE0OS40OCwtMTIyLjE5IC04Mi41NiwtMC4wMDIgLTE0OS40OCw1NC43IC0xNDkuNDgsMTIyLjE5IDAsNjcuNDggNjYuOTIsMTIyLjE4IDE0OS40OCwxMjIuMTggeiBtIDAsMTIuNTggYyAtOTEuMDYsMCAtMTY0Ljg3LC02MC4zNCAtMTY0Ljg3LC0xMzQuNzYgMCwtNzQuNDMgNzMuODEsLTEzNC43NyAxNjQuODcsLTEzNC43NyA5MS4wNSwwIDE2NC44Nyw2MC4zNCAxNjQuODcsMTM0Ljc3IDAsNzQuNDIgLTczLjgyLDEzNC43NiAtMTY0Ljg3LDEzNC43NiB6Ii8+CiAgPHBhdGggc3R5bGU9ImZpbGw6IzhjYzhmNiIgaWQ9InBhdGgyMSIgZD0ibSAyNjAuNTgsMjQ3LjIgYyA4LjI2LC0xLjA4IDE2LjMxLC0zLjM5IDIzLjY5LC02LjkxIGwgNDcuMDUsMTkuMzQgLTI0LjY1LC0zNi45MiBjIDEyLjk3LC0xNi43MSAxMi45NywtMzcuODYgMCwtNTQuNTYgbCAyNC42NSwtMzYuOTIgLTQ0LjM1LDE4LjIyIC0xLjI3LDE5LjUgYyAxOC4xMywxNS4yMSAxNy42NywzOS4yMiAtMS4xNyw1My45IC00LjgxLDMuODQgLTEwLjMzLDYuNTQgLTE2LjEyLDguNDYgbCAtNy44MywxNS44OSB6IG0gLTIyLC0xMDMuNTggYyAtMTMuNjUsMS43OCAtMjMuODgsNi45NSAtMjMuODgsNi45NSBsIC00Ny4wMSwtMTkuMzQgMjQuNjUsMzYuOTIgYyAtMTIuOTcsMTYuNyAtMTIuOTcsMzcuODUgMCw1NC41NiBsIC0yNC42NSwzNi45MiA0NS41NiwtMTguNzMgMC45NywtMTguMjYgYyAtMTkuMDYsLTE1LjE2IC0xOC45MSwtMzkuNDEgMC4zMSwtNTQuNjQgNC40NywtNCAxNC41OCwtOC4yNCAxNi4wOSwtOC40NiBsIDcuOTYsLTE1LjkyIHoiLz4KICA8cGF0aCBpZD0icGF0aDIzIiBkPSJtIDIyNi45NCwxOTEuMyAtNi4yNiwxMTEgNTAuNzEsLTEwMi45NCA3LjI2LC0xMTAuNzkgLTUxLjcxLDEwMi43MyB6IG0gLTAuNjUsODkuNjcgNDAuNjEsLTgyLjQxIC0zNS41OSwtNi40OSAtNS4wMiw4OC45IHoiLz4KPC9nPgo8L3N2Zz4=">
		<br>
		<br>
		<div style="color:#778899; font-weight:bold;">(c) 2017-{$year} unix-world.org</div>
	</div>
</body>
<script>
	document.getElementById('U-A').innerText = String(navigator.userAgent);
</script>
</html>
HTML;
if((string)$_POST['servicename'] == '') {
	die((string)$starter);
} //end if

//-- main handler

$urls = [
//	'https://www.unix-world.org/',
//	'https://www.wikipedia.org/',
	'https://play.google.com/store/apps/category/GAME_PUZZLE/?hl=en&gl=us#8', //.rand(0,3),
//	''
];

header('Access-Control-Allow-Origin: *');
header('Content-type: application/json');
echo json_encode([
	'message' => 'OK',
	'url' => (string) $urls[array_rand($urls)]
]);
file_put_contents('data/post-'.time().'.log', json_encode($_POST));

// end of php code
?>