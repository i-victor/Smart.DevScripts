<?php

namespace SmartModExtLib\Cart;

/**
 * Cart: A very simple PHP cart library.
 * Copyright (c) 2017 Sei Kan
 * Distributed under the terms of the MIT License.
 * Redistributions of files must retain the above copyright notice.
 *
 * Copyright (c) 2018 Radu Ovidiu I. / unix-world.org
 *
 */

final class ecommCart {

	/**
	 * An unique ID for the cart.
	 *
	 * @var string
	 */
	private $cartId;

	/**
	 * Maximum item allowed in the cart.
	 *
	 * @var int
	 */
	private $cartMaxItem = 250;

	/**
	 * Maximum quantity of a item allowed in the cart.
	 *
	 * @var int
	 */
	private $itemMaxQuantity = 999999.99;

	/**
	 * Enable or disable cookie.
	 *
	 * @var bool
	 */
	private $useCookie = false;

	/**
	 * A collection of cart items.
	 *
	 * @var array
	 */
	private $items = [];


	/**
	 * Initialize cart.
	 *
	 * @param array $options
	 */
	public function __construct($options=[]) {
		if(isset($options['cartMaxItem']) && preg_match('/^\d+$/', $options['cartMaxItem'])) {
			$this->cartMaxItem = $options['cartMaxItem'];
		} //end if
		if(isset($options['itemMaxQuantity']) && preg_match('/^\d+$/', $options['itemMaxQuantity'])) {
			$this->itemMaxQuantity = $options['itemMaxQuantity'];
		} //end if
		if(isset($options['useCookie']) && $options['useCookie']) {
			$this->useCookie = true;
		} //end if
		$this->cartId = 'eComm_Cart';
		$this->read();
	} //END FUNCTION


	/**
	 * Get items in  cart.
	 *
	 * @return array
	 */
	public function getItems() {
		return $this->items;
	} //END FUNCTION


	/**
	 * Check if the cart is empty.
	 *
	 * @return bool
	 */
	public function isEmpty() {
		return empty(array_filter($this->items));
	} //END FUNCTION


	/**
	 * Get the total of item in cart.
	 *
	 * @return int
	 */
	public function getTotalItem() {
		$total = 0;
		if(!is_array($this->items)) {
			$this->items = array();
		} //end if
		foreach($this->items as $key => $items) {
			if(is_array($items)) {
				foreach($items as $kk => $item) {
					++$total;
				} //end foreach
			} //end if
		} //end foreach
		return $total;
	} //END FUNCTION


	/**
	 * Get the total of item quantity in cart.
	 *
	 * @return int
	 */
	public function getTotalQuantity() {
		$quantity = 0;
		if(!is_array($this->items)) {
			$this->items = array();
		} //end if
		foreach($this->items as $key => $items) {
			if(is_array($items)) {
				foreach($items as $kk => $item) {
					$quantity += $item['quantity'];
				} //end foreach
			} //end if
		} //end foreach
		return $quantity;
	} //END FUNCTION


	/**
	 * Get the sum of a attribute from cart.
	 *
	 * @return int
	 */
	public function getAttributeTotal() {
		$total = 0;
		if(!is_array($this->items)) {
			$this->items = array();
		} //end if
		foreach($this->items as $key => $items) {
			if(is_array($items)) {
				foreach($items as $kk => $item) {
					$total += $item['sell']['price'] * $item['quantity'];
				} //end foreach
			} //end if
		} //end foreach
		return $total;
	} //END FUNCTION


	/**
	 * Remove all items from cart.
	 */
	public function clear() {
		$this->items = [];
		$this->write();
	} //END FUNCTION


	private function calculateHash($attributes) {
	//	return sha1(json_encode($attributes));
		return sha1(\Smart::json_encode($attributes, false, true, false));
	} //END FUNCTION


	/**
	 * Check if a item exist in cart.
	 *
	 * @param string $id
	 * @param array  $attributes
	 *
	 * @return bool
	 */
	public function isItemExists($id, $attributes=[]) {
		$id = (string) $id;
		$attributes = (is_array($attributes)) ? $attributes : [];
		if(isset($this->items[$id])) {
			$hash = $this->calculateHash($attributes);
			if(is_array($this->items[$id])) {
				foreach($this->items[$id] as $key => $item) {
					if($item['hash'] == $hash) {
						return true;
					} //end if
				} //end foreach
			} //end if
		} //end if
		return false;
	} //END FUNCTION


