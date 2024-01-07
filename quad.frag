#version 330 core

out vec4 out_color;

uniform vec2 u_resolution;
uniform vec2 u_mouse;
uniform float u_time;
uniform int u_frame;

float easeInOutCubic(float x);
vec3 inferno(float t);

vec2 point[5] = vec2[5](
    vec2(0.83,0.75),
    vec2(0.60,0.07),
    vec2(0.28,0.64),
    vec2(0.31,0.26),
    vec2(0.5,0.5)
);

vec3 colors[5] = vec3[5](
    vec3(17.0/255.0, 115.0/255.0, 185.0/255.0),
    vec3(215.0/255.0, 85.0/255.0, 38.0/255.0),
    vec3(236.0/255.0, 176.0/255.0, 53.0/255.0),
    vec3(117.0/255.0, 52.0/255.0, 137.0/255.0),
    vec3(130.0/255.0, 170.0/255.0, 69.0/255.0)
);

void main()
{
    // dummy vars to make sure the uniforms are used and therefore not optimized away
    float dummy0 = u_resolution.x;
    float dummy1 = u_time;
    float dummy2 = u_mouse.x;
    float dummy3 = u_frame;

    vec2 st = gl_FragCoord.xy/u_resolution.xy;
    vec2 mouse = u_mouse/u_resolution;
    mouse.y = 1.0 - mouse.y; // flip y-axis

    // Basic voronoi example from:
    // https://thebookofshaders.com/12/

    vec3 color = vec3(.0);

    // Cell positions
    point[4] = mouse;

    float m_dist = 1.;  // minimum distance
    int m_point = 0;    // index of the closest point

    // Iterate through the points positions
    for (int i = 0; i < 5; i++) {
        // L1 norm
        // float dist = abs(st.x-point[i].x) + abs(st.y-point[i].y);
        // L2 norm
        float dist = distance(st, point[i]);
        // L infinite norm
        // float dist = max(abs(st.x-point[i].x),abs(st.y-point[i].y));

        // Keep the closer distance
        m_dist = min(m_dist, dist);
        if( m_dist == dist ){
            m_point = i;
        }
    }

    // pick a color based on the closest point. sample the inferno colormap
    color = colors[m_point];
    color *= 1.0 - m_dist*2.1;

    // Show isolines
    // color -= step(.7,abs(sin(200.0*m_dist)))*.3;

    out_color = vec4(color,1.0);
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