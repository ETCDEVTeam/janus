class Janus < Formula
  desc "Shared CI deployer + version syntax tool for ETCDEV"
  homepage ""
  url "https://github.com/whilei/janus/releases/download/v0.1.1/janus_v0.1.1_Darwin_x86_64.tar.gz"
  version "0.1.1"
  sha256 "d7ea50ba89ce46ec2bd40c39f32a3f4b2d2989d9bc2b71c7e1f944aed5871e54"

  def install
    bin.install "janus"
  end
end
