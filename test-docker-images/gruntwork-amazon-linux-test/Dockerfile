# TODO: Is it worth referencing a specific tag instead of latest?
FROM amazonlinux:2017.12

# Reduce Docker image size per https://blog.replicated.com/refactoring-a-dockerfile-for-image-size/
# - perl-Digest-SHA: installs shasum
RUN yum update -y && \
    yum upgrade -y && \
    yum install -y \
        hostname \
        jq \
        perl-Digest-SHA \
        rsyslog \
        sudo \
        tar \
        vim \
        wget && \
        yum clean all && rm -rf /var/cache/yum

# Installing pip with yum doesn't actually put it in the PATH, so we use easy_install instead. Pip will now be placed
# in /usr/local/bin, but amazonlinux's sudo uses a sanitzed PATH that does not include /usr/local/bin, so we symlink pip.
# The last line upgrades pip to the latest version.
RUN curl https://bootstrap.pypa.io/ez_setup.py | sudo /usr/bin/python && \
    easy_install pip && \
    pip install --upgrade pip

# Install the AWSCLI (which apparently does not come pre-bundled with Amazon Linux!)
RUN pip install awscli --upgrade

# Ideally, we'd install the latest version of Docker to avoid a conflict between the Docker client in this container
# and the Docker API on your local host, but installing the latest version of Docker yields the error "Requires:
# container-selinux >= 2.9", which indicates that a newer Linux kernel version is required than what comes with Amazon Linux.
# So we settle for the Amazon Linux supported version for now.
RUN yum install -y docker