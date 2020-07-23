# Overview

This repository is only dummy app for learning purpose.
The things that I want to learn here are [casbin](https://casbin.org/) and [scs](https://github.com/alexedwards/scs).

# Casbin 
Casbin is an authorization library that support access controll. With this dummy app, I try to implement it for RBAC (Role-Based Access Control).
The scenario here is to limit access for a certain role into a certain HTTP endpoint.
Using Casbin seems pretty straight forward for simple case, and it has a lot of feauture to explore as well.
It support policy with a `.csv` file and also has adapter for several common databases which will come in handy if you manage a lot of roles.
In case of few roles, I think file based policy still float my boat.
I like this library so far.

# scs
scs is HTTP Session Management for Go. In this dummy app I just barely use it for sample.
But it has few cool features as well since it is considered as easy to extend from their README.md

# Notes
The dummy app is for learning purpose only and not a production ready code. 
This code only for learning the minimum stuff to use Casbin and scs library.
