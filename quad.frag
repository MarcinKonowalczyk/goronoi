#version 330 core

out vec4 color;

uniform vec2 u_resolution;
uniform vec2 u_mouse;
uniform float u_time;

float easeInOutCubic(float x);

void main()
{
    // dummy variabels to make sure the uniforms are used and therefore not optimized away
    float dummy0 = u_resolution.x;
    float dummy1 = u_time;
    float dummy2 = u_mouse.x;

    vec2 frag = gl_FragCoord.xy/u_resolution.xy;
    vec2 mouse = u_mouse/u_resolution;
    mouse.y = 1.0 - mouse.y; // flip y-axis

    float red = frag.x;
    float green = frag.y;

    red = easeInOutCubic(red);
    green = easeInOutCubic(green);

    float dist = distance(frag, mouse);
    float radius = 0.05;
    radius *= 1.0 + sin(u_time * 20.0) * 0.2;
    if (dist < radius) {
        red = 0.0;
        green = 0.0;
    }

    color = vec4(red, green, 0.0, 1.0);
}

float easeInOutCubic(float x)
{
    if (x < 0.5) {
        return 4.0 * x * x * x;
    } else {
        float x1 = -2.0 * x + 2.0;
        float x2 = x1 * x1 * x1;
        return 1.0 - x2 / 2.0;
    }
}