#version 100

attribute vec3 pos;
attribute vec3 texPos;
attribute vec3 norm;

uniform mat4 world;
uniform mat4 camera;
uniform mat4 project;

varying vec3 texCoord;
varying vec3 normal;

void main(void) {
	gl_Position = project * camera * world * vec4(pos, 1.0);
	texCoord = texPos;
	normal = norm;
}
