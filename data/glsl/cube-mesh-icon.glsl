[VERTEX]
#version 100

uniform mat4 uModelMatrix;
uniform mat4 uProjectionMatrix;

attribute vec3 aVertexPosition;
attribute vec2 aAtlasXY;
attribute vec2 aVertexUV;
attribute float aVertexLightLevel;

varying float lightLevel;
varying vec2 uv;
varying vec2 atlasXY;

void main() {
    uv = aVertexUV;
    atlasXY = aAtlasXY;
    lightLevel = aVertexLightLevel;
	gl_Position = uProjectionMatrix * uModelMatrix * vec4(aVertexPosition, 1.0);
    gl_Position[2] = aVertexPosition[2];
}

[FRAGMENT]
#version 100

precision mediump float;

const float cAtlasScale = 1.0 / 128.0;

uniform sampler2D uAtlas;

varying float lightLevel;
varying vec2 uv;
varying vec2 atlasXY;

void main() {
    vec2 auv = atlasXY;
    auv.x += fract(uv.x);
    auv.y += fract(uv.y);
    auv *= cAtlasScale;
    vec3 color = texture2D(uAtlas, auv).rgb;
    gl_FragColor = vec4(color * lightLevel, 1.0);
}
