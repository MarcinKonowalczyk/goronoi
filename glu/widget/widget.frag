#version 330 core

out vec4 color;

uniform vec2 u_resolution;
uniform vec4 u_color;
uniform vec2 u_mouse;
uniform bool u_mouse_down;
uniform bool u_mouse_over;

void main()
{
	vec2 st = gl_FragCoord.xy/u_resolution.xy;
	vec2 mouse = u_mouse/u_resolution;
    mouse.y = 1.0 - mouse.y; // flip y-axis

	float mouse_dist = distance(st, mouse);
	if (mouse_dist < 0.1) {
		if (u_mouse_down) {
			// Yellow
			color = vec4(1.0, 1.0, 0.0, 0.5);
		} else {
			// Red
			color = vec4(1.0, 0.0, 0.0, 0.5);
		}
	} else {
		if (u_mouse_over) {
				color = vec4(u_color.rgb, 1.0);
			} else {
				color = vec4(u_color.rgb, 0.5);
		}
	}

}
