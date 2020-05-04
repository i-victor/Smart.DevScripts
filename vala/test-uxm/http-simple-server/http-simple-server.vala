
// valac --pkg libsoup-2.4 http-simple-server.vala

using GLib;

void default_handler(Soup.Server server, Soup.Message msg, string path, GLib.HashTable? query, Soup.ClientContext client) {

	string response_text = """
		<html>
		  <body>
			<p>Current location: %s</p>
			<p><a href="/xml">Test XML</a></p>
		  </body>
		</html>""".printf (path);

	msg.set_response("text/html", Soup.MemoryUse.COPY, response_text.data);

}

void xml_handler(Soup.Server server, Soup.Message msg, string path, GLib.HashTable? query, Soup.ClientContext client) {

	string response_text = "<node><subnode>test</subnode></node>";

	msg.set_response("text/xml", Soup.MemoryUse.COPY, response_text.data);

}

void main () {

	var server = new Soup.Server (Soup.SERVER_PORT, 18888);
	server.add_handler ("/", default_handler);
	server.add_handler ("/xml", xml_handler);
	server.run ();

}


// #end
