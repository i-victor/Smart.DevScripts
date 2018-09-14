<?php
// Controller: kml.test
// Route: ?/page/kml.test
// Author: Radu Ilies
// r.2015-04-04

//----------------------------------------------------- PREVENT EXECUTION BEFORE RUNTIME READY
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------

define('SMART_APP_MODULE_AREA', 'SHARED'); // INDEX, ADMIN, SHARED

// Custom map markers:
//http://www.iconsdb.com/soylent-red-icons/map-marker-5-icon.html

//http://mapicons.nicolasmollet.com/
// hotels: 1AB1E7, important as eb0e0e


class SmartAppIndexController extends SmartAbstractAppController {


public function Run() {

//--
$mode = $this->RequestVarGet('mode', '', 'string');
//--

//--
if(((string)$mode == 'points') OR ((string)$mode == 'markers')) {
	$the_display_mode = 'points';
} else {
	$the_display_mode = 'path';
} //end if else
//--
$the_default_map = 'openstreetmap';
//--
$the_proxy_cache_url = '';
$the_proxy_cache_url = '?/page/maps-cache.mapnik-proxy'; // uncomment this to use the proxy cache
$the_proxy_buf_level = '1'; // default is zero, load only needed
//--

$the_kml_guides = $this->ControllerGetParam('module-path').'kml/obiective.kml';

//$the_kml_file = $this->ControllerGetParam('module-path').'kml/tarnita-padis.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/tarnita-belis-rachitele-tranis.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/baraj-dragan-around.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/stana_de_vale-valea_draganului.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/stana_de_vale-valea_iadului-lesu.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/alesd-stana_de_vale.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/vladeasa-rachitele-giurcuta-horea-cluj.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/rachitele-vladeasa.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/drafts/rachitele-valea-stanciului.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/drafts/stana-coborare.kml';
//$the_kml_file = $this->ControllerGetParam('module-path').'kml/stana-de-vale-aria-vulturilor.kml';
$the_kml_file = $this->ControllerGetParam('module-path').'kml/stana-de-vale-padis.kml';

//===

$kml_arr = \SmartModExtLib\Kml\KmlParser::kml_to_array($the_kml_file);
$guides_arr = \SmartModExtLib\Kml\KmlParser::kml_to_array($the_kml_guides);

$title = $kml_arr['title'];
$description = $kml_arr['description'];


//--
$path_ratio = 0.2; // this ~ 250 m based on measure by map legend
$newarr = (array) \SmartModExtLib\Kml\KmlParser::parse_arr_coords($this->ControllerGetParam('module-path'), $kml_arr['coordinates'], $path_ratio);
//--
$coordinates = "\n";
//--

//--
if((string)$the_display_mode == 'path') {
	$coordinates .= \SmartModExtLib\Kml\KmlParser::build_js_arr($newarr, 0, 'multiline'); // path
	if(count($newarr) <= 0) {
		$nextarrindex = 0;
	} elseif(count($newarr) == 1) {
		$nextarrindex = 1; // start point
	} elseif(count($newarr) > 1) {
		$nextarrindex = 2; // start and end point
	} //end if else
	$coordinates .= \SmartModExtLib\Kml\KmlParser::build_js_arr((array) \SmartModExtLib\Kml\KmlParser::parse_arr_coords($this->ControllerGetParam('module-path'), $guides_arr['coordinates'], 0), $nextarrindex, 'points'); // spots
} else {
	$coordinates .= \SmartModExtLib\Kml\KmlParser::build_js_arr($newarr, 0, 'points'); // spots
	$coordinates .= \SmartModExtLib\Kml\KmlParser::build_js_arr((array) \SmartModExtLib\Kml\KmlParser::parse_arr_coords($this->ControllerGetParam('module-path'), $guides_arr['coordinates'], 0), count($newarr), 'points'); // spots
} //end if else
//--

//-- map zoom
$zoom = '12';
//-- set the map center
$latitude = 0;
$longitude = 0;
//--
$middlepoint_idx = floor(count($newarr) / 2);
$middlepoint_arr = $newarr[$middlepoint_idx];
if(is_array($middlepoint_arr)) {
	$latitude = trim($middlepoint_arr['lat']);
	$longitude = trim($middlepoint_arr['lon']);
} //end if else
//--

$tmp_main = <<<HTML
	<link rel="stylesheet" href="modules/mod-js-components/views/js/jsmaps/theme/default/style.css"  type="text/css">
	<script type="text/javascript" src="modules/mod-js-components/views/js/jsmaps/openlayers.js"></script>
	<script type="text/javascript" src="modules/mod-js-components/views/js/jsmaps/openlayers-google.js"></script>
	<script type="text/javascript" src="modules/mod-js-components/views/js/jsmaps/smart/smartmaps.js"></script>
	<script type="text/javascript">
		//--
		var map_was_init = false;
		var GoogleMapsLoaded = false;
		//--
		var the_proxy_buf_level = {$the_proxy_buf_level};
		var the_proxy_cache_url = '{$the_proxy_cache_url}';
		var the_proxy_mode = 'mapnik'; // can be: custom-mapnik
		//--
		OpenLayers.ImgPath = 'modules/mod-js-components/views/js/jsmaps/img/';
		var MySmartMap = new Smart_Maps(the_proxy_buf_level, the_proxy_cache_url, the_proxy_mode);
		MySmartMap.SetMapDivID('map1');
		MySmartMap.setDebug(0); // enable debug
		var the_map_type = 'openstreetmap'; // the selected map
		//--
		// Full Open Map(s): openstreetmap
		// Open Map(s): mapquest | mapquest-aerial | cyclemap | cyclemap-transport
		// Commercial Maps (Google v3): google | google-physical | google-hybrid | google-aerial
		var the_map_type = '{$the_default_map}'; // the selected map
		//--
		function Init_SmartMaps() {
			//--
			var MyBrowserHeight = $(window).height() - 5;
			//-- reset the Canvas (this is required for switching from a Map type to another (example: switch from Openstreetmaps to Googlemaps)
			\$('#map1_canvas').html('<div id="map1" style="width: 99%; height: ' + MyBrowserHeight + 'px; border: 1px solid #CCCCCC;"></div>');
			//-- draw the map
			MySmartMap.DrawMap(the_lat, the_lon, the_zoom, the_markers, the_areas, the_map_type);
			//--
		} //END FUNCTION
		//--
	</script>
	<script type="text/javascript">
		//--
		//MySmartMap.setIconOffsetMode(false); // for paths we may want icons centered
		//MySmartMap.setMarkerIcon('{$this->ControllerGetParam('module-path')}img/path-circle.png', 7, 7);
		MySmartMap.SetAreasStyle("#FFCC00", 0.75, 8, '#003399', 0.5);
		//--
		var the_markers = new Array();
		var the_areas = new Array();
		//--
		{$coordinates}
		//--
		var the_lat = {$latitude};
		var the_lon = {$longitude};
		var the_zoom = {$zoom};
		//--
		function Render_SmartMaps() {
			//--
			if(map_was_init) {
				the_zoom = MySmartMap.getCurrentZoom();
				the_lat = MySmartMap.getCurrentLat();
				the_lon = MySmartMap.getCurrentLon();
			} //end if
			//--
			map_was_init = true;
			//--
			\$('#map1_canvas').html('<div style="position:fixed; top:100px; z-index:7101; width: 99%; text-align: center;"><span style="color: #FFFFFF; background-color: #1D0092; opacity:0.75; border-radius: 8px; padding: 5px; font-size:40px;">... Rendering Map Data ...</span></div>');
			//--
			if((the_map_type === 'google') || (the_map_type === 'google-physical') || (the_map_type === 'google-hybrid') || (the_map_type === 'google-aerial')) {
				(function(d) { // load the google maps javascript and callback later the Init_SmartMaps (after loading)
					if(GoogleMapsLoaded) {
						Init_SmartMaps();
					} else {
						GoogleMapsLoaded = true;
						var js;
						var id = 'googlemapsjs';
						var ref = d.getElementsByTagName('script')[0];
						js = d.createElement('script');
						js.id = id;
						js.async = true;
						js.type = "text/javascript";
						js.src = "//maps.google.com/maps/api/js?v=3&callback=Init_SmartMaps";
						ref.parentNode.insertBefore(js, ref);
					} //end if else
				}(document));
			} else {
				Init_SmartMaps();
			} //end if else
			//--
		} //END FUNCTION
		//--
		function draw_the_map() {
			Render_SmartMaps();
			MySmartMap.RestrictExtent(20, 43.5, 30.5, 48.5); // example to restrict map to RO area
		} //END FUNCTION
		//--
	</script>
	<div style="position:fixed; top:5px; left:50px; z-index:7023; width: 700px; padding: 5px; color: #000000; background-color: #FFFFFF; opacity:0.75; border-radius: 8px;">
	<span style="font-size: 20px; font-weight: bold;">Path: {$title}</span>
	<br>
	<span style="font-size: 12px;">{$description}</span>
	<br><span id="map1_data"></span>
	</div>
	<!-- form -->
		<!-- extra #start (this is required only for demo) -->
		<span id="DebugData">
		<select id="SmartMapsControlType" onchange="the_map_type = this.value; Render_SmartMaps();" style="position:fixed; top:3px; right:10px; z-index:7021; color: #FFFFFF; background-color: #1D0092; font-size: 14px; padding: 2px; width: 190px; border: 0px; border-radius: 4px;">
			<option value="openstreetmap">OSM Maps</option>
			<option value="cyclemap">Cycle Maps</option>
			<option value="google">Google Maps</option>
		</select>
		</span>
		<!-- extra #end -->
	<!-- end form -->
	<div id="map1_canvas"></div>
	<div style="position:fixed; bottom:1px; left:1px; z-index:7022; text-align: center; font-size: 10px; width:700px; border: 1px solid #ECECEC; color: #000000; background-color: #FFFFFF; opacity:0.7;">
		<div id="map1_lonlat">&nbsp;</div>
	</div>
	<script type="text/javascript">
		//--
		\$('#SmartMapsControlType').val(the_map_type);
		//--
		var MyMapActionMode = 0;
		MySmartMap.OperationModeSwitch(MyMapActionMode);
		\$('#SmartMapsControlAction').val(MyMapActionMode);
		//--
		$(document).ready(function(){
			draw_the_map();
		});
		$(window).resize(function() {
			draw_the_map()
		});
		//--
	</script>
HTML;

//$tmp_main = $coordinates;

//-- title
$this->PageViewSetVar(
	'title',
	'KML Reader (test)'
);
//--

//-- main content
$this->PageViewSetVar(
	'main',
	$tmp_main
);
//--

} //END FUNCTION

} //END CLASS


class SmartAppAdminController extends SmartAppIndexController {

	// this will clone the IndexAppModule to run exactly the same action in admin.php

} //END CLASS


//end of php code
?>