# Product Requirements Document (PRD)
# Fider - User Feedback & Feature Request Portal

**Version:** 1.0  
**Date:** February 3, 2026  
**Status:** Current State Documentation

---

## 1. Executive Summary

### 1.1 Product Overview
Fider is a self-hosted, open-source feedback portal designed to help organizations collect, organize, and respond to customer feature requests and suggestions. The platform provides a centralized space where customers can share ideas, vote on suggestions, and engage in discussions about product direction.

### 1.2 Product Vision
"Give your customers a voice and let them tell you what they need. Spend less time guessing and more time building the right product."

### 1.3 Target Audience
- **Primary**: SaaS companies, product teams, and organizations seeking structured customer feedback
- **Secondary**: Open-source projects, community-driven initiatives, internal product teams
- **User Personas**:
  - End Users/Visitors: Submit and vote on suggestions
  - Collaborators: Moderate content, manage posts, assign tags
  - Administrators: Full control over site configuration, users, and integrations

### 1.4 Deployment Options
- **Cloud-hosted**: Fully managed service by Fider creators
- **Self-hosted**: Deploy on own infrastructure with full control

---

## 2. Core Features & Capabilities

### 2.1 User Management & Authentication

#### 2.1.1 User Roles
Three-tiered role system:
- **Visitor** (Role = 1): Basic role for all users
  - Submit posts and comments
  - Vote on posts
  - Requires moderation if not trusted
- **Collaborator** (Role = 2): Limited administrative access
  - All visitor permissions
  - Moderate content
  - Manage tags and post statuses
  - Pin posts and comments
  - Access to administrative pages
- **Administrator** (Role = 3): Full administrative access
  - All collaborator permissions
  - Manage users and roles
  - Configure site settings
  - Export data
  - Manage webhooks
  - Billing management (when enabled)

#### 2.1.2 User Status
- **Active** (1): Normal operating status
- **Deleted** (2): User chose to delete account
- **Blocked** (3): Blocked by staff members

#### 2.1.3 Authentication Methods
**Email-based Authentication:**
- Sign in with email and verification code
- Email verification for new users
- Password-less authentication flow

**OAuth Providers:**
- Facebook
- Google
- GitHub
- Configurable per tenant
- System-level and custom OAuth providers

**SAML/SSO:**
- Enterprise SSO support via SAML
- Configurable IdP integration
- Shibboleth/CAS support
- Login, ACS (Assertion Consumer Service), and metadata endpoints

**User Trust System:**
- Trusted users bypass moderation
- Administrators can trust/untrust users
- New visitors require moderation by default

#### 2.1.4 User Features
- Avatar support (Gravatar, letter avatars, custom uploads)
- User profile management
- Email change with verification
- API key generation for API access
- Account deletion

---

### 2.2 Post (Idea) Management

#### 2.2.1 Post Entity
Core fields:
- ID, Number (user-facing), Title, Slug
- Description (markdown supported)
- Author/User
- Vote count, Comment count
- Tags (array)
- Status
- Creation timestamp
- Approval status (for moderation)
- Pin status (pinned posts appear first)

#### 2.2.2 Post Statuses
- **Open** (0): Default status for new posts
- **Started** (1): Work in progress on the feature
- **Completed** (2): Feature implemented
- **Declined** (3): Request declined by team
- **Planned** (4): Accepted and on roadmap
- **Duplicate** (5): Duplicate of another post
- **Deleted** (6): Removed from site

#### 2.2.3 Post Operations
**For All Users:**
- Create new posts
- View posts (respects privacy settings)
- Vote on posts (toggle vote on/off)
- Subscribe to post updates
- Flag inappropriate posts

**For Collaborators & Admins:**
- Update post status with response
- Pin posts to top of lists
- Assign/unassign tags
- Mark as duplicate (links to original)
- Approve/decline flagged posts

**For Administrators Only:**
- Delete posts permanently
- Clear flags
- Export posts to CSV

