{
  inputs.lambda.url = "github:whslabs/rust-lambda-cloudtrail";
  outputs = { lambda, ... }: {
    defaultPackage = lambda.defaultPackage;
  };
}
