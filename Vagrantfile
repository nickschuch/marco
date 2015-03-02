# Backup
#
# Provides a local development environment for Marco.
#

Vagrant.configure("2") do |config|
  config.vm.box       = "ubuntu/trusty64"
  config.vm.host_name = "marco.dev"
  
  config.vm.network :private_network, :ip => "192.168.80.10"
  
  config.vm.synced_folder ".", "/opt/golang/src/github.com/nickschuch/marco"
  
  config.vm.provider :virtualbox do |vb|
    vb.customize ["modifyvm", :id, "--memory", 1024]
  end

  config.vm.provision "shell", path: "scripts/provision.sh"
end
