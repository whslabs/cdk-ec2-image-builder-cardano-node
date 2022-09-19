#! /usr/bin/env nix-shell
#! nix-shell --pure -i sh

set -ex

cdk bootstrap

cdk synth

cdk deploy
