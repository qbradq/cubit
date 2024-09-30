[VERTEX]
#version 100

const float scale = 1.0/16.0;

uniform mat4 uModelViewMatrix;
uniform mat4 uProjectionMatrix;
uniform vec3 uRotationPoint;
uniform float uLightLevels[6*6];
uniform int uFacing;

attribute vec3 aVertexPosition;
attribute vec3 aVertexColor;
attribute float aVertexFacing;

varying vec3 color;
varying float lightLevel;

void main() {
	color = aVertexColor;
	lightLevel = uLightLevels[uFacing*6+int(aVertexFacing)];
	vec3 vp = aVertexPosition-uRotationPoint;
	gl_Position = uProjectionMatrix * uModelViewMatrix * vec4(vp*scale, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

varying vec3 color;
varying float lightLevel;

void main() {
	gl_FragColor = vec4(color * lightLevel, 1.0);
}
