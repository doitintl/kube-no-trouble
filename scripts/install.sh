#!/usr/bin/env sh

# Set strict error checking
set -emou pipefail
LC_CTYPE=C

# Enable debug output if $DEBUG is set to true
[ "${DEBUG:="false"}" = "true" ] && set -x

# Optional vars
APP_NAME="${APP_NAME:="kubent"}"
GITHUB_REPO="${GITHUB_REPO:="doitintl/kube-no-trouble"}"
TARGET_DIR="${TARGET_DIR:="/usr/local/bin"}"
TARGET_ARCH="${TARGET_ARCH:="$(uname -m | sed 's/x86_64/amd64/')"}"
TARGET_OS="${TARGET_OS:="$(uname -s | tr '[:upper:]' '[:lower:]')"}"
REQUIRED_BINARIES=${REQUIRED_BINARIES:='tar curl'}


# Check if we have all the dependencies
check_binaries() {
  for bin in ${REQUIRED_BINARIES}; do
    [ -x "$(command -v "${bin}")" ] || fail "Required dependency ${bin} not found in path"
  done
}

# Check if we have colors available, it looks good
check_colors(){
  if command -v tput > /dev/null; then
    COLORS="$(tput colors)"
    if [ -n "${COLORS}" ] && [ "${COLORS}" -ge 8 ]; then
      GREEN="$(tput setaf 2)"
      RED="$(tput setaf 1)"
      YELLOW="$(tput setaf 3)"
      NOCOL="$(tput sgr0)"
    fi
  else
    GREEN=''
    RED=''
    YELLOW=''
    NOCOL=''
  fi
}

# Print error and exit
fail() {
  echo "${RED}ERROR: $*${NOCOL}"
  exit 1
}

usage()
{
cat << EOF

OPTIONS:
   -h      Show help.
   -d      Directory where kubent will be installed. Default is /usr/local/bin
   -a      Architecture to install (x86_64 only atm.). Default is to auto-detect.
   -o      OS (linux, macos). Default is to auto-detect.
EOF
}

get_latest_release() {
curl -sL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" \
  | grep -oE '"tag_name": "[^"]*' \
  | grep -oE '[^"]*$'
}

download_version() {
version="${1}"
curl -L -o- "https://github.com/${GITHUB_REPO}/releases/download/${version}/${APP_NAME}-${version}-${TARGET_OS}-${TARGET_ARCH}.tar.gz" \
  | tar -xz -C "${TARGET_DIR}"
}

check_colors
check_binaries
echo ">>> ${GREEN}${APP_NAME} installation script${NOCOL} <<<"

while getopts "hd:a:o:" OPTION
do
     case $OPTION in
         h) usage; exit;;
         d) TARGET_DIR="$OPTARG";;
         a) TARGET_ARCH="$OPTARG";;
         o) TARGET_OS="$OPTARG";;
         ?) usage; exit;;
     esac
done

echo "${YELLOW}>${NOCOL} Detecting latest version"
release="$(get_latest_release)"

echo "${YELLOW}>${NOCOL} Downloading version ${release}"
download_version "${release}"

echo "${YELLOW}>${NOCOL} Done. ${GREEN}${APP_NAME}${NOCOL} was installed to ${TARGET_DIR}/."
