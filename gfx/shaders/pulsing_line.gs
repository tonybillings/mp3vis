#version 410 core

layout(lines) in;
layout(line_strip, max_vertices = 200) out;

uniform float step;
uniform float thickness;

void main() {
    float thicknessNorm = step * thickness;

    for (float offset = 0.0; offset <= thicknessNorm; offset += step) {
        for (int i = 0; i < gl_in.length(); i++) {
            gl_Position = gl_in[i].gl_Position + vec4(0.0, offset, 0.0, 0.0);
            EmitVertex();
        }

        EndPrimitive();
    }
}
