Build docker container version of the app, that is runnable on kubernetes.

The go app:  jsonsvalidator needs to be copied into this directory so the Dockerfile will find it.

This examples assumes the app need  a kubernetes access certificate in this directory:  ca-certificates.crt
