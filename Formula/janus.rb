class Janus < Formula
  desc "Shared CI deployer + version syntax tool for ETCDEV"
  homepage ""
  url "https://github.com/whilei/janus/releases/download/v0.1.0/janus_v0.1.0_Darwin_x86_64.tar.gz"
  version "0.1.0"
  sha256 "23e46f30666571ca8f52d6663975ec76c55637ed42d712804d500b8240929fc9"

  def install
    bin.install "janus"
  end
end
