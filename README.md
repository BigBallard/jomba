# JOMBA

JSON Object Mapping By Abstraction

- [Introduction](#introduction)
- [Output](#output)
- [Considerations](#considerations)

## Introduction
Jomba is a tool that will take a valid JSON formatted file and present the abstract JSON structure. The abstraction shows a single composite instance of each property, object, and array with a count for the number of instances.

Why is this useful? Consider creating a database schema and you first need to know the data structures. Jomba will give you an abstract look at fields and structure of the data which can also help with parsing the data, thereby reducing the time to discover and know your data.

Jomba is run as command line utility that at the least needs to be passed a file that contains JSON data.

eg.

```
> jomba -i <path-to-file>
```

The path to the file can be relative or absolute.

## Output

After processing a JSON file, the output is either printed to the console or written to a file as instructed by the command line options given. Either way, the output format is the same.

The root of every JSON file is required to be an object, so the first line will generally look like this:
```
{ 
```
This equates to saying one anonymous object instance at that level. Following will be a list of fields in order of parsing, or as they would appear in the provided JSON file. Following the name of the field will be the number of times it occurs in the aggregate of objects.
```
{ 
    name: 1
```
Generally speaking, the number of field instances should not be less than or equal to the number of objects. So here is an example of a deeper aggregate:
```
...
    orders: { 18
        id: 18    
        options: 3
    }
...
```
In this example, there are a total of 18 instances of the order field objects. The aggregate of an order shows to have 18 objects with the `id` field, and 3 objects with the `options` field.

Arrays are essentially shown the same way, but with an array bracket. The fields and counts in the array notation show the aggregate of objects within that array.

> Arrays within arrays are not currently supported but are currently in the think tank!

## Considerations

Jomba takes the general intent of JSON structures when working with arrays and therefor makes some assumptions on the data that is provided. What is meant by this is that arrays are typically a series of entities that are instances of the same data model. This is not always the case however, and there are times where arrays are used to group data almost like an object is meant to. You knowing the data schema is going to help with interpreting what Jomba outputs.