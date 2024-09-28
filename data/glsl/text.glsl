[VERTEX]
#version 460

uniform mat4 uProjectionMatrix;
uniform vec3 uPosition;

in vec2 aVertexPosition;
in vec2 aVertexUV;

out vec2 uv;

void main() {
	uv = aVertexUV / vec2(64.0, 32.0);
	gl_Position = uProjectionMatrix * vec4(vec3(aVertexPosition, 0.0)+
        uPosition, 1.0);
}

[FRAGMENT]
#version 460

precision mediump float;

uniform sampler2D uFont;

in vec2 uv;

out vec4 color;

void main() {
    color = texture(uFont, uv);
}
