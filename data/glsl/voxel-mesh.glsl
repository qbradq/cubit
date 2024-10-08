[VERTEX]
#version 100

const float cScale = 1.0/16.0;
const vec3 cRotationPoint = vec3(8.0, 8.0, 8.0);

uniform mat4 uModelMatrix;
uniform mat4 uViewMatrix;
uniform mat4 uProjectionMatrix;
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
	vec3 vp = aVertexPosition-cRotationPoint;
	gl_Position = uProjectionMatrix * uViewMatrix * uModelMatrix *
		vec4(vp*cScale, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

varying vec3 color;
varying float lightLevel;

void main() {
	gl_FragColor = vec4(color * lightLevel, 1.0);
}
