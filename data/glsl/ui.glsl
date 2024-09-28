[VERTEX]
#version 460

uniform mat4 uProjectionMatrix;
uniform vec3 uPosition;

in layout(location=0) vec2 aVertexPosition;
in layout(location=1) vec2 aVertexUV;

out vec2 uv;

void main() {
	uv = aVertexUV / 128.0;
    vec3 pos = vec3(aVertexPosition, 0.0)+uPosition;
	gl_Position = uProjectionMatrix * vec4(pos, 1.0);
}

[FRAGMENT]
#version 460

uniform layout(binding=0) sampler2D uAtlas;

in vec2 uv;

out vec4 color;

void main() {
    color = texture(uAtlas, uv);
}
