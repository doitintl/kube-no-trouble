#!/usr/bin/env sh

# Set strict error checking
set -euf
LC_CTYPE=C

# Enable debug output if $DEBUG is set to true
[ "${DEBUG:="false"}" = "true" ] && set -x

# Optional vars
APP_NAME="${APP_NAME:="kubent"}"
GITHUB_REPO="${GITHUB_REPO:="LeMyst/kube-no-trouble"}"
TARGET_DIR="${TARGET_DIR:="/usr/local/bin"}"
TARGET_ARCH="${TARGET_ARCH:="$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')"}"
TARGET_OS="${TARGET_OS:="$(uname -s | tr '[:upper:]' '[:lower:]')"}"
TARGET_NIGHTLY="${TARGET_NIGHTLY:=false}"
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
   -a      Architecture to install (amd64, arm64). Default is to auto-detect.
   -o      OS (linux, macos, windows [amd64 only]). Default is to auto-detect.
   -n      Install latest nightly. Default is to install the latest stable release.
EOF
}

get_latest_release() {
curl -sL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" \
  | grep -oE '"tag_name": "[^"]*' \
  | grep -oE '[^"]*$'
}

get_latest_nightly_release() {
  curl -sL "https://api.github.com/repos/${GITHUB_REPO}/releases"  \
  | grep -oE '"tag_name": "[^"]*' \
  | grep -oE '[^"]*$' \
  | head -1
}

download_version() {
  version="${1}"
  sudo=""

  if [ ! -w "${TARGET_DIR}" ]; then
    if [ -x "$(command -v 'sudo')" ]; then
      echo "Target directory (${TARGET_DIR}) is not writable, trying to use sudo"
      sudo="sudo"
    fi
    ${sudo} [ -w "${TARGET_DIR}" ] \
      || fail "Target directory (${TARGET_DIR}) is not writable (destination can be changed using \$TARGET_DIR variable)"
  fi

  curl -L -o- "https://github.com/${GITHUB_REPO}/releases/download/${version}/${APP_NAME}-${version}-${TARGET_OS}-${TARGET_ARCH}.tar.gz" \
    | ${sudo} tar -xz -C "${TARGET_DIR}"
}

check_colors
check_binaries
echo ">>> ${GREEN}${APP_NAME} installation script${NOCOL} <<<"

while getopts "hd:a:o:n" OPTION
do
     case $OPTION in
         h) usage; exit;;
         d) TARGET_DIR="$OPTARG";;
         a) TARGET_ARCH="$OPTARG";;
         o) TARGET_OS="$OPTARG";;
         n) TARGET_NIGHTLY=true ;;
         ?) usage; exit;;
     esac
done

echo "${YELLOW}>${NOCOL} Detecting latest version"
release=""
if [ "${TARGET_NIGHTLY}" = true ]; then
  release="$(get_latest_nightly_release)"
else
  release="$(get_latest_release)"
fi

echo "${YELLOW}>${NOCOL} Downloading version ${release}"
download_version "${release}"

echo "${YELLOW}>${NOCOL} Done. ${GREEN}${APP_NAME}${NOCOL} was installed to ${TARGET_DIR}/."
