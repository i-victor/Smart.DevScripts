<?php
// [@[#[!SF.DEV-ONLY!]#]@]
// Controller: Samples/BenchMark
// Route: ?/page/samples.benchmark (?page=samples.benchmark)
// Author: unix-world.org
// v.3.7.5 r.2018.03.09 / smart.framework.v.3.7

//----------------------------------------------------- PREVENT EXECUTION BEFORE RUNTIME READY
if(!defined('SMART_FRAMEWORK_RUNTIME_READY')) { // this must be defined in the first line of the application
	@http_response_code(500);
	die('Invalid Runtime Status in PHP Script: '.@basename(__FILE__).' ...');
} //end if
//-----------------------------------------------------

define('SMART_APP_MODULE_AREA', 'INDEX'); // INDEX, ADMIN, SHARED

/**
 * Index Controller
 *
 * @ignore
 *
 */
class SmartAppIndexController extends SmartAbstractAppController {

	public function Run() {

		//-- dissalow run this sample if not test mode enabled
		if(!defined('SMART_FRAMEWORK_TEST_MODE') OR (SMART_FRAMEWORK_TEST_MODE !== true)) {
			$this->PageViewSetErrorStatus(503, 'ERROR: Test mode is disabled ...');
			return;
		} //end if
		//--

		//--
		$op = $this->RequestVarGet('op', '', 'string');
		//--

		//--
		$cart = new \SmartModExtLib\Cart\ecommCart([
			'cartId' 			=> 'eCommCart1',
			'cartMaxItem' 		=> 10, // Maximum item can added to cart, 0 = Unlimited
			'itemMaxQuantity' 	=> 50, // Maximum quantity of a item can be added to cart, 0 = Unlimited
			'useCookie' 		=> false // Do not use cookie, cart items will gone after browser closed
		]);
		//--
		$cart_currency = 'US$';
		//--
		$products = Smart::json_decode('
		[
		   {
			  "id":100,
			  "name":"iPhone SE (32 GB)",
			  "image":{
				 "source":"https:\/\/user-images.githubusercontent.com\/73107\/30917969-cc1b8586-a3cf-11e7-872c-92d98d24afb0.png",
				 "width":200,
				 "height":250
			  },
			  "attributes": {
				  "color":{
					 "silver":"Silver",
					 "gold":"Gold",
					 "space-gray":"Space Gray",
					 "rose-gold":"Rose Gold"
				  },
				  "size":{
					  "std":"Standard",
					  "xs":"XS"
				  }
			  },
			  "price":"349.00",
			  "qtytype":"dec"
		   },
		   {
			  "id":101,
			  "name":"iPhone SE (128 GB)",
			  "image":{
				 "source":"https:\/\/user-images.githubusercontent.com\/73107\/30917969-cc1b8586-a3cf-11e7-872c-92d98d24afb0.png",
				 "width":200,
				 "height":250
			  },
			  "attributes": {
				  "color":{
					 "silver":"Silver",
					 "gold":"Gold",
					 "space-gray":"Space Gray",
					 "rose-gold":"Rose Gold"
				  }
			  },
			  "price":"449.00"
		   },
		   {
			  "id":102,
			  "name":"iPhone 6s (32 GB)",
			  "image":{
				 "source":"https:\/\/user-images.githubusercontent.com\/73107\/30917728-4052e8c8-a3cf-11e7-93df-7ac32ab8dca5.png",
				 "width":157,
				 "height":250
			  },
			  "attributes": {
				  "color":{
					 "silver":"Silver",
					 "gold":"Gold",
					 "space-gray":"Space Gray",
					 "rose-gold":"Rose Gold"
				  }
			  },
			  "price":"449.00"
		   },
		   {
			  "id":103,
			  "name":"iPhone 6s (128 GB)",
			  "image":{
				 "source":"https:\/\/user-images.githubusercontent.com\/73107\/30917728-4052e8c8-a3cf-11e7-93df-7ac32ab8dca5.png",
				 "width":157,
				 "height":250
			  },
			  "attributes": {
				  "color":{
					 "silver":"Silver",
					 "gold":"Gold",
					 "space-gray":"Space Gray",
					 "rose-gold":"Rose Gold"
				  }
			  },
			  "price":"549.00"
		   },
		   {
			  "id":104,
			  "name":"iPhone 6s Plus (32 GB)",
			  "image":{
				 "source":"https:\/\/user-images.githubusercontent.com\/73107\/30917727-405206e2-a3cf-11e7-943c-b9bc44155c24.png",
				 "width":158,
				 "height":250
			  },
			  "attributes": {
				  "color":{
					 "silver":"Silver",
					 "gold":"Gold",
					 "space-gray":"Space Gray",
					 "rose-gold":"Rose Gold"
				  }
			  },
			  "price":"549.00"
		   }
		]
		');
		//--
		if((string)$op == 'cart-json') {
			//--
			$this->PageViewSetCfg('rawpage', true);
			//--
			//print_r($_POST); die();
			$cart_op = $this->RequestVarGet('cart_action', '', 'string');
			$frm = $this->RequestVarGet('frm', [], 'array');
			//--
			$redirect = '';
			$message = '???';
			//-- Empty the cart
			if((string)$cart_op == 'empty') {
				$message = 'Cart cleared';
				$cart->clear();
				$redirect = '?page='.Smart::escape_url($this->ControllerGetParam('controller')).'&op=cart';
			} //end if
			//-- Add item
			if((string)$cart_op == 'add') {
				$cart_item_id = (string) $frm['id'];
				$cart_item_hash = (string) $frm['hash'];
				$cart_item_qty = (string) $frm['qty'];
				$cart_item_atts = (array) $frm['att'];
				$message = 'Product added';
				foreach($products as $key => $product) {
					if((string)$cart_item_id == (string)$product['id']) {
						break;
					} //end if
				} //end foreach
				$cart->add(
					$product['id'],
					(array) $cart_item_atts,
					$cart_item_qty,
					[
						'currency' => (string) $cart_currency,
						'price' => (float) $product['price'],
						'tax' 	=> (int) $product['tax']
					]
				);
			} //end if
			//-- Update item
			if((string)$cart_op == 'update') {
				$redirect = '?page='.Smart::escape_url($this->ControllerGetParam('controller')).'&op=cart';
				if((string)$frm['cart'] == '@cart') {
					$message = 'Cart updated';
					foreach($frm as $key => $val) {
						if(is_array($val)) {
							$cart_item_id = (string) $val['id'];
							$cart_item_qty = (string) $val['qty'];
							$cart_item_hash = (string) $key;
							$cart_item_atts = (array) $val['att'];
							$cart->update(
								(string) $cart_item_id,
								(string) $cart_item_hash, //(array) $cart_item_atts,
								$cart_item_qty,
								false
							);
						} //end if
					} //end foreach
					$cart->write();
				} else {
					$cart_item_id = (string) $frm['id'];
					$cart_item_hash = (string) $frm['hash'];
					$cart_item_qty = (string) $frm['qty'];
				//	$cart_item_atts = (array) $frm['att'];
					$message = 'Product quantity updated';
					$cart->update(
						(string) $cart_item_id,
						(string) $cart_item_hash, //(array) $cart_item_atts,
						$cart_item_qty
					);
				} //end if else
			} //end if
			//-- Remove item
			if((string)$cart_op == 'remove') {
				$cart_item_id = (string) $frm['id'];
				$cart_item_hash = (string) $frm['hash'];
			//	$cart_item_qty = (string) $frm['qty'];
				$message = 'Product removed';
				$redirect = '?page='.Smart::escape_url($this->ControllerGetParam('controller')).'&op=cart';
				foreach($products as $key => $product) {
					if((string)$cart_item_id == (string)$product['id']) {
						break;
					} //end if
				} //end foreach
				$cart->remove(
					$product['id'],
					(string) $cart_item_hash //(array) $cart_item_atts,
				);
			} //end if
			//--
			$answer = 'OK';
			$title = 'Cart';
			//--
			$this->PageViewSetVar(
				'main',
				(string) SmartComponents::js_ajax_replyto_html_form($answer, $title, $message, $redirect)
			);
			return;
			//--
		} elseif((string)$op == 'cart') {
			//--
			$all_items = [];
			$cart_items = [];
			if(!$cart->isEmpty()) {
				$all_items = $cart->getItems();
				//print_r($all_items); die();
				foreach($all_items as $id => $items) {
					foreach($items as $key => $item) {
						foreach($products as $kk => $product) {
							if((string)$id == (string)$product['id']) {
								break;
							} //end if
						} //end if
						$tmp_arr = [];
						$tmp_arr['id'] = $item['id'];
						$tmp_arr['hash'] = $item['hash'];
						if((string)$product['qtytype'] == 'dec') {
							$tmp_arr['quantity'] = 0 + \Smart::format_number_dec($item['quantity'], 2, '.', '');
						} else {
							$tmp_arr['quantity'] = \Smart::format_number_int($item['quantity'], '+');
						} //end if else
						$tmp_arr['name'] = $product['name'];
						$tmp_arr['price'] = $item['sell']['price'];
						$tmp_arr['tax'] = $item['sell']['tax'];
						$tmp_arr['currency'] = $item['sell']['currency'];
						$tmp_arr['attributes'] = (array) $item['attributes'];
//print_r($item); die();
						$cart_items[] = (array) $tmp_arr;
					} //end foreach
				} //end foreach
			} //end if
			//--
			$tpl = 'cart.mtpl.htm';
			$arr = [
				'PAGE-URL' 		=> (string) $this->ControllerGetParam('controller'),
				'CART-CURRENCY' 	=> (string) $cart_currency,
				'CART-TOTAL' 		=> (string) Smart::format_number_dec($cart->getAttributeTotal(), 2, '.', ''),
				'CART-ITEMS' 		=> (array) $cart_items
			];
			//--
		} elseif((string)$op == 'checkout') {
			//--
			return 404; // not implemented
			//--
		} else {
			//--
			$arr = [];
			if(is_array($products)) {
				foreach($products as $key => $val) {
					//print_r($val); die();
//print_r($val['attributes']); die();
					$arr[] = [
						'id' 		=> $val['id'],
						'name' 		=> $val['name'],
						'price' 	=> $val['price'],
						'currency' 	=> $cart_currency,
						'img-src' 	=> $val['image']['source'],
						'img-w' 	=> $val['image']['width'],
						'img-h' 	=> $val['image']['height'],
						'atts' 		=> $val['attributes'],
						'hash' 		=> $cart->calculateHash($val['id'], $val['attributes'])
					];
				} //end foreach
			} //end if
			//--
			$tpl = 'shop.mtpl.htm';
			$arr = [
				'PAGE-URL' 		=> (string) $this->ControllerGetParam('controller'),
				'PRODUCTS-ARR' 	=> (array) $arr
			];
			//--
		} //end if else
		//--
		$this->PageViewSetVars([
			'title' => 'eCommerce Test',
			'main' => SmartMarkersTemplating::render_file_template(
				(string) $this->ControllerGetParam('module-view-path').$tpl,
				(array) $arr,
				'no' // don't use caching (use of caching make sense only if file template is used more than once per execution)
			)
		]);
		//--

	} //END FUNCTION

} //END CLASS

//end of php code
?>