
// valac --pkg json-glib-1.0 sample.vala

using Json;

int main(string[] args) {

	Json.Parser parser = new Json.Parser();
	var data = "{\"foo\":\"123abc\", \"num\":123, \"arr\":[\"1\",2,3]}";

	try {

		parser.load_from_data(data, -1);

		var root_object = parser.get_root().get_object();

		var foo = root_object.get_string_member("foo"); // .get_string_member_with_default("foo", ""); // is available for json-glib-1.0 >= 1.6
		stdout.printf("Foo = %s\n", foo);

		var num = root_object.get_int_member("num"); // .get_int_member_with_default("num", 0); // is available for json-glib-1.0 >= 1.6
		stdout.printf("Num = %lld\n", num);

		var arr = root_object.get_array_member("arr");
		foreach(var elem in arr.get_elements()) {
			var elTest = elem.get_string();
			if(elTest == null) {
				elTest = elem.get_int().to_string();
			}
			stdout.printf("Arr[] = %s\n", elTest);
		}

		// https://valadoc.org/json-glib-1.0/Json.Object.get_array_member.html

	} catch (Error e) {

		stderr.printf("Unable to parse the string as Json: %s\n", e.message);

	}

	return 0;

}

