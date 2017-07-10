class Janus < Formula
  desc "Shared CI deployer + version syntax tool for ETCDEV"
  homepage ""
  url "https://github.com/whilei/janus/releases/download/v0.1.6/janus_0.1.6_darwin_amd64.tar.gz"
  version "0.1.6"
  sha256 "b0345168b4cce549e9665605dde1740184fa2fa7777e5862a786ff92e31c41fe"

  def install
    bin.install "janus"
  end
end
