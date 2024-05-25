#!/bin/bash
#From https://github.com/oneclickvirt/cpuTest
#2024.05.25


rm -rf /tmp/geekbench
arch=$(uname -m)
release_date="20240525"
if [ ! -d "/tmp" ]; then
    mkdir /tmp
fi
if ! command -v wget >/dev/null 2>&1; then
    echo "The wget command is not detected, please download it before executing this script."
fi
if ! command -v tar >/dev/null 2>&1; then
    echo "The tar command is not detected, please download it before executing this script."
fi
if [ "$1" != "-v" ]; then
  echo "Error: the -v option must be used"
  exit 1
fi
if [ -z "$2" ]; then
  echo "Error: a value must be provided (gb4, gb5 or gb6)"
  exit 1
fi
case "$2" in
  gb4|gb5|gb6)
    gbv="$2"
    ;;
  *)
    echo "Error: Invalid value. Must be gb4, gb5 or gb6"
    exit 1
    ;;
esac
case $gbv in
  4)
    case $arch in
      "x86_64" | "x86" | "amd64" | "x64")
        wget -O /tmp/geekbench.tar.gz https://github.com/oneclickvirt/cpuTest/releases/download/${release_date}/Geekbench-4.4.4-Linux.tar.gz
        ;;
      *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
    esac
    ;;
  5)
    case $arch in
      "x86_64" | "x86" | "amd64" | "x64")
        wget -O /tmp/geekbench.tar.gz https://github.com/oneclickvirt/cpuTest/releases/download/${release_date}/Geekbench-5.5.1-Linux.tar.gz
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O /tmp/geekbench.tar.gz https://github.com/oneclickvirt/cpuTest/releases/download/${release_date}/Geekbench-5.5.1-LinuxARMPreview.tar.gz
        ;;
      *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
    esac
    ;;
  6)
    case $arch in
      "x86_64" | "x86" | "amd64" | "x64")
        wget -O /tmp/geekbench.tar.gz https://github.com/oneclickvirt/cpuTest/releases/download/${release_date}/Geekbench-6.3.0-Linux.tar.gz
        ;;
      "armv7l" | "armv8" | "armv8l" | "aarch64" | "arm64")
        wget -O /tmp/geekbench.tar.gz https://github.com/oneclickvirt/cpuTest/releases/download/${release_date}/Geekbench-6.3.0-LinuxARMPreview.tar.gz
        ;;
      *)
        echo "Unsupported architecture: $arch"
        exit 1
        ;;
    esac
    ;;
esac

chmod 777 /tmp/geekbench.tar.gz
tar -xvf /tmp/geekbench.tar.gz
if [ -f /tmp/geekbench ]; then 
  for file in /tmp/geekbench/geekbench[0-9]*; do
    if [ -f "$file" ]; then
      target="/tmp/geekbench/geekbench"
      mv "$file" "$target"
    fi
  done
fi
chmod 777 /tmp/geekbench/geekbench
/tmp/geekbench/geekbench --version
if [ $? -ne 0 ]; then
  echo "Geekbench failed to check the version, please leave an error message in the repository's issues."
fi