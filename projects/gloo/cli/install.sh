#! /bin/sh

set -eu


GLOO_VERSION=$(curl -H"Accept: application/vnd.github.v3+json" https://api.github.com/repos/solo-io/gloo/releases/latest| python -c "import sys, json; print(json.load(sys.stdin)['tag_name'])" )

if [ "$(uname -s)" = "Darwin" ]; then
  OS=darwin
else
  OS=linux
fi

tmp=$(mktemp -d /tmp/gloo.XXXXXX)
filename="glooctl-${OS}-amd64"
url="https://github.com/solo-io/supergloo/releases/download/${GLOO_VERSION}/${filename}"
(
  cd "$tmp"

  echo "Downloading ${filename}..."

  SHA=$(curl -sL "${url}.sha256")
  curl -LO "${url}"
  echo ""
  echo "Download complete!, validating checksum..."
  checksum=$(openssl dgst -sha256 "${filename}" | awk '{ print $2 }')
  if [ "$checksum" != "$SHA" ]; then
    echo "Checksum validation failed." >&2
    exit 1
  fi
  echo "Checksum valid."
  echo ""
)

(
  cd "$HOME"
  mkdir -p ".gloo/bin"
  mv "${tmp}/${filename}" ".gloo/bin/gloo"
  chmod +x ".gloo/bin/gloo"
)

rm -r "$tmp"

echo "Gloo was successfully installed 🎉"
echo ""
echo "Add the gloo CLI to your path with:"
echo ""
echo "  export PATH=\$PATH:\$HOME/.gloo/bin"
echo ""
echo "Now run:"
echo ""
echo "  glooctl install kube        # install gloo into the 'gloo-system' namespace"
echo ""
echo "Looking for more? Visit https://glooe.solo.io/installation/"
echo ""