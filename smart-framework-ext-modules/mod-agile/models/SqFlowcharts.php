<?php
// Class: \SmartModDataModel\Agile\SqFlowcharts
// Author: Radu Ovidiu I.

namespace SmartModDataModel\Agile;

//----------------------------------------------------- PREVENT DIRECT EXECUTION
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------


//=====================================================================================
//===================================================================================== CLASS START
//=====================================================================================


final class SqFlowcharts {

	// ->
	// v.170825

	private $db = null;
	private $sqdb = '#db/flowcharts.sqlite3';


	public function __construct() {
		//--
		$this->db = new \SmartSQliteDb((string)$this->sqdb);
		$this->db->open();
		//--
		if(!\SmartFileSystem::is_type_file((string)$this->sqdb)) {
			if($this->db instanceof \SmartSQliteDb) {
				$this->db->close();
			} //end if
			throw new \Exception('FLOWCHARTS: DB SQLITE File does NOT Exists !');
			return;
		} //end if
		//--
		if(!$this->initDBSchema()) {
			throw new \Exception('FLOWCHARTS: Failed to Initialize the Flowcharts Data Table !');
			return;
		} //end if
		//--
	} //END FUNCTION


	public function __destruct() {
		//--
		if(!$this->db instanceof \SmartSQliteDb) {
			return;
		} //end if
		//--
		$this->db->close();
		//--
	} //END FUNCTION


	public function getOneByUuid($uuid) {
		//--
		if(!$this->checkInstance()) {
			return array();
		} //end if
		//--
		$uuid = (string) trim((string)$uuid);
		if((string)$uuid == '') {
			return array();
		} //end if
		//--
		return (array) $this->db->read_asdata('SELECT * FROM `flow_charts` WHERE (`uuid` = ?) ORDER BY `id` DESC LIMIT 1 OFFSET 0',
			[
				(string) $uuid
			]
		);
		//--
	} //END FUNCTION


	public function getAllByUuid($limit=100) {
		//--
		if(!$this->checkInstance()) {
			return array();
		} //end if
		//--
		return (array) $this->db->read_adata(
			'SELECT * FROM `flow_charts` GROUP BY `uuid` ORDER BY `id` DESC LIMIT '.(int)$limit.' OFFSET 0'
		);
		//--
	} //END FUNCTION


	public function getNewUuid() {
		//--
		return (string) \Smart::uuid_10_seq().'-'.\Smart::uuid_10_num().'-'.\Smart::uuid_10_str(); // str 32 chars, very unique
		//--
	} //END FUNCTION


	public function saveData($data, $user='') {
		//--
		$data = (array) $data;
		//--
		$newdata = array();
		//--
		$newdata['uuid'] = (string) trim((string)$data['uuid']);
		if((string)$newdata['uuid'] == '') {
			return -1; // empty uuid
		} //end if
		//--
		$newdata['dtime'] = (string) date('Y-m-d H:i:s');
		$newdata['user'] = (string) trim((string)$user);
		//--
		$newdata['title'] = (string) trim((string)$data['title']);
		if((string)$newdata['title'] == '') {
			return -2; // invalid title
		} //end if
		//--
		$newdata['saved_data'] = (string) trim((string)$data['saved_data']);
		if((string)$newdata['saved_data'] == '') {
			return -3; // invalid json
		} //end if
		//--
		$compare = (array) $this->getOneByUuid((string)$newdata['uuid']);
		if(((string)$compare['uuid'] === (string)$newdata['uuid']) AND ((string)$compare['title'] === (string)$newdata['title']) AND ((string)$compare['saved_data'] === (string)$newdata['saved_data'])) {
			return 1; // data not changed ...
		} //end if
		//--
		$wr = $this->db->write_data(
			'INSERT INTO `flow_charts` '.$this->db->prepare_statement(
				(array) $newdata,
				'insert'
			)
		);
		//--
		return (int) $wr[1];
		//--
	} //END FUNCTION


//--

	private function checkInstance() {
		//--
		if(!$this->db instanceof \SmartSQliteDb) {
			throw new \Exception('Invalid FLOWCHARTS DB Connection !');
			return 0;
		} //end if
		//--
		return 1;
		//--
	} //END FUNCTION


	private function initDBSchema() {
		//--
		if(!$this->db instanceof \SmartSQliteDb) {
			throw new \Exception('Invalid FLOWCHARTS DB Connection !');
			return 0;
		} //end if
		//--
		if($this->db->check_if_table_exists('flow_charts') != 1) { // create table if not exists
			//--
			$this->db->write_data("CREATE TABLE 'flow_charts' ('id' INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, 'dtime' character varying(23) NOT NULL, 'uuid' character varying(32) NOT NULL, 'title' character varying(255) NOT NULL, 'saved_data' TEXT NOT NULL, 'user' character varying(96) NOT NULL); CREATE INDEX 'uuid' ON 'flow_charts' ('uuid');");
			//--
		} //end if
		//--
		if($this->db->check_if_table_exists('flow_charts') != 1) { // create table if not exists
			//--
			return 0;
			//--
		} //end if
		//--
		return 1;
		//--
	} //END FUNCTION


} //END CLASS


//=====================================================================================
//===================================================================================== CLASS END
//=====================================================================================


// end of php code
?>