# Create docker image for multiarch

# Usage: ./hack/pushMultiArch.sh

docker buildx create --name amdAndArm --driver=docker-container default

docker buildx build --builder=amdAndArm \
  --platform=linux/amd64,linux/arm64 \
  --push \
  -t yvyom/github-gist-backend:v1.0-alpha .
