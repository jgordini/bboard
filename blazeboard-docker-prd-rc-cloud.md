# BlazeBoard Docker Implementation PRD

**Version:** 1.0  
**Last Updated:** January 30, 2026  
**Scope:** Docker containerization and deployment infrastructure for BlazeBoard

---

## 1. Executive Summary

This document specifies the Docker implementation requirements for deploying BlazeBoard—a community feedback and idea management platform—in containerized environments. The implementation supports both self-hosted white-label deployments and UAB internal infrastructure, with primary deployment on UAB Research Computing Cloud (RC Cloud).

**Primary Objectives:**

- Provide portable, reproducible deployments across any Docker-compatible infrastructure
- Enable zero-downtime deployments with health checks and graceful shutdowns
- Support both development and production configurations
- Maintain security best practices for container deployments
- Deploy on UAB RC Cloud ([cloud.rc.uab.edu](https://dashboard.cloud.rc.uab.edu)) using OpenStack virtual machines

---

## 2. Container Architecture Overview

### 2.1 Service Topology

```
┌─────────────────────────────────────────────────────────────────┐
│                        Host Machine / VM                         │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────────────┐   │
│  │   Nginx     │   │  BlazeBoard │   │     PostgreSQL      │   │
│  │   Reverse   │──▶│   Next.js   │──▶│     Database        │   │
│  │   Proxy     │   │   App       │   │     (Optional*)     │   │
│  │   :80/:443  │   │   :3000     │   │     :5432           │   │
│  └─────────────┘   └─────────────┘   └─────────────────────┘   │
│         │                                                       │
│         ▼                                                       │
│  ┌─────────────┐                                               │
│  │  SSL Certs  │                                               │
│  │  (mounted)  │                                               │
│  └─────────────┘                                               │
└─────────────────────────────────────────────────────────────────┘

*PostgreSQL can be external (Azure PostgreSQL, Supabase, etc.)
```

### 2.2 Container Services

| Service            | Base Image           | Purpose                        | Exposed Port    |
| ------------------ | -------------------- | ------------------------------ | --------------- |
| `blazeboard-app`   | `node:20-alpine`     | Next.js application server     | 3000 (internal) |
| `blazeboard-nginx` | `nginx:1.25-alpine`  | Reverse proxy, SSL termination | 80, 443         |
| `blazeboard-db`    | `postgres:16-alpine` | PostgreSQL database (optional) | 5432 (internal) |

---

## 3. Core Application Features

The following features are delivered through the containerized Next.js application:

### 3.1 Idea Management

| Feature             | Description                                                                                    |
| ------------------- | ---------------------------------------------------------------------------------------------- |
| **Create**          | Submit ideas with title (≤200 chars), description (≤2,000 chars), impact statement, categories |
| **Read**            | Paginated list with filtering (category, status, date) and sorting (votes, date, activity)     |
| **Update**          | Edit own ideas within 15-minute window; admins/moderators can edit any                         |
| **Delete**          | Soft delete with audit trail; moderation review for permanent removal                          |
| **Status Tracking** | `new` → `under_review` → `planned` → `completed` / `declined`                                  |
| **Idea Numbering**  | Sequential ID for easy reference (e.g., "IDEA-42")                                             |

### 3.2 Voting Engine

| Feature           | Description                                                             |
| ----------------- | ----------------------------------------------------------------------- |
| **Vote Types**    | Binary upvote/downvote (+1/-1)                                          |
| **Deduplication** | Database constraint `UNIQUE(user_id, idea_id)` prevents duplicate votes |
| **Real-Time**     | Postgres `LISTEN/NOTIFY` pushes vote changes to all connected clients   |
| **Display**       | Net score with visual progress bar                                      |
| **Notifications** | Email alerts when ideas reach configurable vote thresholds              |

### 3.3 Commenting System

| Feature        | Description                             |
| -------------- | --------------------------------------- |
| **Threading**  | Nested reply threads (infinite depth)   |
| **Formatting** | Markdown support with safe rendering    |
| **Voting**     | Up/down voting on comments              |
| **Moderation** | Pinning, flagging, edit window (15 min) |

### 3.4 Moderation Suite

| Feature               | Description                                  |
| --------------------- | -------------------------------------------- |
| **Flagging**          | User/content flagging with reason codes      |
| **Queue**             | Moderation queue with bulk actions           |
| **User Management**   | Temporary/permanent banning                  |
| **Content Filtering** | Profanity filter with configurable blocklist |
| **Audit Trail**       | Full log of all moderation actions           |

---

## 4. Frontend Application Specifications

### 4.1 Technology Stack

| Layer                | Technology                      | Version | Purpose                                              |
| -------------------- | ------------------------------- | ------- | ---------------------------------------------------- |
| **Framework**        | Next.js (App Router)            | 14.x+   | Server-side rendering, API routes, Server Components |
| **Language**         | TypeScript                      | 5.x     | Type safety across frontend and backend              |
| **UI Components**    | shadcn/ui                       | Latest  | Copy-paste component architecture, no vendor lock-in |
| **Styling**          | Tailwind CSS                    | 3.x     | Utility-first CSS, rapid development                 |
| **State Management** | React Server Components + hooks | -       | Minimal client-side state                            |
| **Real-Time**        | Postgres LISTEN/NOTIFY          | -       | Live vote counts without WebSocket services          |
| **Forms**            | React Hook Form + Zod           | -       | Form validation and handling                         |
| **Markdown**         | react-markdown                  | -       | Comment/idea body rendering                          |

#### 4.1.1 shadcn/ui Component System

BlazeBoard uses **[shadcn/ui](https://ui.shadcn.com/)** as the primary UI component library. shadcn/ui is a collection of re-usable components built with **Radix UI** primitives and **Tailwind CSS**, copied into the project (copy-paste architecture) rather than installed as a dependency. This approach provides:

| Benefit | Description |
| -------- | ----------- |
| **No vendor lock-in** | Components live in the codebase; styling and behavior can be customized freely. |
| **Accessibility** | Radix UI primitives provide keyboard navigation, focus management, and ARIA patterns out of the box. |
| **Design consistency** | Components align with the Fider-inspired design system and UAB branding (see §4.2). |
| **Tailwind integration** | Full compatibility with Tailwind utility classes and the project's semantic color tokens. |

Required base components for BlazeBoard include: `button`, `card`, `badge`, `input`, `textarea`, `select`, `dropdown-menu`, `dialog`, `toast`, `avatar`, `separator`, and `progress`. New components are added via the CLI: `npx shadcn-ui@latest add <component>`. Full setup, Tailwind configuration for UAB + Fider design, and component implementation examples are in **§4.2.10 shadcn/ui Implementation Guide**.

### 4.2 Design System (Fider-Inspired)

BlazeBoard's UI follows a clean, modern design language inspired by Fider, emphasizing simplicity, clarity, and user engagement.

#### 4.2.1 Visual Design Principles

| Principle         | Implementation                                               |
| ----------------- | ------------------------------------------------------------ |
| **Minimalism**    | Clean layouts with ample whitespace, reducing cognitive load |
| **Clarity**       | Clear visual hierarchy using typography, spacing, and color  |
| **Engagement**    | Interactive elements with subtle animations and hover states |
| **Accessibility** | High contrast, large touch targets, clear focus states       |
| **Consistency**   | Unified design language across all pages and components      |

#### 4.2.2 Typography System

```typescript
// Design tokens for typography
const typography = {
  fonts: {
    heading:
      '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    body: '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
    mono: '"Fira Code", "Courier New", monospace',
  },
  sizes: {
    hero: '3rem', // 48px - Hero headlines ("Help us build the best feedback platform")
    h1: '2rem', // 32px - Page titles ("Restrict how many ideas an user can vote")
    h2: '1.5rem', // 24px - Section headers
    h3: '1.25rem', // 20px - Card titles
    body: '1rem', // 16px - Main content
    small: '0.875rem', // 14px - Meta information, timestamps
    tiny: '0.75rem', // 12px - Labels, badges
  },
  weights: {
    regular: 400,
    medium: 500,
    semibold: 600,
    bold: 700,
  },
  lineHeights: {
    tight: 1.2,
    normal: 1.5,
    relaxed: 1.75,
  },
};
```

#### 4.2.3 Color Palette (UAB + Fider-Inspired)

```typescript
// Color system combining UAB branding with Fider's clean aesthetic
const colors = {
  // Primary - UAB Green
  primary: {
    50: '#f0f7f4',
    100: '#d9ece3',
    500: '#1A5632', // Main UAB Green
    600: '#154428',
    700: '#10331e',
    900: '#0a1f13',
  },

  // Accent - Campus Green
  accent: {
    50: '#f7fde7',
    500: '#90D408', // Main Campus Green
    600: '#7ab307',
    700: '#659205',
  },

  // Interactive - Blue (Fider-inspired)
  interactive: {
    50: '#eff6ff',
    500: '#3b82f6', // Primary action color (Vote buttons, Sign in)
    600: '#2563eb',
    700: '#1d4ed8',
  },

  // Neutrals - Gray scale
  neutral: {
    50: '#fafafa', // Background
    100: '#f5f5f5', // Card backgrounds
    200: '#e5e5e5', // Borders
    300: '#d4d4d4', // Disabled states
    500: '#737373', // Secondary text
    700: '#404040', // Primary text
    900: '#171717', // Headings
  },

  // Status colors
  status: {
    open: '#3b82f6', // Blue
    discussion: '#f59e0b', // Orange/Amber
    planned: '#8b5cf6', // Purple
    completed: '#10b981', // Green
    declined: '#6b7280', // Gray
  },

  // Semantic colors
  semantic: {
    success: '#10b981',
    warning: '#FDB913', // UAB Gold
    error: '#ef4444',
    info: '#3b82f6',
  },
};
```

#### 4.2.4 Layout Structure (Homepage - Fider Pattern)

```
┌─────────────────────────────────────────────────────────────────┐
│ Header: Logo, RSS, Dark Mode Toggle, Sign In Button            │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ Hero Section (Left-aligned)                               │ │
│  │ ┌─────────────────────────────────────────┐               │ │
│  │ │ "Help us build the                      │   Search      │ │
│  │ │  best feedback platform"                │   Box         │ │
│  │ │                                         │               │ │
│  │ │ Tagline + Call to Action                │               │ │
│  │ └─────────────────────────────────────────┘               │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ [➕ Enter your suggestion here...]                        │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  [🔍 FILTER]  [SORT BY: TRENDING ▼]              [Search 🔎]  │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐ │
│  │ Idea Card                                      1 ♡        │ │
│  │ ┌─────────────────────────────────────────────────────┐   │ │
│  │ │ Add an option to display the privacy policy...      │   │ │
│  │ │ Right now, the privacy policy is only shown...      │   │ │
│  │ │                                                     │   │ │
│  │ │ 4 Votes                                             │   │ │
│  │ └─────────────────────────────────────────────────────┘   │ │
│  └───────────────────────────────────────────────────────────┘ │
│                                                                 │
│  [More idea cards in list format...]                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### 4.2.5 Idea Detail Page Layout (Fider Pattern - Simplified)

```
┌─────────────────────────────────────────────────────────────────┐
│ Header: Logo, RSS, Dark Mode Toggle, Sign In Button            │
├─────────────────────────────────────────────────────────────────┤
│ [← Back to all suggestions]                                     │
│                                                                 │
│ ┌───────────┐  ┌─────────────────────────────────────────────┐ │
│ │   Votes   │  │ Idea Title (Large, Bold)                    │ │
│ │           │  │ "Restrict how many ideas an user can vote"  │ │
│ │    148    │  │                                             │ │
│ │           │  │ 👤 Posted by Author · Nov 2017 · OPEN      │ │
│ │   Votes   │  │                                             │ │
│ │           │  │ Description text...                         │ │
│ │           │  │                                             │ │
│ │           │  │ [🏷️ Under Discussion]                      │ │
│ │           │  │                                             │ │
│ │           │  │ [💙 Vote for this idea]                     │ │
│ │           │  │                                             │ │
│ │           │  │ 148 Votes                                   │ │
│ │           │  │ ████████████████████░░░░░░                  │ │
│ │           │  │                                             │ │
│ │           │  │ [🔗 Copy link] [📡 Comment Feed]            │ │
│ │           │  │                                             │ │
│ │           │  │ ──────────────────────────────────────      │ │
│ │           │  │                                             │ │
│ │           │  │ Discussion   29                             │ │
│ │           │  │                                             │ │
│ │           │  │ ┌─────────────────────────────────────────┐│ │
│ │Powered by │  │ │ [Rich text editor with toolbar]         ││ │
│ │BlazeBoard │  │ │ Leave a comment                         ││ │
│ │           │  │ │                                         ││ │
│ │           │  │ └─────────────────────────────────────────┘│ │
│ │           │  │ [Post]                                      │ │
│ │           │  │                                             │ │
│ │           │  │ ┌─────────────────────────────────────────┐│ │
│ │           │  │ │ 👤 Matt Roberts ✓  5 months ago   ⋯    ││ │
│ │           │  │ │ It's still something we're mulling...   ││ │
│ │           │  │ └─────────────────────────────────────────┘│ │
│ │           │  │                                             │ │
│ │           │  │ [More comments...]                          │ │
│ └───────────┘  └─────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

#### 4.2.6 Component Specifications

##### Idea Card (List View)

```typescript
interface IdeaCard {
  layout: 'horizontal' | 'compact';
  spacing: {
    padding: '1.5rem'; // 24px
    gap: '1rem'; // 16px between elements
  };
  elements: {
    title: {
      fontSize: '1.125rem'; // 18px
      fontWeight: 600;
      color: 'neutral.900';
      maxLines: 2;
      lineHeight: 1.4;
    };
    description: {
      fontSize: '0.875rem'; // 14px
      color: 'neutral.600';
      maxLines: 2;
      lineHeight: 1.5;
    };
    voteCount: {
      fontSize: '1.5rem'; // 24px
      fontWeight: 700;
      position: 'left-aligned';
      label: 'Votes';
    };
    metadata: {
      fontSize: '0.75rem'; // 12px
      color: 'neutral.500';
      items: ['commentCount', 'timestamp'];
    };
  };
  hover: {
    background: 'neutral.50';
    transition: 'all 150ms ease';
  };
}
```

##### Vote Button (Primary Action)

```typescript
interface VoteButton {
  default: {
    background: 'interactive.500'; // Blue #3b82f6
    color: 'white';
    padding: '0.75rem 1.5rem'; // 12px 24px
    borderRadius: '0.375rem'; // 6px
    fontSize: '0.875rem'; // 14px
    fontWeight: 600;
    icon: '💙' | 'thumbs-up';
  };
  voted: {
    background: 'interactive.600';
    icon: '✓';
  };
  hover: {
    background: 'interactive.600';
    transform: 'translateY(-1px)';
    boxShadow: '0 4px 12px rgba(59, 130, 246, 0.25)';
  };
  disabled: {
    background: 'neutral.300';
    cursor: 'not-allowed';
    opacity: 0.6;
  };
}
```

##### Status Badge

```typescript
interface StatusBadge {
  variants: {
    open: {
      background: 'rgba(59, 130, 246, 0.1)';
      color: 'interactive.700';
      text: 'OPEN';
    };
    discussion: {
      background: 'rgba(245, 158, 11, 0.1)';
      color: '#d97706';
      text: 'Under Discussion';
    };
    planned: {
      background: 'rgba(139, 92, 246, 0.1)';
      color: '#7c3aed';
      text: 'Planned';
    };
    completed: {
      background: 'rgba(16, 185, 129, 0.1)';
      color: '#059669';
      text: 'Completed';
    };
    declined: {
      background: 'rgba(107, 114, 128, 0.1)';
      color: '#4b5563';
      text: 'Declined';
    };
  };
  style: {
    padding: '0.25rem 0.75rem'; // 4px 12px
    borderRadius: '0.25rem'; // 4px
    fontSize: '0.75rem'; // 12px
    fontWeight: 600;
    textTransform: 'uppercase';
    letterSpacing: '0.025em';
  };
}
```

##### Vote Count Sidebar Display

```typescript
interface VoteCountDisplay {
  layout: 'vertical-sidebar';
  container: {
    width: '120px'; // Simplified sidebar width
    padding: '1.5rem 1rem';
    textAlign: 'center';
  };
  voteNumber: {
    fontSize: '3rem'; // Large vote count display (48px)
    fontWeight: 700;
    color: 'neutral.900';
    marginBottom: '0.5rem';
  };
  label: {
    fontSize: '0.875rem';
    color: 'neutral.600';
    fontWeight: 500;
    text: 'Votes';
  };
  style: {
    display: 'flex';
    flexDirection: 'column';
    alignItems: 'center';
    gap: '0.25rem';
  };
}
```

##### Comment Thread

```typescript
interface CommentThread {
  comment: {
    padding: '1.5rem';
    borderBottom: '1px solid neutral.200';
    spacing: '1rem';
  };
  author: {
    display: 'flex';
    gap: '0.75rem';
    elements: {
      avatar: {
        size: '40px';
        borderRadius: '50%';
      };
      name: {
        fontSize: '0.875rem';
        fontWeight: 600;
        color: 'neutral.900';
      };
      badge: {
        // For verified/staff users
        color: 'interactive.500';
        icon: '✓';
      };
      timestamp: {
        fontSize: '0.75rem';
        color: 'neutral.500';
      };
    };
  };
  content: {
    fontSize: '0.875rem';
    lineHeight: 1.6;
    color: 'neutral.700';
    marginTop: '0.5rem';
  };
  actions: {
    display: 'flex';
    gap: '1rem';
    fontSize: '0.75rem';
    color: 'neutral.500';
    marginTop: '0.75rem';
    items: ['reply', 'edit', 'delete', 'menu'];
  };
}
```

##### Rich Text Editor

```typescript
interface RichTextEditor {
  toolbar: {
    background: 'neutral.100';
    borderBottom: '1px solid neutral.200';
    padding: '0.5rem';
    buttons: [
      'heading1',
      'heading2',
      'bold',
      'italic',
      'strikethrough',
      'bulletList',
      'numberedList',
      'code',
      'mention',
      'image',
      'link'
    ];
    buttonSize: '32px';
    spacing: '0.25rem';
  };
  textarea: {
    minHeight: '120px';
    padding: '1rem';
    fontSize: '0.875rem';
    lineHeight: 1.6;
    placeholder: 'Leave a comment';
    border: '1px solid neutral.200';
    borderRadius: '0.375rem';
  };
  actions: {
    position: 'bottom-left';
    submitButton: {
      text: 'Post';
      variant: 'interactive.500';
    };
    markdownToggle: {
      text: 'Switch to markdown editor';
      position: 'bottom-right';
      fontSize: '0.75rem';
      color: 'neutral.500';
    };
  };
}
```

#### 4.2.7 Spacing System

```typescript
// Consistent spacing scale (Tailwind-based)
const spacing = {
  0: '0',
  1: '0.25rem', // 4px
  2: '0.5rem', // 8px
  3: '0.75rem', // 12px
  4: '1rem', // 16px
  5: '1.25rem', // 20px
  6: '1.5rem', // 24px
  8: '2rem', // 32px
  10: '2.5rem', // 40px
  12: '3rem', // 48px
  16: '4rem', // 64px
  20: '5rem', // 80px
  24: '6rem', // 96px
};

// Component spacing guidelines
const componentSpacing = {
  cardPadding: spacing[6], // 24px
  sectionGap: spacing[8], // 32px
  headerPadding: spacing[4], // 16px
  buttonPadding: '12px 24px',
  inputPadding: spacing[3], // 12px
  containerMaxWidth: '1280px',
  contentMaxWidth: '900px',
};
```

#### 4.2.8 Interactive States

```typescript
// Hover, focus, active states
const interactiveStates = {
  hover: {
    card: {
      background: 'neutral.50',
      cursor: 'pointer',
      transition: 'background 150ms ease',
    },
    button: {
      primary: {
        background: 'interactive.600',
        transform: 'translateY(-1px)',
        boxShadow: '0 4px 12px rgba(59, 130, 246, 0.25)',
      },
      ghost: {
        background: 'neutral.100',
      },
    },
    link: {
      color: 'interactive.700',
      textDecoration: 'underline',
    },
  },
  focus: {
    outline: '2px solid interactive.500',
    outlineOffset: '2px',
    borderRadius: '0.25rem',
  },
  active: {
    transform: 'scale(0.98)',
    transition: 'transform 100ms ease',
  },
  disabled: {
    opacity: 0.6,
    cursor: 'not-allowed',
    pointerEvents: 'none',
  },
};
```

#### 4.2.9 Animation Guidelines

```typescript
// Subtle, purposeful animations
const animations = {
  transitions: {
    fast: '100ms ease',
    normal: '150ms ease',
    slow: '300ms ease',
  },
  properties: {
    fadeIn: {
      from: { opacity: 0 },
      to: { opacity: 1 },
      duration: '200ms',
    },
    slideIn: {
      from: { transform: 'translateY(8px)', opacity: 0 },
      to: { transform: 'translateY(0)', opacity: 1 },
      duration: '200ms',
    },
    scaleIn: {
      from: { transform: 'scale(0.95)', opacity: 0 },
      to: { transform: 'scale(1)', opacity: 1 },
      duration: '150ms',
    },
  },
  reducedMotion: {
    respectUserPreference: true,
    fallback: 'opacity-only',
  },
};
```

#### 4.2.10 shadcn/ui Implementation Guide

BlazeBoard uses [shadcn/ui](https://ui.shadcn.com/) for component architecture. shadcn/ui provides copy-paste components built with Radix UI and Tailwind CSS, offering full control without package dependencies.

##### Installation and Setup

```bash
# Initialize shadcn/ui in Next.js project
npx shadcn-ui@latest init

# Install required components for Fider design
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add badge
npx shadcn-ui@latest add input
npx shadcn-ui@latest add textarea
npx shadcn-ui@latest add select
npx shadcn-ui@latest add dropdown-menu
npx shadcn-ui@latest add dialog
npx shadcn-ui@latest add toast
npx shadcn-ui@latest add avatar
npx shadcn-ui@latest add separator
npx shadcn-ui@latest add progress
```

##### Tailwind Configuration for UAB + Fider Design

```typescript
// tailwind.config.ts
import type { Config } from "tailwindcss"

const config = {
  darkMode: ["class"],
  content: [
    './pages/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
    './app/**/*.{ts,tsx}',
    './src/**/*.{ts,tsx}',
  ],
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        // UAB Brand Colors
        uab: {
          green: "#1A5632",
          "dragons-lair": "#033319",
        },
        campus: {
          green: "#90D408",
        },
        loyal: {
          yellow: "#FDB913",
        },
        smoke: {
          gray: "#808285",
        },

        // Fider Interactive Colors
        interactive: {
          50: '#eff6ff',
          500: '#3b82f6',  // Primary action color
          600: '#2563eb',
          700: '#1d4ed8',
        },

        // shadcn/ui semantic colors
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      fontFamily: {
        sans: ['"Aktiv Grotesk"', '-apple-system', 'BlinkMacSystemFont', '"Segoe UI"', 'sans-serif'],
        serif: ['"Kulturista"', 'Georgia', '"Times New Roman"', 'serif'],
      },
      fontSize: {
        base: ['18px', { lineHeight: '1.75' }],  // UAB standard
      },
      keyframes: {
        "accordion-down": {
          from: { height: "0" },
          to: { height: "var(--radix-accordion-content-height)" },
        },
        "accordion-up": {
          from: { height: "var(--radix-accordion-content-height)" },
          to: { height: "0" },
        },
      },
      animation: {
        "accordion-down": "accordion-down 0.2s ease-out",
        "accordion-up": "accordion-up 0.2s ease-out",
      },
    },
  },
  plugins: [require("tailwindcss-animate")],
} satisfies Config

export default config
```

##### CSS Variables (globals.css)

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    /* UAB Brand Colors mapped to shadcn variables */
    --background: 0 0% 100%;
    --foreground: 154 53% 22%; /* UAB Green */

    --card: 0 0% 100%;
    --card-foreground: 154 53% 22%;

    --popover: 0 0% 100%;
    --popover-foreground: 154 53% 22%;

    /* Primary = UAB Green */
    --primary: 154 53% 22%;
    --primary-foreground: 0 0% 100%;

    /* Secondary = Campus Green */
    --secondary: 80 92% 43%;
    --secondary-foreground: 154 53% 22%;

    /* Muted = Smoke Gray */
    --muted: 0 0% 96%;
    --muted-foreground: 0 1% 51%;

    /* Accent = Interactive Blue */
    --accent: 217 91% 60%;
    --accent-foreground: 0 0% 100%;

    --destructive: 0 84% 60%;
    --destructive-foreground: 0 0% 98%;

    --border: 0 0% 90%;
    --input: 0 0% 90%;
    --ring: 154 53% 22%; /* UAB Green for focus rings */

    --radius: 0.375rem; /* 6px - Fider style */
  }

  .dark {
    --background: 154 53% 10%;
    --foreground: 0 0% 98%;

    --card: 154 53% 12%;
    --card-foreground: 0 0% 98%;

    --popover: 154 53% 12%;
    --popover-foreground: 0 0% 98%;

    --primary: 80 92% 43%; /* Campus Green in dark mode */
    --primary-foreground: 154 53% 22%;

    --secondary: 154 53% 15%;
    --secondary-foreground: 0 0% 98%;

    --muted: 154 53% 15%;
    --muted-foreground: 0 1% 65%;

    --accent: 217 91% 60%;
    --accent-foreground: 0 0% 98%;

    --destructive: 0 62% 30%;
    --destructive-foreground: 0 0% 98%;

    --border: 154 53% 20%;
    --input: 154 53% 20%;
    --ring: 80 92% 43%;
  }
}

@layer base {
  * {
    @apply border-border;
  }
  body {
    @apply bg-background text-foreground;
    font-size: 18px; /* UAB standard */
  }
}
```

##### Fider Design Component Implementations

###### Vote Button Component

```typescript
// components/voting/vote-button.tsx
'use client';

import { Button } from '@/components/ui/button';
import { ThumbsUp } from 'lucide-react';
import { cn } from '@/lib/utils';

interface VoteButtonProps {
  voted: boolean;
  loading?: boolean;
  disabled?: boolean;
  onClick: () => void;
  className?: string;
}

export function VoteButton({
  voted,
  loading,
  disabled,
  onClick,
  className,
}: VoteButtonProps) {
  return (
    <Button
      onClick={onClick}
      disabled={disabled || loading}
      className={cn(
        'bg-interactive-500 hover:bg-interactive-600',
        'text-white font-semibold',
        'px-6 py-3 rounded-md',
        'transition-all duration-150',
        'hover:-translate-y-0.5',
        'hover:shadow-lg hover:shadow-interactive-500/25',
        'disabled:opacity-60 disabled:cursor-not-allowed',
        voted && 'bg-interactive-600',
        className
      )}
    >
      <ThumbsUp className="mr-2 h-4 w-4" />
      {voted ? 'Voted' : 'Vote for this idea'}
    </Button>
  );
}
```

###### Idea Card Component

```typescript
// components/ideas/idea-card.tsx
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { MessageCircle } from 'lucide-react';
import Link from 'next/link';
import { cn } from '@/lib/utils';

interface IdeaCardProps {
  id: string;
  title: string;
  description: string;
  voteCount: number;
  commentCount: number;
  status:
    | 'open'
    | 'discussion'
    | 'planned'
    | 'completed'
    | 'declined';
  createdAt: Date;
}

const statusConfig = {
  open: {
    label: 'OPEN',
    className: 'bg-interactive-500/10 text-interactive-700',
  },
  discussion: {
    label: 'Under Discussion',
    className: 'bg-amber-500/10 text-amber-700',
  },
  planned: {
    label: 'Planned',
    className: 'bg-purple-500/10 text-purple-700',
  },
  completed: {
    label: 'Completed',
    className: 'bg-green-500/10 text-green-700',
  },
  declined: {
    label: 'Declined',
    className: 'bg-gray-500/10 text-gray-700',
  },
};

export function IdeaCard(props: IdeaCardProps) {
  const status = statusConfig[props.status];

  return (
    <Link href={`/ideas/${props.id}`}>
      <Card className="hover:bg-muted/50 transition-colors cursor-pointer">
        <CardContent className="p-6 flex gap-4">
          {/* Vote Count */}
          <div className="flex flex-col items-center min-w-[60px]">
            <span className="text-3xl font-bold text-foreground">
              {props.voteCount}
            </span>
            <span className="text-sm text-muted-foreground">
              Votes
            </span>
          </div>

          {/* Content */}
          <div className="flex-1 min-w-0">
            <div className="flex items-start justify-between gap-4 mb-2">
              <h3 className="font-semibold text-lg text-foreground line-clamp-2">
                {props.title}
              </h3>
              <Badge
                variant="secondary"
                className={cn('shrink-0', status.className)}
              >
                {status.label}
              </Badge>
            </div>

            <p className="text-sm text-muted-foreground line-clamp-2 mb-3">
              {props.description}
            </p>

            <div className="flex items-center gap-4 text-xs text-muted-foreground">
              <span className="flex items-center gap-1">
                <MessageCircle className="h-3 w-3" />
                {props.commentCount} comments
              </span>
              <span>{props.createdAt.toLocaleDateString()}</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}
```

###### Hero Section Component

```typescript
// components/layout/hero-section.tsx
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Plus } from 'lucide-react';

interface HeroSectionProps {
  headline: string;
  description: string;
  callToAction?: string;
  onSuggestionClick: () => void;
}

export function HeroSection(props: HeroSectionProps) {
  return (
    <div className="bg-gradient-to-b from-muted/30 to-background py-16">
      <div className="container max-w-4xl">
        {/* Headline */}
        <h1 className="text-5xl font-bold text-primary mb-4 font-serif">
          {props.headline}
        </h1>

        {/* Description */}
        <p className="text-lg text-muted-foreground mb-6 leading-relaxed">
          {props.description}
        </p>

        {/* Call to Action */}
        {props.callToAction && (
          <p className="text-sm text-muted-foreground mb-8">
            {props.callToAction}
          </p>
        )}

        {/* Suggestion Input */}
        <div
          onClick={props.onSuggestionClick}
          className="flex items-center gap-3 p-4 bg-card border border-input rounded-lg hover:border-ring cursor-pointer transition-colors"
        >
          <Plus className="h-5 w-5 text-interactive-500" />
          <span className="text-muted-foreground">
            Enter your suggestion here...
          </span>
        </div>
      </div>
    </div>
  );
}
```

###### Status Badge Component

```typescript
// components/ideas/status-badge.tsx
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

type IdeaStatus =
  | 'open'
  | 'discussion'
  | 'planned'
  | 'completed'
  | 'declined';

const statusConfig = {
  open: {
    label: 'OPEN',
    className:
      'bg-interactive-500/10 text-interactive-700 hover:bg-interactive-500/20',
  },
  discussion: {
    label: 'Under Discussion',
    className: 'bg-amber-500/10 text-amber-700 hover:bg-amber-500/20',
  },
  planned: {
    label: 'Planned',
    className:
      'bg-purple-500/10 text-purple-700 hover:bg-purple-500/20',
  },
  completed: {
    label: 'Completed',
    className: 'bg-green-500/10 text-green-700 hover:bg-green-500/20',
  },
  declined: {
    label: 'Declined',
    className: 'bg-gray-500/10 text-gray-700 hover:bg-gray-500/20',
  },
};

interface StatusBadgeProps {
  status: IdeaStatus;
  className?: string;
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const config = statusConfig[status];

  return (
    <Badge
      variant="secondary"
      className={cn(
        'font-semibold uppercase text-xs tracking-wide',
        config.className,
        className
      )}
    >
      {config.label}
    </Badge>
  );
}
```

###### Vote Progress Bar Component

```typescript
// components/voting/vote-progress.tsx
import { Progress } from '@/components/ui/progress';
import { cn } from '@/lib/utils';

interface VoteProgressProps {
  current: number;
  max: number;
  className?: string;
}

export function VoteProgress({
  current,
  max,
  className,
}: VoteProgressProps) {
  const percentage = Math.min((current / max) * 100, 100);

  return (
    <div className={cn('space-y-2', className)}>
      <div className="flex items-center justify-between">
        <span className="text-3xl font-bold text-foreground">
          {current}
        </span>
        <span className="text-sm text-muted-foreground">Votes</span>
      </div>
      <Progress
        value={percentage}
        className="h-2"
        indicatorClassName="bg-interactive-500"
      />
    </div>
  );
}
```

##### Component Usage Examples

```typescript
// app/page.tsx
import { HeroSection } from '@/components/layout/hero-section';
import { IdeaCard } from '@/components/ideas/idea-card';
import { Button } from '@/components/ui/button';
import { Filter } from 'lucide-react';

export default function HomePage() {
  const handleSuggestionClick = () => {
    // Open suggestion modal
  };

  return (
    <div>
      <HeroSection
        headline="Help us improve UAB IT services"
        description="Share your ideas to help us build better technology solutions for the UAB community. Your feedback shapes the future of IT at UAB."
        callToAction="Your ideas make a difference!"
        onSuggestionClick={handleSuggestionClick}
      />

      <div className="container py-8">
        {/* Filter Bar */}
        <div className="flex items-center gap-4 mb-6">
          <Button variant="outline" size="sm">
            <Filter className="mr-2 h-4 w-4" />
            FILTER
          </Button>
          {/* Sort dropdown */}
        </div>

        {/* Idea List */}
        <div className="space-y-4">{/* Map over ideas */}</div>
      </div>
    </div>
  );
}
```

### 4.3 Application Pages & Routes

| Route               | Page               | Description                                    | Auth Required     |
| ------------------- | ------------------ | ---------------------------------------------- | ----------------- |
| `/`                 | Home/Ideas List    | Paginated idea feed with filters, hero section | No (configurable) |
| `/ideas/[id]`       | Idea Detail        | Full idea view with voter sidebar and comments | No (configurable) |
| `/ideas/new`        | Submit Idea        | Inline suggestion input or full form modal     | Yes               |
| `/profile`          | User Profile       | User's ideas, votes, settings                  | Yes               |
| `/admin`            | Admin Dashboard    | Categories, users, settings                    | Admin only        |
| `/admin/moderation` | Moderation Queue   | Flagged content review                         | Moderator+        |
| `/auth/login`       | Login              | Magic link or SSO redirect                     | No                |
| `/auth/verify`      | Email Verification | Magic link callback                            | No                |

### 4.4 Core UI Components (Fider-Inspired Architecture)

```
components/
├── layout/
│   ├── Header.tsx                 # Clean header: Logo, RSS, Dark Mode, Sign In
│   ├── HeroSection.tsx            # Homepage hero with tagline and CTA
│   ├── Footer.tsx                 # Minimal footer with "Powered by" attribution
│   └── Container.tsx              # Max-width container wrapper (1280px)
│
├── ideas/
│   ├── IdeaList.tsx               # Paginated list of idea cards
│   ├── IdeaCard.tsx               # Horizontal card with vote count, title, description
│   ├── IdeaDetail.tsx             # Full idea page with two-column layout
│   ├── IdeaForm.tsx               # Inline quick-submit or modal form
│   ├── IdeaFilters.tsx            # Filter and sort controls
│   ├── IdeaStatusBadge.tsx        # Colored status pills (Open, Planned, etc.)
│   └── IdeaHero.tsx               # Large idea title and metadata section
│
├── voting/
│   ├── VoteButton.tsx             # Primary blue action button with icon
│   ├── VoteCount.tsx              # Large vote number display
│   ├── VoteProgress.tsx           # Horizontal progress bar
│   ├── VoteSidebar.tsx            # Left sidebar with vote count
│   └── useRealTimeVotes.ts        # Hook for live vote updates
│
├── comments/
│   ├── CommentThread.tsx          # Flat comment list (no deep nesting)
│   ├── CommentItem.tsx            # Single comment with avatar, author, timestamp
│   ├── CommentForm.tsx            # Rich text editor with toolbar
│   ├── CommentEditor.tsx          # Markdown/WYSIWYG toggle editor
│   ├── CommentActions.tsx         # Reply, Edit, Delete menu (⋯)
│   └── CommentMeta.tsx            # Author badge (✓), timestamp, "edited"
│
├── input/
│   ├── SearchBar.tsx              # Top-right search input
│   ├── SuggestionInput.tsx        # Hero "Enter your suggestion here..." input
│   ├── RichTextToolbar.tsx        # Editor formatting toolbar
│   └── MarkdownEditor.tsx         # Plain markdown textarea
│
├── ui/                            # shadcn/ui base components
│   ├── button.tsx                 # Primary, ghost, outline variants
│   ├── card.tsx                   # White background, subtle shadow
│   ├── badge.tsx                  # Status pills and labels
│   ├── avatar.tsx                 # Circular user images
│   ├── dropdown-menu.tsx          # Comment actions, user menu
│   ├── dialog.tsx                 # Modals for forms
│   ├── input.tsx                  # Text inputs with focus styles
│   ├── textarea.tsx               # Multi-line inputs
│   ├── select.tsx                 # Sort dropdown
│   └── toast.tsx                  # Success/error notifications
│
├── meta/
│   ├── UserBadge.tsx              # Author name with verified checkmark
│   ├── Timestamp.tsx              # Relative time display ("5 months ago")
│   ├── MetaInfo.tsx               # Posted by, date, status row
│   └── AttributionFooter.tsx      # "Powered by Fider" (customizable)
│
└── navigation/
    ├── FilterBar.tsx              # [FILTER] [SORT BY: ▼] row
    ├── Breadcrumb.tsx             # "← Back to all suggestions"
    ├── SortDropdown.tsx           # Trending, Top, Recent, etc.
    └── ThemeToggle.tsx            # Dark/light mode switch
```

### 4.5 White-Label Branding System

All branding is controlled via environment variables and a JSON configuration file, enabling complete customization without code changes. The system supports Fider-style customization.

**Configuration File:** `public/config/branding.json`

```json
{
  "organizationName": "Your Organization",
  "productName": "BlazeBoard",
  "tagline": "Share. Vote. Innovate.",
  "hero": {
    "headline": "Help us build the best feedback platform",
    "description": "We're on a mission to build the best open source feedback platform, but we know this is not possible without your help. We need your feedback on how to improve it!",
    "callToAction": "Please don't use this as a test site - if you want to play around there is a demo site ✅"
  },
  "logo": {
    "url": "/assets/logo.svg",
    "height": "40px",
    "showInHeader": true
  },
  "colors": {
    "primary": "#1A5632",
    "accent": "#90D408",
    "interactive": "#3b82f6",
    "background": "#fafafa",
    "card": "#ffffff"
  },
  "fonts": {
    "heading": "Inter, -apple-system, BlinkMacSystemFont, Segoe UI, sans-serif",
    "body": "Inter, -apple-system, BlinkMacSystemFont, Segoe UI, sans-serif"
  },
  "features": {
    "anonymousViewing": true,
    "downvotesEnabled": false,
    "markdownComments": true,
    "showHeroSection": true,
    "showVoteDisplay": true,
    "inlineSuggestionInput": true,
    "commentFeed": true,
    "darkModeToggle": true,
    "rssFeeds": true
  },
  "ui": {
    "voteButtonStyle": "filled",
    "cardStyle": "minimal",
    "statusBadgeStyle": "pill",
    "layoutDensity": "comfortable"
  },
  "attribution": {
    "show": true,
    "text": "Powered by BlazeBoard",
    "url": "https://github.com/your-org/blazeboard"
  }
}
```

**CSS Variable Injection:**

```typescript
// lib/branding.ts
export function injectBrandingCSS(config: BrandingConfig) {
  return `
    :root {
      /* Colors */
      --color-primary: ${config.colors.primary};
      --color-accent: ${config.colors.accent};
      --color-interactive: ${config.colors.interactive};
      --color-background: ${config.colors.background};
      --color-card: ${config.colors.card};
      
      /* Derived colors (HSL for Tailwind) */
      --primary: ${hexToHSL(config.colors.primary)};
      --accent: ${hexToHSL(config.colors.accent)};
      --interactive: ${hexToHSL(config.colors.interactive)};
      
      /* Typography */
      --font-heading: ${config.fonts.heading};
      --font-body: ${config.fonts.body};
      
      /* Spacing (based on density) */
      --card-padding: ${
        config.ui.layoutDensity === 'compact' ? '1rem' : '1.5rem'
      };
      --section-gap: ${
        config.ui.layoutDensity === 'compact' ? '1.5rem' : '2rem'
      };
    }
  `;
}
```

**Docker Environment Overrides:**

```yaml
# Branding can be overridden via environment variables
environment:
  # Organization
  - BRANDING_ORG_NAME=UAB IT
  - BRANDING_PRODUCT_NAME=BlazeBoard
  - BRANDING_TAGLINE=Share. Vote. Innovate.

  # Hero Section
  - BRANDING_HERO_HEADLINE=Help us improve UAB IT services
  - BRANDING_HERO_DESCRIPTION=Share your ideas to help us build better technology solutions for the UAB community.

  # Colors
  - BRANDING_PRIMARY_COLOR=#1A5632
  - BRANDING_ACCENT_COLOR=#90D408
  - BRANDING_INTERACTIVE_COLOR=#3b82f6

  # Assets
  - BRANDING_LOGO_URL=/assets/uab-logo.svg

  # Features
  - BRANDING_SHOW_HERO=true
  - BRANDING_ENABLE_DARK_MODE=true
  - BRANDING_ENABLE_RSS=true

  # Attribution
  - BRANDING_ATTRIBUTION_TEXT=Powered by BlazeBoard
  - BRANDING_ATTRIBUTION_URL=https://github.com/uab-it/blazeboard
```

### 4.6 UAB Internal Branding (Fixed Configuration)

For UAB deployments, branding is locked to institutional standards:

| Element           | Value                                      |
| ----------------- | ------------------------------------------ |
| **Primary Color** | UAB Green `#1A5632` (HSL: `154 53% 22%`)   |
| **Accent Color**  | Campus Green `#90D408` (HSL: `80 92% 43%`) |
| **Warning Color** | UAB Gold `#FDB913`                         |
| **Typography**    | Aktiv Grotesk (primary), Inter (fallback)  |
| **Header Text**   | "BlazeBoard - UAB IT Ideas"                |
| **Tagline**       | "Share. Vote. Innovate."                   |
| **Logo**          | UAB IT logo (provided by UAB Marketing)    |

### 4.6 UAB Internal Branding (Fixed Configuration)

For UAB deployments, branding is locked to institutional standards with Fider-style presentation:

| Element                  | Value                                                                                                                                | Notes                                         |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------- |
| **Primary Color**        | UAB Green `#1A5632`                                                                                                                  | Main brand color for headings, links          |
| **Accent Color**         | Campus Green `#90D408`                                                                                                               | Secondary brand color for highlights          |
| **Interactive Color**    | Blue `#3b82f6`                                                                                                                       | Vote buttons, primary actions (Fider pattern) |
| **Warning/Alert Color**  | Loyal Yellow `#FDB913`                                                                                                               | Also called "UAB Gold"                        |
| **Neutral Color**        | Smoke Gray `#808285`                                                                                                                 | Secondary text, borders                       |
| **Background Tints**     | Campus Green 15% `#EEF9DA`, Smoke Gray 7% `#F5F5F5`                                                                                  | Light backgrounds for sections                |
| **Typography - Primary** | Aktiv Grotesk                                                                                                                        | Sans-serif for body text, headings, buttons   |
| **Typography - Accent**  | Kulturista                                                                                                                           | Serif for decorative headings (optional)      |
| **Base Font Size**       | 18px                                                                                                                                 | Body text (`p`, `li`)                         |
| **Hero Headline**        | "Help us improve UAB IT services"                                                                                                    |                                               |
| **Hero Description**     | "Share your ideas to help us build better technology solutions for the UAB community. Your feedback shapes the future of IT at UAB." |                                               |
| **Tagline**              | "Share. Vote. Innovate."                                                                                                             |                                               |
| **Logo**                 | UAB IT logo (SVG format)                                                                                                             | Provided by UAB Marketing                     |
| **Vote Button Style**    | Filled blue (#3b82f6) with white text                                                                                                | Fider pattern                                 |
| **Attribution**          | "Powered by BlazeBoard"                                                                                                              | Footer attribution                            |

**UAB Color Tints Available:**

- All brand colors support tints at 7%, 10%, 15%, 33%, 45%, and 50%
- Example: Campus Green 15% = `#EEF9DA` (light green background)
- Smoke Gray 7% = light gray for alternating sections

**UAB-Specific Features:**

- Dark mode toggle enabled
- RSS feeds for idea updates
- Inline suggestion input on homepage
- Comment feed with threading
- Status badges: Open (blue), Under Discussion (amber), Planned (purple), Completed (green), Declined (gray)

**Font Loading Notes:**

- Aktiv Grotesk loaded via UAB web font service
- Kulturista requires Safari compatibility fix when used
- Fallback fonts: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif

### 4.7 Responsive Design Breakpoints

````css
/* Tailwind default breakpoints used */
/* sm: 640px  - Mobile landscape */
/* md: 768px  - Tablets */
/* lg: 1024px - Desktop */
/* xl: 1280px - Large desktop */

### 4.7 Responsive Design Breakpoints

```css
/* Tailwind default breakpoints used */
/* sm: 640px  - Mobile landscape */
/* md: 768px  - Tablets */
/* lg: 1024px - Desktop */
/* xl: 1280px - Large desktop */
````

**Layout Behaviors (Fider-Inspired):**

| Breakpoint              | Layout                     | Components                                                                                                                                         |
| ----------------------- | -------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Mobile (<768px)**     | Single column              | - Stacked vote sidebar<br>- Full-width suggestion input<br>- Collapsed filters<br>- Mobile-optimized comment editor<br>- Bottom sticky vote button |
| **Tablet (768-1024px)** | Two-column                 | - Side-by-side content and sidebar<br>- Collapsible filter panel<br>- Inline editor toolbar                                                        |
| **Desktop (>1024px)**   | Two-column + fixed sidebar | - Left: Vote sidebar (120px fixed)<br>- Center: Main content (max 900px)<br>- Full rich text editor<br>- Sticky navigation                         |

**Fider-Specific Responsive Patterns:**

- Hero section scales from 48px headline on desktop to 32px on mobile
- Suggestion input transitions from inline to modal on mobile
- Vote button becomes sticky footer on mobile idea detail pages
- Comment actions (⋯) menu adapts from dropdown to slide-up sheet on mobile
- Filter/Sort bar stacks vertically on mobile

### 4.8 Accessibility Requirements (WCAG 2.1 AA)

| Requirement             | Implementation                                         |
| ----------------------- | ------------------------------------------------------ |
| **Keyboard Navigation** | Full tab navigation for voting, commenting, forms      |
| **Screen Readers**      | ARIA labels on vote buttons, status badges, forms      |
| **Color Contrast**      | Minimum 4.5:1 for text, 3:1 for UI components          |
| **Focus Indicators**    | Visible focus rings on all interactive elements        |
| **Reduced Motion**      | Respect `prefers-reduced-motion` media query           |
| **Alt Text**            | Required for all images, user avatars                  |
| **Form Labels**         | Associated labels for all form inputs                  |
| **Error Messages**      | Clear, descriptive error states with ARIA live regions |

### 4.9 Real-Time Updates Architecture

```typescript
// hooks/useRealTimeVotes.ts
// Leverages Postgres LISTEN/NOTIFY without external services

import { useEffect, useState } from 'react';

export function useRealTimeVotes(ideaId: string) {
  const [voteCount, setVoteCount] = useState<number>(0);

  useEffect(() => {
    // Server-Sent Events connection to Next.js API route
    const eventSource = new EventSource(
      `/api/ideas/${ideaId}/votes/stream`
    );

    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setVoteCount(data.voteCount);
    };

    return () => eventSource.close();
  }, [ideaId]);

  return voteCount;
}
```

```typescript
// app/api/ideas/[id]/votes/stream/route.ts
// SSE endpoint backed by Postgres LISTEN/NOTIFY

export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const encoder = new TextEncoder();

  const stream = new ReadableStream({
    async start(controller) {
      const client = await pool.connect();
      await client.query(`LISTEN vote_update`);

      client.on('notification', (msg) => {
        const payload = JSON.parse(msg.payload);
        if (payload.idea_id === params.id) {
          controller.enqueue(
            encoder.encode(`data: ${JSON.stringify(payload)}\n\n`)
          );
        }
      });
    },
  });

  return new Response(stream, {
    headers: {
      'Content-Type': 'text/event-stream',
      'Cache-Control': 'no-cache',
      Connection: 'keep-alive',
    },
  });
}
```

### 4.9 Static Asset Handling

**Asset Organization:**

```
public/
├── assets/
│   ├── logo.svg              # Default/placeholder logo
│   ├── logo-dark.svg         # Dark mode variant
│   └── favicon.ico
├── config/
│   └── branding.json         # Runtime branding config
└── fonts/
    ├── inter-var.woff2       # Default font
    └── aktiv-grotesk.woff2   # UAB font (if licensed)
```

**Docker Volume for Custom Assets:**

```yaml
# Mount custom branding assets
volumes:
  - ./custom-assets/logo.svg:/app/public/assets/logo.svg:ro
  - ./custom-assets/branding.json:/app/public/config/branding.json:ro
  - ./custom-assets/favicon.ico:/app/public/favicon.ico:ro
```

### 4.10 Frontend Build Optimization

**Next.js Production Build Settings:**

```javascript
// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',

  // Enable static optimization where possible
  experimental: {
    optimizeCss: true,
  },

  // Image optimization (disable if using external CDN)
  images: {
    unoptimized: process.env.DISABLE_IMAGE_OPTIMIZATION === 'true',
    remotePatterns: [{ protocol: 'https', hostname: '*.uab.edu' }],
  },

  // Bundle analyzer (dev only)
  ...(process.env.ANALYZE === 'true' && {
    webpack: (config, { isServer }) => {
      if (!isServer) {
        const {
          BundleAnalyzerPlugin,
        } = require('webpack-bundle-analyzer');
        config.plugins.push(new BundleAnalyzerPlugin());
      }
      return config;
    },
  }),
};

