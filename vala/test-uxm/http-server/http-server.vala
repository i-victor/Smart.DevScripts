
// valac --pkg libsoup-2.4 http-server.vala

using GLib;

public class NoodleSoupServer : Soup.Server {

	private int access_counter = 0;

	public NoodleSoupServer () {

		Object (port: 18888);
		GLib.assert (this != null);

		// Links:
		//   http://localhost:18888/about.html
		//   http://localhost:18888/index.html
		//   http://localhost:18888/
		this.add_handler ("/quit.html", quit_handler);
		this.add_handler ("/about.html", about_handler);
		this.add_handler ("/index.html", root_handler);
		this.add_handler ("/", root_handler);

		// Links:
		//   http://localhost:18888/*
		//   http://localhost:18888/foo.html
		this.add_handler (null, default_handler);
	}

	private static void root_handler (Soup.Server server, Soup.Message msg, string path, GLib.HashTable? query, Soup.ClientContext client) {
		string html_head = "<head><title>Index</title></head>";
		string html_body = "<body><h1>Index:</h1></body>";
		msg.set_response ("text/html", Soup.MemoryUse.COPY, "<html>%s%s</html>".printf (html_head, html_body).data);
	}

	private static void quit_handler (Soup.Server server, Soup.Message msg, string path, GLib.HashTable? query, Soup.ClientContext client) {
		Process.exit (0);
	}

	private static void about_handler (Soup.Server server, Soup.Message msg, string path, GLib.HashTable? query, Soup.ClientContext client) {
		string html_head = "<head><title>About</title></head>";
		string html_body = "<body><h1>About:</h1></body>";
		msg.set_response ("text/html", Soup.MemoryUse.COPY, "<html>%s%s</html>".printf (html_head, html_body).data);
	}

	private static void default_handler (Soup.Server server, Soup.Message msg, string path, GLib.HashTable? query, Soup.ClientContext client) {
		NoodleSoupServer self = server as NoodleSoupServer;
		GLib.assert (self != null);
		if (msg.uri.get_path () == "/foo.html") {
			// http://localhost:18888/foo.html
			string html_head = "<head><title>Default</title></head>";
			string html_body = "<body><h1>Default:</h1><p>%s</p><p>%u</p></body>".printf (msg.uri.to_string (false), ++self.access_counter);
			msg.set_response ("text/html", Soup.MemoryUse.COPY, "<html>%s%s</html>".printf (html_head, html_body).data);
		} else {
			// 404:
			msg.set_response ("text/html", Soup.MemoryUse.COPY, "<html><head><title>404</title></head><body><h1>404</h1></body></html>".data);
			msg.status_code = 404;
		}
	}

	public static int main (string[] args) {
		NoodleSoupServer server = new NoodleSoupServer ();
		server.run ();
		return 0;
	}
}

// #end
