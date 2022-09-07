# Netcupcloud
## Disclaimer
This project is an unofficial go module, changes to APIs can destroy functionalities to this module / change  the behavior of it. Use this module for education purposes only, I am not responsible for possible term and conditions violations or unintentionally caused costs. This project is not affiliated with the company netcup GmbH. Netcup is a registered trademark of netcup GmbH, Karlsruhe,Germany.

## Description
A go module, providing a cloud-like API for netcup. 
https://www.netcup.eu/

## Problem
Netcup currently does not support the management of cloud resources via an API. 
This includes f.e. the automatic provisioning of VPS.


## Test the package:
1. Use the local.env.example to create own local.env file with your credentials
2. Run ```go test``` to test the package.


## TODOs
Ideas:
hash the whole html pages, to check if the web pages have changed, make it as a warning during tests