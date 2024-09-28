[VERTEX]
#version 460

const float scale = 1.0 / 16.0;

uniform layout(location=0) mat4 uModelViewMatrix;
uniform layout(location=1) mat4 uProjectionMatrix;

in layout(location=0) ivec3 aVertexPosition;
in layout(location=1) ivec3 aXYZ;

out vec3 xyz;

void main() {
	gl_Position = uProjectionMatrix * uModelViewMatrix *
		vec4(vec3(aVertexPosition)*scale, 1.0);
	xyz = vec3(aXYZ);
}

[FRAGMENT]
#version 460

uniform layout(binding=0) sampler2DArray uVolume;

in vec3 xyz;

out vec4 fragColor;

void main() {
	vec4 c = texelFetch(uVolume, ivec3(floor(xyz.x), floor(xyz.y),
		floor(xyz.z)), 0);
	if(c.a < 1.0) {
		discard;
	} else {
		fragColor = c;
	}
}
