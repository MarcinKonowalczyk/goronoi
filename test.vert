#version 330 core

in vec4 vertexData;

out vec4 vertexColor;

void main()
{
    gl_Position = vec4(vertexData.x, vertexData.y, 0.0f, 1.0f);
    vertexColor = vec4(vertexData.z, vertexData.w, 0.5f, 0.5f);
}