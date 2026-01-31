# UAB Design System for Blazeboard

## Overview

Blazeboard uses the University of Alabama at Birmingham (UAB) visual identity system. This document describes how UAB branding is implemented in the codebase.

## Colors

### Primary Colors

| Color | Hex | CSS Variable | Usage |
|-------|-----|--------------|-------|
| UAB Green | `#1A5632` | `--uab-green` | Primary brand color, headings, buttons |
| Campus Green | `#90D408` | `--campus-green` | Accents, hover states, icons |
| Loyal Yellow | `#FDB913` | `--loyal-yellow` | Highlights, warnings, accents |
| Smoke Gray | `#808285` | `--smoke-gray` | Body text, borders |

### Tint Colors

Colors have tint variations at 7%, 15%, 33%, and 50% opacity:
- `--campus-green-7`, `--campus-green-15`, etc.
- `--smoke-gray-7`, `--smoke-gray-15`, etc.

### Dark Mode

Dark mode adjusts colors for contrast:
- UAB Green: Lightened 20%
- Campus Green: Stays vibrant
- Loyal Yellow: Slightly muted

## Typography

### Fonts

- **Primary**: Aktiv Grotesk (body text, UI)
- **Accent**: Kulturista (h1, h2, decorative headings)

### Scale

| Element | Font | Size | Weight | Color |
|---------|------|------|--------|-------|
| H1 | Kulturista | 42-48px | Bold | UAB Green |
| H2 | Kulturista | 28-36px | Bold | UAB Green |
| H3-H6 | Aktiv Grotesk | 16-24px | Semi-bold | UAB Green |
| Body | Aktiv Grotesk | 18px | Regular | Smoke Gray |

### Loading Fonts

Fonts are loaded from Adobe Fonts:
```html
<link rel="stylesheet" href="https://use.typekit.net/[KIT_ID].css">
```

Safari requires the `.font-kulturista` class for proper rendering.

## Components

### Buttons

**Primary**: UAB Green background, white text
**Secondary**: White background, UAB Green border
**Link**: Campus Green text

### Forms

**Focus state**: 2px Campus Green border with subtle glow
**Error state**: Red border

### Cards

**Base**: White background, rounded corners, shadow
**Accent**: 4px left border rotating through UAB colors

### Tags

Pill-shaped with 15% tint backgrounds, rotating colors

## Utility Classes

### Color Utilities

```scss
.uab-green           // UAB Green text
.bg-uab-green        // UAB Green background
.campus-green        // Campus Green text
.bg-campus-green-15  // Campus Green 15% tint background
```

### Grid Utilities

```scss
.uab-grid-two-across    // 2-column responsive grid
.uab-grid-three-across  // 3-column responsive grid
.uab-grid-four-across   // 4-column responsive grid
```

### Card Utilities

```scss
.uab-card               // Base card style
.uab-card-accent        // Card with left border
.border-campus-green    // Campus Green border
.border-loyal-yellow    // Loyal Yellow border
```

## Accessibility

### Color Contrast

All color combinations meet WCAG AA standards:
- UAB Green on white: 8.35:1 (AAA)
- Smoke Gray on white: 4.54:1 (AA)

**Note**: Campus Green (2.8:1) is only used for accents, not body text.

### Focus States

All interactive elements have 2px Campus Green focus outlines with 2px offset.

### Motion

The design respects `prefers-reduced-motion` preference.

## File Structure

### Variables
- `public/assets/styles/variables/_colors.scss` - UAB colors
- `public/assets/styles/variables/_text.scss` - Typography
- `public/assets/styles/variables/_dark-colors.scss` - Dark mode

### Utilities
- `public/assets/styles/utility/uab.scss` - UAB-specific classes

## Resources

- [UAB Brand Guidelines](https://www.uab.edu/toolkit/branding)
- Design document: `docs/plans/2026-01-31-blazeboard-uab-redesign.md`
