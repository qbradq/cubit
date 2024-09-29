[VERTEX]
#version 100

const float scale = 1.0/16.0;

uniform mat4 uModelViewMatrix;
uniform mat4 uProjectionMatrix;

attribute vec3 aVertexPosition;
attribute vec3 aVertexColor;
attribute float aVertexLightLevel;

varying vec3 color;
varying float lightLevel;

void main() {
	color = aVertexColor;
	lightLevel = aVertexLightLevel;
	gl_Position = uProjectionMatrix * uModelViewMatrix *
		vec4(aVertexPosition*scale, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

varying vec3 color;
varying float lightLevel;

void main() {
	gl_FragColor = vec4(color * lightLevel, 1.0);
}
