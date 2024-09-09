#version 100

precision mediump float;
precision mediump int;

uniform int drawMode;
uniform sampler2D tex;
uniform sampler2D atlas;

varying vec3 texCoord;
varying vec3 normal;

void main(void) {
    if (drawMode == 0) {
        gl_FragColor = vec4(texCoord.xyz, 1.0);
    } else if (drawMode == 1) {
        gl_FragColor = texture2D(atlas, texCoord.xy);
    } else {
        gl_FragColor = texture2D(tex, texCoord.xy);
    }
}
