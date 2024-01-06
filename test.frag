#version 330 core

uniform vec4 u_color;

in vec2 TexCoords;
out vec4 color;

uniform sampler2D text;

void main()
{    
    // dummy vars to make sure the uniforms are used and therefore not optimized away
    float dummy1 = u_color.x;

    vec3 textColor = vec3(1.0, 1.0, 1.0); // Fix to red for now

    vec4 sampled = vec4(1.0, 1.0, 1.0, texture(text, TexCoords).r);
    // color = vec4(textColor, 1.0) * sampled;
    // debug. output white
    color = vec4(textColor, 1.0);
} 