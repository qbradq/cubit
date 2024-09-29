[VERTEX]
#version 100

uniform mat4 uModelViewMatrix;
uniform mat4 uProjectionMatrix;
uniform mat4 uNormalMatrix;

attribute vec3 aVertexPosition;
attribute vec2 aVertexUV;
attribute float aVertexLightLevel;

varying float lightLevel;
varying vec2 uv;

void main() {
	uv = aVertexUV / 128.0;
    lightLevel = aVertexLightLevel;
	gl_Position = uProjectionMatrix * uModelViewMatrix *
        vec4(aVertexPosition, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

uniform sampler2D uAtlas;

varying float lightLevel;
varying vec2 uv;

void main() {
    vec3 color = texture2D(uAtlas, uv).rgb;
    gl_FragColor = vec4(color * lightLevel, 1.0);
}
