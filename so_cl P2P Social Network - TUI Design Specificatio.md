<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# so_cl P2P Social Network - TUI Design Specification (Terminal Symbols)

## 1. Overall Layout Structure

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ {＊} so_cl │ Social P2P Network v0.1.5                                        │
├─────────────┬─────────────────────────────────┬─────────────────────────────┤
│   NAV       │           MAIN CONTENT          │          SIDEBAR            │
│ (20% width) │         (50% width)            │       (30% width)           │
│             │                                 │                             │
│ [Menu Items]│        Page Content            │        Context Info         │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
│             │                                 │                             │
└─────────────┴─────────────────────────────────┴─────────────────────────────┘
```


## 2. Navigation Menu (Left Panel)

```
NAV
┌─────────────┐
│ [Home]      │  ← Current page highlighted
│ [Discover]  │
│ [Peers]     │
│ [Profile]   │
│ [Settings]  │
└─────────────┘
```

**Navigation Features:**

- Current page: Reverse video (background green, text black)
- Menu items in brackets `[ ]`
- Vertical spacing between items
- Single pixel border around panel


## 3. Home Page Layout

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ {＊} so_cl │ Social P2P Network v0.1.5                                        │
├─────────────┬─────────────────────────────────┬─────────────────────────────┤
│   NAV       │             FEED                │          NETWORK            │
│             │                                 │                             │
│ [Home]      │ ┌─────────────────────────────┐ │ Status: Connected           │
│ [Discover]  │ │ @alice: Just deployed my   │ │ Peers: 42                   │
│ [Peers]     │ │ new DHT node! *            │ │ Speed: 1.2 NB/s             │
│ [Profile]   │ │ <reply> <share> <like> 23   │ │ ┌─────────────────────────┐ │
│ [Settings]  │ │ 2 min ago                  │ │ │    ASCII ART GALLERY     │ │
│             │ └─────────────────────────────┘ │ │                         │ │
│             │                                 │ │   [preview 1]          │ │
│             │ ┌─────────────────────────────┐ │ │   [preview 2]          │ │
│             │ │ @bob: Check out this ASCII  │ │ │   [preview 3]          │ │
│             │ │ art I made! ~              │ │ └─────────────────────────┘ │
│             │ │ <reply> <share> <like> 42   │ │                             │
│             │ │ 15 min ago                 │ │                             │
│             │ └─────────────────────────────┘ │                             │
│             │                                 │                             │
│             │ [New Post] [Load More]          │                             │
└─────────────┴─────────────────────────────────┴─────────────────────────────┘
```

**Post Format:**

```
┌─────────────────────────────┐
│ @username: Post content...   │  ← Max 280 characters
│ <reply> <share> <like> ##    │  ← Interaction buttons
│ timestamp ago                │  ← Relative time
└─────────────────────────────┘
```


