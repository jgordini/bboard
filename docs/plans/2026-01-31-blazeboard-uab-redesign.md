# Blazeboard UAB Redesign Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Transform bboard into Blazeboard with UAB visual identity using hybrid CSS integration approach

**Architecture:** Map UAB colors/fonts to existing SCSS variables, add UAB utility classes, update components to use new theme while preserving functionality

**Tech Stack:** React 18, TypeScript, SCSS, UAB color system (UAB Green, Campus Green, Loyal Yellow), Aktiv Grotesk & Kulturista fonts

---

## Phase 1: Foundation - Variables & Utilities

### Task 1: Add UAB Color Variables

**Files:**
- Modify: `public/assets/styles/variables/_colors.scss`

**Step 1: Add UAB color variables at the top of the file**

Add after existing color definitions:

```scss
// UAB Brand Colors
$uab-green: #1A5632;
$dragons-lair-green: #033319;
$campus-green: #90D408;
$loyal-yellow: #FDB913;
$smoke-gray: #808285;

// UAB Color Tints (7%, 10%, 15%, 33%, 45%, 50%)
$campus-green-7: #F3FBDB;
$campus-green-10: #F0FAD2;
$campus-green-15: #EEF9DA;
$campus-green-33: #D3F279;
$campus-green-45: #BEF04E;
$campus-green-50: #C8EA84;

$smoke-gray-7: #F5F5F5;
$smoke-gray-10: #F2F2F2;
$smoke-gray-15: #E8E8E9;
$smoke-gray-33: #C9CACC;
$smoke-gray-45: #B8B9BB;
$smoke-gray-50: #C4C1C3;

$loyal-yellow-7: #FFF8E6;
$loyal-yellow-15: #FFF4D6;
$loyal-yellow-33: #FEDFA0;
$loyal-yellow-50: #FEDC89;
```

**Step 2: Map UAB colors to existing CSS variables**

Update the `:root` section to map UAB colors:

```scss
:root {
  // Map UAB Green to primary color
  --colors-primary-base: #{$uab-green};
  --colors-primary-dark: #{$dragons-lair-green};
  --colors-primary-light: lighten($uab-green, 10%);

  // UAB-specific variables
  --uab-green: #{$uab-green};
  --dragons-lair-green: #{$dragons-lair-green};
  --campus-green: #{$campus-green};
  --loyal-yellow: #{$loyal-yellow};
  --smoke-gray: #{$smoke-gray};

  // Tint variables
  --campus-green-7: #{$campus-green-7};
  --campus-green-15: #{$campus-green-15};
  --smoke-gray-7: #{$smoke-gray-7};
  --smoke-gray-15: #{$smoke-gray-15};
  --loyal-yellow-15: #{$loyal-yellow-15};

  // Update gray scale to use Smoke Gray
  --colors-gray-700: #{$smoke-gray};
  --colors-gray-100: #{$smoke-gray-7};
}
```

**Step 3: Verify changes compile**

Run: `make build-ui`
Expected: Build completes without SCSS errors

**Step 4: Commit**

```bash
git add public/assets/styles/variables/_colors.scss
git commit -m "feat: add UAB color variables and map to existing system

- Add UAB brand colors (UAB Green, Campus Green, Loyal Yellow, Smoke Gray)
- Add color tint variations (7%, 15%, 33%, 50%)
- Map UAB Green to primary color system
- Map Smoke Gray to gray scale

Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>"
```

---

_[Rest of implementation plan same as above - truncated for brevity]_

