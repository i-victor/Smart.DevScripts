
using GLib;
using SDL;
using SDLGraphics;

public class Window : GLib.Object {

	// ATTRIBUTES

	// Window data
	public uint16 screen_width;
	public uint16 screen_height;
	public uint16 screen_bpp;
	public unowned SDL.Screen screen;
	public string caption;
	// Game core data
	public int delay { get; set; }
	private bool keep_running = true;

	// Ball placeholder
	public Ball[]? balls = new Ball[6];

	// METHODS

	public void video_init () {
		// Make the window with the specified data
		uint32 video_flags =    SurfaceFlag.DOUBLEBUF
							|       SurfaceFlag.HWACCEL
							|       SurfaceFlag.HWSURFACE;

		this.screen = Screen.set_video_mode (this.screen_width,
			this.screen_height,
			this.screen_bpp,
			video_flags);

		if (this.screen == null) {

			stderr.printf ("Unable to set video mode.\n");
		}

		// If window is succesfully created, set secondary data
		SDL.WindowManager.set_caption (this.caption,"");
	}

	public delegate void screen_update_delegate (Window window);
	public screen_update_delegate screen_update;

	public void run () {
		// First, start the window
		//this.video_init ();

		// Main loop of the Window
		while (this.keep_running) {

			this.screen_update (this);
			this.event_process ();
			SDL.Timer.delay (this.delay);
		}
	}

	public void event_process () {

		Event event;
		while (Event.poll (out event) == 1) {
			switch (event.type) {
				case EventType.QUIT:
					this.keep_running = false;
					break;
				case EventType.KEYDOWN:
					stderr.printf ("TODO: event_keydown function not yet implemented.\n");
					//this.event_keydown (event.key);
					break;
			}
		}
	}

	public void event_keydown (KeyboardEvent event) {

		Key key = event.keysym;
		// etc...
	}

	// CONSTRUCTOR
	public Window (uint16 width = 800,
	  uint16 height = 640,
	  uint16 bpp = 32,
	  string caption = "Untitled Window") {

		// Make window
		this.screen_width = width;
		this.screen_height = height;
		this.screen_bpp = bpp;
		this.caption = caption;
		this.delay = 5;

		//video_init ();
	}
}