#### 2.2.4 Post Features
- Markdown support for descriptions
- Duplicate detection via similarity search
- SEO-friendly URLs with slugs
- Pinning capability (by staff)
- Vote restrictions (can't vote on completed/declined/duplicate)
- RSS/Atom feed for posts

---

### 2.3 Comments & Discussions

#### 2.3.1 Comment Entity
Core fields:
- ID, Content (markdown)
- Author/User
- Timestamps (created, edited)
- Editor information (if edited)
- Reaction counts (emoji reactions)
- Attachments (file references)
- Approval status
- Flags count (for moderators)
- Pin status

#### 2.3.2 Comment Operations
**For Authenticated Users:**
- Post comments on posts
- Edit own comments
- Delete own comments
- Add emoji reactions to comments
- Flag inappropriate comments

**For Collaborators & Admins:**
- Pin comments
- View flag counts
- Approve/decline flagged comments
- Moderate content

**For Administrators:**
- Delete any comment
- Clear comment flags

#### 2.3.3 Comment Features
- Markdown support
- File attachments
- Emoji reactions (multiple types)
- Edit history tracking
- Pinned comments appear first
- Real-time reaction counts

---

### 2.4 Voting System

#### 2.4.1 Vote Entity
- User-to-Post relationship
- Vote count aggregation
- Voted status per user

#### 2.4.2 Voting Rules
- One vote per user per post
- Toggle vote on/off
- Cannot vote on completed/declined/duplicate posts
- Vote counts visible to all users
- Leaderboard tracking for top posts

#### 2.4.3 Vote Operations
- Add vote (POST)
- Remove vote (DELETE)
- Toggle vote (POST)
- List votes for a post (GET)

---

### 2.5 Tag System

#### 2.5.1 Tag Entity
Core fields:
- ID, Name, Slug
- Color (for visual distinction)
- Public/Private visibility flag

#### 2.5.2 Tag Operations
**For All Users:**
- View public tags
- Filter posts by tags

**For Collaborators:**
- Assign tags to posts
- Unassign tags from posts

**For Administrators:**
- Create new tags
- Edit tag properties
- Delete tags
- Set tag visibility (public/private)

#### 2.5.3 Tag Features
- Color-coded visual display
- Public tags visible to all
- Private tags for internal organization
- Multiple tags per post
- Tag-based filtering and search

---

### 2.6 Tenant (Site) Management

#### 2.6.1 Tenant Entity
Core configuration:
- Name, Subdomain, Custom CNAME
- Logo (blob storage)
- Welcome message and header
- Invitation text (call-to-action)
- Custom CSS
- Locale settings
- Status (Active, Pending, Disabled)

#### 2.6.2 Privacy & Access Control
- **Public sites**: Anyone can view
- **Private sites**: Requires authentication
- Prevention of search engine indexing option
- Email authentication toggle

#### 2.6.3 Tenant Features
- Multi-tenant architecture
- Subdomain-based isolation
- Custom domain (CNAME) support
- Custom branding (logo, CSS)
- Localization support (multiple languages)
- RSS/Atom feed toggle
- Allowed authentication schemes configuration

#### 2.6.4 Tenant Status
- **Active**: Normal operation
- **Pending**: Awaiting activation
- **Disabled**: Locked/suspended

---

### 2.7 Moderation System

#### 2.7.1 Content Moderation
**Automatic Moderation:**
- New visitors' posts/comments require approval
- Trusted users bypass moderation
- Flagging system for inappropriate content

**Manual Moderation:**
- Review queue for pending items
- Approve and verify user (grants trust)
- Decline and block user
- Approve without trust
- Decline without block

#### 2.7.2 Moderation Features
- Flag posts for review
- Flag comments for review
- Moderation count tracking
- Separate queues for posts and comments
- Bulk moderation actions
- User trust management

#### 2.7.3 Commercial Feature
Content moderation is a **commercial feature** requiring a license key:
- `COMMERCIAL_KEY` environment variable
- License validation system
- Graceful degradation if unlicensed

---

### 2.8 Notification System

#### 2.8.1 Notification Entity
Core fields:
- ID, Title, Link
- Read/unread status
- Created timestamp
- Author information
- Avatar data

#### 2.8.2 Notification Types
- New comments on subscribed posts
- Status changes on user's posts
- Responses from staff
- Mentions (via taggable users)

#### 2.8.3 Notification Operations
- List all notifications
- Get unread count
- Mark as read (individual)
- Mark all as read
- View notification details (redirects)

#### 2.8.4 Notification Features
- Real-time unread count
- Email notifications (configurable)
- In-app notification center
- Automatic subscription on vote/comment
- Manual subscription toggle

---

### 2.9 Search & Discovery

#### 2.9.1 Search Capabilities
- Full-text search on posts
- Similar post detection (duplicate prevention)
- Tag-based filtering
- Status-based filtering
- User-based filtering

#### 2.9.2 Search Features
- Noise word filtering
- Similarity scoring
- Search API endpoint
- Autocomplete suggestions

#### 2.9.3 Discovery Features
- Homepage post listing
- Status-based views
- Tag-based views
- Leaderboards (top posts, top users)
- RSS/Atom feeds

---

### 2.10 Webhooks & Integrations

#### 2.10.1 Webhook Entity
Core fields:
- ID, Name, Type
- URL, HTTP Method
- Custom headers (key-value map)
- Content template
- Status (Enabled/Disabled)

#### 2.10.2 Webhook Types
Configurable triggers for various events:
- New post created
- Post status changed
- New comment added
- Post voted
- User invited
- (Extensible event system)

#### 2.10.3 Webhook Operations
**For Administrators:**
- Create webhooks
- Update webhook configuration
- Delete webhooks
- Test webhook (send test payload)
- Preview webhook payload
- Get webhook props by type

#### 2.10.4 Webhook Features
- Custom HTTP headers
- Template-based payload customization
- Event filtering
- Status management (enable/disable)
- Test/preview functionality
- Delivery tracking

#### 2.10.5 Incoming Webhooks
- **Stripe webhooks**: For billing events (when billing enabled)
- Webhook signature verification
- Event processing

---

### 2.11 Billing & Subscriptions

#### 2.11.1 Billing Features (Commercial)
Available when `BILLING_ENABLED`:
- Stripe integration
- Customer portal access
- Checkout session creation
- Subscription management
- Plan selection
- Payment method management

#### 2.11.2 Billing Operations
**For Administrators:**
- Access billing page
- Create Stripe portal session
- Create Stripe checkout session
- View current plan status

#### 2.11.3 Commercial Features
Unlocked with valid license:
- Content moderation
- Advanced features (future)

---

### 2.12 Email System

#### 2.12.1 Email Providers
**Mailgun:**
- API-based sending
- Regional support (US/EU)
- Configuration via environment

**SMTP:**
- Generic SMTP support
- Configurable host, port, credentials
- TLS/SSL support

#### 2.12.2 Email Types
- Verification emails (sign-up, sign-in)
- Email change verification
- Invitation emails
- Notification emails
- Sample invite emails

#### 2.12.3 Email Features
- Template-based emails
- Branded emails (tenant logo/colors)
- No-reply sender configuration
- Multi-language support
- MailHog support (development)

---

### 2.13 API

#### 2.13.1 API Structure
**Public API** (`/api/v1/*`):
- No authentication required
- Read-only operations
- Search posts, get post details
- List tags, comments, votes
- Leaderboards

**Members API** (`/api/v1/*`):
- Authentication required
- Create/update posts
- Add/edit/delete comments
- Vote operations
- Subscription management
- Flag content

**Staff API** (`/api/v1/*`):
- Collaborator/Admin only
- List users
- Send invitations
- Pin posts/comments
- Tag assignment
- Flag management

**Admin API** (`/api/v1/*`):
- Administrator only
- User management
- Tag management
- Delete operations
- Moderation actions
- Clear flags

#### 2.13.2 API Features
- RESTful design
- JSON responses
- API key authentication (user-specific)
- Rate limiting (via middleware)
- CORS support
- Comprehensive error handling

---

### 2.14 User Interface

#### 2.14.1 Frontend Technology
- React 18 with TypeScript
- Server-side rendering (SSR)
- Lazy-loaded page components
- SCSS with utility classes (BEM-style)
- LinguiJS for internationalization

#### 2.14.2 Key Pages
**Public Pages:**
- Home (post listing with filters)
- Post details
- Leaderboard (top ideas, top users)
- Sign in/Sign up
- Legal (Terms, Privacy)

**Authenticated Pages:**
- User settings
- Notifications center
- Change email verification

**Administrative Pages:**
- General settings
- Advanced settings
- Privacy settings
- User management
- Tag management
- Authentication management
- Webhooks
- Invitations
- Content moderation
- Flagged posts/comments
- Export (CSV, backup)
- Billing (when enabled)

**Special Pages:**
- Design system preview
- OAuth callback/echo pages
- SAML endpoints
- Error pages (401, 403, 404, 410, 500)
- Maintenance page

#### 2.14.3 UI Features
- Responsive design
- Dark mode support (via custom CSS)
- Markdown rendering
- Real-time updates
- Keyboard shortcuts
- Accessibility features
- Progressive enhancement

#### 2.14.4 Design System
- Reusable component library
- Utility-first CSS classes
- Consistent spacing/typography
- Icon system (Heroicons)
- Illustration library (Undraw)

---

## 3. Technical Architecture

### 3.1 Backend Stack
- **Language**: Go 1.22+
- **Architecture**: CQRS pattern (Command Query Responsibility Segregation)
- **Service Layer**: Bus-based dependency injection
- **Database**: PostgreSQL
- **Migrations**: SQL-based (numbered files)

### 3.2 Frontend Stack
- **Framework**: React 18
- **Language**: TypeScript
- **Bundler**: Webpack
- **SSR**: Custom server-side rendering
- **Styling**: SCSS with utility classes
- **i18n**: LinguiJS

### 3.3 Infrastructure
- **Deployment**: Docker support
- **Database**: PostgreSQL 13+
- **Storage**: Blob storage for avatars, images, attachments
- **Email**: SMTP or Mailgun
- **Development**: Air (hot reload), Webpack watch

### 3.4 Code Organization

#### Backend Structure
```
app/
├── cmd/            # Application entry point, routes
├── handlers/       # HTTP request handlers
│   ├── apiv1/      # API v1 handlers
│   └── webhooks/   # Webhook handlers
├── models/         # Data models
│   ├── entity/     # Database entities
│   ├── cmd/        # Commands
│   ├── query/      # Queries
│   ├── action/     # User input actions
│   ├── dto/        # Data transfer objects
│   └── enum/       # Enumerations
├── services/       # Business logic
│   ├── sqlstore/   # PostgreSQL implementation
│   ├── oauth/      # OAuth services
│   └── saml/       # SAML services
├── middlewares/    # HTTP middlewares
└── pkg/            # Shared packages
    ├── bus/        # Service registry
    ├── web/        # Web framework
    ├── dbx/        # Database utilities
    └── env/        # Environment configuration
```

#### Frontend Structure
```
public/
├── pages/          # Page components (lazy-loaded)
│   ├── Home/
│   ├── ShowPost/
│   ├── Administration/
│   └── ...
├── components/     # Reusable components
├── services/       # API clients
├── hooks/          # Custom React hooks
├── models/         # TypeScript interfaces
└── assets/
    ├── styles/     # SCSS files
    │   └── utility/  # Utility classes
    └── images/     # Icons, illustrations
```

### 3.5 Database Schema (Key Tables)
- `tenants` - Site/tenant configuration
- `users` - User accounts
- `posts` - Feature requests/ideas
- `comments` - Post comments
- `tags` - Organizational tags
- `votes` - User votes on posts
- `notifications` - User notifications
- `webhooks` - Outbound webhook configurations
- `user_providers` - OAuth/SAML provider links

### 3.6 Security
- JWT-based sessions
- CSRF protection (middleware)
- CORS configuration
- Role-based access control (RBAC)
- Content sanitization
- SQL injection protection (parameterized queries)
- XSS prevention (markdown sanitization)
- Secure headers (via middleware)
- SAML/OAuth security flows

---

## 4. Configuration & Environment

### 4.1 Core Configuration
```bash
BASE_URL              # Application base URL
GO_ENV                # Environment (development/production)
DATABASE_URL          # PostgreSQL connection string
JWT_SECRET            # Secret for JWT signing
```

### 4.2 Logging
```bash
LOG_LEVEL             # DEBUG, INFO, WARN, ERROR
LOG_CONSOLE           # Console logging toggle
LOG_SQL               # SQL query logging
LOG_FILE              # File logging toggle
LOG_FILE_OUTPUT       # Log file path
```

### 4.3 Authentication
```bash
# OAuth Providers
OAUTH_FACEBOOK_APPID
OAUTH_FACEBOOK_SECRET
OAUTH_GOOGLE_CLIENTID
OAUTH_GOOGLE_SECRET
OAUTH_GITHUB_CLIENTID
OAUTH_GITHUB_SECRET

# SAML/SSO
SAML_ENTITY_ID
SAML_IDP_ENTITY_ID
SAML_IDP_SSO_URL
SAML_IDP_CERT
SAML_SP_CERT_PATH
SAML_SP_KEY_PATH
```

### 4.4 Email
```bash
EMAIL_NOREPLY         # Sender email address

# Mailgun
EMAIL_MAILGUN_API
EMAIL_MAILGUN_DOMAIN
EMAIL_MAILGUN_REGION

# SMTP
EMAIL_SMTP_HOST
EMAIL_SMTP_PORT
EMAIL_SMTP_USERNAME
EMAIL_SMTP_PASSWORD
```

### 4.5 Commercial Features
```bash
COMMERCIAL_KEY        # License key for commercial features
LICENSE_PUBLIC_KEY    # Public key for verification (auto-updated)
LICENSE_PRIVATE_KEY   # Private key (hosted platform only)
```

### 4.6 Maintenance Mode
```bash
MAINTENANCE           # Enable maintenance mode
MAINTENANCE_MESSAGE   # Custom maintenance message
MAINTENANCE_UNTIL     # Expected uptime message
```

---

## 5. Development Workflow

### 5.1 Local Development
```bash
make watch            # Hot reload (server + UI)
make migrate          # Run database migrations
make build            # Production build
```

### 5.2 Testing
```bash
make test             # All tests (server + UI)
make test-server      # Go tests
make test-ui          # Jest tests
make coverage-server  # With coverage
make test-e2e-ui      # E2E tests (Cucumber)
```

### 5.3 Linting
```bash
make lint             # Lint all code
make lint-server      # golangci-lint
make lint-ui          # ESLint
```

### 5.4 Localization
```bash
make locale-extract   # Extract translations
make locale-reset     # Reset specific keys
```

### 5.5 Docker Environment
- PostgreSQL on port 5432
- Application on port 3000
- MailHog on port 8025 (email preview)

---

## 6. Key User Flows

### 6.1 New User Registration
1. User visits signup page
2. User provides email and name
3. System sends verification email
4. User clicks verification link
5. Account activated with Visitor role

### 6.2 Submit Feature Request
1. User clicks "Share Feedback" button
2. User enters title and description
3. System suggests similar posts (duplicate detection)
4. User confirms uniqueness
5. User optionally adds tags (if collaborator)
6. Post created (pending approval if user not trusted)
7. User automatically subscribed to post

### 6.3 Vote on Post
1. User views post
2. User clicks vote button
3. Vote count increments
4. User subscribed to post (receives notifications)
5. User can toggle vote off

### 6.4 Staff Response Workflow
1. Collaborator reviews post
2. Collaborator sets status (Started/Planned/etc.)
3. Collaborator adds response text
4. Response saved with timestamp
5. All subscribers notified

### 6.5 Moderation Workflow
1. New post from untrusted user appears in queue
2. Moderator reviews post
3. Moderator approves and verifies user (grants trust), OR
4. Moderator declines and blocks user
5. User marked as trusted/blocked
6. Future posts auto-approved/auto-blocked

---

## 7. Internationalization

### 7.1 Supported Languages
- English (default)
- Arabic (ar)
- Czech (cs)
- German (de)
- Greek (el)
- Spanish (es-ES)
- Persian (fa)
- French (fr)
- Italian (it)
- Japanese (ja)
- Korean (ko)
- Dutch (nl)
- Polish (pl)
- Portuguese Brazil (pt-BR)
- Russian (ru)
- Sinhala Sri Lanka (si-LK)
- Slovak (sk)
- Swedish (sv-SE)
- Turkish (tr)
- Chinese China (zh-CN)

### 7.2 Translation System
- LinguiJS framework
- JSON-based translation files
- Per-tenant locale configuration
- Automatic language detection
- Fallback to English
- Admin pages forced to English (temporary)

---

## 8. Performance & Scalability

### 8.1 Performance Features
- Server-side rendering (SSR) for fast initial load
- Code splitting (lazy-loaded pages)
- Asset caching (CDN-ready)
- Database indexing
- Query optimization
- Gzip compression

### 8.2 Caching Strategy
- Static assets: 365 days
- Tenant assets: 30 days
- Avatars: 5 days
- Feed: 5 minutes
- Custom CSS: Versioned with MD5

### 8.3 Scalability Considerations
- Multi-tenant architecture
- Horizontal scaling capability
- Database connection pooling
- Stateless application design
- CDN integration for assets

---

## 9. Limitations & Constraints

### 9.1 Current Limitations
- Single database (PostgreSQL only)
- No real-time WebSocket updates
- Limited file attachment types
- Manual content moderation (no AI)
- English-only admin interface (temporary)

### 9.2 Known Constraints
- Requires PostgreSQL 13+
- Requires Go 1.22+
- Requires Node 21/22
- No Windows native support (WSL recommended)

---

## 10. Success Metrics

### 10.1 User Engagement
- Post creation rate
- Voting activity
- Comment engagement
- User retention
- Session duration

### 10.2 Content Quality
- Post approval rate
- Duplicate detection accuracy
- Moderation queue size
- Flag resolution time

### 10.3 Platform Health
- Response time (API)
- Uptime percentage
- Error rate
- Database query performance

---

## 11. Future Enhancements

### 11.1 Planned Features
- Real-time updates (WebSockets)
- Advanced analytics dashboard
- AI-powered duplicate detection
- Roadmap visualization
- Mobile applications
- Advanced role permissions
- Custom fields for posts
- Integration marketplace

### 11.2 Under Consideration
- Multi-language admin interface
- Advanced reporting
- Custom workflows
- SLA tracking
- Customer segmentation
- A/B testing support

---

## 12. Support & Resources

### 12.1 Documentation
- Official docs: https://docs.fider.io
- Self-hosted guide: https://docs.fider.io/self-hosted/
- Contributing guide: CONTRIBUTING.md

### 12.2 Community
- Demo site: https://demo.fider.io
- Feedback site: https://feedback.fider.io
- GitHub: https://github.com/getfider/fider

### 12.3 Commercial Support
- Hosted platform: https://fider.io
- Enterprise support available
- Custom development services

---

## 13. Compliance & Legal

### 13.1 Privacy
- GDPR compliance ready
- User data export
- Account deletion
- Consent management
- Privacy policy support

### 13.2 Terms of Service
- Customizable per tenant
- Legal page framework
- Content moderation policies

### 13.3 Data Retention
- User data retained until account deletion
- Post history maintained
- Audit trail for administrative actions
- Export capabilities (CSV, ZIP backup)

---

## 14. Appendix

### 14.1 Glossary
- **Post**: User-submitted feature request or suggestion (also called "idea")
- **Tenant**: Individual site/organization instance
- **Collaborator**: Moderator with limited admin access
- **Visitor**: Regular authenticated user
- **Bus**: Service registry and command/query dispatch system
- **CQRS**: Command Query Responsibility Segregation pattern

### 14.2 Model Naming Conventions
- `entity.*` - Database table entities
- `cmd.*` - Commands that execute actions (may return values)
- `query.*` - Read-only data queries (with Result field)
- `action.*` - User input for POST/PUT/PATCH (maps to commands)
- `dto.*` - Data transfer objects between packages
- `enum.*` - Enumeration types and constants

### 14.3 CSS Naming Conventions
- `#p-<page>` - Page ID
- `.c-<component>` - Component class
- `.c-<component>__<element>` - Component element (BEM)
- `.c-<component>--<modifier>` - Component modifier (BEM)
- Utility classes (no prefix) - Reusable utilities (similar to Tailwind)

---

**Document End**

*This PRD reflects the current state of the Fider application as of February 3, 2026. For the latest updates, refer to the official documentation and repository.*
