
// valac --pkg gtk+-3.0 hello-world.vala

public class MyApplication : Gtk.Application {
	protected override void activate () {
		var window = new Gtk.ApplicationWindow (this);
		var label = new Gtk.Label ("Hello World !");
		window.add (label);
		window.set_title ("Vala Sample");
		window.set_default_size (320, 200);
		window.show_all ();
	}
}

public int main (string[] args) {
	return new MyApplication ().run (args);
}

