class Janus < Formula
  desc "Shared CI deployer + version syntax tool for ETCDEV"
  homepage ""
  url "https://github.com/whilei/janus/releases/download/v0.1.8/janus_0.1.8_darwin_amd64.tar.gz"
  version "0.1.8"
  sha256 "ed0990064eb234d7468a1367e98d476f7f6f93040ce2e976d5a42dcc6d0bc434"

  def install
    bin.install "janus"
  end
end
