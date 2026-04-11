# Build Phases — The Sanctuary Lodge

A step-by-step plan for building the Vue 3 lodge management website.
Each phase is self-contained and shippable before moving to the next.

---

## Phase 1 — Project Scaffold & Design Foundation

Stand up the project with all tooling and design tokens in place before writing a single feature.

- [ ] Init Vue 3 project with Vite (`npm create vue@latest`)
- [ ] Install and configure Tailwind CSS v4
- [ ] Wire up OKLCH color tokens in `tailwind.config.js`
  - `primary` — `oklch(0.55 0.12 45)` warm terracotta
  - `secondary` — `oklch(0.94 0.02 60)` cream background
  - `accent` — `oklch(0.45 0.08 145)` muted forest green
  - Full surface scale, text tokens, outline, error
- [ ] Add Google Fonts — Noto Serif + Manrope
- [ ] Add Material Symbols Outlined icon font
- [ ] Install Vue Router 4 and Pinia
- [ ] Install Axios
- [ ] Set up base folder structure (`views/`, `components/`, `stores/`, `composables/`)
- [ ] Create global CSS reset and base styles (scroll-smooth, font defaults)
- [ ] Create `BaseButton.vue`, `BaseInput.vue`, `BaseCard.vue` — core UI primitives
- [ ] Create `useScrollReveal.js` composable (Intersection Observer)

**Exit criteria:** App runs, color tokens render correctly, base components are styled to spec.

---

## Phase 2 — Layout Shell & Navigation

Shared chrome that wraps every page.

- [x] `AppNavbar.vue` — glassmorphic sticky nav
  - Logo (italic serif)
  - Nav links: Rooms / Experiences / About / Concierge
  - Sign In button + Reserve button (primary gradient)
  - Collapses to hamburger on mobile
  - User avatar + initials with dropdown (My Bookings / Sign Out) when authenticated
- [x] `MobileBottomNav.vue` — icon tab bar (Explore / Bookings / Saved / Profile)
- [x] `AppFooter.vue` — multi-column links + copyright
- [x] `AppLayout.vue` — default layout slot wrapping navbar + footer
- [x] `AuthLayout.vue` — minimal layout for login/register (no footer nav clutter)
- [x] Wire up router with placeholder views for all routes

**Exit criteria:** Can navigate between all routes, nav and footer render on every page.

---

## Phase 3 — Home / Landing Page

The first impression. Closest to a complete HTML sample in the docs.

- [x] `HeroSection.vue` — full-height image, serif headline, italic accent word, gradient overlay
- [x] `SearchBar.vue` — overlapping widget (-120px margin-top editorial overlap)
  - Fields: Destination, Check-in/out, Guests ± stepper
  - Glassmorphic card, Material icons per field
  - Pushes query params to /rooms on submit
- [x] `PropertyGrid.vue` — 3-col grid of lodge/room cards
  - 4:5 aspect ratio images, hover scale 1.05 over 1000ms
  - Serif property name, star rating + review count, price on right
  - Glassmorphic location badge (top-left)
  - Scroll-reveal staggered entrance
- [x] `WhyChooseUs.vue` — asymmetric bento layout (5-col text + 7-col 2×2 offset image grid)
- [x] `NewsletterCta.vue` — centered section, bottom-border email input, success state
- [x] `HomeView.vue` — thin composition of all home sections

**Exit criteria:** Landing page matches the `the_sanctuary_landing_page_updated` HTML sample.

---

## Phase 4 — Authentication

Gate for all booking features.

- [x] Pinia `auth` store — `user`, `token`, `isAuthenticated`, `login()`, `logout()`, `register()`, `fetchUser()`
- [x] `LoginView.vue` — email + password, inline validation, show/hide toggle, animated error banner
- [x] `RegisterView.vue` — name, email, password + confirm, live password strength meter, inline validation
- [x] Router navigation guards — unauthenticated → login, authenticated → redirect away from login/register
- [x] Persist auth token + user to `localStorage`; `fetchUser()` rehydrates on boot
- [x] User avatar with initials + dropdown in navbar (done in Phase 2)

**Exit criteria:** Can register, log in, log out. Protected routes redirect to login.

---

## Phase 5 — Rooms Browsing

Let guests explore what's available.

- [x] `RoomsView.vue` — filterable grid of all room types
  - Filter bar: room type chips, capacity slider, max price slider, available-only toggle
  - Reads `?guests=` query param from SearchBar
  - Reactive filtered list with empty state
  - Responsive grid (1 col mobile → 3 col desktop)
- [x] `RoomCard.vue` — image with overlay on unavailable, capacity badge, amenity icon row, price
- [x] `RoomDetailView.vue` — full room page
  - 3-image gallery with prev/next arrows, dot indicators, thumbnail strip, fade transition
  - Room description, amenities grid (icon + label)
  - Capacity, bed type, size, location
  - Sticky sidebar: date pickers, guest stepper, meal plan radio selector, live price breakdown, Reserve CTA with date validation
