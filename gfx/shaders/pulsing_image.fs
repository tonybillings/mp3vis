#version 410

in vec2 TexCoords;
out vec4 color;
uniform sampler2D backgroundTexture;
uniform float angle;

void main() {
    vec2 centeredCoords = TexCoords - vec2(0.5, 0.5);

    vec2 rotatedCoords;
    rotatedCoords.x = centeredCoords.x * cos(angle) - centeredCoords.y * sin(angle);
    rotatedCoords.y = centeredCoords.x * sin(angle) + centeredCoords.y * cos(angle);

    rotatedCoords += vec2(0.5, 0.5);

    color = texture(backgroundTexture, rotatedCoords);
    color.rgb *= 0.3;
    color.r += angle * 0.4;
}
