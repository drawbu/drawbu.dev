---
title: How to install NixOS on a Contabo Server
date: 2023-08-02
uri: install-nixos-contabo-vps
author:
  name: drawbu
  email: contact@drawbu.dev
description: Step by step installation of NixOS on a contabo server. This is mainly a guide for myself
---

In this guide, we'll see how to setup NixOS on a Ubuntu 22.04 Contabo VPS.
For this, we'll use a script known as [NixOS-Infect](https://github.com/elitak/nixos-infect) that is gonna replace the installed OS with NixOS.

## Rent a server
Should not be a very technical step, just go on [contabo.com](https://contabo.com) and purchase a VPS. I am just using the cheaper version, so it should works on anything.

Wait for the installation to be complete, it could take some times (from 5 seconds to 3 hours, really random).

## Connect to the server
To connect to the server, just run:
```bash
ssh root@vps_ip_address
```
Then, enter the password you set when you purchased it, or the one you received by mail. If none of them works, don't worry you can set it again on Contabo's website.

## Setup SSH Key
The root user will not have a password when nixos-infect runs to completion. To enable root login, you **must** have an SSH key configured.

To do that, just follow I made: [here](Connect%20to%20a%20server%20over%20SSH.md)

## Change the hostname
Contabo sets the hostname to something like `vmi######.contaboserver.net`, NixOS only allows RFC 1035 compliant hostnames ([see here](https://search.nixos.org/options?show=networking.hostName&query=hostname)). Run `hostname something_without_dots` before running the script. If you run the script before changing the hostname - remove the `/etc/nixos/configuration.nix` so it's regenerated with the new hostname.

Run this set the new hostname to `nixos`, or anything you want (just, no `.`).
```bash
hostname nixos
```

## Run the script
Connect to the server, and just run the script, and wait for completion. You can always inspect it before hand, but keep in mind that it is supposed to run on an almost empty VPS, just purchased.
```bash
curl https://raw.githubusercontent.com/elitak/nixos-infect/master/nixos-infect | NIX_CHANNEL=nixos-22.11 bash -x
```

When it's over, it will reboot and eject you. Just wait a little, reconnect to the server, and voilà, you're on NixOS baby!
