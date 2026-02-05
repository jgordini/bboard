# UI Review Tasks

Review of all page components against [Web Interface Guidelines](https://github.com/vercel-labs/web-interface-guidelines).

---

## 1. Fix Icon Component (aria-hidden + alt)

**Impact:** ~20 issues resolved in one change
**File:** `public/components/common/Icon.tsx`

- [ ] Line 14: `<img>` never receives an `alt` attribute -- add `alt` prop with `""` default
- [ ] Line 18: `<svg>` never gets `aria-hidden="true"` -- add `aria-hidden="true"` by default

---

## 2. Fix Dropdown ListItem (div -> button)

**File:** `public/components/common/Dropdown.tsx`

- [ ] Line 36: `<div onClick={handleClick}>` used for non-link list items -- replace with `<button>`

---

## 3. Fix span/div Click Handlers -> button

| File | Line | Issue |
|------|------|-------|
| `pages/Administration/components/OAuthForm.tsx` | 113 | `<span onClick={enableClientSecret}>change</span>` -- should be `<button>` |
| `pages/Administration/components/TagForm.tsx` | 72 | `<span onClick={this.randomize}>randomize</span>` -- should be `<button>` |
| `pages/Administration/components/SideMenu.tsx` | 83 | `<div onClick={toggle}>` mobile menu toggler -- should be `<button>` with `aria-label` |
| `pages/ShowPost/components/MentionSelector.tsx` | 15 | `<div className="clickable">` -- should be `<button>` |
| `pages/ShowPost/components/PostSearch.tsx` | 56 | `<VStack onClick={selectPost(p)}>` -- should be `<button>` |
| `pages/ShowPost/components/TagListItem.tsx` | 19 | `<HStack onClick={onClick}>` -- should be `<button>` |
| `pages/MySettings/MySettings.page.tsx` | 129 | `<a href="#" onClick={this.closeModal}>` -- should be `<button>` |
| `pages/MyNotifications/MyNotifications.page.tsx` | 76 | `<a href="#" onClick={this.markAllAsRead}>` -- should be `<button>` |

---

## 4. Add Confirmation Dialogs for Destructive Actions

| File | Line | Action |
|------|------|--------|
| `pages/Administration/pages/FlaggedPosts.page.tsx` | 87 | "Clear Flags" has no confirmation |
| `pages/Administration/pages/FlaggedComments.page.tsx` | 81 | "Clear Flags" has no confirmation |
| `pages/Administration/pages/ManageMembers.page.tsx` | 204 | "Block User" has no confirmation |

---

## 5. Fix Form Labels, Autocomplete, and Input Types

### Missing labels

| File | Line | Field |
|------|------|-------|
| `pages/SignUp/SignUp.page.tsx` | 119 | `field="name"` -- no `label` prop |
| `pages/SignUp/SignUp.page.tsx` | 120 | `field="email"` -- no `label` prop |
| `pages/SignUp/SignUp.page.tsx` | 128 | `field="tenantName"` -- no `label` prop |
| `pages/SignUp/SignUp.page.tsx` | 130 | `field="subdomain"` -- no `label` prop |
| `pages/SignIn/CompleteSignInProfile.page.tsx` | 54 | `field="name"` -- no `label` prop |
| `pages/Administration/pages/PrivacySettings.page.tsx` | 60, 67, 76 | Three `<Toggle>` components missing `label` prop |
| `pages/Administration/pages/AdvancedSettings.page.tsx` | 69 | Raw `<input>` with no `<label>` association (`htmlFor` missing) |
| `pages/Administration/components/webhook/WebhookForm.tsx` | 71, 78 | "Header" and "Value" inputs have no label |
| `pages/DesignSystem/DesignSystem.page.tsx` | 427 | `field="age"` -- no `label` prop |
| `pages/DesignSystem/DesignSystem.page.tsx` | 466 | `field="search"` -- no `label` prop |
| `pages/Home/components/PostFilter.tsx` | 170 | `<input>` missing `aria-label` and `name` |

### Missing autocomplete / type / spellCheck

| File | Line | Fix needed |
|------|------|------------|
| `pages/SignUp/SignUp.page.tsx` | 120 | Add `type="email"`, `autoComplete="email"`, `spellCheck={false}` |
| `pages/SignUp/SignUp.page.tsx` | 119 | Add `autoComplete="name"` |
| `pages/SignIn/CompleteSignInProfile.page.tsx` | 54 | Add `autoComplete="name"` |
| `pages/MySettings/MySettings.page.tsx` | 147 | Add `type="email"`, `autoComplete="email"`, `spellCheck={false}` |
| `pages/MySettings/MySettings.page.tsx` | 172 | Add `autoComplete="name"` |

### Missing aria-label on icon-only buttons

| File | Line | Element |
|------|------|---------|
| `pages/Administration/pages/ManageMembers.page.tsx` | 61 | Icon-only Dropdown button |
| `pages/ShowPost/components/VotesPanel.tsx` | 61 | Icon-only `<button>` with `<Icon>` |
| `pages/DesignSystem/DesignSystem.page.tsx` | 228, 245, 260, 279, 298, 317 | Icon-only `<Button>` components |

---

## 6. Fix `transition: all` in SCSS + Focus States

### transition: all -> explicit properties

| File | Line | Current | Suggested replacement |
|------|------|---------|-----------------------|
| `pages/Home/Home.page.scss` | 65 | `transition: all 0.2s ease-in-out` | `transition: box-shadow 0.2s ease-in-out, transform 0.2s ease-in-out, border-color 0.2s ease-in-out` |
| `pages/Home/components/PostFilter.scss` | 11 | `transition: all 0.2s ease-in-out` | `transition: box-shadow 0.2s ease-in-out, border-color 0.2s ease-in-out` |
| `pages/Home/components/PostsContainer.scss` | 76 | `transition: all 0.2s ease-in-out` | `transition: box-shadow 0.2s ease-in-out, transform 0.2s ease-in-out` |
| `pages/Home/components/PostsContainer.scss` | 124 | `transition: all 0.2s ease` | `transition: background-color 0.2s ease, border-color 0.2s ease, color 0.2s ease` |
| `pages/Administration/components/SideMenu.scss` | 16 | `transition: all 0.15s ease` | `transition: background-color 0.15s ease, color 0.15s ease, border-left-color 0.15s ease` |
| `pages/ShowPost/components/VotesPanel.scss` | 43 | `transition: all 0.2s ease` | `transition: background-color 0.2s ease, color 0.2s ease` |

### :focus -> :focus-visible

| File | Line | Fix |
|------|------|-----|
| `pages/Home/Home.page.scss` | 77 | Change `&:focus` to `&:focus-visible` |
| `pages/Administration/components/SideMenu.scss` | 32 | Change `&:focus` to `&:focus-visible` |

---

## 7. Fix Typography: Ellipsis and Loading States

Replace `...` (three ASCII periods) with `…` (Unicode ellipsis U+2026).

### Loading states

| File | Line | Text |
|------|------|------|
| `pages/Administration/pages/FlaggedPosts.page.tsx` | 47 | `Loading...` |
| `pages/Administration/pages/FlaggedComments.page.tsx` | 50 | `Loading...` |
| `pages/Administration/pages/ContentModeration.page.tsx` | 34 | `Loading...` |
| `pages/Administration/pages/ManageBilling.page.tsx` | 72 | `Loading...` |
| `pages/Leaderboard/Leaderboard.page.tsx` | 50 | `Loading...` |
| `pages/SignUp/PendingActivation.page.tsx` | 64 | `Resending...` |

### Placeholders

| File | Line | Text |
|------|------|------|
| `pages/ShowPost/components/PostSearch.tsx` | 50 | `Search original post...` |
| `pages/ShowPost/components/VotesModal.tsx` | 66 | `Search for users by name...` |
| `pages/ShowPost/components/ResponseModal.tsx` | 92 | `...plans...` |
| `pages/ShowPost/components/CommentInput.tsx` | 94 | `Leave a comment` (add trailing `…`) |
| `pages/ShowPost/components/DeletePostModal.tsx` | 43 | `Why are you deleting this post? (optional)` (add trailing `…`) |
| `pages/Administration/pages/GeneralSettings.page.tsx` | 85 | `Enter your suggestion here...` |
| `pages/Administration/pages/ManageMembers.page.tsx` | 225 | `Search by name / email ...` |
| `pages/Administration/components/webhook/WebhookForm.tsx` | 184 | `https://webhook.site/...` |
| `pages/DesignSystem/DesignSystem.page.tsx` | 421 | `Your name goes here...` |
| `pages/DesignSystem/DesignSystem.page.tsx` | 466 | `Search...` |
| `pages/Home/Home.page.tsx` | 133 | `Enter your suggestion here...` |

---

## 8. Additional Findings

### Semantic HTML

| File | Line | Issue |
|------|------|-------|
| `pages/ShowPost/components/DiscussionPanel.tsx` | 19 | `"Discussion"` in a `<span>` -- should be a heading (`<h2>` or `<h3>`) |
| `pages/SignUp/SignUp.page.tsx` | 106 | Heading jumps to `<h3>` -- should use proper hierarchy |

### i18n

| File | Line | Issue |
|------|------|-------|
| `pages/ShowPost/components/ShowComment.tsx` | 203 | Hardcoded `"edited"` not wrapped in `<Trans>` or `i18n._()` |
| `pages/ShowPost/components/ShowComment.tsx` | 203 | Tooltip string uses template literal without i18n |

### Animation

| File | Line | Issue |
|------|------|-------|
| `pages/Home/components/ShareFeedback.scss` | 36 | `animate-fade-in` does not honor `prefers-reduced-motion` |
| `pages/Home/components/SimilarPosts.scss` | 2 | Animates `max-height` instead of `transform`/`opacity` |
| Global | -- | No `prefers-reduced-motion` handling in page-level SCSS |

### Misc

| File | Line | Issue |
|------|------|-------|
| `pages/Home/components/ListPosts.tsx` | 58 | Vote/comment counts missing `font-variant-numeric: tabular-nums` |
| `pages/Home/components/PostFilter.tsx` | 170 | `<input>` missing `autocomplete="off"` |

---

## Summary

| Category | Count | Priority |
|----------|-------|----------|
| Icon component (systemic) | ~20 downstream | **High** |
| `<div>`/`<span>`/`<a href="#">` as buttons | 9 | **High** |
| Form inputs without labels | 11 | **High** |
| Destructive actions without confirmation | 3 | **High** |
| Missing `aria-label` on icon-only buttons | 8 | **High** |
| `transition: all` in SCSS | 6 | Medium |
| `:focus` instead of `:focus-visible` | 2 | Medium |
| Missing `autocomplete`/`type` on inputs | 5 | Medium |
| `...` instead of `…` (ellipsis) | 17 | Low |
| Semantic heading issues | 2 | Low |
| i18n gaps | 2 | Low |
| Animation `prefers-reduced-motion` | 3 | Low |
