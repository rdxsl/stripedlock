# Striped Lock

## Usage

This repo implements a Golang striped locking library for use in concurrency. Based on ideas from Maya Raviv and also the following discussion: https://plus.google.com/+googleguava/posts/DAH1HKhfuoE

A striped lock is initialized with an array of sync.Mutex of a fixed size, and resources are striped across the locks based on the hashcode of their resource id.

The goal of striped locking is to give a developer control over the memory vs speed tradeoff when needing to lock a large number of resources.