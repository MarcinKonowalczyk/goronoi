#version 330 core

in vec4 vertexColor;
out vec4 fragmentColor;

void main()
{
    fragmentColor = vertexColor;
    // fragmentColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}