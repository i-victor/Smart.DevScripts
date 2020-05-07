
// valac -g --pkg sdl --pkg sdl-gfx -X -lSDL_gfx -X -I/usr/include/SDL World.vala Window.vala Ball.vala main.vala -o BouncingBall

using GLib;
using SDL;
using SDLGraphics;

// Re-write of BBWindow's screen_update (linked in the run void)
void screen_update_definition (Window window) {

	window.screen.fill (null,0);

	if (window.balls != null) {
		// If balls exist, iterate through them
		foreach (Ball ball in window.balls) {
			if (ball != null) {
				ball.represent ();
			}
		}
	}

	window.screen.flip ();
}

void main (string[] args) {

	// Make a world for the ball
	World BBWorld = new World ();
	// Make a window to display the ball
	Window BBWindow = new Window ();
		// Set the window's properties
		BBWindow.delay = 5;
		BBWindow.caption = "Wolter's SDL Bouncing Ball";
		// Display the window
		BBWindow.video_init ();
	// Make a ball
	Ball BBall = new Ball (BBWorld, BBWindow);
	// Assign the screen_update delegate
	BBWindow.screen_update = screen_update_definition;
	// Run the window
	BBWindow.run ();
}
