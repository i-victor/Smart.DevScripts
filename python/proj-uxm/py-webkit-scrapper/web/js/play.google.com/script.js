
// javascript parser for play.google.com
// (c) 2017-2019 Radu.I

var uxmJScriptVersion = 'v.20190919.1059#play.google.com';
var uxmJSUrlHash = null;

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
	var mybttns = [];
	jQuery('a.JC71ub').each(function(index){
		mylinks.push(jQuery(this).attr('href'));
	});
	myRawLinks = mylinks.join('\\n');
	//alert('Links: ' + myRawLinks);
	jQuery('a.LkLjZd.ScJHi.U8Ww7d.xjAeve.nMZKrb.id-track-click').each(function(index){
		if(jQuery(this).is(':visible')) {
			mybttns.push(index);
		}
	});
	myRawBttns = mybttns.join('\\n');
	//--
	var myData = {
		servicename: 'Webkit',
		servicetype: 'play.google.com',
		crrlink: String(window.location.href),
		buttons: String(myRawBttns),
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
	//--
} //END FUNCTION


function uxmRunFollowButton(hash) {
	//--
	if(hash) {
		hash = parseInt(hash);
		if(!isNaN(hash) && isFinite(hash) && (hash >= 0)) {
			var u = jQuery('a.LkLjZd.ScJHi.U8Ww7d.xjAeve.nMZKrb.id-track-click').eq(hash).attr('href');
			if(u) {
				rndTime = Math.floor(Math.random() * 7001) + 500;
				u += '&hl=en&gl=us';
				setTimeout(function(){ window.location = String(u); }, rndTime); // jQuery('a.LkLjZd.ScJHi.U8Ww7d.xjAeve.nMZKrb.id-track-click').eq(hash)[0].click();
			} else {
				alert('FAILED to follow See More button #' + hash);
			} //end if else
		} //end if else
	} //end if
	//--
} //END FUNCTION


function uxmRunScript() {
	//--
	try {
		uxmJSUrlHash = String(window.location.hash);
	} catch(err){
		uxmJSUrlHash = '';
	} //end try catch
	if(uxmJSUrlHash) {
		uxmJSUrlHash = parseInt(uxmJSUrlHash.substr(1));
		if(isNaN(uxmJSUrlHash) || !isFinite(uxmJSUrlHash) || (uxmJSUrlHash < 0)) {
			uxmJSUrlHash = false;
		} //end if
	} else {
		uxmJSUrlHash = false;
	} //end if else
	//--
	var rndTime;
	//--
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
										if(uxmJSUrlHash) {
											uxmRunFollowButton(uxmJSUrlHash);
										} else {
											uxmRunSpider(uxmScriptUrl);
										} //end if else
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

uxmRunScript();

// #END
