[VERTEX]
#version 100

uniform mat4 uModelViewMatrix;
uniform mat4 uProjectionMatrix;

attribute vec3 aVertexPosition;
attribute vec3 aVertexColor;

varying vec3 color;

void main() {
	gl_PointSize = 32.0;
	color = aVertexColor;
	gl_Position = uProjectionMatrix * uModelViewMatrix *
		vec4(aVertexPosition, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

varying vec3 color;

void main() {
	gl_FragColor = vec4(color.rgb, 1.0);
}
