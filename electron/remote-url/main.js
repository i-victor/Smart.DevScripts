// main.js

const electron = require ('electron');

const app = electron.app; // electron module
app.allowRendererProcessReuse = true; // get rid of message: (electron) The default value of app.allowRendererProcessReuse is deprecated

const BrowserWindow = electron.BrowserWindow; //enables UI
const Menu = electron.Menu; // menu module

let win;

app.on('ready', _ => {

	win = new BrowserWindow({
		width: 800,
		height: 600,
	});

	const template = [
		{
			label: 'Help',
			submenu: [{ // adds submenu
					label: 'About',
				}, {
					type: 'separator' // horizontal line between submenu items
				},{
					label: 'Quit',
					role: 'quit' // closes app when clicked

				}]
		},
		{
			label: 'Refresh', // Refreshes or reloads the page when clicked
			role: 'reload'
		},
	];
	const menu = Menu.buildFromTemplate(template); // sets the menu
	Menu.setApplicationMenu(menu);

	win.loadURL('http://demo.unix-world.org/smart-framework/')    // loads this URL

});


// #END
