#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <version-tag> <homebrew-tap-directory>" >&2
  exit 1
fi

tag="$1"
tap_dir="$2"
version="${tag#v}"
archive_url="https://github.com/anton-fuji/dibo/archive/refs/tags/${tag}.tar.gz"
archive_file="$(mktemp)"
trap 'rm -f "$archive_file"' EXIT

curl --fail --location --silent --show-error --output "$archive_file" "$archive_url"
checksum="$(shasum -a 256 "$archive_file" | awk '{print $1}')"

mkdir -p "$tap_dir/Formula"
cat >"$tap_dir/Formula/dibo.rb" <<FORMULA
class Dibo < Formula
  desc "A CLI tool for generating .dockerignore files from templates"
  homepage "https://github.com/anton-fuji/dibo"
  url "${archive_url}"
  sha256 "${checksum}"
  license "MIT"

  depends_on "go" => :build

  def install
    ldflags = "-s -w -X github.com/anton-fuji/dibo/cmd.version=#{version}"
    system "go", "build", *std_go_args(ldflags: ldflags), "."
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/dibo --version")
  end
end
FORMULA
