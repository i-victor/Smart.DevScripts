
// valac --pkg gtk+-3.0 hello-world.vala

// Sample Vala GTK : Object Mode

public class MyApplication : Gtk.Application {

	protected override void activate () {

		var window = new Gtk.ApplicationWindow(this); // create new window
		window.set_title("Vala Sample"); // add window title
		window.set_default_size(640, 400); // set default window size(width and height)

		var label = new Gtk.Label("Hello World !"); // create a label
		window.add(label); // add a label to the window

		window.show_all(); // display the window

	}

}

public int main (string[] args) {

	return new MyApplication ().run (args);

}
