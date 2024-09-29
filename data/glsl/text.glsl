[VERTEX]
#version 100

uniform mat4 uProjectionMatrix;
uniform vec3 uPosition;

attribute vec2 aVertexPosition;
attribute vec2 aVertexUV;
attribute vec3 aVertexColor;

varying vec2 uv;
varying vec3 color;

void main() {
	uv = aVertexUV / vec2(64.0, 32.0);
    color = aVertexColor;
	gl_Position = uProjectionMatrix * vec4(vec3(aVertexPosition, 0.0)+
        uPosition, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

uniform sampler2D uFont;

varying vec2 uv;
varying vec3 color;

void main() {
    vec4 f = texture2D(uFont, uv);
    gl_FragColor = f*vec4(color, 1.0);
}