## 4. Discover Page Layout

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ {＊} so_cl │ Social P2P Network v0.1.5                                        │
├─────────────┬─────────────────────────────────┬─────────────────────────────┤
│   NAV       │            DISCOVER              │          NETWORK            │
│             │                                 │                             │
│ [Home]      │ ┌─ TRENDING TOPICS ───────────┐ │ Status: Connected           │
│ [Discover]  │ │ #retro-terminal  * 234      │ │ Peers: 42                   │
│ [Peers]     │ │ #ascii-art     ^ 156        │ │ Speed: 1.2 NB/s             │
│ [Profile]   │ │ #p2p-social    + 89         │ │ ┌─────────────────────────┐ │
│ [Settings]  │ │ #dht-network   + 45         │ │ │    ASCII ART GALLERY     │ │
│             │ └─────────────────────────────┘ │ │                         │ │
│             │                                 │ │   [trending art 1]      │ │
│             │ ┌─ POPULAR POSTS ─────────────┐ │ │   [trending art 2]      │ │
│             │ │ @artist: New ASCII masterpiece│ │ │   [trending art 3]      │ │
│             │ │ <like> 234 <share> 45       │ │ └─────────────────────────┘ │
│             │ └─────────────────────────────┘ │                             │
│             │                                 │                             │
│             │ ┌─ NEW PEERS TO FOLLOW ───────┐ │                             │
│             │ │ [@coder42]  •  127 posts    │ │                             │
│             │ │ [@pixelart] •  89 posts     │ │                             │
│             │ │ [@retrodev] •  156 posts    │ │                             │
│             │ └─────────────────────────────┘ │                             │
└─────────────┴─────────────────────────────────┴─────────────────────────────┘
```

**Trending Indicators:**

- `*` = Hot/trending (replaces 🔥)
- `^` = Rising (replaces 📈)
- `+` = Growing (replaces ⬆)


## 5. Peers Page Layout

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ {＊} so_cl │ Social P2P Network v0.1.5                                        │
├─────────────┬─────────────────────────────────┬─────────────────────────────┤
│   NAV       │              PEERS               │        NETWORK DETAILS      │
│             │                                 │                             │
│ [Home]      │ ┌─ CONNECTED PEERS (42) ──────┐ │ Node ID: 0x7f3a...9b2c     │
│ [Discover]  │ │ ● @alice      • Active      │ │ Status: Direct Connection   │
│ [Peers]     │ │   192.168.1.42:8080         │ │ Uptime: 2h 34m 12s         │
│ [Profile]   │ │   ↑ 1.2 MB/s  ↓ 0.8 MB/s    │ │ Protocol: libp2p v1.2.3     │
│ [Settings]  │ │   [Message] [Follow] [Info] │ │                             │
│             │ └─────────────────────────────┘ │ ┌─ CONNECTION METRICS ────┐ │
│             │                                 │ │ Latency: 23ms             │ │
│             │ ┌─ RECENT PEERS ───────────────┐ │ │ Bandwidth: 2.0 MB/s       │ │
│             │ │ ○ @charlie    • 5 min ago   │ │ │ Packets: 1,247/s         │ │
│             │ │   Disconnected              │ │ │ Loss: 0.1%                │ │
│             │ │   [Reconnect]               │ │ └─────────────────────────┘ │
│             │ └─────────────────────────────┘ │                             │
│             │                                 │                             │
│             │ ┌─ PEER STATISTICS ────────────┐ │                             │
│             │ │ Total Connections: 127       │ │                             │
│             │ │ Active Sessions: 42          │ │                             │
│             │ │ Data Transferred: 2.3 GB     │ │                             │
│             │ │ Network Health: 98%          │ │                             │
│             │ └─────────────────────────────┘ │                             │
└─────────────┴─────────────────────────────────┴─────────────────────────────┘
```

**Peer Status Indicators:**

- ● = Active/Online
- ○ = Offline/Disconnected
- `>>` = High-speed connection (replaces ⚡)
- `..` = Slow connection (replaces 🐢)


## 6. Profile Page Layout

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ {＊} so_cl │ Social P2P Network v0.1.5                                        │
├─────────────┬─────────────────────────────────┬─────────────────────────────┤
│   NAV       │              PROFILE             │          NETWORK            │
│             │                                 │                             │
│ [Home]      │ ┌─────────────────────────────┐ │ Status: Connected           │
│ [Discover]  │ │  ╔═══════════════════════╗  │ │ Peers: 42                   │
│ [Peers]     │ │  ║  ▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄  ║  │ │ Speed: 1.2 NB/s             │
│ [Profile]   │ │  ║  █▀▄ ▄▀▄ ▀█ ▀▄▀▄ ▄▀▄ █  ║  │ │ ┌─────────────────────────┐ │
│ [Settings]  │ │  ║  █ █ █ █ █▄▄▄▄▄▄▄▄▄▄█  ║  │ │ │    ASCII ART GALLERY     │ │
│             │ │  ║  █▄▀ ▀▄▀ ▄█ ▄▀▄▀ ▀▄▀ █  ║  │ │ │                         │ │
│             │ │  ║  ▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀  ║  │ │ │   [your art 1]         │ │
│             │ │  ║    so_cl               ║  │ │ │   [your art 2]         │ │
│             │ │  ╚═══════════════════════╝  │ │ │   [your art 3]         │ │
│             │ │  @bobanekoso_cl            │ │ └─────────────────────────┘ │
│             │ │  Digital Artist & Retro     │ │                             │
│             │ │  Terminal Enthusiast        │ │                             │
│             │ │  Joined: 2023-03-15         │ │                             │
│             │ └─────────────────────────────┘ │                             │
│             │                                 │                             │
│             │ ┌─ PROFILE STATISTICS ────────┐ │                             │
│             │ │ Posts: 127    Following: 42 │ │                             │
│             │ │ Followers: 89  Likes: 1,247 │ │                             │
│             │ │ ASCII Arts: 15 Reposts: 234 │ │                             │
│             │ └─────────────────────────────┘ │                             │
└─────────────┴─────────────────────────────────┴─────────────────────────────┘
```

**Profile Components:**

- **ASCII Avatar**: Boxed ASCII art (20x10 characters)
- **Username**: Below avatar
- **Display Name**: Under username

