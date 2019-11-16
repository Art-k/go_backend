export GOPATH=/home/art-k/go/src/go_backend/src
cd $HOME/go/src/go_backend
git pull
go build
sudo -su
service go_backend stop
service go_backend start
