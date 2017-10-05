# kafka-management-tool

## Introduction
The kafka-management-tool is a command line `Go` tool that has been created for the sole purpose of being able to republish kafka messages in a given topic. Additionally it can also print the deserialised JSON into the terminal window for the user to read.

## Prerequisites 
To use this tool you will need to:
- Git clone this repo onto a local environment
- Have Go installed

## In-depth
As described in the introduction, this tool will print and republish specific messages into the specified Kafka topics. Given that the mandatory flags are entered into the tool, the tool will go through a number of steps in order to print and republish the chosen message. Below is the basic flow of the tool.

![alt text](https://user-images.githubusercontent.com/29541485/31216979-5a986170-a9ad-11e7-8117-795084401a56.png)
