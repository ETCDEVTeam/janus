class Janus < Formula
  desc "Shared CI deployer + version syntax tool for ETCDEV"
  homepage ""
  url "https://github.com/ethereumproject/janus/releases/download/v0.1.2/janus_v0.1.2_Darwin_x86_64.tar.gz"
  version "0.1.2"
  sha256 "a3a8b14e39462350eafd935430d7a4bd28c66ebede053ae137b35652485f1d73"

  def install
    bin.install "janus"
  end
end
