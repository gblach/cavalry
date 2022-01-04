# Cavalry

Unittest framework for containerized applications.

## Abstract

The Cavalry is a tool that builds one or more containers, then runs unittests in them.  
If all tests pass, it can push images to the container registry.  
At the end it removes the containers and optionally the images.

## Command line options

	cavalry [-c dir] [-e <podman|docker>] [-f <oci|docker>] [-np] [-nr] [-p] [Cavalryfile]

**-c directory**  
Change working directory.  
Cavalry will look for Cavalryfile in this directory.
Also this directory will be used as the default value for the DIR directive.

**-e &lt;podman | docker&gt;**  
Choose the engine: podman or docker.  
Cavalry checks if podman or docker is installed on the system and uses the command it found. This option is useful if you have both podman and docker installed.

**-f &lt;oci | docker&gt;**  
Choose the image format: oci or docker.  
This option is useful if you are using podman, because docker always uses the docker format.

**-m &lt;email address&gt;**  
Send an email to this address in case of failure.  
Emails are sent using /usr/sbin/sendmail or the command pointed by SENDMAIL_CMD env variable if defined.

**-h**	 
Show help message.

**-np**  
Do not push any images into container registries.
In other words, omit all PUSH directives.

**-nr**  
Do not remove any images and containers built by Cavalry and keep them running.

**-p**  
Show commands that Cavalry plans to execute instead of executing them.

**-ma**  
In conjunction with the -m option, this means to always send the e-mail not only on failure.

**-v**  
Show version and exit.

**Cavalryfile**  
Point the Cavalryfile to read. The default is "Cavalryfile".

## Cavalryfile

The Cavalryfile is line oriented. Each line contains the directive and parameters. A valid file should contain at least one CONTAINER and one EXEC directive.

	CONTAINER tag
This directive instructs Cavalry to build new image and run the container. This image will be tagged with tag parameter.

The following directives apply to the latest CONTAINER directive: DIR, FILE, PUSH, and KEEP.

	DIR directory
Points to the directory where the container's image is to be built. Defaults to the current working directory.

	FILE Dockerfile
Points the Dockerfile to be used during building the container's image. The default is "Dockerfile".

	PUSH registry_addr
Points the registry addess where image will be pushed after passing the tests. You should be logged in to the registry prior to running Cavalry.

	KEEP
Mark that the container's image cannot be removed during the cleanup step.

	EXEC tag command
Defines the test command to run. This command will be executed in the container with the image tagged by the tag parameter.
