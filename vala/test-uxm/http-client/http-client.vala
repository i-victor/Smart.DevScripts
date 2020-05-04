
// valac --pkg libsoup-2.4 http-client.vala

void main () {

	var session = new Soup.Session ();
	session.ssl_strict = false; // allow self-signed certificates

	var message = new Soup.Message ("GET", "https://127.0.0.1/");

	/* see if we need HTTP auth */
	session.authenticate.connect ((sess, msg, auth, retrying) => {
		if (!retrying) {
			stdout.printf("Authentication required\n");
			auth.authenticate ("user", "password");
		}
	});

	/* send a sync request */
	session.send_message(message);

	var status_code = message.status_code;
	stdout.printf("Status Code: %s\n", status_code.to_string());

	message.response_headers.foreach ((name, val) => {
		stdout.printf("Name: %s -> Value: %s\n", name, val);
	});

	stdout.printf("Message length: %lld\n%s\n", message.response_body.length, message.response_body.data);

}

