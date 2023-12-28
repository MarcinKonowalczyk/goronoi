#version 330 core

out vec4 color;

uniform vec2 u_resolution;
uniform vec2 u_mouse;
uniform float u_time;

void main()
{
    // dummy variabels to make sure the uniforms are used and therefore not optimized away
    float dummy0 = u_resolution.x;
    float dummy1 = u_time;
    float dummy2 = u_mouse.x;

    vec2 st = gl_FragCoord.xy/u_resolution.xy;
    st.x *= u_resolution.x/u_resolution.y;

    // vec2 mouse = u_mouse/u_resolution;
    // mouse.y = 0.5; // pretrend mouse is in the middle of the screen
    // // float dist = distance(st, mouse);
    // float dist = abs(gl_FragCoord.x - u_mouse.x);

    float red = st.x;
    float green = st.y;

    // if (dist < 0.1)
    // {
    //     red = 0.0;
    //     green = 0.0;
    // }
    

    color = vec4(red, green, 0.0, 1.0);
}