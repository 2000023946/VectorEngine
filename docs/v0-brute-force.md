# 01 — Brute-Force Vector Search

## Goal

Build the simplest correct vector search engine possible.

This implementation will serve as the baseline for every optimization
added later to VectorEngine.

---

## Problem

VectorEngine needs to store vectors and find the vectors most similar
to a given query vector.

The simplest possible approach is to compare the query against every
vector in the database.

For every stored vector:

1. Calculate its distance from the query.
2. Keep track of the closest vectors.
3. Return the top-k closest vectors.

This is called brute-force or exact nearest-neighbor search.

---

## Initial Architecture

```text
Client
  |
  v
VectorEngine
  |
  v
Vector Storage
  |
  v
Linear Scan
  |
  v
Distance Calculation
  |
  v
Top-K Results