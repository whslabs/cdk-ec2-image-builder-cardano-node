#! /bin/sh

set -ex

cdk bootstrap

cdk synth

cdk deploy --require-approval=never