- [x] `usePricing.js` composable — reactive `{ nightCount, baseTotal, mealCost, taxes, grandTotal }`

**Exit criteria:** Can browse rooms, view detail, see a live price estimate update as dates/guests change.

---

## Phase 6 — Reservation Flow

The core booking experience — 3-section form with sticky summary.

- [x] Pinia `booking` store — holds draft booking state
  - `roomId`, `roomType`, `checkIn`, `checkOut`, `guestCount`
  - `mealPlan` — `'full_board' | 'half_board' | 'breakfast' | 'none'`
  - `guestInfo` — `{ firstName, lastName, email, phone, nationality, passportId }`
  - `specialRequests`
- [x] `ReservationView.vue` — two-column layout (form left, summary right)
  - Staggered section fadeIn animations matching HTML sample
  - Prefills email/name from logged-in user on mount
  - Success overlay before redirect to bookings
- [x] `GuestInfoForm.vue` — Section 01: numbered header, white card panel, all fields
- [x] `StayDetailsForm.vue` — Section 02: room type (locked display), guest stepper, date pickers with validation
- [x] `PreferencesForm.vue` — Section 03: composes MealPlanSelector + special requests
- [x] `MealPlanSelector.vue` — styled radio cards with icon, description, per-person rate
- [x] `ReservationSummary.vue` — sticky panel: room image with inner-glow, formatted dates, full price breakdown, Platinum Guarantee badge
- [x] Form validation — required fields, email format, date logic (check-out after check-in)
- [x] Submit → POST to API → success overlay → redirect to BookingsView

**Exit criteria:** Full booking flow works end-to-end. Summary updates live as form changes.

---

## Phase 7 — My Bookings Dashboard

Authenticated user's booking history and management.

- [ ] Pinia `reservations` store — `active[]`, `past[]`, `loading`, `fetchAll()`, `cancel(id)`, `rebook(id)`
- [ ] `BookingsView.vue` — hero + two sections
- [ ] `BookingCard.vue` — active reservation card
  - Room image, property name, dates, guest count, meal plan
  - Status badge (accent green = confirmed, amber = pending)
  - "View Details" and "Cancel" action buttons
  - Hover: card scale 1.01 over 700ms, image zoom 1.1 over 1000ms
- [ ] `PastBookingCard.vue` — grayscale on load, full colour on hover
  - Verification badge, "Rebook" button
- [ ] Cancel reservation confirmation modal
- [ ] Empty state when no bookings exist

**Exit criteria:** Dashboard matches `my_bookings_updated` HTML sample. Can cancel and rebook.

---

## Phase 8 — Polish & Production Readiness

Tie up loose ends before deployment.

- [ ] Dark mode — audit all components for `dark:` Tailwind classes
- [ ] Full responsive QA — mobile (375px), tablet (768px), desktop (1280px+)
- [ ] Accessibility audit — semantic HTML, ARIA labels, keyboard navigation, focus states
- [ ] Loading skeletons for async data (rooms grid, bookings dashboard)
- [ ] Error states — API failures, empty search results, 404 page
- [ ] Page `<title>` and meta tags per route (vue-meta or `useHead`)
- [ ] Image optimisation — lazy loading, `srcset` for responsive images
- [ ] Environment variables — `VITE_API_BASE_URL`
- [ ] Production build test (`npm run build && npm run preview`)
- [ ] Deploy

**Exit criteria:** Lighthouse scores ≥ 90 performance, 100 accessibility. Zero console errors in prod build.

---

## Route Map (for reference)

| Path | View | Auth required |
|---|---|---|
| `/` | `HomeView` | No |
| `/rooms` | `RoomsView` | No |
| `/rooms/:id` | `RoomDetailView` | No |
| `/reserve/:roomId` | `ReservationView` | Yes |
| `/bookings` | `BookingsView` | Yes |
| `/login` | `LoginView` | No |
| `/register` | `RegisterView` | No |

---

## Design Rules (quick reference)

| Rule | Do | Don't |
|---|---|---|
| Borders | Background color shifts | 1px solid lines |
| Text | `on-surface oklch(0.18 0.02 45)` | Pure `#000000` |
| Radius | `4px` default, `8px` large | `0px` or `999px` |
| CTAs | Primary gradient terracotta | Flat solid fills |
| Active / confirmed | Accent forest green | Bright `#00ff00` |
| Inputs | Bottom border only, transparent bg | Boxed input fields |
| Shadows | Diffused soft `oklch(.../ 0.06)` | Hard drop shadows |
| Dividers | Never inside cards/lists | — |
| Glassmorphism | 80% opacity + `backdrop-blur-xl` | Solid opaque nav |
