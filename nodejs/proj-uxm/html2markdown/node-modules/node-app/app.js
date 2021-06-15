
// app.js # nodejs
// convert HTML to Markdown :: r.20210613
// (c) 2021 unix-world.org

// HINT: if get the message: `Trace/BPT trap`, use the --jitless option to run nodejs
// Ex: node --jitless ./modules/mod-docs/node-modules/node-app/app.js ./wpub/devdocs/jquery

(() => {

	const fs = require('fs');

	const fNameJsonDB = 'db.optimized.json';
	let dirPath = null;
	let theJsonDbFile = null;
	let theJsonMdFile = null
	const cmdArgs = process.argv.slice(2);
	if(cmdArgs && Array.isArray(cmdArgs) && cmdArgs.length && cmdArgs[0]) {
		dirPath = String(cmdArgs[0] || '').replace(/(\/)+$/, '').trim();
		if(!dirPath || (dirPath == '/')) {
			console.error('ERR: First argument is empty ... Expected a path to a dir that contains a `' + fNameJsonDB + '` file ...');
			process.exit(2);
			return false;
		}
		dirPath = String(dirPath) + '/'
		const isDirExists = fs.existsSync(dirPath) && fs.lstatSync(dirPath).isDirectory();
		if(!isDirExists) {
			console.error('ERR: First argument is invalid ... Expected a path to a dir that contains a `' + fNameJsonDB + '` file ...');
			process.exit(3);
			return false;
		}
		theJsonDbFile = String(dirPath) + String(fNameJsonDB);
		const isSrcFileExists = fs.existsSync(theJsonDbFile) && fs.lstatSync(theJsonDbFile).isFile();
		if(!isSrcFileExists) {
			console.error('ERR: First argument is wrong ... Expected a path to a dir that contains a `' + fNameJsonDB + '` file ...');
			process.exit(4);
			return false;
		}
	} else {
		console.error('ERR: No arguments ... Expected a path to a dir that contains a `' + fNameJsonDB + '` file ...');
		process.exit(1);
		return false;
	} //end if

	const TurndownService = require('../turndown/lib/turndown');
	const turndownService = new TurndownService();
	const turndownPluginGfm = require('../turndown/lib/turndown-plugin-gfm');
	const gfm = turndownPluginGfm.gfm;
	turndownService.use(gfm);

	theJsonMdFile = dirPath + 'db-md.json';

	fs.access(theJsonMdFile, fs.F_OK, (err) => {
		if(err) {
			console.log('CLEANUP: OK, JSON MD DB does not exist, NO CLEANUP IS NECESSARY ...');
			return true;
		}
		try {
			console.log('CLEANUP: Trying to DELETE JSON MD DB ...');
			fs.unlinkSync(theJsonMdFile);
		} catch(err) {
			console.error('ERR: CLEANUP: FAILED to delete JSON MD DB ...', err);
			process.exit(11);
			return false;
		}
		console.log('CLEANUP: OK, JSON MD DB was DELETED ...');
		return true;
	});

	console.log('Processing DB File', theJsonDbFile);
	fs.readFile(theJsonDbFile, null, (err, theSource) => {

		if(err) {
			console.error('ERR: READ JSON: FAILED to READ ' + theJsonDbFile + ' ...', err);
			process.exit(21);
			return false;
		}

		let json = null;
		try {
			json = JSON.parse(theSource);
		} catch(err) {
			console.error('ERR: READ JSON: FAILED to PARSE ' + theJsonDbFile + ' ...', err);
			process.exit(22);
			return false;
		}

		if((typeof(json) !== 'object') || (json === null)) {
			console.error('ERR: READ JSON: INVALID FORMAT / NOT ASSOCIATIVE ARRAY {OBJECT} # ' + theJsonDbFile + ' ...');
			process.exit(23);
			return false;
		}

		let currentKey = null;
		let currentHTML = null;
		let markdown = null;
		let mDocs = {};
		let loops = 0;
		let convertedOk = 0;
		try {
			Object.keys(json).forEach(key => {
			//	console.log(process.memoryUsage());
				loops++;
				currentKey = String(key || '');
				currentHTML = String(json[key] || '');
				json[key] = null; // free memory
				markdown = null;
				if(!currentKey || !currentHTML) {
					console.warn('WARN: CONVERT JSON: SKIP: INVALID KEY OR INVALID HTML Data at key: `' + currentKey + '` # ' + theJsonDbFile + ' ...');
				} else {
					currentHTML = '<!DOCTYPE html><html><head><meta charset="UTF-8"></head><body>' + currentHTML + '</body></html>';
					markdown = String(turndownService.turndown(currentHTML) || '');
					if(!markdown) {
						console.warn('WARN: CONVERT JSON: EMPTY MARKDOWN at key: `' + currentKey + '` # ' + theJsonDbFile + ' ...');
					} else {
						convertedOk++;
						console.log('CONVERT JSON: OK [PROCESSED] [' + loops + '/' + convertedOk + '] key: `' + currentKey + '` ; Markdown length is:', markdown.length);
						mDocs[currentKey] = markdown;
					}
				} //end if else
			});
		} catch(err) {
			console.error('ERR: CONVERT JSON: FAILED at key: `' + currentKey + '` # ' + theJsonDbFile + ' ...', err);
			process.exit(24);
			return false;
		}

		try {
			mDocs = String(JSON.stringify(mDocs, null, 2) || '');
		} catch(err) {
			mDocs = '';
			console.error('ERR: CONVERT JSON: FAILED to compose JSON MARKDOWN # `' + theJsonMdFile + '` ...', err);
			process.exit(25);
			return false;
		}

		fs.writeFile(theJsonMdFile, mDocs, (err) => {
			if(err) {
				console.error('ERR: WRITE MARKDOWN: FAILED to WRITE `' + theJsonMdFile + '` ...', err);
				process.exit(26);
				return false;
			}
		});

		console.log('WRITE MARKDOWN: DONE, SAVED as `' + theJsonMdFile + '` ... OK');
		return true;

	});

})();

// #END
