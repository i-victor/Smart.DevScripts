
using GLib;
using SDL;
using SDLGraphics;

public class Ball : GLib.Object {

	// ATTRIBUTES

	// Vital attributes (Var. | Prop.)
	public float mass       { get; set; }
	public float elast      { get; set; }
	// Position:
	public float posx { get; set; }
	public float posy { get; set; }
	// Velocity:
	public float velx { get; set; }
	public float vely { get; set; }
		private float vely_term;
	// Acceleration:
	public float accelx { get; set; }
	public float accely { get; set; }
	// Appearance
	public uint radius      { get; set; }
	public uint32 color     { get; set; }
	// World attributes
	unowned World world;
		// Rendered values
		private float gravity_rend;
		private float air_res_rend;
		private float wind_rend;
		private float friction_rend;
	// Window data
	unowned Window window;

	// METHODS

	// Re-render values (temporary)
	private void recalc_values () {

		// Vt = sqrt ([2 × m × g] ÷ [density × surface × drag])
		this.vely_term = (float) Math.sqrt ( (2*this.mass*world.gravity) / (1.2*(2*(float) this.radius/100)*0.47) );
		stdout.printf ("Terminal velocity: %f\n", this.vely_term);
	}

	// Calculate velocities based on vectors
	private void cvel () {
		// Add acceleration
		this.velx += this.accelx;
		this.vely += this.accely;
		// Perform other worldly calculations
			// Gravity
			this.vely += this.gravity_rend;
			// Air resistance
			if (this.velx > 0) { this.velx -= this.air_res_rend; }
			else { this.velx += this.air_res_rend; }
		// Limit (terminal velocity)
		if (this.vely > this.vely_term) { this.vely = this.vely_term; }
	}

	// Bounce in X
	private void bouncex () {
		this.velx *= -1;
		// Use elasticity
		stdout.printf ("Ball bounced in X\n");
		if (this.velx > 0) {
			this.velx -= this.elast;
			this.velx -= this.friction_rend;

		}
		else {
			this.velx += this.elast;
		}
	}

	// Bounce in Y
	private void bouncey () {
		this.vely *= -1;
		// Use elasticity
		// NOTE: Its obvius that the vely is negative when the ball bounces off
		this.vely += this.elast;

		// Use friction (in X)
		if (this.velx > 0) {
			this.velx -= this.friction_rend;
		}
		else {
			this.velx += this.friction_rend;
		}
	}

	// Make the ball collide (and bounce)
	public void collision_detection () {

		// Screen level
		// For Y
		if ( (this.posy+this.radius) >= (float) window.screen_height) {
			this.posy = window.screen_height - this.radius;
			this.bouncey ();
		}
		// For X
		if (this.posx-this.radius <= 0) {
			this.posx = this.radius;
			bouncex ();
		}
		else if ((this.posx+this.radius) >= window.screen_width) {
			this.posx = window.screen_width-this.radius;
			bouncex ();
		}
	}

	// Gluepoint of all motion functions
	private void move () {
		// Detect collisions
		this.collision_detection ();
		// Calculate velocities
		this.cvel ();
		// Move in X
		this.posx += this.velx;
		// Move in Y
		this.posy += this.vely;
	}

	// Draws the ball on the screen
	private void draw () {

		// Draw as a circle on the screen
		Circle.fill_color (window.screen,
		  (int16) this.posx,
		  (int16) this.posy,
		  (int16) this.radius,
		  this.color);
	}

	// Gluepoint of previous 2 functions
	public void represent () {
		this.move ();
		this.draw ();
	}

	// CONSTRUCTOR

	public Ball (World world, Window window) {

		// Select world
		this.world = world;
		// Select window
		this.window = window;
		// Assign attributes
		this.mass = 1f;
		this.elast = 0.5f;
		this.radius = 5;
		this.color = 0xAFFF00FFU;
		// Define vectors
		this.posx = (float) (window.screen_width-this.radius) / 2;
		this.posy = (float) 10;
		//this.posy = (float) (window.screen_height-this.radius) / 2;
		this.velx = 0f;
		this.vely = 0f;
		this.accelx = 0f;
		this.accely = 0f;
		// Calculate variables
		this.recalc_values ();
		this.gravity_rend = (world.gravity*(float) window.delay/200);
		this.air_res_rend = (world.air_res*(float) window.delay/200);
		this.friction_rend = (world.friction*(float) window.delay/200);
		// Assign self to window
		for (uint iter = 0; iter <= window.balls.length; iter++) {
			if (window.balls[iter] == null) {
				window.balls[iter] = this;
				break;
			}
		}
		// Testing options
		this.radius = 5;
		this.velx = 2;
	}
}
