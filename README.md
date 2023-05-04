# Werewolf Game

<p align="center">
<a href="https://tp-o.itch.io/werewolf" target="_blank">
    <img src="./assets/play-btn.gif" width="200" height="200" />
</a>
<h2 align="center">Click to play</h2>
</p>

# Introduction

## Overview

This multiplayer game is based on a very well-known board card game Werewolf. In the original game, it only allows players to play in order while everyone else have to wait for their turns. For this reason, we decided to make some changes to gameplay to reduce the waiting time and make it more enjoyable. The player must control their chacracter to complete their job instead of just selecting boring options. Basically, the gameplay is enhanced, but the game concept is kept the same as the original.

## Gameplay

The two main factions are [Villager](#role) vs [Werewolf](#role), but there also have third factions to make the game more interesting. Each faction has its own win condition, and the game is over if one side wins. Day, dusk and night are 3 [implemented phases](#phase). In each phase, there is one or more turns relied on the played [roles](#role).

There are 2 types of turn: private and public. In private turn, only roles played in that turn can play; otherwise, everyone can play. One turn is also played by one or more roles, and if they finish or skip their work, they can do whatever they want until the end of the current phase.

In the moring, Villager gathers together and finds the Werewolf. The special roles do their job at dusk, and Werewolf and the rest do their job at night.

# Technical

## Challenge

...

## Platform & Tool

- Go
- Node.js
- Unity
- Redis cluster
- PostgreSQL

# Demo

Click [here](https://tp-o.itch.io/werewolf).

# License

...
