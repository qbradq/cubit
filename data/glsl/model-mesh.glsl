[VERTEX]
#version 100

const float scale = 1.0/16.0;

uniform mat4 uModelMatrix;
uniform mat4 uViewMatrix;
uniform mat4 uProjectionMatrix;
uniform vec3 uRotationPoint;
uniform float uLightLevels[6];

attribute vec3 aVertexPosition;
attribute vec3 aVertexColor;
attribute float aVertexFacing;

varying vec3 color;
varying float lightLevel;

void main() {
	color = aVertexColor;
	lightLevel = uLightLevels[int(aVertexFacing)];
	vec3 vp = aVertexPosition-uRotationPoint;
	vp = vec3(vec4(vp, 1.0) * uModelMatrix);
	vp = vp*scale;
	gl_Position = uProjectionMatrix * uViewMatrix * vec4(vp, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

varying vec3 color;
varying float lightLevel;

void main() {
	gl_FragColor = vec4(color * lightLevel, 1.0);
}
