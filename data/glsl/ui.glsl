[VERTEX]
#version 100

uniform mat4 uProjectionMatrix;
uniform vec3 uPosition;

attribute vec2 aVertexPosition;
attribute vec2 aVertexUV;

varying vec2 uv;

void main() {
	uv = aVertexUV / 128.0;
    vec3 pos = vec3(aVertexPosition, 0.0)+uPosition;
	gl_Position = uProjectionMatrix * vec4(pos, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

uniform sampler2D uAtlas;

varying vec2 uv;

void main() {
    vec4 color = texture2D(uAtlas, uv);
    gl_FragColor = color;
}
