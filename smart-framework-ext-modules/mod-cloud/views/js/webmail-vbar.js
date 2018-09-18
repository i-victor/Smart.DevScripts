
// NetVision WebMail - Vertical Dynamic Bar
// (c) 2006-2018 unix-world.org

// v.180208

//==================================================================
//==================================================================

var NetVision_WebMail_Allow_DinamicBar = true;

var NetVision_WebMail_Vertical_DynamicBar = new function() { // START CLASS

var min_left = 1;
var max_left = 500;

// read previous stored value (from cookie) and set the width of div
this.restoreWidth = function() {
	//--
	var xwidth = 175; // this is by default
	//--
	var the_cookie = parseInt(SmartJS_BrowserUtils.getCookie('NetVision_WebMail_LeftArea_Size'));
	if((!isNaN(the_cookie)) && (the_cookie >= min_left) && (the_cookie <= max_left)) {
		xwidth = the_cookie;
	} //end if
	jQuery('#nv_webmail_left_div').css({ 'width': xwidth });
	//--
} //END FUNCTION

this.handleDrag = function() {
	//--
	if(NetVision_WebMail_Allow_DinamicBar === true) {
		//--
		var area = jQuery('#nv_webmail_left_div');
		var pwidth = parseInt(area.width());
		//--
		jQuery('#nv_webmail_resizer_div').bind('dragstart', function(event) {
			//console.log('drag start');
		}).bind('drag', function(event) {
			//console.log('drag: ' + event.pageX);
			pwidth = Math.round(event.pageX);
			if(pwidth < min_left) {
				pwidth = min_left;
			} else if(pwidth > max_left) {
				pwidth = max_left;
			} //end if
			area.css({
				width: pwidth + 'px'
			});
		}).bind('dragend', function(event) {
			//console.log('drag end');
			SmartJS_BrowserUtils.setCookie('NetVision_WebMail_LeftArea_Size', pwidth, 0, '/');
		});
		//--
	} //end if
	//--
} //END FUNCTION

} //END CLASS

//==================================================================
//==================================================================

// #END
