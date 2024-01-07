#version 150 core

//vertex position
in vec2 vert;

// pass through to fragTexCoord
in vec2 vertTexCoord;

// window resolution
uniform vec2 u_resolution;

// pass to frag
out vec2 fragTexCoord;

void main() {
   fragTexCoord = vertTexCoord;

   vec2 clipSpace = (vert / u_resolution * 2.0) - 1.0;
   gl_Position = vec4(clipSpace * vec2(1, -1), 0, 1);
}