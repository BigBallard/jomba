# Jomba

JSON Object Mapping By Abstraction

<!-- TOC -->
[Introduction](#introduction)
<!-- /TOC -->

## Introduction
Jomba is a tool that will take a valid JSON formatted file and present the abstract JSON structure. The abstraction shows a single composite instance of each property, object, and array with a count for the number of instances.

Why is this useful? Consider creating a database schema and you first need to know the data structures. Jomba will give you an abstract look at fields and structure of the data which can also help with parsing the data, thereby reducing the time to discover and know your data.

Jomba is run as command line utility that at the least needs to be passed a file that contains JSON data.

eg.

```
> jomba -f <path-to-file>
```

The path to the file can be relative or absolute.