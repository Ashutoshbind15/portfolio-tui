---
title: React Compiler — a practical guide
date: 2026-05-23
summary: Write normal React. Let the compiler skip the work you used to memoize by hand.
---

# React Compiler — a practical guide

The React Compiler is a build-time tool that optimizes your React code
automatically. You write normal components and hooks; the compiler
analyzes them and rewrites the output so your app skips unnecessary
work on updates.

Compile-time optimization is not new to frontend — Svelte, Vue, Solid,
and others have done this for a while. React’s version arrived later,
but it is built on deep static analysis and is widely considered one of
the best optimizing compilers in the space. You keep writing React
as-is.

## What it does (beyond useMemo)

Most people hear “React Compiler” and think automatic memoization. That
is the main user-facing win, but a compiler does more than cache values.

The React Compiler lowers your code into an internal representation
(HIR), runs multiple analysis passes over it, and emits optimized
JavaScript. Along the way it:

- Tracks data flow and mutability to understand which values change
  together and which are safe to reuse
- Validates Rules of React, catching things like conditional hook
  calls, setState during render, and invalid mutations of props/state
- Inserts memoization automatically, including conditional memoization
  that is hard to write by hand
- Rewrites render logic at build time so stable subtrees and values
  are cached and children do not re-render when nothing relevant changed

You do not add annotations in the default setup. The compiler works on
plain JavaScript/TypeScript that follows normal React conventions.

## Why this matters for how you write code

Before the compiler, performance often meant sprinkling `React.memo`,
`useMemo`, and `useCallback` through hot paths. That code is noisy, easy
to get wrong, and makes components harder to read.

With the compiler, you can write straightforward React:

- Derived values computed directly in the component body
- Event handlers defined as regular functions and passed to children
- Multiple `useState` hooks and `useEffect` calls without wrapping
  everything defensively

That clarity helps humans and coding agents. When components read like
plain logic — state at the top, derived values next, handlers below,
JSX at the bottom — agents generate fewer incorrect memo dependencies,
miss fewer stale-closure bugs, and produce code that is easier to
review. Less boilerplate means less surface area for both you and the
agent to get wrong.

## How parent-to-child memoization actually works

With `useCallback` and `React.memo`, passing something like
`onClick={() => handleMove(id)}` directly in JSX creates a new function
reference every render. The child sees a “new” prop every time and
re-renders anyway.

The React Compiler handles this in the parent. It memoizes callbacks,
object literals, and JSX elements inside the component body. When the
values those depend on have not changed, the child receives the same
prop references as last render and can skip re-rendering.

Both the parent and the child need to be compiled for this to work end
to end. In infer mode that happens automatically as long as both follow
naming conventions (PascalCase components, `use*` hooks).

## A real example from a chess app

I run the compiler in a chess training app — puzzle solving, game
analysis, opening trees. The whole client has zero `useMemo`,
`useCallback`, or `React.memo`.

**Puzzle view** (parent) holds move progress, a review/scrub index,
display position, and an async mutation for submitting moves. On each
render it derives board state from FEN and move history (live position
vs. review position). It passes an `onMove` handler and the current FEN
down to a chessboard child, which computes legal-move highlights when a
piece is selected.

The parent updates often — async pending state, step-through review,
animated move playback. Without the compiler you would normally memoize
the handler and derived board props so the board child does not repaint
on every tick. With the compiler, both components compile cleanly and
the board only updates when FEN, orientation, or interactivity actually
change.

**Game analysis view** (parent) is heavier: a branching analysis tree
with add/delete mutations, debounced engine evaluation updating every
~180ms, and tabbed panels for commentary, recommendations, and a move
tree graph. It renders a mini board child synced to the current tree
position, plus panel children for engine lines and graph nodes.

Engine eval ticks cause frequent parent re-renders. The panels and board
children only re-render when the slice of tree state they depend on
changes, because the compiler stabilizes the props and JSX each child
receives from the parent.

The pattern is the same everywhere in the app: stateful parent, derived
values in the render body, handlers passed straight to children, no
manual memo layer.

## Compilation modes

Two you need to know:

**infer** (default). The compiler finds components and hooks by
convention: PascalCase names that return JSX, functions prefixed with
`use` that call hooks. No extra code required. This is where most
projects should end up.

**annotation**. Only compiles functions marked with `"use memo"` at the
top of the body. Useful for trying the compiler on a few components in
an existing app before going all-in.

You can opt out of any function with `"use no memo"` while debugging.

## Common gotchas

- Focuses on update performance (fewer unnecessary re-renders), not
  faster initial loads
- Only optimizes function components and hooks, not class components
- Manual `useMemo` / `useCallback` still work as escape hatches when
  you need explicit control (e.g. effect dependencies)
- Parent and child both need to be compiled for prop stabilization to
  carry across the boundary

## Getting started

Works best with React 19 (runtime built into React). Also supports
React 17 and 18: install `react-compiler-runtime` and set
`target: '18'` or `'17'` in config. Either way, install the Babel
plugin:

```bash
pnpm add -D babel-plugin-react-compiler
```

Wire it into your build (must run before other transforms):

- **Next.js** — `reactCompiler: true` in config
- **Vite** (`@vitejs/plugin-react` v6+) — `@rolldown/plugin-babel` +
  `reactCompilerPreset()`
- **Expo** (SDK 54+) — enable the `reactCompiler` experiment in app
  config

Enable `eslint-plugin-react-hooks` (latest) and fix Rules of React
violations. Start with infer mode, or annotation + `"use memo"` for a
smaller first step. Check React DevTools v5+ for the Memo ✨ badge on
compiled components.

Write normal React. Turn it on, follow naming conventions, fix lint
violations, and let the compiler handle the rest.
