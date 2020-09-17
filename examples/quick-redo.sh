set -ex

pushd ../
make install
popd
terraform init
terraform plan