	/**
	 * Add item to cart.
	 *
	 * @param string $id
	 * @param int    $quantity
	 * @param array  $attributes
	 *
	 * @return bool
	 */
	public function add($id, $attributes, $quantity=1, $sell=[]) {
		$id = (string) $id;
		$attributes = (is_array($attributes)) ? array_filter($attributes) : []; // must filter out non-existent keys
		$quantity = (preg_match('/^\d+$/', $quantity)) ? $quantity : 1;
		if($quantity <= 0) {
			$quantity = 1;
		} //end if
		$sell = (array) $sell;
		$sell['price'] = (float) $sell['price'];
		$sell['tax'] = (float) $sell['tax'];
		$hash = $this->calculateHash($attributes);
		if(count($this->items) >= $this->cartMaxItem && $this->cartMaxItem != 0) {
			return false;
		} //end if
		if(is_array($this->items[$id])) {
			foreach($this->items[$id] as $index => $item) {
				if((string)$item['hash'] == (string)$hash) {
					$this->items[$id][$index]['quantity'] += $quantity;
					$this->items[$id][$index]['quantity'] = ($this->itemMaxQuantity < $this->items[$id][$index]['quantity'] && $this->itemMaxQuantity != 0) ? $this->itemMaxQuantity : $this->items[$id][$index]['quantity'];
					$this->write();
					return true;
				} //end if
			} //end foreach
		} //end if
		$this->items[$id][] = [
			'hash'       => (string) $hash,
			'id'         => (string) $id,
			'attributes' => (array) $attributes,
			'quantity'   => ($quantity > $this->itemMaxQuantity && $this->itemMaxQuantity != 0) ? $this->itemMaxQuantity : $quantity,
			'sell'       => (array) $sell
		];
		$this->write();
		return true;
	} //END FUNCTION


	/**
	 * Update item quantity.
	 *
	 * @param string $id
	 * @param int    $quantity
	 * @param array/string  $attributes
	 *
	 * @return bool
	 */
	public function update($id, $attributes, $quantity=1) {
		$id = (string) $id;
		$quantity = (preg_match('/^\d+$/', $quantity)) ? $quantity : 1;
		if($quantity == 0) {
			$this->remove($id, $attributes);
			return true;
		} //end if
		if(is_array($this->items[$id])) {
			if(is_array($attributes)) {
				$hash = $this->calculateHash(array_keys($attributes));
			} else {
				$hash = (string) $attributes;
				$attributes = [];
			} //end if else
			if((string)$hash == '') {
				return false;
			} //end if
			foreach($this->items[$id] as $index => $item) {
				if((string)$item['hash'] == (string)$hash) {
					$this->items[$id][$index]['quantity'] = $quantity;
					$this->items[$id][$index]['quantity'] = ($this->itemMaxQuantity < $this->items[$id][$index]['quantity'] && $this->itemMaxQuantity != 0) ? $this->itemMaxQuantity : $this->items[$id][$index]['quantity'];
					$this->write();
					return true;
				} //end if
			} //end foreach
		} //end if
		return false;
	} //END FUNCTION


	/**
	 * Remove item from cart.
	 *
	 * @param string $id
	 * @param array/string  $attributes
	 *
	 * @return bool
	 */
	public function remove($id, $attributes) {
		$id = (string) $id;
		if(!is_array($this->items[$id])) {
			return false;
		} //end if
		if(is_array($attributes)) {
			if(empty($attributes)) {
				unset($this->items[$id]);
				$this->write();
				return true;
			} //end if
			$hash = $this->calculateHash(array_keys($attributes));
		} else {
			$hash = (string) $attributes;
				$attributes = [];
		} //end if else
		if((string)$hash == '') {
			return false;
		} //end if
		foreach($this->items[$id] as $index => $item) {
			if((string)$item['hash'] == (string)$hash) {
				unset($this->items[$id][$index]);
				$this->write();
				return true;
			} //end if
		} //end foreach
		return false;
	} //END FUNCTION


	/**
	 * Destroy cart session.
	 */
	public function destroy() {
		$this->items = [];
		if($this->useCookie) {
			\SmartUtils::set_cookie((string)$this->cartId, '', -1);
		} else {
			\SmartSession::set((string)$this->cartId, null);
		} //end if else
		return true;
	} //END FUNCTION


	/**
	 * Read items from cart session.
	 */
	private function read() {
		if($this->useCookie) {
			$this->items = \Smart::json_decode(\SmartUtils::data_unarchive(\SmartFrameworkRegistry::getCookieVar((string)$this->cartId)));
		} else { // session
			$this->items = \SmartSession::get((string)$this->cartId);
		} //end if else
		if(!is_array($this->items)) {
			$this->items = [];
		} //end if
		return true;
	} //END FUNCTION


	/**
	 * Write changes into cart session.
	 */
	private function write() {
		if($this->useCookie) {
			\SmartUtils::set_cookie($this->cartId, (string)\SmartUtils::data_archive((string)\Smart::json_encode((array)$this->items)), time() + 604800);
		} else {
			\SmartSession::set((string)$this->cartId, (array)$this->items);
		} //end if else
		return true;
	} //END FUNCTION


} //END CLASS


//end of php code
?>