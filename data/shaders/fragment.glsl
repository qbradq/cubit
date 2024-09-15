#version 100

precision mediump float;
precision mediump int;

uniform int drawMode;
uniform sampler2D tex;
uniform sampler2D atlas;
uniform vec3 lightPos;
uniform vec3 lightColor;
uniform vec3 cameraPos;

varying vec3 fragPos;
varying vec3 texCoord;
varying vec3 normal;

void main(void) {
    // Color sampling
    vec3 color;
    if (drawMode == 0) {
        color = texCoord;
    } else if (drawMode == 1) {
        color = texture2D(atlas, texCoord.xy).xyz;
    } else if (drawMode == 2) {
        color = texture2D(tex, texCoord.xy).xyz;
    } else {
        gl_FragColor = vec4(texCoord, 1.0);
        return;
    }
    // Ambient
    vec3 ambient = 0.1 * lightColor;
    // Diffuse
    vec3 norm = normalize(normal);
    vec3 lightDir = normalize(lightPos - fragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * lightColor;
    // Specular
    vec3 viewDir = normalize(cameraPos - fragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 2.0);
    vec3 specular = 0.5 * spec * lightColor;
    // Color mixing
    vec3 result = (ambient + diffuse + specular) * vec3(color.xyz);
    gl_FragColor = vec4(result, 1.0);
}
