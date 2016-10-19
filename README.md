# CPVT-CloudNativeLab

## Lab Objective
The objective of this lab is to introduce some common tools, methodologies and design considerations when you start to build out a Cloud Native application. It will also allow you to become famailiar with Docker, interacting with the services and an introduction to the Golang language.

## What are the lab requirements
There are only two requirements for this lab, 1) is to have Docker installed and running. If configured correctly, you should be able to run native Docker commands at the terminal or command window and not need to use Sudo or Administrator for it to succeed. 2) is to have an active Internet connection to enable you to download some images from Docker Hub.

## What components are we going to use
There are a large number of different tools we could use in this lab, for example Server Load Balancers and each one of them has many different flavours both Open Source and Commercial offerings. This lab has taken some of the most popular Open Source offerings at the time of writing to focus in on. You can however swap individual components out for alternatives, however this lab will not cover the integration points for other tools.

We will use;
+  Docker
+  Consul
+  Fabio
+  Redis
+  Golang
+  RabbitMQ
+  Others

## What is the lab outcome?
At the end of this lab, you will have hopefully gained a good understanding of Docker and Docker parameters, how to setup supporting tools such as Consul and how automatic and manual registeration occurs, how to create and compile GO code and insert into a container.
If you would like to see a graphical output of what you are going to produce, you can see a quick spoiler at the end of step 6 in this document.

![alt text](https://github.com/robjporter/CPVT-CloudNativeLab/raw/master/images/final.png "Output display 1")

## Any assumptions about the lab?
The lab has been designed that if you have no or little experience of these components, you will be able to follow through step­by­step and get to the same output as someone who has good knowledge of all of these tools.

**This lab has been written on and for a MAC. Other *nix platforms should be able to follow the commands exactly, however Windows users will need to amend the paths used.**

## The layout of the lab PDF looks like the following;
Title | Page
------|-----
Overview | 1
Step 1 ­ Docker Setup | 6
Step 2 ­ Consul Setup | 11
Step 3 ­ Fabio Setup | 16
Step 4 ­ Redis Setup | 20
Step 5 ­ Golang Setup | 23
Step 6 ­ Bringing it all together | 28
Step 7 ­ Where to now? | 37
Step 8 ­ Other things to try | 42