module.exports = nextConfig;
```

**Build Performance Targets:**

| Metric                         | Target               |
| ------------------------------ | -------------------- |
| Build time                     | < 3 minutes          |
| First Contentful Paint (FCP)   | < 1.5s               |
| Largest Contentful Paint (LCP) | < 2.5s               |
| Cumulative Layout Shift (CLS)  | < 0.1                |
| Time to Interactive (TTI)      | < 3.5s               |
| Bundle size (gzipped)          | < 150KB (initial JS) |

---

## 5. Docker Image Specifications

### 5.1 Application Image (`blazeboard-app`)

**Base Image:** `node:20-alpine`

**Rationale:**

- Alpine Linux minimizes attack surface (~5MB base)
- Node.js 20 LTS provides long-term support through April 2026
- Compatible with Next.js 14+ requirements

**Build Stages:**

```dockerfile
# Stage 1: Dependencies
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --only=production

# Stage 2: Builder
FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
ENV NEXT_TELEMETRY_DISABLED=1
RUN npm run build

# Stage 3: Production Runner
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:3000/api/health || exit 1

CMD ["node", "server.js"]
```

**Image Size Target:** < 200MB (production image)

### 5.2 Next.js Configuration for Standalone Output

```javascript
// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  experimental: {
    // Optimize for container deployments
    outputFileTracingRoot: require('path').join(__dirname, '../'),
  },
  // Disable image optimization if using external CDN
  images: {
    unoptimized: process.env.DISABLE_IMAGE_OPTIMIZATION === 'true',
  },
};

