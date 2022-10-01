#! /bin/sh

set -ex

dockerd &

cdk bootstrap

cdk synth

cdk deploy --require-approval=never
