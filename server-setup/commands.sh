#!/bin/bash

# Create git user
useradd --home-dir /home/git --create-home --shell /bin/bash git

# Login as git user
su - git

# Install gitolite
git clone https://github.com/sitaramc/gitolite

mkdir -p $HOME/bin

gitolite/install -to $HOME/bin

# Add gitolite to PATH

echo "export PATH=$HOME/bin:$PATH" >> $HOME/.bashrc

source $HOME/.bashrc

# Initialize gitolite

gitolite setup -pk admin.pub

# Exit git user
exit

exit

# Clone gitolite-admin repository 
git clone git@localhost:gitolite-admin.git
