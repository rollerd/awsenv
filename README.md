# awsenv

awsenv is a small tool for quickly switching between aws credentials and optionally copying credentials to the clipboard as 'export commands' for use in other terminal windows/ssh sessions.

### Purpose
If you use aws-cli, boto3, or similar set of tools, you probably already know that you can configure your credentials (and config) in `.ini` style config files that live in your home directory.
Each section in these config files contain the different 'profiles' that you can use to authenticate to AWS or assume various roles (see [AWS named profiles](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-profiles.html) for more info).

Normally you would need to set the **AWS_PROFILE** environment variable to switch between these profiles:

```sh
export AWS_PROFILE=prod
```

and the profile env var would not persist across different local shell sessions.
`awsenv` helps with this issue by always using the `default` profile and swapping the underlying credentials values.

### Usage
Step 1: Add the following to your `~/.bashrc` (or equivalent):

```
export AWS_PROFILE=default
```

Step 2: If you have an existing `[default]` section in your `~/.aws/credentials` file, you should rename it to something else.

Step 3: Copy the awsenv binary to your /usr/local/bin directory and make it executable:

```
tar -xf darwin_amd64.tar
mv awsenv /usr/local/bin
chmod 755 /usr/local/bin/awsenv
```

Step 4: Start a new shell and try it out by running:

```
awsenv <existing credential section>
```

You should be able to view your `~/.aws/credentials` file and see that there is a new `[default]` section filled in with the values from the `<existing credential section>` you specified in the command.

The script also creates a file called `.awsenv` in your home directory that is populated with the name of the new active aws_profile. This can be useful for something like a tmux status bar script that can poll the file and display the value.

#### Copying Credentials
Run:

```
awsenv source <existing credential section>
```

The credentials will be copied to your clipboard as commands to set environment variables. eg "export AWS_ACCESS_KEY_ID=<some access key>"

#### Building/Modifying
A docker image built for crosscompiling can be used to build the binary (default macos below for CGO): https://github.com/elastic/golang-crossbuild

```
docker run  -it -e CGO_ENABLED=1 -v `pwd`:/app -w /app docker.elastic.co/beats-dev/golang-crossbuild:1.21.12-darwin-arm64-debian10 --build-cmd "go mod init github.com/rollerd/awsenv;go mod tidy;go build" -p "darwin/arm64"
```

