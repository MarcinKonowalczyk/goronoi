#version 330 core

out vec4 out_color;

uniform vec4 u_color;
uniform sampler2D u_texture;

void main()
{
    // dummy vars to make sure the uniforms are used and therefore not optimized away
    float dummy1 = u_color.x;

    vec4 text_color = vec4(1.0, 0.0, 0.0, 1.0);

    vec4 sampled = vec4(1.0, 0.0, 0.0, texture(u_texture, vec2(0.0, 0.0)).r);
    out_color = sampled * text_color;
}