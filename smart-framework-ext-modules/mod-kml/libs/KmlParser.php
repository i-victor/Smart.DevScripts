<?php
// Module Lib: KmlParser
// Framework: Smart.Framework
// Author: Radu Ilies
// r.2015-04-04

namespace SmartModExtLib\Kml;


//----------------------------------------------------- PREVENT DIRECT EXECUTION
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------


class KmlParser {


//=======================================================
public static function geo_distance($lat1, $lon1, $lat2, $lon2, $unit='km') {
	//--
	// Calculate distance between two points on a globe (Haversine Formula)
	// http://en.wikipedia.org/wiki/Haversine_formula
	//--
	$earth_radius_km = 6371; // kilometers
	$earth_radius_miles = 3960; // miles
	//--
    $x_lat = deg2rad($lat2 - $lat1);
    $x_lon = deg2rad($lon2 - $lon1);
	//--
    $a = sin($x_lat / 2) * sin($x_lat / 2) + cos(deg2rad($lat1)) * cos(deg2rad($lat2)) * sin($x_lon / 2) * sin($x_lon / 2);
    $c = 2 * asin(sqrt($a));
    //--
    $d_km = $earth_radius_km * $c;
    $d_miles = $earth_radius_miles * $c;
    $d_nmiles = $d_miles * 0.8684;
	//--
	$out = $d_km;
	//--
	return $out;
	//--
} //END FUNCTION
//=======================================================


//=======================================================
public static function kml_to_array($y_kml_file) {
	//--
	$kml_content = \SmartFileSystem::read($y_kml_file);
	if(strlen($kml_content) <= 0) {
		@http_response_code(500);
		die('Map read failed: '.$y_kml_file);
	} //end if
	//--
	$tmp_arr = array();
	$tmp_arr = explode('<Document>', $kml_content);
	$kml_content = '<kml>'."\n".'<Document>'."\n".trim($tmp_arr[1]);
	$tmp_arr = array();
	//--
	//echo $kml_content;
	//print_r($tmp_arr);
	//--
	$arr = (new \SmartXmlParser())->transform($kml_content);
	//print_r($arr);
	//--
	$tmp_title = htmlspecialchars(''.$arr['kml']['Document']['name']);
	$tmp_desc = htmlspecialchars(''.$arr['kml']['Document']['description']);
	//--
	$tmp_placemarks = array();
	if(is_array($arr['kml']['Document']['Placemark'][0])) {
		$tmp_placemarks = (array) $arr['kml']['Document']['Placemark'];
	} else {
		$tmp_placemarks[0] = $arr['kml']['Document']['Placemark'];
	} //end if else
	//--
	$tmp_coords = '';
	//print_r($tmp_placemarks);
	for($i=0; $i<count($tmp_placemarks); $i++) {
		//--
		//print_r($tmp_placemarks);
		if(strlen(trim($tmp_placemarks[$i]['MultiGeometry']['LineString']['coordinates'])) > 0) { // google maps multigeometry
			$tmp_coords .= trim($tmp_placemarks[$i]['MultiGeometry']['LineString']['coordinates'])."\n";
		} else { // default is iPhone where I'd go
			$tmp_coords .= trim($tmp_placemarks[$i]['LineString']['coordinates'])."\n";
		} //end if else
		//--
	} //end for
	//echo $tmp_coords;
	$tmp_coords = @str_replace(array("\r\n", "\r"), array("\n", "\n"), $tmp_coords);
	//--
	$tmp_arr_coords = @explode("\n", $tmp_coords);
	//--
	return array('title' => $tmp_title, 'description' => $tmp_desc, 'coordinates' => (array) $tmp_arr_coords);
	//--
} //END FUNCTION
//=======================================================


//=======================================================
public static function parse_arr_coords($path_to_images, $y_arr_coords, $y_path_delta) {
	//--
	$y_arr_coords = (array) $y_arr_coords;
	//print_r($y_arr_coords);
	//--
	$idx = 0;
	$crr_lat = 0;
	$crr_lon = 0;
	$crr_ttl = '';
	$crr_icon_file = '';
	$crr_icon_w = '32';
	$crr_icon_h = '37';
	//--
	$crr_delta = 0 + $y_path_delta; // delta distance in km (0 to display all)
	//--
	$out_arr = array();
	$max_coords = count($y_arr_coords);
	//--
	for($i=0; $i<$max_coords; $i++) {
		//--
		//echo 'Line #'.$i.'<br>';
		//--
		$y_arr_coords[$i] = trim($y_arr_coords[$i]);
		//--
		if(stristr($y_arr_coords[$i], ',') !== false) { // if we find a line with lat,lon
			//--
			$tmp_arr = @explode(',', $y_arr_coords[$i]);
			//--
			$draw = true;
			if($crr_delta > 0) {
				//--
				if(($crr_lat != 0) OR ($crr_lon != 0)) {
					//--
					$calc_dist = self::geo_distance(trim($tmp_arr[1]), trim($tmp_arr[0]), $crr_lat, $crr_lon, 'km');
					//echo $calc_dist.'<br>';
					//--
					if($calc_dist < $crr_delta) {
						$draw = false; // points are too close, skip
					} //end if
					//--
				} //end if
				//--
			} //end if
			//echo $i.'#Calc/'.$draw.'/'.$calc_dist.'['.$crr_lat.','.$crr_lon.']'.' @ ['.trim($tmp_arr[1]).','.trim($tmp_arr[0]).']'.'<br>';
			//--
			if($draw) {
				//--
				if(strlen($tmp_arr[3]) > 0) {
					$tmp_title = ''.$tmp_arr[3];
				} else {
					$tmp_title = 'Marker #'.($idx+1); //.' ('.trim($tmp_arr[1]).', '.trim($tmp_arr[0]).')';
				} //end if else
				//--
				if(strlen($tmp_arr[4]) > 0) {
					$tmp_icon_file = $path_to_images.'img/custom/'.trim($tmp_arr[4]).'.png';
				} else {
					$tmp_icon_file = '';
				} //end if else
				//--
				$idx++;
				//--
				$crr_lat = trim($tmp_arr[1]);
				$crr_lon = trim($tmp_arr[0]);
				$crr_ttl = trim($tmp_title);
				$crr_icon_file = trim($tmp_icon_file);
				//--
				$out_arr[] = array('lat' => $crr_lat, 'lon' => $crr_lon, 'title' => $crr_ttl, 'icon_file' => $crr_icon_file, 'icon_w' => $crr_icon_w, 'icon_h' => $crr_icon_h);
				//--
			} //end if
			//--
		} //end if
		//--
	} //end for
	//--
	//print_r($out_arr);
	return $out_arr;
	//--
} //END FUNCTION
//=======================================================


//=======================================================
public static function build_js_arr($y_arr_parsed_coords, $idx, $y_mode) {
	//--
	$y_arr_parsed_coords = (array) $y_arr_parsed_coords;
	//--
	$max = count($y_arr_parsed_coords);
	//--
	$out_ext_js = '';
	$out_js = "\n";
	//--
	if((string)$y_mode == 'multiline') {
		$out_js .= "the_areas[0] = new Array('PATH_1', 'multiline', 'MULTILINESTRING((";
	} //end if
	//--
	for($i=0; $i<$max; $i++) {
		//--
		$tmp_arr = (array) $y_arr_parsed_coords[$i];
		//--
		if((string)$y_mode == 'points') {
			$out_js .= 'the_markers['.$idx.'] = new Array(\''.sha1(print_r($tmp_arr,1).time().microtime().rand(1000,9999)).'\', '.trim($tmp_arr['lat']).', '.trim($tmp_arr['lon']).', \''.addslashes(trim($tmp_arr['title'])).'\', \'\', \''.addslashes($tmp_arr['icon_file']).'\', '.((int)0+$tmp_arr['icon_w']).', '.((int)0+$tmp_arr['icon_h']).');'."\n";
		} else {
			if($i == 0) {
				$out_ext_js .= 'the_markers[0] = new Array(\''.sha1(print_r($tmp_arr,1).time().microtime().rand(1000,9999)).'\', '.trim($tmp_arr['lat']).', '.trim($tmp_arr['lon']).', \''.addslashes('START').'\', \'\', \''.addslashes($tmp_arr['icon_file']).'\', '.((int)0+$tmp_arr['icon_w']).', '.((int)0+$tmp_arr['icon_h']).');'."\n";
			} //end if
			$out_js .= ''.trim($tmp_arr['lon']).' '.trim($tmp_arr['lat']);
			if($i < ($max - 1)) {
				$out_js .= ', ';
			} else {
				$out_ext_js .= 'the_markers[1] = new Array(\''.sha1(print_r($tmp_arr,1).time().microtime().rand(1000,9999)).'\', '.trim($tmp_arr['lat']).', '.trim($tmp_arr['lon']).', \''.addslashes('END').'\', \'\', \''.addslashes($tmp_arr['icon_file']).'\', '.((int)0+$tmp_arr['icon_w']).', '.((int)0+$tmp_arr['icon_h']).');'."\n";
			} //end if
		} //end if else
		//--
		$idx++;
		//--
	} //end for
	//--
	if((string)$y_mode == 'multiline') {
		$out_js .= "))');"."\n";
		$out_js .= $out_ext_js."\n";
	} //end if
	//--
	return $out_js;
	//--
} //END FUNCTION
//=======================================================


} //END CLASS


//end of php code
?>