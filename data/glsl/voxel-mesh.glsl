[VERTEX]
#version 100

uniform mat4 uModelViewMatrix;
uniform mat4 uProjectionMatrix;
uniform mat4 uNormalMatrix;

attribute vec3 aVertexPosition;
attribute vec3 aVertexNormal;
attribute vec3 aVertexColor;

varying vec3 normal;
varying vec3 color;

void main() {
	color = aVertexColor;
	vec3 norm = aVertexNormal - vec3(1.0, 1.0, 1.0);
	normal = (uNormalMatrix * vec4(norm, 1.0)).xyz;
	gl_Position = uProjectionMatrix * uModelViewMatrix *
		vec4(aVertexPosition, 1.0);
}

[FRAGMENT]
#version 100

precision mediump float;

varying vec3 normal;
varying vec3 color;

const vec3 ambientLightColor = vec3(1.0, 1.0, 1.0);
const float ambientIntensity = 0.5;
const vec3 diffuseLightColor = vec3(1.0, 1.0, 1.0);
const vec3 lightDirection = normalize(vec3(1.0, 1.0, 1.0));

void main() {
	vec3 ambient = ambientLightColor * ambientIntensity;
	float diffuseIntensity = max(dot(normal, lightDirection), 0.0);
	vec3 diffuse = diffuseLightColor * diffuseIntensity;
	gl_FragColor = vec4((ambient + diffuse) * color, 1.0);
}
