#version 100

uniform mat4 world;
uniform mat4 camera;
uniform mat4 project;

attribute vec3 pos;
attribute vec3 texPos;
attribute vec3 norm;

varying vec3 texCoord;
varying vec3 fragPos;
varying vec3 normal;

void main(void) {
	gl_Position = project * camera * world * vec4(pos, 1.0);
	texCoord = texPos;
	fragPos = vec3(world * vec4(pos, 1.0));
	normal = norm;
}
