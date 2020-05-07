
// valac --pkg gtk+-3.0 hello-world2.vala

// Sample Vala GTK : Programatic Mode

using Gtk;

int main (string[] args) {

	Gtk.init(ref args);

	var window = new Window(); // create a window
	window.title = "First GTK+ Program"; // set window title
	window.set_default_size(350, 70); // set default window size(width and height)

	window.window_position = WindowPosition.CENTER; // position the window to the center of screen
	window.border_width = 10; // set window border size

	var button = new Button.with_label ("Click me!"); // create a button
	button.clicked.connect (() => {
		button.label = "Thank you"; // action for click the button
	});
	window.add(button); // add the button to the window

	window.destroy.connect(Gtk.main_quit); // on window close, exit the program

	window.show_all(); // display window
	Gtk.main(); // gtk main loop

	return 0;

}
