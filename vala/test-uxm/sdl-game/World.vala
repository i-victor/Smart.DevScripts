
using GLib;

public class World : GLib.Object {

	public float gravity    { get; set; }
	public float air_res    { get; set; }
	public float wind               { get; set; }
	public float friction   { get; set; }

	public World () {
		// Set default values
		this.gravity = 9.8f;
		this.air_res = 0.05f;
		this.wind = 0.2f;
		this.friction = 0.03f;
	}
}

