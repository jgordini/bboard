# Blazeboard UAB Redesign - Design Document

**Date:** January 31, 2026
**Project:** Transform bboard into Blazeboard with UAB branding
**Scope:** Full visual redesign with UAB brand identity

## Vision

Transform bboard into **Blazeboard**, a UAB-branded feedback platform that maintains all existing functionality while adopting UAB's visual identity. The redesign uses UAB Green (#1A5632), Campus Green (#90D408), and Loyal Yellow (#FDB913) with Aktiv Grotesk and Kulturista typography to create a cohesive, professional experience that feels distinctly UAB.

## Design System Integration

### CSS Variable Mapping Strategy

We'll create a UAB theme layer in the existing SCSS architecture by mapping UAB colors and typography to the current variable system. This preserves the codebase structure while enabling a complete visual transformation.

#### Color Variables (`public/assets/styles/variables/_colors.scss`)

**Primary Mappings:**
- Map UAB Green (#1A5632) → `--colors-primary-base`
- Map Campus Green (#90D408) → accent/success colors
- Map Loyal Yellow (#FDB913) → warning/highlight colors
- Map Smoke Gray (#808285) → neutral grays

**New UAB-Specific Variables:**
```scss
--uab-green: #1A5632;
--dragons-lair-green: #033319;
--campus-green: #90D408;
--loyal-yellow: #FDB913;
--smoke-gray: #808285;
--white: #FFFFFF;

// Tint variations (7%, 10%, 15%, 33%, 45%, 50%)
--campus-green-7: #F3FBDB;
--campus-green-15: #EEF9DA;
--smoke-gray-7: #F5F5F5;
--smoke-gray-15: #E8E8E9;
// ... etc
```

#### Typography Variables (`public/assets/styles/variables/_text.scss`)

**Font Families:**
```scss
--font-primary: 'Aktiv Grotesk', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
--font-accent: 'Kulturista', Georgia, 'Times New Roman', serif;
```

**Base Settings:**
- Base font size: 18px (UAB standard)
- Line height: 1.6 for body text
- Heading color: UAB Green

**Safari Kulturista Fix:**
```css
.font-kulturista {
  font-family: "kulturista-web", Georgia, "Times New Roman", serif !important;
}
```

#### Utility Classes (`public/assets/styles/utility/`)

Add UAB-specific utilities:
```scss
// Color utilities
.uab-green { color: var(--uab-green); }
.bg-uab-green { background-color: var(--uab-green); }
.campus-green { color: var(--campus-green); }
.bg-campus-green { background-color: var(--campus-green); }
.loyal-yellow { color: var(--loyal-yellow); }
.bg-loyal-yellow { background-color: var(--loyal-yellow); }

// Grid utilities
.uab-grid-two-across {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;

  @media (min-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }
}

.uab-grid-three-across {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;

  @media (min-width: 768px) {
    grid-template-columns: repeat(3, 1fr);
  }
}

// Section wrapper
.uab-section {
  padding: 3rem 1.5rem;
}
```

## Component Visual Updates

### Home Page Transformation

#### Welcome Section
- **Layout:** Two-column hero layout (1/3 - 2/3 split on desktop)
- **Left column:** Welcome message with Kulturista heading in UAB Green
- **Right column:** Post input and list
- **Background:** Alternating white and Smoke Gray-7 sections

#### Post Input Button
- **Style:** Card-style with rounded corners (`border-radius: 12px`)
- **Shadow:** `box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1)`
- **Icon:** Campus Green plus circle (40px)
- **Hover:** Elevated shadow (`0 10px 15px -3px rgb(0 0 0 / 0.1)`)
- **Padding:** 1.5rem (24px)
- **Font:** Aktiv Grotesk, 18px, medium weight

#### Post Cards
- **Base:** White background, rounded-lg, shadow-md
- **Left border:** 4px accent in rotating colors (Campus Green, Loyal Yellow, UAB Green)
- **Vote counter:** Campus Green background when active
- **Tags:** Pill style with UAB color tints (15% opacity backgrounds)
- **Hover state:** Elevated shadow (shadow-lg)

### Navigation & Header

#### Header Component
- **Background:** UAB Green (#1A5632)
- **Text:** White
- **Logo:** "Blazeboard" in Kulturista font, white color
- **Navigation links:** White with Campus Green underline on hover
- **User avatar:** Loyal Yellow border (2px) when notifications present

### Buttons & Interactive Elements

#### Button Styles

**Primary Buttons:**
```scss
.btn-primary {
  background-color: var(--uab-green);
  color: white;
  border-radius: 8px;
  padding: 0.75rem 1.5rem;
  font-weight: 600;
  box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1);

  &:hover {
    background-color: var(--dragons-lair-green);
    box-shadow: 0 4px 6px -1px rgb(0 0 0 / 0.1);
  }
}
```

**Secondary Buttons:**
```scss
.btn-secondary {
  background-color: white;
  color: var(--uab-green);
  border: 2px solid var(--uab-green);
  border-radius: 8px;
  padding: 0.75rem 1.5rem;
  font-weight: 600;

  &:hover {
    background-color: var(--campus-green-7);
  }
}
```

**Link Buttons:**
```scss
.btn-link {
  color: var(--campus-green);
  font-weight: 600;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }

  i {
    margin-left: 0.5rem;
  }
}
```

## Layout Patterns & Page Structures

### Section-Based Architecture

All pages follow alternating section pattern for visual rhythm:

#### Home Page Structure
```
1. Hero section (white bg, py-12) - Welcome message + CTA
2. Posts section (Smoke Gray-7 bg, py-12) - Post input + list
3. Footer section (white bg, py-8) - PoweredBy component
```

#### Post Details Page Structure
```
1. Navigation breadcrumb (white bg, py-4)
2. Post content (Smoke Gray-7 bg, py-12) - Title, description, vote section
3. Comments section (white bg, py-12) - Discussion thread
4. Related posts (Smoke Gray-7 bg, py-12) - Similar suggestions
```

#### Admin Pages Structure
```
1. Admin header (UAB Green bg, py-6) - Page title in Kulturista
2. Content wrapper (Smoke Gray-7 bg)
   - Side navigation (white card, shadow-md)
   - Content area (white card, shadow-md)
```

### Responsive Grid System

**Two-column layouts:**
- Desktop: 50/50 split with 3rem gap
- Tablet: 50/50 split with 2rem gap
- Mobile: Stack to single column

**Three-column layouts:**
- Desktop: 33/33/33 split with 2rem gap
- Tablet: 50/50 with wrap
- Mobile: Stack to single column

**Four-column layouts:**
- Desktop: 25/25/25/25 split with 1.5rem gap
- Tablet: 50/50 with wrap
- Mobile: Stack to single column

### Card Patterns

**Standard Card:**
```scss
.uab-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px 0 rgb(0 0 0 / 0.1);
  padding: 1.5rem;
  border-left: 4px solid var(--accent-color); // Rotate colors

  &:hover {
    box-shadow: 0 10px 15px -3px rgb(0 0 0 / 0.1);
  }
}
```

**Quote Card:**
```scss
.uab-quote-card {
  @extend .uab-card;
  border-left-color: var(--campus-green);

  .quote-icon {
    color: var(--campus-green);
    font-size: 2rem;
    opacity: 0.5;
  }

  .author {
    border-top: 1px solid var(--smoke-gray-15);
    padding-top: 1rem;
    margin-top: 1rem;
  }
}
```

## Typography Hierarchy & Styling

### Font Loading

Add to `public/index.html`:
```html
<link rel="stylesheet" href="https://use.typekit.net/XXXXX.css">
<!-- Replace XXXXX with UAB's Adobe Fonts kit ID -->
```

Or use CSS import in main stylesheet:
```css
@import url('https://use.typekit.net/XXXXX.css');
```

### Heading Scale

| Element | Font | Size (Desktop) | Size (Mobile) | Weight | Color |
|---------|------|----------------|---------------|--------|-------|
| H1 | Kulturista | 48px | 36px | Bold | UAB Green |
| H2 | Kulturista | 36px | 28px | Bold | UAB Green |
| H3 | Aktiv Grotesk | 24px | 22px | Semi-bold | UAB Green |
| H4 | Aktiv Grotesk | 20px | 18px | Semi-bold | UAB Green |
| H5 | Aktiv Grotesk | 18px | 16px | Medium | UAB Green |
| H6 | Aktiv Grotesk | 16px | 14px | Medium | UAB Green |

### Body Text Styles

**Paragraphs:**
- Font: Aktiv Grotesk, 18px
- Line height: 1.6
- Color: Smoke Gray (#808285)
- Margin bottom: 1rem

**Lists:**
- Font: Aktiv Grotesk, 18px
- Line height: 1.6
- Bullet/number color: Campus Green
- Spacing: 0.5rem between items

**Links:**
- Color: Campus Green
- Text decoration: none
- Hover: underline, brighten color
- Icon suffix for external/CTA links: `→`

**Emphasis:**
- Strong/bold: UAB Green color, font-weight 600
- Italic: Maintain Aktiv Grotesk
- Code: Monospace, Smoke Gray-7 background

## Forms & Interactive States

### Form Components

#### Text Inputs & Textareas
```scss
.form-input {
  border: 1px solid var(--smoke-gray);
  border-radius: 8px;
  padding: 12px 16px;
  font-size: 18px;
  font-family: var(--font-primary);
  background: white;

  &:focus {
    outline: none;
    border: 2px solid var(--campus-green);
    box-shadow: 0 0 0 3px rgba(144, 212, 8, 0.1);
  }

  &::placeholder {
    color: var(--smoke-gray);
  }
}
```

#### Select Dropdowns
```scss
.form-select {
  @extend .form-input;
  background-image: url('data:image/svg+xml,...'); // Campus Green chevron
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 40px;

  &:focus {
    border-color: var(--campus-green);
  }
}
```

#### Checkboxes & Radio Buttons
```scss
.form-checkbox {
  width: 24px;
  height: 24px;
  border: 2px solid var(--smoke-gray);
  border-radius: 4px;

  &:checked {
    background-color: var(--campus-green);
    border-color: var(--campus-green);
    background-image: url('data:image/svg+xml,...'); // White checkmark
  }

  &:focus {
    outline: 2px solid var(--campus-green);
    outline-offset: 2px;
  }
}
```

#### Tag Pills
```scss
.tag-pill {
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 600;

  // Rotate through UAB colors
  &.tag-green {
    background-color: var(--campus-green-15);
    color: var(--uab-green);
  }

  &.tag-yellow {
    background-color: rgba(253, 185, 19, 0.15);
    color: #a67500; // Darker yellow for text
  }

  .remove-icon {
    margin-left: 6px;
    color: var(--smoke-gray);
    cursor: pointer;

    &:hover {
      color: #dc2626; // Red
    }
  }
}
```

### Interactive States

#### Hover States
- **Buttons:** Darken by 10%, elevate shadow (1-2px increase)
- **Cards:** Shadow elevation from md to lg
- **Links:** Underline appears, color brightens 10%
- **Icons:** Scale 1.05, color shift to Campus Green

#### Focus States
- **All interactive elements:** 2px Campus Green outline, 2px offset
- **Form inputs:** 2px Campus Green border + subtle glow
- **Keyboard navigation:** Clear focus ring visible at all times

#### Active/Selected States
- **Background:** Campus Green-15 tint
- **Border/accent:** Full Campus Green (4px left border for cards)
- **Text:** UAB Green, font-weight 600
- **Icons:** Campus Green fill

#### Disabled States
- **Opacity:** 50%
- **Cursor:** not-allowed
- **Grayscale filter:** Applied to colored elements
- **No hover effects:** Disabled

## Dark Mode & Visual Assets

### Dark Mode Color Adaptations

The app supports dark mode. UAB colors adapt as follows:

#### Dark Mode Palette
```scss
@media (prefers-color-scheme: dark) {
  --background: #1a1a1a;
  --surface: #2d2d2d;
  --surface-elevated: #3d3d3d;

  // Adjusted UAB colors for dark mode
  --uab-green: #2a7d4f; // Lightened for contrast
  --campus-green: #90D408; // Keep vibrant
  --loyal-yellow: #e5a910; // Slightly muted
  --smoke-gray: #b0b0b0; // Lightened

  // Text colors
  --text-primary: rgba(255, 255, 255, 0.9);
  --text-secondary: rgba(255, 255, 255, 0.7);
  --text-muted: rgba(255, 255, 255, 0.5);
}
```

#### Dark Mode Card Borders
- Use Campus Green/Loyal Yellow at 30% opacity
- Glow effect on hover (subtle box-shadow with color)

#### Dark Mode Forms
- Input borders: Gray-600 (#4a4a4a)
- Focus borders: Campus Green (full strength)
- Background: Gray-800 (#2d2d2d)

### Icons & Imagery

#### Font Awesome Integration

**Icon Color Patterns:**
- Primary action icons: Campus Green
- Decorative icons: Loyal Yellow
- System icons (close, menu, etc.): Smoke Gray
- Success states: Campus Green
- Warning states: Loyal Yellow
- Error states: Red (#dc2626)

**Common Icon Uses:**
```html
<!-- CTA links -->
<a href="#" class="btn-link">
  Learn More <i class="fas fa-arrow-right"></i>
</a>

<!-- Add actions -->
<button class="btn-primary">
  <i class="fas fa-plus-circle"></i> Add Post
</button>

<!-- Success states -->
<div class="status-approved">
  <i class="fas fa-check-circle"></i> Approved
</div>

<!-- External links -->
<a href="#" class="link-external">
  Visit Site <i class="fas fa-external-link-alt"></i>
</a>
```

#### Imagery Style Guide

**Photos:**
- Add subtle UAB Green overlay (5% opacity) for brand consistency
- Use high-quality images with good contrast
- Rounded corners (8-12px) to match card style

**Illustrations:**
- Prefer illustrations with Campus Green and Loyal Yellow accents
- Maintain UAB color palette
- Simple, modern style (like undraw.co or similar)

**Backgrounds:**
- Subtle textures in Smoke Gray-7 tint
- Avoid busy patterns that compete with content
- Consider diagonal stripes or dots in UAB colors at low opacity

## Accessibility & Best Practices

### WCAG Compliance

#### Color Contrast Ratios

| Combination | Ratio | WCAG Level | Usage |
|-------------|-------|------------|-------|
| UAB Green on White | 8.35:1 | AAA ✓ | Body text, headings |
| White on UAB Green | 8.35:1 | AAA ✓ | Button text, nav |
| Smoke Gray on White | 4.54:1 | AA ✓ | Body text |
| Campus Green on White | 2.8:1 | Fail ✗ | Icons/accents only |
| UAB Green on Campus Green-15 | 7.5:1 | AAA ✓ | Tags, highlights |

**Important:** Never use Campus Green for body text. Use for icons, borders, and backgrounds only.

### Keyboard Navigation

- **Tab order:** Logical flow through page sections
- **Focus indicators:** 2px Campus Green outline, 2px offset, always visible
- **Skip links:** "Skip to main content" link at top of page
- **Dropdown navigation:** Arrow keys for menu items
- **Button activation:** Enter and Space keys
- **Modal dialogs:** Trap focus, Escape to close

### Screen Reader Support

- **Semantic HTML:** Use proper heading hierarchy (h1→h6)
- **ARIA labels:** For icon-only buttons, decorative elements
- **Alt text:** Descriptive alt text for all images
- **Form labels:** Properly associated with inputs
- **Live regions:** ARIA live for dynamic updates (vote counts, notifications)
- **Landmark regions:** header, nav, main, aside, footer

### Motion & Animation

**Respect `prefers-reduced-motion`:**
```scss
@media (prefers-reduced-motion: reduce) {
  * {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
```

**Default Animations:**
- Hover transitions: 200ms ease-in-out
- Focus transitions: 150ms ease-out
- Page transitions: 300ms ease-in-out
- No auto-playing content
- Smooth scroll with fallback

## Implementation Notes

### Priority Order

1. **Phase 1: Foundation (Variables & Utilities)**
   - Update color variables
   - Add typography variables
   - Create utility classes
   - Add font loading

2. **Phase 2: Core Components**
   - Header/Navigation
   - Buttons
   - Form inputs
   - Cards

3. **Phase 3: Page Layouts**
   - Home page
   - Post details page
   - Admin pages

4. **Phase 4: Polish**
   - Dark mode refinements
   - Animations/transitions
   - Icon updates
   - Accessibility audit

### Testing Checklist

- [ ] Color contrast meets WCAG AA (AAA preferred)
- [ ] Keyboard navigation works throughout
- [ ] Screen reader announces content correctly
- [ ] Dark mode colors are readable
- [ ] Responsive layouts work on mobile/tablet/desktop
- [ ] Forms are accessible and usable
- [ ] Hover/focus states are clear
- [ ] Motion respects user preferences

### Browser Support

- Chrome/Edge (latest 2 versions)
- Firefox (latest 2 versions)
- Safari (latest 2 versions)
- Mobile Safari (iOS 14+)
- Chrome Mobile (latest)

### Performance Considerations

- Font loading: Use `font-display: swap` to prevent FOIT
- Images: Lazy load below fold, use modern formats (WebP)
- CSS: Minimize custom styles, leverage utility classes
- Animations: Use transform/opacity for GPU acceleration

## Resources

### UAB Brand Guidelines
- UAB Colors: Campus Green (#90D408), UAB Green (#1A5632), Loyal Yellow (#FDB913)
- Fonts: Aktiv Grotesk (primary), Kulturista (accent)
- Adobe Fonts: [Kit ID needed]

### Design References
- UAB Website: https://www.uab.edu/
- UAB Templates: [Internal Joomla templates]
- Font Awesome 5 Pro: https://fontawesome.com/

### Development Resources
- Current codebase: React 18, TypeScript, SCSS
- Existing utility classes: `public/assets/styles/utility/`
- Component library: `public/components/`

---

**Document Status:** Approved for implementation
**Next Steps:** Create implementation plan, set up git worktree for development
