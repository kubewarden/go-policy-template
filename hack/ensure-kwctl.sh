#!/usr/bin/env bash
set -o errexit
set -o nounset
set -o pipefail

arch=x86_64
os="unknown"
bin_dir=bin

if [[ "$OSTYPE" == "linux"* ]]; then
  os="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  os="darwin"
fi

if [[ "$os" == "unknown" ]]; then
  echo "OS '$OSTYPE' not supported. Aborting." >&2
  exit 1
fi

if [ -z "${1}" ]; then
  echo "must provide kwctl version as first and only parameter"
  exit 1
fi

# Turn colors in this script off by setting the NO_COLOR variable in your
# environment to any value:
#
# $ NO_COLOR=1 test.sh
NO_COLOR=${NO_COLOR:-""}
if [ -z "$NO_COLOR" ]; then
  header=$'\e[1;33m'
  reset=$'\e[0m'
else
  header=''
  reset=''
fi

function header_text {
  echo "$header$*$reset"
}

header_text "downloading kwctl"

kwctl_version=${1}
kwctl_bin_name="kwctl-$os-$arch"
kwctl_archive_name="$kwctl_bin_name.zip"
kwctl_download_url=" https://github.com/kubewarden/kwctl/releases/download/$kwctl_version/$kwctl_archive_name"
kwctl_archive_path="$bin_dir/$kwctl_archive_name"

mkdir -p $bin_dir
if [ ! -f $kwctl_archive_path ]; then
  echo "download url: ${kwctl_download_url}"
  echo "archive path: $kwctl_archive_path"
  curl -fsL ${kwctl_download_url} -o "$kwctl_archive_path"
fi
header_text "kwctl downloaded"
tar -zvxf "$kwctl_archive_path" -C "$bin_dir/"
mv $bin_dir/$kwctl_bin_name $bin_dir/kwctl
rm $kwctl_archive_path
