#! /bin/sh

cdk bootstrap

cdk synth

cdk deploy --require-approval=never
