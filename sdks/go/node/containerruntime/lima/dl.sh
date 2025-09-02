mkdir -p embeds/amd64 embeds/arm64

curl -L https://github.com/lima-vm/lima/releases/download/v1.2.1/lima-1.2.1-Darwin-arm64.tar.gz | tar -xzv ./bin/limactl ./share/lima/lima-guestagent.Linux-aarch64.gz
gunzip ./share/lima/lima-guestagent.Linux-aarch64.gz
mv ./bin/limactl embeds/arm64/limactl
mv ./share/lima/lima-guestagent.Linux-aarch64 embeds/arm64/lima-guestagent.Linux-aarch64

curl -L https://github.com/lima-vm/lima/releases/download/v1.2.1/lima-1.2.1-Darwin-x86_64.tar.gz | tar -xzv ./bin/limactl ./share/lima/lima-guestagent.Linux-x86_64.gz
gunzip ./share/lima/lima-guestagent.Linux-x86_64.gz
mv ./bin/limactl embeds/amd64/limactl
mv ./share/lima/lima-guestagent.Linux-x86_64 embeds/amd64/lima-guestagent.Linux-amd64