module.exports = nextConfig;
```

### 5.3 Nginx Reverse Proxy Image

**Base Image:** `nginx:1.25-alpine`

**Configuration Template:**

```nginx
# nginx/nginx.conf
upstream blazeboard {
    server blazeboard-app:3000;
    keepalive 32;
}

server {
    listen 80;
    server_name _;

    # Redirect HTTP to HTTPS in production
    location / {
        return 301 https://$host$request_uri;
    }

    # Health check endpoint (no redirect)
    location /health {
        return 200 'OK';
        add_header Content-Type text/plain;
    }
}

server {
    listen 443 ssl http2;
    server_name _;

    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Gzip compression
    gzip on;
    gzip_types text/plain text/css application/json application/javascript text/xml;
    gzip_min_length 1000;

    location / {
        proxy_pass http://blazeboard;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Static assets with long cache
    location /_next/static {
        proxy_pass http://blazeboard;
        proxy_cache_valid 200 365d;
        add_header Cache-Control "public, max-age=31536000, immutable";
    }
}
```

---

## 6. Docker Compose Configurations

### 6.1 Development Configuration

```yaml
# docker-compose.dev.yml
version: '3.9'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
      target: deps # Use deps stage for hot reload
    container_name: blazeboard-dev
    ports:
      - '3000:3000'
    volumes:
      - .:/app
      - /app/node_modules
      - /app/.next
    environment:
      - NODE_ENV=development
      - DATABASE_URL=postgresql://blazeboard:devpassword@db:5432/blazeboard
      - NEXTAUTH_URL=http://localhost:3000
      - NEXTAUTH_SECRET=dev-secret-change-in-production
    depends_on:
      db:
        condition: service_healthy
    networks:
      - blazeboard-network

  db:
    image: postgres:16-alpine
    container_name: blazeboard-db-dev
    environment:
      POSTGRES_USER: blazeboard
      POSTGRES_PASSWORD: devpassword
      POSTGRES_DB: blazeboard
    volumes:
      - postgres-dev-data:/var/lib/postgresql/data
      - ./db/init:/docker-entrypoint-initdb.d:ro
    ports:
      - '5432:5432'
    healthcheck:
      test: ['CMD-SHELL', 'pg_isready -U blazeboard -d blazeboard']
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - blazeboard-network

volumes:
  postgres-dev-data:

networks:
  blazeboard-network:
    driver: bridge
```

### 6.2 Production Configuration

```yaml
# docker-compose.prod.yml
version: '3.9'

services:
  app:
    image: ${REGISTRY:-ghcr.io}/blazeboard/blazeboard:${TAG:-latest}
    container_name: blazeboard-app
    restart: unless-stopped
    expose:
      - '3000'
    environment:
      - NODE_ENV=production
      - DATABASE_URL=${DATABASE_URL}
      - NEXTAUTH_URL=${NEXTAUTH_URL}
      - NEXTAUTH_SECRET=${NEXTAUTH_SECRET}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASS=${SMTP_PASS}
    healthcheck:
      test:
        [
          'CMD',
          'wget',
          '--no-verbose',
          '--tries=1',
          '--spider',
          'http://localhost:3000/api/health',
        ]
      interval: 30s
      timeout: 10s
      start_period: 40s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    networks:
      - blazeboard-network
    logging:
      driver: json-file
      options:
        max-size: '10m'
        max-file: '3'

  nginx:
    image: nginx:1.25-alpine
    container_name: blazeboard-nginx
    restart: unless-stopped
    ports:
      - '80:80'
      - '443:443'
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      app:
        condition: service_healthy
    healthcheck:
      test:
        [
          'CMD',
          'wget',
          '--no-verbose',
          '--tries=1',
          '--spider',
          'http://localhost/health',
        ]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - blazeboard-network
    logging:
      driver: json-file
      options:
        max-size: '10m'
        max-file: '3'

  # Optional: Include if using containerized PostgreSQL
  db:
    image: postgres:16-alpine
    container_name: blazeboard-db
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./db/backups:/backups
    healthcheck:
      test: ['CMD-SHELL', 'pg_isready -U ${DB_USER} -d ${DB_NAME}']
      interval: 10s
      timeout: 5s
      retries: 5
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
    networks:
      - blazeboard-network
    # Not exposed externally - accessed only via internal network

volumes:
  postgres-data:
    driver: local

networks:
  blazeboard-network:
    driver: bridge
```

### 6.3 UAB Internal Configuration

```yaml
# docker-compose.uab.yml
version: '3.9'

services:
  app:
    image: ${UAB_REGISTRY}/blazeboard:${TAG:-latest}
    container_name: blazeboard-uab
    restart: unless-stopped
    expose:
      - '3000'
    environment:
      - NODE_ENV=production
      - DATABASE_URL=postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=require
      - NEXTAUTH_URL=https://spark.it.uab.edu
      - NEXTAUTH_SECRET=${NEXTAUTH_SECRET}
      # UAB SAML Configuration
      - SAML_ISSUER=https://spark.it.uab.edu
      - SAML_ENTRY_POINT=https://padlock.idm.uab.edu/cas/idp/profile/SAML2/Redirect/SSO
      - SAML_IDP_CERT=${UAB_IDP_CERT}
      # UAB SMTP
      - SMTP_HOST=smtp.office365.com
      - SMTP_PORT=587
      - SMTP_USER=${UAB_SMTP_USER}
      - SMTP_PASS=${UAB_SMTP_PASS}
      # UAB Branding
      - BRANDING_ORG_NAME=UAB IT
      - BRANDING_PRODUCT_NAME=BlazeBoard
      - BRANDING_PRIMARY_COLOR=#1A5632
    healthcheck:
      test:
        [
          'CMD',
          'wget',
          '--no-verbose',
          '--tries=1',
          '--spider',
          'http://localhost:3000/api/health',
        ]
      interval: 30s
      timeout: 10s
      start_period: 40s
      retries: 3
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '0.5'
          memory: 512M
    networks:
      - blazeboard-uab-network
    logging:
      driver: json-file
      options:
        max-size: '50m'
        max-file: '5'

  nginx:
    image: nginx:1.25-alpine
    container_name: blazeboard-nginx-uab
    restart: unless-stopped
    ports:
      - '80:80'
      - '443:443'
    volumes:
      - ./nginx/nginx.uab.conf:/etc/nginx/nginx.conf:ro
      - /etc/ssl/uab:/etc/nginx/ssl:ro # UAB SSL certificates
    depends_on:
      app:
        condition: service_healthy
    networks:
      - blazeboard-uab-network

networks:
  blazeboard-uab-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.28.0.0/16 # Isolated network for UAB deployment
```

---

## 7. Environment Configuration

### 7.1 Required Environment Variables

| Variable          | Description                        | Example                               | Required |
| ----------------- | ---------------------------------- | ------------------------------------- | -------- |
| `DATABASE_URL`    | PostgreSQL connection string       | `postgresql://user:pass@host:5432/db` | Yes      |
| `NEXTAUTH_URL`    | Application URL                    | `https://blazeboard.example.com`      | Yes      |
| `NEXTAUTH_SECRET` | Session encryption key (32+ chars) | `openssl rand -base64 32`             | Yes      |
| `NODE_ENV`        | Environment mode                   | `production`                          | Yes      |

### 7.2 Optional Environment Variables

| Variable                 | Description                          | Default      |
| ------------------------ | ------------------------------------ | ------------ |
| `SMTP_HOST`              | Email server hostname                | -            |
| `SMTP_PORT`              | Email server port                    | `587`        |
| `SMTP_USER`              | Email authentication user            | -            |
| `SMTP_PASS`              | Email authentication password        | -            |
| `BRANDING_ORG_NAME`      | Organization name for white-labeling | `BlazeBoard` |
| `BRANDING_PRIMARY_COLOR` | Primary theme color (hex)            | `#1A5632`    |
| `LOG_LEVEL`              | Application log verbosity            | `info`       |

### 7.3 Environment File Template

```bash
# .env.production.template
# Database
DATABASE_URL=postgresql://blazeboard:CHANGE_ME@db:5432/blazeboard?sslmode=require
DB_USER=blazeboard
DB_PASSWORD=CHANGE_ME
DB_NAME=blazeboard
DB_HOST=db

# Authentication
NEXTAUTH_URL=https://your-domain.com
NEXTAUTH_SECRET=generate-with-openssl-rand-base64-32

# Email (Optional)
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=noreply@example.com
SMTP_PASS=CHANGE_ME

# Branding (Optional)
BRANDING_ORG_NAME=Your Organization
BRANDING_PRODUCT_NAME=BlazeBoard
BRANDING_PRIMARY_COLOR=#1A5632

# Logging
LOG_LEVEL=info
```

---

## 8. Health Checks & Monitoring

### 8.1 Application Health Endpoint

```typescript
// app/api/health/route.ts
import { NextResponse } from 'next/server';
import { db } from '@/lib/db';

export async function GET() {
  const health = {
    status: 'healthy',
    timestamp: new Date().toISOString(),
    version: process.env.npm_package_version || 'unknown',
    checks: {
      database: 'unknown',
      memory: 'unknown',
    },
  };

  try {
    // Database connectivity check
    await db.$queryRaw`SELECT 1`;
    health.checks.database = 'healthy';
  } catch (error) {
    health.checks.database = 'unhealthy';
    health.status = 'unhealthy';
  }

  // Memory check
  const used = process.memoryUsage();
  const heapUsedMB = Math.round(used.heapUsed / 1024 / 1024);
  health.checks.memory = heapUsedMB < 400 ? 'healthy' : 'warning';

  const statusCode = health.status === 'healthy' ? 200 : 503;
  return NextResponse.json(health, { status: statusCode });
}
```

### 8.2 Docker Health Check Configuration

```yaml
# Health check parameters
healthcheck:
  test:
    [
      'CMD',
      'wget',
      '--no-verbose',
      '--tries=1',
      '--spider',
      'http://localhost:3000/api/health',
    ]
  interval: 30s # Check every 30 seconds
  timeout: 10s # Fail if no response in 10 seconds
  start_period: 40s # Grace period for container startup
  retries: 3 # Mark unhealthy after 3 consecutive failures
```

### 8.3 Monitoring Integration (Optional)

```yaml
# Prometheus metrics endpoint (add to app service)
labels:
  prometheus.io/scrape: 'true'
  prometheus.io/port: '3000'
  prometheus.io/path: '/api/metrics'
```

---

## 9. Security Requirements

### 9.1 Container Security

| Requirement          | Implementation                             |
| -------------------- | ------------------------------------------ |
| Non-root user        | Run as `nextjs` user (UID 1001)            |
| Read-only filesystem | Mount app code as read-only where possible |
| No privileged mode   | Never use `--privileged` flag              |
| Resource limits      | Set CPU/memory limits in compose           |
| Secret management    | Use Docker secrets or external vault       |

### 9.2 Network Security

| Requirement            | Implementation                             |
| ---------------------- | ------------------------------------------ |
| Internal-only database | Database container not exposed to host     |
| TLS 1.2+ only          | Nginx configured with modern cipher suites |
| Security headers       | X-Frame-Options, CSP, HSTS via Nginx       |
| Rate limiting          | Nginx rate limiting on auth endpoints      |

### 9.3 Secrets Management

**Option 1: Docker Secrets (Swarm Mode)**

```yaml
secrets:
  db_password:
    external: true
  nextauth_secret:
    external: true

services:
  app:
    secrets:
      - db_password
      - nextauth_secret
```

**Option 2: External Secret Management**

```bash
# Using HashiCorp Vault or Azure Key Vault
export DATABASE_URL=$(vault kv get -field=url secret/blazeboard/database)
docker-compose up -d
```

---

## 10. Build & Deployment Pipeline

### 10.1 CI/CD Workflow (GitHub Actions)

```yaml
# .github/workflows/docker-build.yml
name: Build and Push Docker Image

on:
  push:
    branches: [main]
    tags: ['v*']
  pull_request:
    branches: [main]

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Container Registry
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=ref,event=branch
            type=semver,pattern={{version}}
            type=sha,prefix=

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          push: ${{ github.event_name != 'pull_request' }}
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Run Trivy vulnerability scan
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
          exit-code: '1'
          severity: 'CRITICAL,HIGH'
```

### 10.2 UAB GitLab Container Registry (Alternative)

For UAB-internal deployments, BlazeBoard can be hosted in the UAB GitLab Container Registry for secure, on-premises container storage.

**Setup Steps:**

1. Create a UAB GitLab account ([UAB GitLab Account Setup](https://docs.rc.uab.edu/account_management/gitlab_account))
2. Create a new project in UAB GitLab
3. Create an access token:
   - Navigate to **Edit Profile** → **Access Tokens**
   - Token name: `blazeboard-container-registry`
   - Expiration: 90 days (recommended)
   - Scopes: Check `read_registry` and `write_registry`
   - Click **Create personal access token** and save the token

**Login to UAB GitLab Registry:**

```bash
# Using access token (recommended)
sudo docker login code.rc.uab.edu:4567 -u <your-blazerid> -p <access-token>
```

**Tag and Push BlazeBoard Image:**

```bash
# Tag the image for UAB GitLab registry
sudo docker tag blazeboard-app:latest code.rc.uab.edu:4567/<gitlab-group>/<project-name>/blazeboard-app:latest

# Push to UAB GitLab registry
sudo docker push code.rc.uab.edu:4567/<gitlab-group>/<project-name>/blazeboard-app:latest
```

**Reference:** [UAB RC Containers Documentation](https://docs.rc.uab.edu/workflow_solutions/getting_containers)

---

## 10.3 RC Cloud Deployment Guide

BlazeBoard is deployed on UAB Research Computing Cloud (RC Cloud) using OpenStack virtual machines with Docker installed. This section provides step-by-step deployment instructions.

### 10.3.1 Prerequisites

- **RC Cloud Account**: Request access at [UAB RC Support](https://docs.rc.uab.edu/help/support) with intended use case
- **UAB Campus Network/VPN**: RC Cloud requires UAB Campus Network or VPN ([UAB VPN Setup](https://www.uab.edu/it/home/tech-solutions/network/vpn))
- **SSH Client**: Required for instance access ([SSH Client Setup](https://docs.rc.uab.edu/uab_cloud/remote_access))

**Reference:** [RC Cloud Documentation](https://docs.rc.uab.edu/uab_cloud/)

### 10.3.2 Instance Creation Workflow

Follow the complete [RC Cloud Tutorial](https://docs.rc.uab.edu/uab_cloud/tutorial/) in order to set up all required components.

#### Step 1: Network Setup

Create network infrastructure before launching instances. See [Network Setup Guide](https://docs.rc.uab.edu/uab_cloud/tutorial/networks).

1. Navigate to **Network** → **Networks** → **Create Network**
2. Create subnet with appropriate CIDR (e.g., `192.168.1.0/24`)
3. Create router and attach to subnet
4. Create floating IP for external access

#### Step 2: Security Configuration

Configure SSH access and security groups. See [Security Setup Guide](https://docs.rc.uab.edu/uab_cloud/tutorial/security).

**Create Key Pair:**

1. Navigate to **Compute** → **Key Pairs** → **Create Key Pair**
2. Name: `blazeboard-key`
3. Download private key (`.pem` file) to `~/.ssh/`
4. Set permissions: `chmod 400 ~/.ssh/blazeboard-key.pem`

**Create Security Groups:**

| Security Group | Purpose      | Rules                                     |
| -------------- | ------------ | ----------------------------------------- |
| `ssh`          | SSH access   | Ingress TCP 22 from UAB network           |
| `http`         | HTTP access  | Ingress TCP 80 from anywhere (0.0.0.0/0)  |
| `https`        | HTTPS access | Ingress TCP 443 from anywhere (0.0.0.0/0) |

#### Step 3: Launch Instance

Create a virtual machine instance. See [Instance Setup Guide](https://docs.rc.uab.edu/uab_cloud/tutorial/instances).

**Instance Configuration:**

| Field                 | Value                             | Notes                                                                              |
| --------------------- | --------------------------------- | ---------------------------------------------------------------------------------- |
| **Instance Name**     | `blazeboard-prod`                 | Follow [naming conventions](https://docs.rc.uab.edu/uab_cloud/#naming-conventions) |
| **Description**       | BlazeBoard production server      |                                                                                    |
| **Availability Zone** | nova                              |                                                                                    |
| **Source**            | Image: Ubuntu 22.04 LTS           | Other options: Ubuntu 20.04                                                        |
| **Volume Size**       | 40 GB                             | Minimum 20 GB, 40 GB recommended                                                   |
| **Delete Volume**     | No                                | Keep OS volume for reuse                                                           |
| **Flavor**            | `m1.medium` or larger             | 2 vCPU, 4 GB RAM minimum                                                           |
| **Network**           | Your created network              | From Step 1                                                                        |
| **Security Groups**   | `default`, `ssh`, `http`, `https` | Multiple groups allowed                                                            |
| **Key Pair**          | `blazeboard-key`                  | From Step 2                                                                        |

**Launch Instance:**

1. Complete all tabs in the Launch Instance dialog
2. Click **Launch Instance**
3. Wait for Status to show **Active** (typically 2-3 minutes)
4. Associate floating IP from Step 1

#### Step 4: SSH Connection

Access your instance via SSH. See [Remote Access Guide](https://docs.rc.uab.edu/uab_cloud/remote_access).

```bash
# Connect to instance (replace <floating-ip> with your assigned IP)
ssh ubuntu@<floating-ip> -i ~/.ssh/blazeboard-key.pem

# Optional: Set up SSH config for easier access
cat >> ~/.ssh/config << EOF
Host blazeboard
  HostName <floating-ip>
  User ubuntu
  IdentityFile ~/.ssh/blazeboard-key.pem
EOF

# Then connect with: ssh blazeboard
```

### 10.3.3 Docker Installation on RC Cloud Instance

Install Docker on your Ubuntu instance. See [Docker Installation Guide](https://docs.rc.uab.edu/workflow_solutions/getting_containers#docker-installation-on-uab-rc-cloud).

```bash
# Update system packages
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
sudo apt install -y docker.io docker-compose

# Add user to docker group (optional, avoids needing sudo)
sudo usermod -aG docker ubuntu

# Log out and back in for group changes to take effect
exit
# Then reconnect via SSH

# Verify Docker installation
docker --version
docker-compose --version
```

### 10.3.4 Deploy BlazeBoard on RC Cloud

#### Option A: Using Docker Compose (Recommended)

**1. Clone Repository:**

```bash
cd ~
git clone https://github.com/your-org/blazeboard.git
cd blazeboard
```

**2. Configure Environment:**

```bash
# Copy environment template
cp .env.template .env

# Edit environment variables
nano .env
```

**Required Environment Variables:**

```bash
# Database
DATABASE_URL=postgresql://blazeboard:SECURE_PASSWORD@db:5432/blazeboard

# NextAuth
NEXTAUTH_URL=https://<your-floating-ip-or-domain>
NEXTAUTH_SECRET=$(openssl rand -base64 32)

# UAB SAML (if using UAB SSO)
SAML_ENTITY_ID=https://<your-domain>/api/auth
SAML_SSO_URL=https://websso.auth.uab.edu/idp/profile/SAML2/Redirect/SSO
SAML_CERT_PATH=/path/to/uab-idp-cert.pem

# Application
NODE_ENV=production
PORT=3000
```

**3. Deploy Services:**

```bash
# Production deployment
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Check status
docker-compose -f docker-compose.prod.yml ps
```

#### Option B: Using UAB GitLab Container Registry

**1. Pull from UAB GitLab:**

```bash
# Login to UAB GitLab registry
docker login code.rc.uab.edu:4567

# Pull image
docker pull code.rc.uab.edu:4567/<gitlab-group>/blazeboard/blazeboard-app:latest
```

**2. Create docker-compose.yml:**

```yaml
version: '3.8'

services:
  app:
    image: code.rc.uab.edu:4567/<gitlab-group>/blazeboard/blazeboard-app:latest
    container_name: blazeboard-app
    restart: unless-stopped
    ports:
      - '3000:3000'
    env_file: .env
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    container_name: blazeboard-db
    restart: unless-stopped
    environment:
      POSTGRES_USER: blazeboard
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: blazeboard
    volumes:
      - postgres_data:/var/lib/postgresql/data

  nginx:
    image: nginx:1.25-alpine
    container_name: blazeboard-nginx
    restart: unless-stopped
    ports:
      - '80:80'
      - '443:443'
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - app

volumes:
  postgres_data:
```

**3. Deploy:**

```bash
docker-compose up -d
```

### 10.3.5 SSL/TLS Certificate Setup

**Option 1: Let's Encrypt (Recommended for Production)**

```bash
# Install Certbot
sudo apt install -y certbot python3-certbot-nginx

# Obtain certificate (replace with your domain)
sudo certbot --nginx -d blazeboard.rc.uab.edu

# Certificates auto-renew via cron
sudo certbot renew --dry-run
```

**Option 2: Self-Signed Certificate (Development/Testing)**

```bash
# Create directory
mkdir -p nginx/ssl

# Generate self-signed certificate
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout nginx/ssl/privkey.pem \
  -out nginx/ssl/fullchain.pem \
  -subj "/C=US/ST=Alabama/L=Birmingham/O=UAB/CN=<floating-ip>"
```

### 10.3.6 Firewall Configuration

RC Cloud security groups act as the firewall. Ensure the following are configured in your security groups:

| Port | Protocol | Source         | Purpose                   |
| ---- | -------- | -------------- | ------------------------- |
| 22   | TCP      | UAB Network    | SSH access                |
| 80   | TCP      | 0.0.0.0/0      | HTTP (redirects to HTTPS) |
| 443  | TCP      | 0.0.0.0/0      | HTTPS for web access      |
| 3000 | TCP      | localhost only | Application (internal)    |
| 5432 | TCP      | localhost only | PostgreSQL (internal)     |

### 10.3.7 Post-Deployment Verification

```bash
# Check container status
docker ps

# Check application health
curl http://localhost:3000/api/health

# Check database connection
docker exec blazeboard-db psql -U blazeboard -c "SELECT 1"

# Check logs
docker logs blazeboard-app
docker logs blazeboard-nginx

# Access application
# Navigate to: https://<floating-ip> or https://<your-domain>
```

### 10.3.8 Backup and Maintenance

**Database Backups:**

```bash
# Manual backup
docker exec blazeboard-db pg_dump -U blazeboard blazeboard | gzip > backup_$(date +%Y%m%d).sql.gz

# Automated daily backups (add to crontab)
crontab -e
# Add line:
# 0 2 * * * /home/ubuntu/blazeboard/scripts/backup.sh
```

**System Updates:**

```bash
# Update Docker images
docker-compose pull
docker-compose up -d

# Update Ubuntu packages
sudo apt update && sudo apt upgrade -y
sudo reboot  # If kernel updates are installed
```

**Reference:** [RC Cloud Tutorial](https://docs.rc.uab.edu/uab_cloud/tutorial/)

---

### 10.4 Deployment Commands

```bash
# Development
docker-compose -f docker-compose.dev.yml up --build

# Production (new deployment)
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d

# Production (zero-downtime update)
docker-compose -f docker-compose.prod.yml pull
docker-compose -f docker-compose.prod.yml up -d --no-deps --build app
docker-compose -f docker-compose.prod.yml exec nginx nginx -s reload

# UAB Internal
docker-compose -f docker-compose.uab.yml --env-file .env.uab up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f app

# Database backup
docker-compose -f docker-compose.prod.yml exec db pg_dump -U blazeboard blazeboard > backup.sql
```

---

## 11. Database Initialization

### 11.1 Init Script Structure

```
db/
├── init/
│   ├── 01-extensions.sql      # Enable required extensions
│   ├── 02-schema.sql          # Create tables
│   ├── 03-rls-policies.sql    # Row-level security
│   └── 04-seed.sql            # Initial data (optional)
└── migrations/
    └── *.sql                  # Version-controlled migrations
```

### 11.2 Example Init Scripts

```sql
-- db/init/01-extensions.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

```sql
-- db/init/02-schema.sql
-- Users table
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  blazerid TEXT UNIQUE,
  email TEXT UNIQUE NOT NULL,
  first_name TEXT,
  last_name TEXT,
  display_name TEXT NOT NULL,
  department TEXT,
  role TEXT DEFAULT 'contributor'
    CHECK (role IN ('viewer','contributor','moderator','admin')),
  is_banned BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Ideas table
CREATE TABLE IF NOT EXISTS ideas (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  idea_number SERIAL UNIQUE,
  author_id UUID REFERENCES users(id) NOT NULL,
  title TEXT NOT NULL CHECK (char_length(title) <= 200),
  body TEXT NOT NULL CHECK (char_length(body) <= 2000),
  impact_statement TEXT,
  status TEXT DEFAULT 'new'
    CHECK (status IN ('new','under_review','planned','completed','declined')),
  tags TEXT[],
  attributes JSONB DEFAULT '{}',
  vote_count INTEGER DEFAULT 0,
  comment_count INTEGER DEFAULT 0,
  is_deleted BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Additional tables: votes, comments, flags
-- (See full schema in main PRD Section 6)
```

---

## 12. Backup & Recovery

### 12.1 Automated Backup Script

```bash
#!/bin/bash
# scripts/backup.sh

set -e

BACKUP_DIR="/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="${BACKUP_DIR}/blazeboard_${TIMESTAMP}.sql.gz"

# Create backup
docker-compose -f docker-compose.prod.yml exec -T db \
  pg_dump -U ${DB_USER} ${DB_NAME} | gzip > ${BACKUP_FILE}

# Retain only last 30 days
find ${BACKUP_DIR} -name "*.sql.gz" -mtime +30 -delete

echo "Backup completed: ${BACKUP_FILE}"
```

### 12.2 Restore Procedure

```bash
#!/bin/bash
# scripts/restore.sh

BACKUP_FILE=$1

if [ -z "$BACKUP_FILE" ]; then
  echo "Usage: ./restore.sh <backup-file.sql.gz>"
  exit 1
fi

# Stop application (prevent writes)
docker-compose -f docker-compose.prod.yml stop app

# Restore database
gunzip -c ${BACKUP_FILE} | docker-compose -f docker-compose.prod.yml exec -T db \
  psql -U ${DB_USER} ${DB_NAME}

# Restart application
docker-compose -f docker-compose.prod.yml start app

echo "Restore completed from: ${BACKUP_FILE}"
```

### 12.3 Backup Schedule (cron)

```bash
# /etc/cron.d/blazeboard-backup
# Daily backup at 2 AM
0 2 * * * root /opt/blazeboard/scripts/backup.sh >> /var/log/blazeboard-backup.log 2>&1
```

---

## 13. Resource Requirements

### 13.1 Minimum Requirements

| Component         | CPU            | Memory     | Storage     |
| ----------------- | -------------- | ---------- | ----------- |
| Application       | 0.25 cores     | 256 MB     | 500 MB      |
| Nginx             | 0.1 cores      | 64 MB      | 50 MB       |
| PostgreSQL        | 0.5 cores      | 512 MB     | 5 GB        |
| **Total Minimum** | **0.85 cores** | **832 MB** | **5.55 GB** |

### 13.2 Recommended Production Requirements

| Component             | CPU                 | Memory         | Storage    |
| --------------------- | ------------------- | -------------- | ---------- |
| Application           | 1-2 cores           | 512 MB - 1 GB  | 1 GB       |
| Nginx                 | 0.25 cores          | 128 MB         | 100 MB     |
| PostgreSQL            | 1-2 cores           | 1-2 GB         | 20+ GB     |
| **Total Recommended** | **2.25-4.25 cores** | **1.6-3.1 GB** | **21+ GB** |

### 13.3 Scaling Considerations

- **Horizontal scaling**: Add app container replicas behind Nginx load balancer
- **Database scaling**: Use managed PostgreSQL (Azure, AWS RDS) for production
- **CDN**: Offload static assets to CDN for high-traffic deployments

### 13.4 RC Cloud Instance Flavors

For RC Cloud deployments, select instance flavors based on expected load:

| Deployment Type         | Recommended Flavor | vCPU | RAM   | Notes                              |
| ----------------------- | ------------------ | ---- | ----- | ---------------------------------- |
| **Development/Testing** | m1.small           | 1    | 2 GB  | Minimum viable for testing         |
| **Small Production**    | m1.medium          | 2    | 4 GB  | Recommended minimum for production |
| **Medium Production**   | m1.large           | 4    | 8 GB  | Handles moderate traffic           |
| **Large Production**    | m1.xlarge          | 8    | 16 GB | High traffic, multiple containers  |

**Storage Recommendations:**

- **OS Volume**: 40 GB (minimum 20 GB)
- **Data Volume** (optional): Create separate persistent volume for database if data exceeds 20 GB
- See [Volumes Guide](https://docs.rc.uab.edu/uab_cloud/tutorial/volumes) for persistent storage

**Note**: For databases, consider using managed services (Azure PostgreSQL, Supabase) rather than self-hosted on RC Cloud for production deployments requiring high availability.

---

## 14. Troubleshooting Guide

### 14.1 Common Issues

| Issue                      | Symptoms                  | Resolution                                     |
| -------------------------- | ------------------------- | ---------------------------------------------- |
| Container won't start      | Exit code 1, no logs      | Check `docker logs blazeboard-app` for errors  |
| Database connection failed | "ECONNREFUSED" errors     | Verify DATABASE_URL, check db container health |
| SSL certificate errors     | Browser security warnings | Verify cert paths in nginx volumes             |
| Out of memory              | Container restarts        | Increase memory limits in compose              |
| Slow performance           | High response times       | Check resource usage with `docker stats`       |

### 14.2 Debug Commands

```bash
# View container status
docker-compose -f docker-compose.prod.yml ps

# View application logs
docker-compose -f docker-compose.prod.yml logs -f app

# Shell into container
docker-compose -f docker-compose.prod.yml exec app sh

# Check resource usage
docker stats

# Inspect container
docker inspect blazeboard-app

# Test database connection
docker-compose -f docker-compose.prod.yml exec db psql -U blazeboard -c "SELECT 1"

# Check nginx config
docker-compose -f docker-compose.prod.yml exec nginx nginx -t
```

### 14.3 RC Cloud Specific Issues

| Issue                                   | Symptoms                                    | Resolution                                                                                                                                                                                               |
| --------------------------------------- | ------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Cannot access RC Cloud dashboard        | Connection timeout                          | Ensure connected to UAB Campus Network or VPN ([UAB VPN Setup](https://www.uab.edu/it/home/tech-solutions/network/vpn))                                                                                  |
| SSH connection refused                  | "Connection refused" or timeout             | 1. Check security group has SSH (port 22) rule<br>2. Verify floating IP is correctly associated<br>3. Check instance status is "Active"                                                                  |
| Instance fails to start                 | Status shows "Error"                        | Check instance fault details on instance page, contact [RC Support](https://docs.rc.uab.edu/help/support) with Instance ID                                                                               |
| Remote Host Identification Changed      | SSH error after recreating instance         | Run `ssh-keygen -R <floating-ip>` to remove old fingerprint ([Guide](https://docs.rc.uab.edu/uab_cloud/remote_access#remove-an-invalid-host-fingerprint))                                                |
| Docker permission denied                | "permission denied" when running docker     | 1. Add user to docker group: `sudo usermod -aG docker ubuntu`<br>2. Log out and back in<br>3. Verify with `groups` command                                                                               |
| Cannot pull from GitLab registry        | Authentication failed                       | 1. Check access token is valid<br>2. Ensure token has `read_registry` scope<br>3. Re-login: `docker login code.rc.uab.edu:4567`                                                                          |
| Application not accessible from browser | Timeout or connection refused               | 1. Check security groups allow HTTP (80) and HTTPS (443)<br>2. Verify nginx container is running: `docker ps`<br>3. Check nginx logs: `docker logs blazeboard-nginx`                                     |
| Database out of disk space              | Application crashes, "no space left" errors | 1. Check volume usage: `df -h`<br>2. Clean old logs: `docker system prune -a`<br>3. Increase volume size via RC Cloud dashboard                                                                          |
| Instance stuck during deletion          | Instance shows "deleting" for >5 minutes    | Contact [RC Support](https://docs.rc.uab.edu/help/support) with Instance ID from Overview tab ([Where is my Instance ID?](https://docs.rc.uab.edu/uab_cloud/tutorial/instances#where-is-my-instance-id)) |

**Common RC Cloud Commands:**

```bash
# Get Instance ID (needed for support requests)
# Navigate to: Compute → Instances → Click instance name → Overview tab → Copy ID

# Check RC Cloud instance resources
free -h              # Memory usage
df -h                # Disk usage
top                  # CPU usage
docker stats         # Container resource usage

# Check network connectivity
ping 8.8.8.8         # Test internet connectivity
curl ifconfig.me     # Check external IP
ss -tuln             # List open ports

# Firewall (security groups handle this, but for reference)
sudo ufw status      # Check firewall status (should be inactive on RC Cloud)
```

**Getting Help:**

- **RC Cloud Documentation**: [https://docs.rc.uab.edu/uab_cloud/](https://docs.rc.uab.edu/uab_cloud/)
- **RC Support**: Email [support@listserv.uab.edu](mailto:support@listserv.uab.edu)
- **Support Page**: [https://docs.rc.uab.edu/help/support](https://docs.rc.uab.edu/help/support)

---

## 15. Success Criteria

### 15.1 Deployment Checklist

**General Deployment:**

- [ ] Docker images build successfully (< 5 minutes)
- [ ] Production image size < 200 MB
- [ ] Health checks pass within 60 seconds of startup
- [ ] Application responds on port 443 with valid SSL
- [ ] Database migrations complete without errors
- [ ] Environment variables properly injected
- [ ] Logs accessible via `docker-compose logs`
- [ ] Backup script executes successfully
- [ ] Container restarts automatically on failure

**RC Cloud Specific:**

- [ ] RC Cloud account created and access confirmed
- [ ] Connected to UAB Campus Network or VPN
- [ ] Network, subnet, and router configured
- [ ] Floating IP created and associated
- [ ] Security groups configured (ssh, http, https)
- [ ] SSH key pair created and downloaded
- [ ] Instance created with appropriate flavor (m1.medium or larger)
- [ ] SSH connection successful to instance
- [ ] Docker and docker-compose installed on instance
- [ ] Application accessible via floating IP
- [ ] SSL certificate configured (Let's Encrypt or self-signed)
- [ ] Backup automation configured via cron

### 15.2 Performance Benchmarks

| Metric                | Target       |
| --------------------- | ------------ |
| Cold start time       | < 30 seconds |
| Health check response | < 500 ms     |
| API response (p95)    | < 200 ms     |
| Memory usage (idle)   | < 300 MB     |
| Container image pull  | < 60 seconds |

---

## Appendix A: File Structure

```
blazeboard/
├── Dockerfile
├── docker-compose.dev.yml
├── docker-compose.prod.yml
├── docker-compose.uab.yml
├── .dockerignore
├── .env.template
├── nginx/
│   ├── nginx.conf
│   ├── nginx.uab.conf
│   └── ssl/
│       ├── fullchain.pem
│       └── privkey.pem
├── db/
│   ├── init/
│   │   ├── 01-extensions.sql
│   │   ├── 02-schema.sql
│   │   └── 03-rls-policies.sql
│   └── backups/
├── scripts/
│   ├── backup.sh
│   ├── restore.sh
│   └── healthcheck.sh
└── .github/
    └── workflows/
        └── docker-build.yml
```

## Appendix B: .dockerignore

```
# .dockerignore
node_modules
.next
.git
.gitignore
*.md
!README.md
.env*
!.env.template
docker-compose*.yml
Dockerfile*
.dockerignore
coverage
.nyc_output
*.log
.DS_Store
```

---

_Document Owner: BlazeBoard Development Team_  
_Review Cycle: As needed with infrastructure changes_
