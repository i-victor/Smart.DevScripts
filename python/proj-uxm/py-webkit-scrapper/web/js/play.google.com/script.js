
// javascript parser for play.google.com
// (c) 2017-2019 Radu.I

var uxmJScriptVersion = 'v.20190916.1637#play.google.com'

function uxmAjaxCall(url, method, data) {
	//--
	return jQuery.ajax({
		cache: false,
		dataType: 'json',
		type: String(method),
		url: String(url),
		data: data
	});
	//--
} //END FUNCTION


function uxmRunSpider(scriptUrl) {
	//--
	var myRawLinks = '';
	var mylinks = [];
	var myNumBttns = 0;
	jQuery('a.JC71ub').each(function(index){
		mylinks.push(jQuery(this).attr('href'));
	});
	myRawLinks = mylinks.join('\\n');
	//alert('Links: ' + myRawLinks);
	jQuery('a.LkLjZd.ScJHi.U8Ww7d.xjAeve.nMZKrb.id-track-click').each(function(index){
		myNumBttns++;
	});
	//--
	var myData = {
		servicename: 'Webkit',
		servicetype: 'play.google.com',
		buttons: Number(myNumBttns),
		crrlink: String(window.location.href),
		urldataset: String(myRawLinks)
	};
	var myDivOK = '<div style="z-index:2147483647; position:fixed; top:10px; left:5px; width:96vw; min-height:200px; background:#E6F3FF; border:1px solid #729FCF; padding:10px; font-size:20px; font-weight:bold; opacity:0.85;">Service POST: REPLY OK</div>';
	var myDivERR = '<div style="z-index:2147483647; position:fixed; top:10px; left:5px; width:96vw; min-height:200px; background:#FFE6E6; border:1px solid #EF2929; padding:10px; font-size:20px; font-weight:bold; opacity:0.85;">Service POST: REPLY ERROR</div>';
	var myDivFAIL = '<div style="z-index:2147483647; position:fixed; top:10px; left:5px; width:96vw; min-height:200px; background:#FFE6E6; border:1px solid #EF2929; padding:10px; font-size:20px; font-weight:bold; opacity:0.85;">Service POST: REPLY FAILED</div>';
	//alert(scriptUrl);
	myAjxCall = uxmAjaxCall(scriptUrl, 'POST', myData);
	myAjxCall.always(function(msg){
		//alert('Sending ...');
	}).done(function(msg){
		//alert('Ajax Done !');
		if(msg.url) {
			setTimeout(function(){
				document.location.href = 'data:text/plain;base64,' + btoa(String(msg.url));
			}, 5000);
		} else {
			setTimeout(function(){
				if(typeof uxmScriptUrl != 'undefined') {
					document.location.href = String(uxmScriptUrl); // return to the start
				} else {
					try {
						window.webkit.messageHandlers.COMM_CHANNEL_JS_WK.postMessage('UxmSignalOk'); // emit done signal to exit
					} catch(err){}
				} //end if else
			}, 5000);
		}
		if(msg.message === 'OK') {
			//alert('POST: REPLY OK');
			jQuery('body').append(myDivOK);
		} else {
			//alert('POST: REPLY ERR');
			jQuery('body').append(myDivERR);
		}
	}).fail(function(msg){
		//alert('POST: REPLY FAIL');
		jQuery('body').append(myDivFAIL);
	});
} //END FUNCTION


function uxmRunScript() {
	//--
	//alert('Run (1)');
	var hash = null;
	try {
		hash = window.location.hash;
	} catch(err){}
	//alert('Hash (1): ' + hash);
	var rndTime;
	if(hash) {
		hash = parseInt(hash.substr(1));
		//alert('Hash (2): ' + hash);
		if(!isNaN(hash) && isFinite(hash) && (hash >= 0)) {
			//alert('Hash (3): ' + hash);
			//alert('BtnText: ' + jQuery('a.LkLjZd.ScJHi.U8Ww7d.xjAeve.nMZKrb.id-track-click').eq(hash).text());
			var u = jQuery('a.LkLjZd.ScJHi.U8Ww7d.xjAeve.nMZKrb.id-track-click').eq(hash).attr('href');
			if(u) {
				rndTime = Math.floor(Math.random() * 7001) + 500;
				setTimeout(function(){ window.location = String(u); }, rndTime);
			//	jQuery('a.LkLjZd.ScJHi.U8Ww7d.xjAeve.nMZKrb.id-track-click').eq(hash)[0].click();
			} else {
				alert('FAILED to follow See More button #' + hash);
			} //end if else
			return;
		} else {
			hash = -1;
		} //end if else
	} //end if
	//--
	//alert('Run (2)');
	var uxmIsInitUrl = false;
	if(jQuery('h1#PyWebkitGTK-INIT-Page').text()) {
		uxmIsInitUrl = true;
	} //end if
	//--
	if(uxmIsInitUrl === true) {
		var confirmRun = confirm('\n' + 'Webkit.Crawler.Js: ' + uxmJScriptVersion + '\n\n' + 'Crawler.Backend: ' + uxmScriptUrl + '\n\n\n' + 'Press OK to CONTINUE or CANCEL to EXIT ...' + '\n');
		if(!confirmRun) {
			window.webkit.messageHandlers.COMM_CHANNEL_JS_WK.postMessage('UxmSignalOk'); // emit done signal to exit
			return;
		} //end if
		setTimeout(function(){
			uxmRunSpider(uxmScriptUrl);
		}, 3000);
	} else {
		rndTime = Math.floor(Math.random() * 2001) + 2500;
		window.scrollTo(0, document.body.scrollHeight);
		setTimeout(function(){
			rndTime = Math.floor(Math.random() * 2001) + 1500;
			window.scrollTo(0, document.body.scrollHeight);
			setTimeout(function(){
				rndTime = Math.floor(Math.random() * 3001) + 2500;
				window.scrollTo(0, document.body.scrollHeight);
				setTimeout(function(){
					rndTime = Math.floor(Math.random() * 3001) + 1500;
					window.scrollTo(0, document.body.scrollHeight);
					setTimeout(function(){
						rndTime = Math.floor(Math.random() * 4001) + 2500;
						window.scrollTo(0, document.body.scrollHeight);
						setTimeout(function(){
							jQuery('span.RveJvd.snByac').trigger('click'); // show more
							rndTime = Math.floor(Math.random() * 4001) + 1500;
							setTimeout(function(){
								rndTime = Math.floor(Math.random() * 1001) + 1500;
								setTimeout(function(){
									window.scrollTo(0, document.body.scrollHeight);
									setTimeout(function(){
										uxmRunSpider(uxmScriptUrl);
									}, 1750);
								}, rndTime);
							}, rndTime);
						}, rndTime);
					}, rndTime);
				}, rndTime);
			}, rndTime);
		}, rndTime);
	} //end if else
	//--
} //END FUNCTION

//alert('Loaded ...');
uxmRunScript();

// #END
