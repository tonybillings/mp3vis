#version 410

uniform float scale;
in vec2 position;

void main() {
    gl_Position = vec4(position * scale, 0.0, 1.0);
}
